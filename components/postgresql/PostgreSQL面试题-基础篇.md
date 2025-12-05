# PostgreSQL面试题 - 基础篇

[TOC]

## 基础概念类

### 1. 什么是PostgreSQL？它有哪些特点？

**答案：**
PostgreSQL是一个功能强大的开源对象关系型数据库管理系统(ORDBMS)，具有30多年的发展历史。

**主要特点：**
- **开源免费**：采用PostgreSQL License，类似BSD/MIT许可
- **ACID完全兼容**：支持完整的事务特性
- **支持复杂数据类型**：JSON、JSONB、数组、hstore、XML等
- **强大的扩展性**：支持自定义函数、数据类型、索引类型
- **高级特性**：窗口函数、CTE、全文搜索、GIS支持(PostGIS)
- **多版本并发控制(MVCC)**：提高并发性能
- **丰富的索引类型**：B-tree、Hash、GiST、SP-GiST、GIN、BRIN
- **支持继承**：表继承功能
- **外部数据包装器(FDW)**：可以访问外部数据源

### 2. PostgreSQL与MySQL的主要区别是什么？

**答案：**

| 特性 | PostgreSQL | MySQL |
|------|-----------|-------|
| **数据类型** | 支持更多复杂类型(数组、JSON、几何等) | 类型相对简单 |
| **事务支持** | 完全支持ACID，所有引擎都支持事务 | InnoDB支持，MyISAM不支持 |
| **并发控制** | MVCC，读写不阻塞 | InnoDB使用MVCC，但实现不同 |
| **子查询** | 性能优秀，支持更复杂的子查询 | 早期版本性能较差 |
| **窗口函数** | 支持完整的窗口函数 | 8.0+才支持 |
| **全文搜索** | 内置全文搜索 | 需要额外配置 |
| **GIS支持** | PostGIS扩展，功能强大 | 基础GIS支持 |
| **复制** | 流复制、逻辑复制 | 主从复制、组复制 |
| **扩展性** | 高度可扩展 | 相对受限 |
| **性能** | 复杂查询性能好 | 简单查询性能好 |

### 3. 什么是MVCC？PostgreSQL如何实现MVCC？

**答案：**
MVCC(Multi-Version Concurrency Control，多版本并发控制)是一种并发控制方法，允许数据库在不使用锁的情况下提供并发访问。

**PostgreSQL的MVCC实现：**
- **版本号**：每个事务都有唯一的事务ID(XID)
- **元组版本**：每行数据(tuple)包含隐藏字段：
  - `xmin`：插入该行的事务ID
  - `xmax`：删除该行的事务ID(0表示未删除)
  - `cmin/cmax`：命令ID，用于同一事务内的可见性判断
- **可见性规则**：根据事务快照判断哪个版本的数据对当前事务可见
- **VACUUM**：清理旧版本数据，回收空间

**优点：**
- 读操作不阻塞写操作
- 写操作不阻塞读操作
- 提高并发性能

**缺点：**
- 需要额外的存储空间
- 需要定期VACUUM清理

### 4. 什么是VACUUM？为什么需要VACUUM？

**答案：**
VACUUM是PostgreSQL的垃圾回收机制，用于清理死元组(dead tuples)并回收存储空间。

**为什么需要VACUUM：**
- **MVCC产生的死元组**：UPDATE和DELETE操作不会立即删除旧版本数据
- **事务ID回卷**：防止事务ID耗尽导致的数据丢失
- **更新统计信息**：帮助查询优化器生成更好的执行计划
- **更新可见性映射**：提高查询性能

**VACUUM类型：**
```sql
-- 普通VACUUM：清理死元组，但不归还空间给操作系统
VACUUM table_name;

-- VACUUM FULL：完全重建表，归还空间，但会锁表
VACUUM FULL table_name;

-- VACUUM ANALYZE：清理并更新统计信息
VACUUM ANALYZE table_name;

-- 自动VACUUM：PostgreSQL会自动运行
-- 配置参数：autovacuum = on
```

**配置自动VACUUM：**
```sql
-- 在postgresql.conf中配置
autovacuum = on
autovacuum_max_workers = 3
autovacuum_naptime = 1min
autovacuum_vacuum_threshold = 50
autovacuum_analyze_threshold = 50
autovacuum_vacuum_scale_factor = 0.2
autovacuum_analyze_scale_factor = 0.1
```

### 5. 解释PostgreSQL的事务隔离级别

**答案：**
PostgreSQL支持SQL标准定义的4种隔离级别，但实际实现了3种：

| 隔离级别 | 脏读 | 不可重复读 | 幻读 | PostgreSQL实现 |
|---------|------|-----------|------|---------------|
| **Read Uncommitted** | 可能 | 可能 | 可能 | 实际等同于Read Committed |
| **Read Committed** | 不可能 | 可能 | 可能 | ✓ 默认级别 |
| **Repeatable Read** | 不可能 | 不可能 | 可能 | ✓ 通过MVCC防止幻读 |
| **Serializable** | 不可能 | 不可能 | 不可能 | ✓ 使用SSI技术 |

**设置隔离级别：**
```sql
-- 设置当前事务隔离级别
BEGIN TRANSACTION ISOLATION LEVEL REPEATABLE READ;

-- 设置会话默认隔离级别
SET SESSION CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL SERIALIZABLE;

-- 查看当前隔离级别
SHOW transaction_isolation;
```

**各级别示例：**

```sql
-- Read Committed（默认）
-- 事务1
BEGIN;
SELECT * FROM accounts WHERE id = 1;  -- balance = 100
-- 此时事务2修改并提交
SELECT * FROM accounts WHERE id = 1;  -- balance = 200（看到了新值）
COMMIT;

-- Repeatable Read
-- 事务1
BEGIN TRANSACTION ISOLATION LEVEL REPEATABLE READ;
SELECT * FROM accounts WHERE id = 1;  -- balance = 100
-- 此时事务2修改并提交
SELECT * FROM accounts WHERE id = 1;  -- balance = 100（仍是旧值）
COMMIT;

-- Serializable
-- 防止序列化异常，如果检测到冲突会回滚
BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE;
-- 执行操作
COMMIT;  -- 可能抛出serialization failure错误
```

### 6. PostgreSQL中的锁有哪些类型？

**答案：**

**表级锁：**

| 锁模式 | 说明 | 冲突 |
|--------|------|------|
| **ACCESS SHARE** | SELECT获取 | 只与ACCESS EXCLUSIVE冲突 |
| **ROW SHARE** | SELECT FOR UPDATE获取 | 与EXCLUSIVE和ACCESS EXCLUSIVE冲突 |
| **ROW EXCLUSIVE** | INSERT/UPDATE/DELETE获取 | 与SHARE及更高级别冲突 |
| **SHARE UPDATE EXCLUSIVE** | VACUUM/CREATE INDEX CONCURRENTLY | 与自己及更高级别冲突 |
| **SHARE** | CREATE INDEX | 与ROW EXCLUSIVE及更高级别冲突 |
| **SHARE ROW EXCLUSIVE** | 很少使用 | 与ROW EXCLUSIVE及更高级别冲突 |
| **EXCLUSIVE** | 阻止并发修改 | 与ROW SHARE及更高级别冲突 |
| **ACCESS EXCLUSIVE** | ALTER TABLE/DROP TABLE等 | 与所有锁冲突 |

**行级锁：**

```sql
-- FOR UPDATE：排他锁，阻止其他事务修改或锁定
SELECT * FROM accounts WHERE id = 1 FOR UPDATE;

-- FOR NO KEY UPDATE：类似FOR UPDATE，但允许其他事务获取FOR KEY SHARE锁
SELECT * FROM accounts WHERE id = 1 FOR NO KEY UPDATE;

-- FOR SHARE：共享锁，阻止其他事务修改，但允许读取
SELECT * FROM accounts WHERE id = 1 FOR SHARE;

-- FOR KEY SHARE：最弱的锁，只阻止FOR UPDATE
SELECT * FROM accounts WHERE id = 1 FOR KEY SHARE;
```

**显式锁定表：**

```sql
-- 显式获取表锁
BEGIN;
LOCK TABLE accounts IN ACCESS EXCLUSIVE MODE;
-- 执行操作
COMMIT;
```

**查看锁信息：**

```sql
-- 查看当前锁
SELECT 
    locktype,
    relation::regclass,
    mode,
    granted,
    pid
FROM pg_locks
WHERE NOT granted;

-- 查看阻塞关系
SELECT 
    blocked_locks.pid AS blocked_pid,
    blocked_activity.usename AS blocked_user,
    blocking_locks.pid AS blocking_pid,
    blocking_activity.usename AS blocking_user,
    blocked_activity.query AS blocked_statement,
    blocking_activity.query AS blocking_statement
FROM pg_catalog.pg_locks blocked_locks
JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
JOIN pg_catalog.pg_locks blocking_locks 
    ON blocking_locks.locktype = blocked_locks.locktype
    AND blocking_locks.database IS NOT DISTINCT FROM blocked_locks.database
    AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
    AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
    AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
    AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
    AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
    AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
    AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
    AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
    AND blocking_locks.pid != blocked_locks.pid
JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
WHERE NOT blocked_locks.granted;
```

**死锁处理：**

```sql
-- PostgreSQL会自动检测死锁并回滚其中一个事务
-- 配置死锁超时时间
SET deadlock_timeout = '1s';

-- 查看死锁日志
-- 在postgresql.conf中配置
log_lock_waits = on
deadlock_timeout = 1s
```

### 7. 什么是连接池？为什么需要连接池？

**答案：**

**连接池**是预先创建并维护一组数据库连接的技术，应用程序可以重复使用这些连接，而不是每次都创建新连接。

**为什么需要连接池：**
- **减少连接开销**：创建PostgreSQL连接需要fork进程，开销较大
- **限制连接数**：防止过多连接耗尽数据库资源
- **提高性能**：复用连接，减少延迟
- **连接管理**：统一管理连接的生命周期

**常用连接池：**

**1. PgBouncer（推荐）**
```ini
# pgbouncer.ini配置
[databases]
mydb = host=localhost port=5432 dbname=mydb

[pgbouncer]
listen_addr = *
listen_port = 6432
auth_type = md5
auth_file = /etc/pgbouncer/userlist.txt
pool_mode = transaction  # session, transaction, statement
max_client_conn = 1000
default_pool_size = 20
reserve_pool_size = 5
reserve_pool_timeout = 3
```

**池模式：**
- **session**：连接在客户端会话期间保持
- **transaction**：连接在事务期间保持（推荐）
- **statement**：连接在语句期间保持（最激进）

**2. Pgpool-II**
- 支持连接池
- 支持负载均衡
- 支持查询缓存
- 支持读写分离

**应用层连接池：**

```python
# Python - psycopg2连接池
from psycopg2 import pool

connection_pool = pool.SimpleConnectionPool(
    minconn=1,
    maxconn=20,
    host='localhost',
    database='mydb',
    user='user',
    password='password'
)

# 获取连接
conn = connection_pool.getconn()
try:
    cursor = conn.cursor()
    cursor.execute("SELECT * FROM users")
    results = cursor.fetchall()
finally:
    # 归还连接
    connection_pool.putconn(conn)
```

```java
// Java - HikariCP连接池
HikariConfig config = new HikariConfig();
config.setJdbcUrl("jdbc:postgresql://localhost:5432/mydb");
config.setUsername("user");
config.setPassword("password");
config.setMaximumPoolSize(20);
config.setMinimumIdle(5);
config.setConnectionTimeout(30000);

HikariDataSource ds = new HikariDataSource(config);
Connection conn = ds.getConnection();
```

**监控连接：**

```sql
-- 查看当前连接数
SELECT count(*) FROM pg_stat_activity;

-- 查看各数据库连接数
SELECT datname, count(*) 
FROM pg_stat_activity 
GROUP BY datname;

-- 查看连接详情
SELECT 
    pid,
    usename,
    application_name,
    client_addr,
    state,
    query
FROM pg_stat_activity
WHERE state != 'idle';

-- 查看最大连接数配置
SHOW max_connections;

-- 终止空闲连接
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE state = 'idle'
  AND state_change < NOW() - INTERVAL '1 hour';
```
