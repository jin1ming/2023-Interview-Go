# PostgreSQL面试题 - 高级特性篇

[TOC]

## 数据类型

### 1. JSONB类型的使用

**基本操作：**
```sql
CREATE TABLE users (id SERIAL, data JSONB);

-- 操作符
data->'key'     -- 返回JSONB
data->>'key'    -- 返回TEXT
data @> '{"key": "value"}'  -- 包含
data ? 'key'    -- 键存在

-- 索引
CREATE INDEX idx_data ON users USING GIN(data);
```

### 2. 数组类型的使用

**基本操作：**
```sql
CREATE TABLE articles (id SERIAL, tags TEXT[]);

-- 操作符
tags @> ARRAY['tag1']  -- 包含
tags && ARRAY['tag1', 'tag2']  -- 重叠

-- 索引
CREATE INDEX idx_tags ON articles USING GIN(tags);
```

### 3. 范围类型的使用

**基本操作：**
```sql
CREATE TABLE bookings (id SERIAL, period DATERANGE);

-- 操作符
period @> '2024-01-01'::date  -- 包含日期
period && '[2024-01-01, 2024-01-10)'::daterange  -- 重叠

-- 排他约束
CREATE EXTENSION btree_gist;
ALTER TABLE bookings 
ADD CONSTRAINT no_overlap 
EXCLUDE USING GIST (room_id WITH =, period WITH &&);
```

## CTE和递归查询

### 4. CTE（公共表表达式）

**基本CTE：**
```sql
WITH recent_orders AS (
    SELECT * FROM orders 
    WHERE created_at > CURRENT_DATE - INTERVAL '7 days'
)
SELECT user_id, COUNT(*) FROM recent_orders GROUP BY user_id;
```

**递归CTE：**
```sql
-- 组织架构树
WITH RECURSIVE subordinates AS (
    SELECT id, name, manager_id, 1 as level
    FROM employees WHERE id = 1
    UNION ALL
    SELECT e.id, e.name, e.manager_id, s.level + 1
    FROM employees e
    JOIN subordinates s ON e.manager_id = s.id
)
SELECT * FROM subordinates;
```

## 窗口函数

### 5. 常用窗口函数

```sql
-- ROW_NUMBER：行号
SELECT name, salary, 
       ROW_NUMBER() OVER (ORDER BY salary DESC) as rank
FROM employees;

-- RANK：排名（有并列）
SELECT name, salary,
       RANK() OVER (ORDER BY salary DESC) as rank
FROM employees;

-- DENSE_RANK：密集排名
SELECT name, salary,
       DENSE_RANK() OVER (ORDER BY salary DESC) as rank
FROM employees;

-- PARTITION BY：分组
SELECT dept, name, salary,
       ROW_NUMBER() OVER (PARTITION BY dept ORDER BY salary DESC) as rank
FROM employees;

-- LAG/LEAD：前后行
SELECT name, salary,
       LAG(salary) OVER (ORDER BY salary) as prev_salary,
       LEAD(salary) OVER (ORDER BY salary) as next_salary
FROM employees;

-- FIRST_VALUE/LAST_VALUE：首尾值
SELECT name, salary,
       FIRST_VALUE(salary) OVER (ORDER BY salary) as min_salary,
       LAST_VALUE(salary) OVER (ORDER BY salary 
           ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING) as max_salary
FROM employees;

-- 移动平均
SELECT date, amount,
       AVG(amount) OVER (ORDER BY date ROWS BETWEEN 6 PRECEDING AND CURRENT ROW) as moving_avg
FROM sales;
```

## 全文搜索

### 6. 全文搜索功能

```sql
-- 创建tsvector
SELECT to_tsvector('english', 'The quick brown fox');

-- 创建tsquery
SELECT to_tsquery('english', 'quick & fox');

-- 全文搜索
SELECT * FROM documents
WHERE to_tsvector('english', content) @@ to_tsquery('english', 'postgresql & database');

-- 创建索引
CREATE INDEX idx_content_fts ON documents 
USING GIN(to_tsvector('english', content));

-- 添加tsvector列（推荐）
ALTER TABLE documents ADD COLUMN content_tsv TSVECTOR;
UPDATE documents SET content_tsv = to_tsvector('english', content);
CREATE INDEX idx_content_tsv ON documents USING GIN(content_tsv);

-- 自动更新tsvector
CREATE TRIGGER tsvector_update BEFORE INSERT OR UPDATE ON documents
FOR EACH ROW EXECUTE FUNCTION
tsvector_update_trigger(content_tsv, 'pg_catalog.english', content);

-- 排名
SELECT title, ts_rank(content_tsv, query) as rank
FROM documents, to_tsquery('english', 'postgresql') query
WHERE content_tsv @@ query
ORDER BY rank DESC;
```

## 分区表

### 7. 表分区

**范围分区：**
```sql
-- 创建主表
CREATE TABLE orders (
    id SERIAL,
    order_date DATE NOT NULL,
    amount NUMERIC
) PARTITION BY RANGE (order_date);

-- 创建分区
CREATE TABLE orders_2023 PARTITION OF orders
FOR VALUES FROM ('2023-01-01') TO ('2024-01-01');

CREATE TABLE orders_2024 PARTITION OF orders
FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

-- 创建默认分区
CREATE TABLE orders_default PARTITION OF orders DEFAULT;
```

**列表分区：**
```sql
CREATE TABLE users (
    id SERIAL,
    country VARCHAR(2),
    name VARCHAR(100)
) PARTITION BY LIST (country);

CREATE TABLE users_cn PARTITION OF users FOR VALUES IN ('CN');
CREATE TABLE users_us PARTITION OF users FOR VALUES IN ('US');
```

**哈希分区：**
```sql
CREATE TABLE logs (
    id SERIAL,
    user_id INTEGER,
    message TEXT
) PARTITION BY HASH (user_id);

CREATE TABLE logs_0 PARTITION OF logs FOR VALUES WITH (MODULUS 4, REMAINDER 0);
CREATE TABLE logs_1 PARTITION OF logs FOR VALUES WITH (MODULUS 4, REMAINDER 1);
CREATE TABLE logs_2 PARTITION OF logs FOR VALUES WITH (MODULUS 4, REMAINDER 2);
CREATE TABLE logs_3 PARTITION OF logs FOR VALUES WITH (MODULUS 4, REMAINDER 3);
```

## 外部数据包装器

### 8. FDW（Foreign Data Wrapper）

```sql
-- 安装postgres_fdw扩展
CREATE EXTENSION postgres_fdw;

-- 创建外部服务器
CREATE SERVER foreign_server
FOREIGN DATA WRAPPER postgres_fdw
OPTIONS (host 'remote-host', port '5432', dbname 'remote_db');

-- 创建用户映射
CREATE USER MAPPING FOR local_user
SERVER foreign_server
OPTIONS (user 'remote_user', password 'password');

-- 创建外部表
CREATE FOREIGN TABLE foreign_orders (
    id INTEGER,
    user_id INTEGER,
    total NUMERIC
) SERVER foreign_server
OPTIONS (schema_name 'public', table_name 'orders');

-- 查询外部表
SELECT * FROM foreign_orders WHERE user_id = 100;

-- 导入外部schema
IMPORT FOREIGN SCHEMA public
FROM SERVER foreign_server
INTO local_schema;
```

## 扩展

### 9. 常用扩展

```sql
-- PostGIS：地理空间数据
CREATE EXTENSION postgis;

-- pg_trgm：模糊搜索
CREATE EXTENSION pg_trgm;
CREATE INDEX idx_name_trgm ON users USING GIN(name gin_trgm_ops);
SELECT * FROM users WHERE name % 'alice';  -- 相似度搜索

-- uuid-ossp：UUID生成
CREATE EXTENSION "uuid-ossp";
SELECT uuid_generate_v4();

-- hstore：键值对存储
CREATE EXTENSION hstore;
CREATE TABLE products (id SERIAL, attributes HSTORE);

-- pg_stat_statements：查询统计
CREATE EXTENSION pg_stat_statements;
SELECT query, calls, total_time FROM pg_stat_statements
ORDER BY total_time DESC LIMIT 10;

-- pgcrypto：加密函数
CREATE EXTENSION pgcrypto;
SELECT crypt('password', gen_salt('bf'));
```

## 性能监控

### 10. 性能监控查询

```sql
-- 查看活动连接
SELECT pid, usename, application_name, client_addr, state, query
FROM pg_stat_activity
WHERE state != 'idle';

-- 查看表大小
SELECT schemaname, tablename,
       pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- 查看索引使用情况
SELECT schemaname, tablename, indexname, idx_scan
FROM pg_stat_user_indexes
ORDER BY idx_scan;

-- 查看缓存命中率
SELECT 
    sum(heap_blks_read) as heap_read,
    sum(heap_blks_hit) as heap_hit,
    sum(heap_blks_hit) / (sum(heap_blks_hit) + sum(heap_blks_read)) as ratio
FROM pg_statio_user_tables;

-- 查看长时间运行的查询
SELECT pid, now() - query_start as duration, query
FROM pg_stat_activity
WHERE state = 'active'
  AND now() - query_start > interval '5 minutes';

-- 终止查询
SELECT pg_terminate_backend(pid);
```
