# 分布式ID生成器设计

## 一、需求分析

| 需求 | 说明 |
|------|------|
| 全局唯一 | 多节点生成不冲突 |
| 趋势递增 | 便于数据库索引 |
| 高性能 | 单机 10万+ QPS |
| 高可用 | 服务不能单点 |
| 信息安全 | 不暴露业务量（可选） |

---

## 二、方案对比

| 方案 | 优点 | 缺点 | QPS |
|------|------|------|-----|
| **UUID** | 简单，本地生成 | 无序，太长（36位） | 极高 |
| **数据库自增** | 简单，有序 | 单点瓶颈，性能低 | 1K |
| **Redis INCR** | 简单，性能好 | 依赖Redis，持久化风险 | 10K |
| **Snowflake** | 有序，高性能，本地生成 | 时钟回拨问题 | 400K |
| **号段模式** | 高可用，DB压力小 | 实现复杂 | 100K |

---

## 三、Snowflake 算法（推荐）

### 结构（64位）

```
0 - 41位时间戳 - 10位机器ID - 12位序列号

| 1bit | 41bit      | 10bit     | 12bit    |
|------|------------|-----------|----------|
| 符号  | 时间戳(ms)  | 机器ID    | 序列号    |
| 0    | 毫秒级时间  | 1024台机器 | 每毫秒4096个 |
```

### 容量计算

- **时间戳**：41位 → 2^41 ms ≈ 69年
- **机器ID**：10位 → 1024台机器
- **序列号**：12位 → 每毫秒 4096 个
- **总QPS**：1024 × 4096 × 1000 = **41亿/秒**

### Go 实现

```go
type Snowflake struct {
    mu        sync.Mutex
    epoch     int64  // 起始时间戳（2020-01-01）
    machineID int64  // 机器ID（0-1023）
    sequence  int64  // 序列号（0-4095）
    lastTime  int64  // 上次生成时间
}

const (
    machineBits   = 10
    sequenceBits  = 12
    machineMax    = 1<<machineBits - 1  // 1023
    sequenceMax   = 1<<sequenceBits - 1 // 4095
    machineShift  = sequenceBits        // 12
    timestampShift = machineBits + sequenceBits // 22
)

func NewSnowflake(machineID int64) *Snowflake {
    return &Snowflake{
        epoch:     1577836800000, // 2020-01-01 00:00:00 UTC
        machineID: machineID & machineMax,
    }
}

func (s *Snowflake) NextID() int64 {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    now := time.Now().UnixMilli()
    
    if now == s.lastTime {
        // 同一毫秒，序列号+1
        s.sequence = (s.sequence + 1) & sequenceMax
        if s.sequence == 0 {
            // 序列号用完，等待下一毫秒
            for now <= s.lastTime {
                now = time.Now().UnixMilli()
            }
        }
    } else if now > s.lastTime {
        // 新的毫秒，序列号归零
        s.sequence = 0
    } else {
        // 时钟回拨，抛异常或等待
        panic("clock moved backwards")
    }
    
    s.lastTime = now
    
    // 组装ID
    id := ((now - s.epoch) << timestampShift) |
          (s.machineID << machineShift) |
          s.sequence
    
    return id
}
```

### 时钟回拨处理

```go
func (s *Snowflake) NextIDSafe() (int64, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    now := time.Now().UnixMilli()
    
    if now < s.lastTime {
        offset := s.lastTime - now
        if offset <= 5 {
            // 回拨5ms内，等待
            time.Sleep(time.Duration(offset) * time.Millisecond)
            now = time.Now().UnixMilli()
        } else {
            // 回拨太大，报错
            return 0, errors.New("clock moved backwards")
        }
    }
    // ... 后续逻辑
}
```

---

## 四、号段模式（美团Leaf）

### 原理

```
1. 服务启动时从DB获取一个号段（如 1-1000）
2. 本地内存发号，用完再取下一段
3. 双Buffer：当前号段用到一定比例时，异步加载下一段
```

### 数据库表

```sql
CREATE TABLE id_segment (
    biz_tag VARCHAR(64) PRIMARY KEY,  -- 业务标识
    max_id BIGINT NOT NULL,           -- 当前最大ID
    step INT NOT NULL,                -- 号段步长
    version INT NOT NULL,             -- 乐观锁版本
    updated_at TIMESTAMP
);

-- 示例数据
INSERT INTO id_segment VALUES ('order', 0, 1000, 0, NOW());
```

### 实现

```go
type Segment struct {
    currentID int64
    maxID     int64
}

type LeafIDGenerator struct {
    bizTag   string
    step     int64
    current  *Segment
    next     *Segment
    loading  bool
    mu       sync.Mutex
}

func (g *LeafIDGenerator) NextID() (int64, error) {
    g.mu.Lock()
    defer g.mu.Unlock()
    
    // 当前号段用完，切换到下一段
    if g.current.currentID >= g.current.maxID {
        if g.next == nil {
            return 0, errors.New("no available segment")
        }
        g.current = g.next
        g.next = nil
    }
    
    // 使用率达到一定比例，异步加载下一段
    usage := float64(g.current.currentID-g.current.maxID+g.step) / float64(g.step)
    if usage > 0.7 && g.next == nil && !g.loading {
        go g.loadNextSegment()
    }
    
    g.current.currentID++
    return g.current.currentID, nil
}

func (g *LeafIDGenerator) loadNextSegment() {
    g.loading = true
    defer func() { g.loading = false }()
    
    // 乐观锁更新
    for {
        var seg IDSegment
        db.Where("biz_tag = ?", g.bizTag).First(&seg)
        
        result := db.Model(&IDSegment{}).
            Where("biz_tag = ? AND version = ?", g.bizTag, seg.Version).
            Updates(map[string]interface{}{
                "max_id":  seg.MaxID + g.step,
                "version": seg.Version + 1,
            })
        
        if result.RowsAffected > 0 {
            g.mu.Lock()
            g.next = &Segment{
                currentID: seg.MaxID,
                maxID:     seg.MaxID + g.step,
            }
            g.mu.Unlock()
            return
        }
        // 乐观锁冲突，重试
    }
}
```

### 双Buffer优化

```
Buffer1: [1, 1000]     ← 当前使用
Buffer2: [1001, 2000]  ← 预加载

当 Buffer1 用到 70% 时，异步加载 Buffer2
Buffer1 用完后无缝切换到 Buffer2
```

---

## 五、Redis 方案

```go
func NextID(bizTag string) (int64, error) {
    // INCR 原子自增
    id, err := redis.Incr("id:" + bizTag).Result()
    if err != nil {
        return 0, err
    }
    return id, nil
}

// 批量获取（减少网络开销）
func NextIDBatch(bizTag string, count int64) (int64, int64, error) {
    // INCRBY 批量获取
    maxID, err := redis.IncrBy("id:"+bizTag, count).Result()
    if err != nil {
        return 0, 0, err
    }
    return maxID - count + 1, maxID, nil
}
```

**风险**：Redis 持久化可能丢数据，需配合 AOF always 或号段模式兜底。

---

## 六、方案选型

| 场景 | 推荐方案 |
|------|----------|
| 简单业务，不要求有序 | UUID |
| 单机服务 | 数据库自增 |
| 高并发，可接受时钟依赖 | **Snowflake** |
| 高可用，DB兜底 | **号段模式（Leaf）** |
| 已有Redis，量不大 | Redis INCR |

---

## 七、面试追问

| 问题 | 回答 |
|------|------|
| Snowflake 时钟回拨怎么办？ | 小回拨等待，大回拨报错或切换机器ID |
| 机器ID如何分配？ | ZooKeeper/etcd 注册，或配置文件 |
| 号段模式DB挂了怎么办？ | 双Buffer撑一段时间，DB主从切换 |
| ID能否反解？ | 可以，从ID提取时间戳、机器ID |
| 如何保证趋势递增？ | Snowflake 时间戳在高位，天然递增 |
| 为什么不用UUID？ | 无序影响B+树索引性能，且太长 |
