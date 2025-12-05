# 即时通讯(IM)系统设计

## 一、核心功能

| 功能 | 说明 |
|------|------|
| 单聊 | 一对一消息 |
| 群聊 | 多人消息（扩散问题） |
| 消息可靠投递 | 不丢、不重、有序 |
| 已读回执 | 对方是否已读 |
| 在线状态 | 用户是否在线 |
| 历史消息 | 消息漫游、多端同步 |

---

## 二、整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                      客户端 (App/Web)                        │
└─────────────────────────────┬───────────────────────────────┘
                              │ WebSocket/TCP 长连接
┌─────────────────────────────▼───────────────────────────────┐
│                     接入层 (Gateway)                         │
│           管理长连接、协议解析、心跳、路由                       │
└─────────────────────────────┬───────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────┐
│                     逻辑层 (Logic)                           │
│           消息处理、鉴权、群组管理、在线状态                     │
└─────────────────────────────┬───────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
┌───────▼───────┐     ┌───────▼───────┐     ┌───────▼───────┐
│    Kafka      │     │    Redis      │     │    MySQL      │
│  消息队列      │     │ 在线状态/路由  │     │  消息持久化    │
└───────────────┘     └───────────────┘     └───────────────┘
```

---

## 三、长连接管理

### 连接建立

```go
// Gateway 服务
func HandleWebSocket(conn *websocket.Conn, userID int64) {
    // 1. 注册连接
    connManager.Add(userID, conn)
    
    // 2. 记录路由信息（用户在哪台Gateway）
    redis.HSet("user:gateway", userID, gatewayAddr)
    redis.SAdd("online:users", userID)
    
    // 3. 心跳检测
    go heartbeat(conn, userID)
    
    // 4. 读取消息
    for {
        msg := conn.ReadMessage()
        handleMessage(userID, msg)
    }
}

// 心跳机制
func heartbeat(conn *websocket.Conn, userID int64) {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        if err := conn.Ping(); err != nil {
            // 连接断开，清理
            connManager.Remove(userID)
            redis.HDel("user:gateway", userID)
            redis.SRem("online:users", userID)
            return
        }
    }
}
```

---

## 四、消息投递模型

### 单聊消息流程

```
发送方 → Gateway A → Logic → 查询接收方Gateway → Gateway B → 接收方

1. 客户端发送消息到 Gateway A
2. Gateway A 转发到 Logic 服务
3. Logic 生成消息ID，存储到 MySQL
4. Logic 查询 Redis 获取接收方所在 Gateway B
5. Logic 推送消息到 Gateway B
6. Gateway B 通过长连接推送给接收方
7. 接收方 ACK 确认
```

```go
func SendMessage(from, to int64, content string) error {
    // 1. 生成消息ID（保证有序）
    msgID := snowflake.NextID()
    
    // 2. 存储消息
    msg := &Message{
        ID:        msgID,
        FromUser:  from,
        ToUser:    to,
        Content:   content,
        Timestamp: time.Now(),
        Status:    "sent",
    }
    db.Create(msg)
    
    // 3. 查询接收方是否在线
    gatewayAddr := redis.HGet("user:gateway", to)
    
    if gatewayAddr != "" {
        // 在线：直接推送
        pushToGateway(gatewayAddr, to, msg)
    } else {
        // 离线：存入离线消息队列
        redis.LPush("offline:msg:"+to, msg)
    }
    
    return nil
}
```

### 群聊消息流程（写扩散 vs 读扩散）

| 方案 | 实现 | 适用场景 |
|------|------|----------|
| **写扩散** | 消息写入每个群成员的收件箱 | 小群（<500人） |
| **读扩散** | 消息只写一份，读取时拉取 | 大群（>500人） |

```go
// 写扩散（小群）
func SendGroupMessage(groupID int64, from int64, content string) {
    members := getGroupMembers(groupID)
    for _, memberID := range members {
        // 写入每个成员的消息队列
        redis.LPush("inbox:"+memberID, msg)
        // 推送在线成员
        if isOnline(memberID) {
            pushMessage(memberID, msg)
        }
    }
}

// 读扩散（大群）
func SendGroupMessageLarge(groupID int64, from int64, content string) {
    // 只写一份到群消息表
    db.Create(&GroupMessage{GroupID: groupID, ...})
    // 通知在线成员有新消息
    notifyGroupMembers(groupID, "new_message")
}
```

---

## 五、消息可靠性

### 1. 消息不丢

```
客户端 → 服务端：
  发送消息 → 服务端存储 → 返回 ACK → 客户端标记已发送
  超时未收到 ACK → 重试

服务端 → 客户端：
  推送消息 → 客户端 ACK → 服务端标记已送达
  超时未收到 ACK → 重推
```

### 2. 消息不重（幂等）

```go
// 客户端生成 clientMsgID，服务端去重
func SendMessage(clientMsgID string, msg *Message) error {
    // 检查是否已处理
    if redis.SetNX("msg:dedup:"+clientMsgID, 1, 24*time.Hour) == false {
        return nil  // 重复消息，忽略
    }
    // 正常处理...
}
```

### 3. 消息有序

```go
// 全局 ID (Snowflake) 只能保证趋势递增，无法保证连续性（无法检测丢消息）
// 优化方案：会话级序列号 (Seq ID)

type Conversation struct {
    UserA   int64
    UserB   int64
    LastSeq int64  // 当前会话最大序号，严格递增：1, 2, 3...
}

// 发送时
func SendMessage(from, to int64, content string) {
    // 利用 Redis 原子自增生成 Seq
    seq := redis.Incr("seq:" + conversationID)
    msg.Seq = seq
    // ...
}

// 接收端（补洞逻辑）
// 收到 Seq=1, Seq=3，发现少了 Seq=2，主动拉取
if msg.Seq > localMaxSeq + 1 {
    pullMissingMessages(localMaxSeq + 1, msg.Seq - 1)
}
localMaxSeq = msg.Seq
```

### 4. 大群消息风暴优化

```
场景：万人群，一条消息触发 10000 次推送，网关压力极大。

优化策略：
1. 消息合并：服务端将 100ms 内的多条消息合并为一个包推送。
2. 智能降频：群消息过快时，对非活跃用户只推“有新消息”通知，不推具体内容，让用户主动拉取。
```

---

## 六、已读回执

```go
// 发送已读回执
func SendReadReceipt(userID, conversationID, lastReadMsgID int64) {
    // 更新已读位置
    redis.HSet("read:position:"+conversationID, userID, lastReadMsgID)
    
    // 通知对方
    otherUser := getOtherUser(conversationID, userID)
    pushMessage(otherUser, &ReadReceipt{
        ConversationID: conversationID,
        ReaderID:       userID,
        LastReadMsgID:  lastReadMsgID,
    })
}

// 查询未读数
func GetUnreadCount(userID, conversationID int64) int {
    lastReadMsgID := redis.HGet("read:position:"+conversationID, userID)
    return db.Model(&Message{}).
        Where("conversation_id = ? AND id > ?", conversationID, lastReadMsgID).
        Count()
}
```

---

## 七、在线状态

```go
// 查询在线状态
func IsOnline(userID int64) bool {
    return redis.SIsMember("online:users", userID)
}

// 批量查询好友在线状态
func GetFriendsOnlineStatus(userID int64) map[int64]bool {
    friends := getFriends(userID)
    result := make(map[int64]bool)
    for _, friendID := range friends {
        result[friendID] = redis.SIsMember("online:users", friendID)
    }
    return result
}

// 状态变更通知（上线/下线）
func NotifyStatusChange(userID int64, online bool) {
    friends := getFriends(userID)
    for _, friendID := range friends {
        if isOnline(friendID) {
            pushMessage(friendID, &StatusChange{UserID: userID, Online: online})
        }
    }
}
```

---

## 八、多端同步

```go
// 用户可能多端登录（手机、电脑、iPad）
// 每个设备一个连接

func HandleMultiDevice(userID int64, deviceID string, conn *websocket.Conn) {
    // 记录：用户 → 多个设备
    redis.HSet("user:devices:"+userID, deviceID, gatewayAddr)
}

// 消息同步到所有设备
func SyncToAllDevices(userID int64, msg *Message) {
    devices := redis.HGetAll("user:devices:" + userID)
    for deviceID, gatewayAddr := range devices {
        pushToGateway(gatewayAddr, userID, deviceID, msg)
    }
}
```

---

## 九、数据模型

```sql
-- 消息表（按会话分表）
CREATE TABLE messages (
    id BIGINT PRIMARY KEY,
    conversation_id BIGINT NOT NULL,
    from_user BIGINT NOT NULL,
    to_user BIGINT,
    group_id BIGINT,
    content TEXT,
    msg_type TINYINT,  -- 1文本 2图片 3语音
    created_at TIMESTAMP,
    INDEX idx_conversation (conversation_id, id)
);

-- 会话表
CREATE TABLE conversations (
    id BIGINT PRIMARY KEY,
    type TINYINT,  -- 1单聊 2群聊
    last_msg_id BIGINT,
    updated_at TIMESTAMP
);

-- 用户会话关系
CREATE TABLE user_conversations (
    user_id BIGINT,
    conversation_id BIGINT,
    last_read_msg_id BIGINT,
    PRIMARY KEY (user_id, conversation_id)
);
```

---

## 十、面试追问

| 问题 | 回答 |
|------|------|
| 如何保证消息不丢？ | ACK 机制 + 重试 + 离线消息存储 |
| 如何保证消息有序？ | Snowflake ID + 会话级 seq |
| 群聊如何扩散？ | 小群写扩散，大群读扩散 |
| 如何处理离线消息？ | Redis 队列存储，上线后拉取 |
| 长连接如何保活？ | 心跳包（30s），超时断开重连 |
| 如何支持百万连接？ | 多 Gateway 水平扩展，单机优化（epoll） |
