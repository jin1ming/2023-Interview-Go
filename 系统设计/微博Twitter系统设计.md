# 微博/Twitter 系统设计

## 一、核心功能

| 功能 | 说明 |
|------|------|
| 发推/发微博 | 文字、图片、视频、转发 |
| 关注/粉丝 | 单向关注关系 |
| Feed 流 | 查看关注人的动态 |
| 互动 | 点赞、评论、转发 |
| 热搜 | 热门话题排行 |

---

## 二、与朋友圈的核心区别

| 对比项 | 微信朋友圈 | 微博/Twitter |
|--------|-----------|--------------|
| 关系 | 双向好友（上限5000） | 单向关注（无上限） |
| 大V问题 | 无 | 有（千万粉丝） |
| Feed模式 | 纯推模式 | **推拉结合** |
| 可见性 | 仅好友可见 | 公开 |

---

## 三、Feed 流设计（核心难点）

### 推拉结合方案

```
普通用户（粉丝 < 1000）：推模式
  发推 → 写入所有粉丝的收件箱

大V用户（粉丝 > 1000）：拉模式
  发推 → 只写自己的发件箱
  粉丝刷Feed时 → 拉取关注的大V动态 + 自己收件箱

读取时：
  收件箱（已推送）+ 拉取大V发件箱 → 合并排序
```

### 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                        发布微博                              │
└─────────────────────────────┬───────────────────────────────┘
                              │
                    ┌─────────▼─────────┐
                    │   判断用户类型     │
                    └─────────┬─────────┘
                              │
          ┌───────────────────┴───────────────────┐
          │                                       │
    ┌─────▼─────┐                           ┌─────▼─────┐
    │  普通用户  │                           │   大V     │
    │  推模式   │                           │  拉模式   │
    └─────┬─────┘                           └─────┬─────┘
          │                                       │
          ▼                                       ▼
    写入粉丝收件箱                            只写发件箱
    (Kafka异步扩散)                          (读时拉取)
```

### 代码实现

```go
const BigVThreshold = 1000  // 大V阈值

func PublishTweet(userID int64, content string) error {
    // 1. 写入发件箱（所有用户都写）
    tweet := &Tweet{
        ID:        snowflake.NextID(),
        UserID:    userID,
        Content:   content,
        CreatedAt: time.Now(),
    }
    db.Create(tweet)
    redis.ZAdd("outbox:"+userID, tweet.CreatedAt.Unix(), tweet.ID)
    
    // 2. 判断是否大V
    followerCount := redis.SCard("followers:" + userID)
    
    if followerCount < BigVThreshold {
        // 普通用户：推模式，异步写入粉丝收件箱
        kafka.Send("tweet-fanout", &FanoutMessage{
            TweetID:   tweet.ID,
            AuthorID:  userID,
            CreatedAt: tweet.CreatedAt,
        })
    }
    // 大V不推送，粉丝读时拉取
    
    return nil
}

// 消费者：写扩散
func ConsumeFanout(msg *FanoutMessage) {
    followers := redis.SMembers("followers:" + msg.AuthorID)
    for _, followerID := range followers {
        redis.ZAdd("inbox:"+followerID, msg.CreatedAt.Unix(), msg.TweetID)
        // 控制收件箱大小
        redis.ZRemRangeByRank("inbox:"+followerID, 0, -1001)
    }
}
```

### 读取 Feed

```go
func GetFeed(userID int64, cursor int64, limit int) []*Tweet {
    // 1. 获取收件箱（已推送的普通用户推文）
    inboxTweets := redis.ZRevRangeByScore("inbox:"+userID, cursor, limit)
    
    // 2. 获取关注的大V列表
    followings := redis.SMembers("following:" + userID)
    bigVs := filterBigV(followings)
    
    // 3. 拉取大V的发件箱
    var bigVTweets []int64
    for _, bigVID := range bigVs {
        tweets := redis.ZRevRangeByScore("outbox:"+bigVID, cursor, limit/len(bigVs))
        bigVTweets = append(bigVTweets, tweets...)
    }
    
    // 4. 合并排序
    allTweetIDs := merge(inboxTweets, bigVTweets)
    sort.Slice(allTweetIDs, func(i, j int) bool {
        return allTweetIDs[i] > allTweetIDs[j]  // 按ID倒序（时间倒序）
    })
    
    // 5. 批量获取推文详情
    return batchGetTweets(allTweetIDs[:limit])
}
```

---

## 四、数据模型

```sql
-- 推文表
CREATE TABLE tweets (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    content TEXT,
    media_urls JSON,
    retweet_id BIGINT,  -- 转发的原推文
    created_at TIMESTAMP,
    INDEX idx_user_time (user_id, created_at DESC)
);

-- 关注关系表
CREATE TABLE follows (
    follower_id BIGINT,   -- 粉丝
    following_id BIGINT,  -- 被关注者
    created_at TIMESTAMP,
    PRIMARY KEY (follower_id, following_id),
    INDEX idx_following (following_id)
);

-- 时间线表（收件箱，可选持久化）
CREATE TABLE timeline (
    user_id BIGINT,
    tweet_id BIGINT,
    created_at TIMESTAMP,
    PRIMARY KEY (user_id, created_at, tweet_id)
);
```

---

## 五、Redis 数据结构

```redis
# 发件箱（用户发的推文）
ZADD outbox:{user_id} {timestamp} {tweet_id}

# 收件箱（推送给用户的推文）
ZADD inbox:{user_id} {timestamp} {tweet_id}

# 关注列表
SADD following:{user_id} {following_id1} {following_id2}

# 粉丝列表
SADD followers:{user_id} {follower_id1} {follower_id2}

# 大V标记
SADD bigv:users {user_id}

# 推文详情缓存
HSET tweet:{id} content "xxx" user_id 123 created_at 1701234567
```

---

## 六、热搜系统

### 实现思路

```go
// 1. 统计话题热度（滑动窗口）
func IncrTopicHeat(topic string) {
    now := time.Now().Unix()
    windowKey := "topic:heat:" + strconv.FormatInt(now/60, 10)  // 每分钟一个窗口
    redis.ZIncrBy(windowKey, 1, topic)
    redis.Expire(windowKey, 10*time.Minute)
}

// 2. 计算热搜榜（定时任务）
func CalcHotSearch() {
    now := time.Now().Unix()
    // 聚合最近5分钟的数据
    keys := []string{}
    for i := 0; i < 5; i++ {
        keys = append(keys, "topic:heat:"+strconv.FormatInt((now-int64(i*60))/60, 10))
    }
    
    // ZUNIONSTORE 合并
    redis.ZUnionStore("hotsearch", keys, redis.ZStore{Aggregate: "SUM"})
    
    // 取 Top 50
    hotTopics := redis.ZRevRange("hotsearch", 0, 49)
}

// 3. 热度衰减（时间越久权重越低）
score = 点赞数 + 评论数*2 + 转发数*3 + 时间衰减因子
```

---

## 七、关注/取关

```go
func Follow(followerID, followingID int64) error {
    // 1. 写入关系表
    db.Create(&Follow{FollowerID: followerID, FollowingID: followingID})
    
    // 2. 更新 Redis
    redis.SAdd("following:"+followerID, followingID)
    redis.SAdd("followers:"+followingID, followerID)
    
    // 3. 如果关注的是普通用户，需要补推历史推文到收件箱
    if !isBigV(followingID) {
        recentTweets := redis.ZRevRange("outbox:"+followingID, 0, 100)
        for _, tweetID := range recentTweets {
            redis.ZAdd("inbox:"+followerID, getTweetTime(tweetID), tweetID)
        }
    }
    
    return nil
}

func Unfollow(followerID, followingID int64) error {
    // 1. 删除关系
    db.Delete(&Follow{}, "follower_id = ? AND following_id = ?", followerID, followingID)
    
    // 2. 更新 Redis
    redis.SRem("following:"+followerID, followingID)
    redis.SRem("followers:"+followingID, followerID)
    
    // 3. 清理收件箱中该用户的推文（异步）
    kafka.Send("unfollow-cleanup", &CleanupMessage{
        UserID:      followerID,
        UnfollowID:  followingID,
    })
    
    return nil
}
```

---

## 八、面试追问

| 问题 | 回答 |
|------|------|
| 为什么用推拉结合？ | 大V纯推会导致写放大，千万粉丝写入太慢 |
| 大V阈值怎么定？ | 根据业务调整，一般 1000-10000 |
| 关注大V后历史推文怎么办？ | 读时拉取，不补推 |
| 取关后收件箱怎么清理？ | 异步清理，或读时过滤 |
| 如何处理热点推文？ | 本地缓存 + Redis 多副本 |
| Feed 流如何分页？ | cursor 分页（用 tweet_id 或 timestamp） |
