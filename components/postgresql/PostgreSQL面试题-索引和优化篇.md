# PostgreSQL面试题 - 索引和优化篇

[TOC]

## 索引相关

### 1. PostgreSQL支持哪些索引类型？各有什么特点？

**答案：**

**1. B-tree索引（默认）**
```sql
CREATE INDEX idx_name ON table_name(column_name);
```
- 最常用的索引类型
- 适用于等值查询和范围查询
- 支持排序操作
- 适用于大部分数据类型
- 支持唯一约束

**2. Hash索引**
```sql
CREATE INDEX idx_name ON table_name USING HASH(column_name);
```
- 仅支持等值查询（=）
- 不支持范围查询
- PostgreSQL 10+才支持WAL日志
- 通常B-tree性能更好，很少使用

**3. GiST索引（Generalized Search Tree）**
```sql
CREATE INDEX idx_name ON table_name USING GIST(column_name);
```
- 通用搜索树
- 支持几何数据类型、全文搜索
- 用于PostGIS地理数据
- 支持最近邻搜索(KNN)
- 更新友好

**4. SP-GiST索引（Space-Partitioned GiST）**
```sql
CREATE INDEX idx_name ON table_name USING SPGIST(column_name);
```
- 空间分区GiST
- 适用于非平衡数据结构
- 如：四叉树、k-d树、前缀树
- 适合电话号码、IP地址等

**5. GIN索引（Generalized Inverted Index）**
```sql
CREATE INDEX idx_name ON table_name USING GIN(column_name);
```
- 倒排索引
- 适用于数组、JSONB、全文搜索
- 查询快，但更新慢
- 索引体积较大

**6. BRIN索引（Block Range Index）**
```sql
CREATE INDEX idx_name ON table_name USING BRIN(column_name);
```
- 块范围索引
- 适用于大表中有序数据
- 索引体积小（相比B-tree可小100倍）
- 适合时间序列数据、日志表
- 查询性能不如B-tree，但维护成本低

**索引类型对比：**

| 索引类型 | 适用场景 | 查询性能 | 更新性能 | 索引大小 |
|---------|---------|---------|---------|---------|
| B-tree | 通用，等值和范围查询 | 好 | 好 | 中等 |
| Hash | 仅等值查询 | 好 | 好 | 中等 |
| GiST | 几何、全文搜索 | 中等 | 好 | 较大 |
| SP-GiST | 非平衡数据 | 中等 | 好 | 中等 |
| GIN | 数组、JSONB、全文 | 很好 | 差 | 大 |
| BRIN | 大表有序数据 | 中等 | 很好 | 很小 |

### 2. 什么时候使用GIN索引？什么时候使用GiST索引？

**答案：**

**GIN索引（倒排索引）：**

**适用场景：**
- JSONB字段查询
- 数组包含查询
- 全文搜索
- 查询频繁，更新较少的场景
- 需要精确匹配

**示例：**

```sql
-- JSONB索引
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    data JSONB
);

CREATE INDEX idx_data ON users USING GIN(data);

-- 查询示例
SELECT * FROM users WHERE data @> '{"city": "Beijing"}';
SELECT * FROM users WHERE data ? 'email';
SELECT * FROM users WHERE data @> '{"tags": ["postgresql"]}';

-- 数组索引
CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    tags TEXT[]
);

CREATE INDEX idx_tags ON articles USING GIN(tags);

-- 查询示例
SELECT * FROM articles WHERE tags @> ARRAY['postgresql'];
SELECT * FROM articles WHERE tags && ARRAY['python', 'java'];

-- 全文搜索索引
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    content TEXT
);

CREATE INDEX idx_content_fts ON documents 
USING GIN(to_tsvector('english', content));

-- 查询示例
SELECT * FROM documents 
WHERE to_tsvector('english', content) @@ to_tsquery('postgresql & database');
```

**GiST索引（通用搜索树）：**

**适用场景：**
- 地理空间数据(PostGIS)
- 范围类型查询
- 更新频繁的场景
- 需要支持多种操作符
- 最近邻搜索(KNN)

**示例：**

```sql
-- 地理空间索引
CREATE EXTENSION postgis;

CREATE TABLE places (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    location GEOMETRY(Point, 4326)
);

CREATE INDEX idx_location ON places USING GIST(location);

-- 查询示例
-- 查找距离某点1000米内的地点
SELECT * FROM places 
WHERE ST_DWithin(location, ST_MakePoint(116.4, 39.9)::geography, 1000);

-- 查找最近的10个地点
SELECT * FROM places 
ORDER BY location <-> ST_MakePoint(116.4, 39.9)::geometry 
LIMIT 10;

-- 范围类型索引
CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    room_id INTEGER,
    period DATERANGE
);

CREATE INDEX idx_period ON bookings USING GIST(period);

-- 查询示例
SELECT * FROM bookings 
WHERE period && '[2024-01-01, 2024-01-31]'::daterange;

-- 全文搜索（GiST也支持，但通常GIN更快）
CREATE INDEX idx_content_gist ON documents 
USING GIST(to_tsvector('english', content));
```

**GIN vs GiST对比：**

| 特性 | GIN | GiST |
|------|-----|------|
| **查询速度** | 更快（3-10倍） | 较快 |
| **更新速度** | 慢 | 快 |
| **索引大小** | 大（2-3倍） | 小 |
| **构建时间** | 长 | 短 |
| **适用场景** | 查询多更新少 | 更新频繁 |
| **全文搜索** | 推荐 | 可用 |
| **JSONB** | 推荐 | 不支持 |
| **几何类型** | 不支持 | 推荐 |
| **KNN搜索** | 不支持 | 支持 |

**选择建议：**
- **全文搜索、JSONB、数组**：优先GIN
- **地理空间、范围类型**：使用GiST
- **更新频繁**：考虑GiST
- **查询性能优先**：选择GIN

### 3. 如何优化PostgreSQL的索引？

**答案：**

**1. 选择合适的索引类型**

```sql
-- 等值查询：B-tree
CREATE INDEX idx_status ON orders(status);

-- 范围查询：B-tree
CREATE INDEX idx_created ON orders(created_at);

-- 全文搜索：GIN
CREATE INDEX idx_content ON articles 
USING GIN(to_tsvector('english', content));

-- 大表有序数据：BRIN
CREATE INDEX idx_timestamp ON logs USING BRIN(timestamp);

-- 地理数据：GiST
CREATE INDEX idx_location ON places USING GIST(location);
```

**2. 使用部分索引（Partial Index）**

```sql
-- 只为活跃用户创建索引
CREATE INDEX idx_active_users ON users(email) 
WHERE status = 'active';

-- 只为未删除的记录创建索引
CREATE INDEX idx_valid_orders ON orders(order_date) 
WHERE deleted_at IS NULL;

-- 只为高价值订单创建索引
CREATE INDEX idx_high_value_orders ON orders(created_at) 
WHERE total > 1000;

-- 优点：
-- 1. 减小索引大小
-- 2. 提高索引维护速度
-- 3. 提高查询性能
```

**3. 使用表达式索引（Expression Index）**

```sql
-- 为函数结果创建索引
CREATE INDEX idx_lower_email ON users(LOWER(email));
SELECT * FROM users WHERE LOWER(email) = 'user@example.com';

-- 为计算结果创建索引
CREATE INDEX idx_total ON orders((price * quantity));
SELECT * FROM orders WHERE price * quantity > 1000;

-- 为JSON路径创建索引
CREATE INDEX idx_city ON users((data->>'city'));
SELECT * FROM users WHERE data->>'city' = 'Beijing';

-- 为日期部分创建索引
CREATE INDEX idx_year ON orders(EXTRACT(YEAR FROM created_at));
SELECT * FROM orders WHERE EXTRACT(YEAR FROM created_at) = 2024;
```

**4. 使用复合索引（Composite Index）**

```sql
-- 注意字段顺序：选择性高的字段放前面
CREATE INDEX idx_user_status_date ON orders(user_id, status, created_at);

-- 可以支持的查询：
SELECT * FROM orders WHERE user_id = 100;
SELECT * FROM orders WHERE user_id = 100 AND status = 'pending';
SELECT * FROM orders WHERE user_id = 100 AND status = 'pending' 
    AND created_at > '2024-01-01';

-- 不能使用索引的查询：
SELECT * FROM orders WHERE status = 'pending';  -- 跳过了第一列
SELECT * FROM orders WHERE created_at > '2024-01-01';  -- 跳过了前两列
```

**5. 使用INCLUDE列（覆盖索引）**

```sql
-- PostgreSQL 11+
CREATE INDEX idx_user_include ON orders(user_id) 
INCLUDE (status, total, created_at);

-- 查询可以只扫描索引，不需要回表（Index Only Scan）
SELECT user_id, status, total, created_at 
FROM orders 
WHERE user_id = 100;
```

**6. 使用唯一索引**

```sql
-- 唯一索引比普通索引更高效
CREATE UNIQUE INDEX idx_email ON users(email);

-- 部分唯一索引
CREATE UNIQUE INDEX idx_active_email ON users(email) 
WHERE status = 'active';
```

**7. 并发创建索引**

```sql
-- 不锁表创建索引（推荐用于生产环境）
CREATE INDEX CONCURRENTLY idx_name ON table_name(column_name);

-- 注意：
-- 1. 创建时间更长
-- 2. 不能在事务中使用
-- 3. 如果失败，会留下INVALID索引，需要手动删除
```

**8. 定期维护索引**

```sql
-- 重建索引
REINDEX INDEX idx_name;
REINDEX TABLE table_name;
REINDEX DATABASE database_name;

-- 并发重建索引（不锁表）
CREATE INDEX CONCURRENTLY idx_new ON table_name(column_name);
DROP INDEX CONCURRENTLY idx_old;
ALTER INDEX idx_new RENAME TO idx_old;

-- 更新统计信息
ANALYZE table_name;
VACUUM ANALYZE table_name;
```

**9. 监控索引使用情况**

```sql
-- 查看未使用的索引
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    pg_size_pretty(pg_relation_size(indexrelid)) as size
FROM pg_stat_user_indexes
WHERE idx_scan = 0 
  AND indexrelname NOT LIKE 'pg_toast%'
ORDER BY pg_relation_size(indexrelid) DESC;

-- 查看索引大小
SELECT 
    indexrelname,
    pg_size_pretty(pg_relation_size(indexrelid)) as size
FROM pg_stat_user_indexes
ORDER BY pg_relation_size(indexrelid) DESC;

-- 查看索引缓存命中率
SELECT 
    indexrelname,
    idx_blks_hit,
    idx_blks_read,
    CASE WHEN (idx_blks_hit + idx_blks_read) = 0 THEN 0
         ELSE round(100.0 * idx_blks_hit / (idx_blks_hit + idx_blks_read), 2)
    END as cache_hit_ratio
FROM pg_statio_user_indexes
ORDER BY cache_hit_ratio;

-- 查看重复索引
SELECT 
    pg_size_pretty(SUM(pg_relation_size(idx))::BIGINT) AS size,
    (array_agg(idx))[1] AS idx1,
    (array_agg(idx))[2] AS idx2,
    (array_agg(idx))[3] AS idx3,
    (array_agg(idx))[4] AS idx4
FROM (
    SELECT 
        indexrelid::regclass AS idx,
        (indrelid::text ||E'\n'|| indclass::text ||E'\n'|| 
         indkey::text ||E'\n'|| COALESCE(indexprs::text,'')||E'\n' || 
         COALESCE(indpred::text,'')) AS key
    FROM pg_index
) sub
GROUP BY key 
HAVING COUNT(*) > 1
ORDER BY SUM(pg_relation_size(idx)) DESC;
```

## 查询优化

### 4. 如何分析和优化慢查询？

**答案：**

**1. 使用EXPLAIN分析查询计划**

```sql
-- 查看查询计划
EXPLAIN SELECT * FROM orders WHERE user_id = 100;

-- 查看实际执行统计
EXPLAIN ANALYZE SELECT * FROM orders WHERE user_id = 100;

-- 查看详细信息（包括缓冲区使用）
EXPLAIN (ANALYZE, BUFFERS, VERBOSE, COSTS, TIMING) 
SELECT * FROM orders WHERE user_id = 100;

-- 查看JSON格式输出
EXPLAIN (ANALYZE, FORMAT JSON) 
SELECT * FROM orders WHERE user_id = 100;
```

**2. 理解EXPLAIN输出**

**常见扫描类型：**
- **Seq Scan**：全表扫描（通常需要优化）
- **Index Scan**：索引扫描，然后回表
- **Index Only Scan**：仅索引扫描（最优）
- **Bitmap Heap Scan**：位图堆扫描（多个索引条件）
- **Bitmap Index Scan**：位图索引扫描

**常见连接类型：**
- **Nested Loop**：嵌套循环（小表）
- **Hash Join**：哈希连接（大表等值连接）
- **Merge Join**：归并连接（已排序数据）

**关键指标：**
```sql
-- cost=0.00..100.00  -- 启动成本..总成本
-- rows=1000          -- 预计返回行数
-- width=50           -- 预计每行字节数
-- actual time=0.1..10.5  -- 实际执行时间（毫秒）
-- loops=1            -- 执行次数
-- Buffers: shared hit=100 read=50  -- 缓冲区命中和读取
```

**3. 常见优化方法**

**添加索引：**
```sql
-- 为WHERE条件添加索引
CREATE INDEX idx_user_id ON orders(user_id);

-- 为JOIN列添加索引
CREATE INDEX idx_order_id ON order_items(order_id);

-- 为ORDER BY列添加索引
CREATE INDEX idx_created_at ON orders(created_at DESC);

-- 为GROUP BY列添加索引
CREATE INDEX idx_status ON orders(status);
```

**优化查询语句：**
```sql
-- 避免SELECT *，只查询需要的列
-- 不好
SELECT * FROM orders WHERE user_id = 100;
-- 好
SELECT id, user_id, total, created_at FROM orders WHERE user_id = 100;

-- 使用LIMIT限制返回行数
SELECT * FROM orders ORDER BY created_at DESC LIMIT 100;

-- 使用EXISTS代替IN（大数据集）
-- 不好
SELECT * FROM users WHERE id IN (SELECT user_id FROM orders);
-- 好
SELECT * FROM users u WHERE EXISTS (
    SELECT 1 FROM orders o WHERE o.user_id = u.id
);

-- 使用JOIN代替子查询
-- 不好
SELECT *, (SELECT COUNT(*) FROM orders WHERE user_id = users.id) as order_count
FROM users;
-- 好
SELECT u.*, COUNT(o.id) as order_count
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id;

-- 避免在WHERE中使用函数
-- 不好（不能使用索引）
SELECT * FROM users WHERE LOWER(email) = 'user@example.com';
-- 好（创建函数索引）
CREATE INDEX idx_lower_email ON users(LOWER(email));
SELECT * FROM users WHERE LOWER(email) = 'user@example.com';

-- 使用UNION ALL代替UNION（如果不需要去重）
SELECT * FROM orders WHERE status = 'pending'
UNION ALL
SELECT * FROM orders WHERE status = 'processing';
```

**使用CTE优化复杂查询：**
```sql
-- 使用CTE分解复杂查询
WITH recent_orders AS (
    SELECT user_id, COUNT(*) as order_count
    FROM orders
    WHERE created_at > CURRENT_DATE - INTERVAL '30 days'
    GROUP BY user_id
),
active_users AS (
    SELECT id, name, email
    FROM users
    WHERE status = 'active'
)
SELECT u.*, COALESCE(o.order_count, 0) as recent_orders
FROM active_users u
LEFT JOIN recent_orders o ON u.id = o.user_id;
```

**4. 配置优化**

```sql
-- 增加工作内存（用于排序、哈希等）
SET work_mem = '256MB';  -- 会话级别
-- 或在postgresql.conf中设置
work_mem = 64MB  -- 全局默认值

-- 增加维护工作内存（用于VACUUM、CREATE INDEX等）
SET maintenance_work_mem = '1GB';

-- 调整查询规划器参数
SET random_page_cost = 1.1;  -- SSD使用较小值
SET seq_page_cost = 1.0;
SET effective_cache_size = '16GB';  -- 操作系统缓存大小

-- 调整并行查询
SET max_parallel_workers_per_gather = 4;
SET parallel_tuple_cost = 0.1;
SET parallel_setup_cost = 1000;
```

**5. 启用慢查询日志**

```sql
-- 在postgresql.conf中配置
logging_collector = on
log_directory = 'pg_log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_min_duration_statement = 1000  -- 记录超过1秒的查询
log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '
log_statement = 'all'  -- 记录所有SQL（可选）
log_duration = on
log_lock_waits = on
log_temp_files = 0  -- 记录所有临时文件
```

**6. 使用pg_stat_statements扩展**

```sql
-- 安装扩展
CREATE EXTENSION pg_stat_statements;

-- 在postgresql.conf中配置
shared_preload_libraries = 'pg_stat_statements'
pg_stat_statements.track = all

-- 查看最慢的查询
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    max_time,
    stddev_time,
    rows
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 20;

-- 查看最频繁的查询
SELECT 
    query,
    calls,
    total_time,
    mean_time
FROM pg_stat_statements
ORDER BY calls DESC
LIMIT 20;

-- 重置统计信息
SELECT pg_stat_statements_reset();
```

### 5. PostgreSQL的JOIN类型和优化策略是什么？

**答案：**

**JOIN类型：**

**1. INNER JOIN（内连接）**
```sql
-- 只返回匹配的行
SELECT u.name, o.order_id, o.total
FROM users u
INNER JOIN orders o ON u.id = o.user_id;
```

**2. LEFT JOIN（左外连接）**
```sql
-- 返回左表所有行，右表不匹配则为NULL
SELECT u.name, o.order_id, o.total
FROM users u
LEFT JOIN orders o ON u.id = o.user_id;

-- 查找没有订单的用户
SELECT u.name
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.id IS NULL;
```

**3. RIGHT JOIN（右外连接）**
```sql
-- 返回右表所有行，左表不匹配则为NULL
SELECT u.name, o.order_id, o.total
FROM users u
RIGHT JOIN orders o ON u.id = o.user_id;
```

**4. FULL OUTER JOIN（全外连接）**
```sql
-- 返回两表所有行，不匹配则为NULL
SELECT u.name, o.order_id
FROM users u
FULL OUTER JOIN orders o ON u.id = o.user_id;
```

**5. CROSS JOIN（交叉连接）**
```sql
-- 笛卡尔积
SELECT u.name, p.product_name
FROM users u
CROSS JOIN products p;

-- 等价于
SELECT u.name, p.product_name
FROM users u, products p;
```

**6. SELF JOIN（自连接）**
```sql
-- 查找同一城市的用户对
SELECT u1.name as user1, u2.name as user2, u1.city
FROM users u1
JOIN users u2 ON u1.city = u2.city AND u1.id < u2.id;
```

**JOIN执行策略：**

**1. Nested Loop Join（嵌套循环）**
- 外表每行都扫描内表
- 适用于小表连接
- 适用于内表有索引的情况

```sql
-- 示例
SELECT * FROM small_table s
JOIN large_table l ON s.id = l.ref_id;
-- 如果large_table.ref_id有索引，可能使用Nested Loop
```

**2. Hash Join（哈希连接）**
- 构建哈希表，然后探测
- 适用于大表等值连接
- 需要足够的work_mem

```sql
-- 增加work_mem提高Hash Join性能
SET work_mem = '256MB';

SELECT * FROM large_table1 l1
JOIN large_table2 l2 ON l1.id = l2.ref_id;
```

**3. Merge Join（归并连接）**
- 两表同时扫描
- 适用于已排序的数据
- 需要JOIN列有序

```sql
-- 如果有索引，可能使用Merge Join
CREATE INDEX idx_user_id ON orders(user_id);
CREATE INDEX idx_id ON users(id);

SELECT * FROM users u
JOIN orders o ON u.id = o.user_id;
```

**JOIN优化策略：**

**1. 在JOIN列上创建索引**
```sql
CREATE INDEX idx_user_id ON orders(user_id);
CREATE INDEX idx_id ON users(id);
```

**2. 使用合适的JOIN顺序**
```sql
-- 小表在前，大表在后
SELECT * FROM small_table s
JOIN large_table l ON s.id = l.ref_id;

-- 使用JOIN提示（PostgreSQL不直接支持，但可以通过配置影响）
SET join_collapse_limit = 1;  -- 禁止重排JOIN顺序
```

**3. 避免JOIN过多表**
```sql
-- 不好：JOIN太多表
SELECT * FROM t1
JOIN t2 ON t1.id = t2.t1_id
JOIN t3 ON t2.id = t3.t2_id
JOIN t4 ON t3.id = t4.t3_id
JOIN t5 ON t4.id = t5.t4_id
JOIN t6 ON t5.id = t6.t5_id;

-- 好：使用CTE分解
WITH step1 AS (
    SELECT * FROM t1
    JOIN t2 ON t1.id = t2.t1_id
    JOIN t3 ON t2.id = t3.t2_id
),
step2 AS (
    SELECT * FROM t4
    JOIN t5 ON t4.id = t5.t4_id
    JOIN t6 ON t5.id = t6.t5_id
)
SELECT * FROM step1
JOIN step2 ON step1.id = step2.ref_id;
```

**4. 使用WHERE过滤减少JOIN数据量**
```sql
-- 先过滤再JOIN
SELECT u.name, o.order_id
FROM users u
JOIN orders o ON u.id = o.user_id
WHERE o.created_at > '2024-01-01'  -- 先过滤
  AND u.status = 'active';
```

**5. 使用LATERAL JOIN处理相关子查询**
```sql
-- 不好：相关子查询
SELECT u.name,
       (SELECT COUNT(*) FROM orders WHERE user_id = u.id) as order_count,
       (SELECT MAX(total) FROM orders WHERE user_id = u.id) as max_order
FROM users u;

-- 好：使用LATERAL JOIN
SELECT u.name, o.order_count, o.max_order
FROM users u
LEFT JOIN LATERAL (
    SELECT 
        COUNT(*) as order_count,
        MAX(total) as max_order
    FROM orders
    WHERE user_id = u.id
) o ON true;
```

**6. 使用半连接和反连接**
```sql
-- 半连接：EXISTS
SELECT * FROM users u
WHERE EXISTS (
    SELECT 1 FROM orders o WHERE o.user_id = u.id
);

-- 反连接：NOT EXISTS
SELECT * FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM orders o WHERE o.user_id = u.id
);

-- 使用LEFT JOIN + IS NULL（反连接）
SELECT u.* FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.id IS NULL;
```

**7. 监控JOIN性能**
```sql
-- 使用EXPLAIN ANALYZE查看JOIN策略
EXPLAIN (ANALYZE, BUFFERS) 
SELECT * FROM users u
JOIN orders o ON u.id = o.user_id;

-- 查看JOIN相关的配置
SHOW join_collapse_limit;
SHOW from_collapse_limit;
SHOW geqo_threshold;
```
