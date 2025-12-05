# ClickHouse命令篇

[TOC]

## 常用客户端命令

- 连接数据库：`clickhouse-client -h <host> --port 9000 -u <user> --password <password>`
- 交互式模式：直接输入 `clickhouse-client`
- 执行 SQL 文件：`clickhouse-client --multiquery < script.sql`
- 导出数据为 CSV：`clickhouse-client --query="SELECT * FROM table" --format=CSV > data.csv`

## 数据库与表管理

- 创建数据库：`CREATE DATABASE IF NOT EXISTS mydb`
- 切换数据库：`USE mydb`
- 查看建表语句：`SHOW CREATE TABLE mytable`
- 删除表：`DROP TABLE IF EXISTS mytable`

## 常用建表语句 (MergeTree)

```sql
CREATE TABLE hits_v1 (
    EventDate Date,
    EventTime DateTime,
    UserID UInt64,
    URL String,
    Referer String
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(EventDate)
ORDER BY (EventDate, intHash32(UserID))
SAMPLE BY intHash32(UserID)
SETTINGS index_granularity = 8192;
```

## 数据操作 (DML)

> 注意：ClickHouse 的更新和删除操作（Mutation）是重操作，应谨慎使用。

- 插入数据：`INSERT INTO mytable (col1, col2) VALUES (1, 'a'), (2, 'b')`
- 异步更新：`ALTER TABLE mytable UPDATE col1 = 10 WHERE id = 1`
- 异步删除：`ALTER TABLE mytable DELETE WHERE id = 1`
- 优化表（强制合并）：`OPTIMIZE TABLE mytable FINAL`

## 运维与监控

- 查看正在执行的查询：`SHOW PROCESSLIST`
- 杀死查询：`KILL QUERY WHERE query_id = '...'`
- 查看表占用空间：
  ```sql
  SELECT 
      table, 
      formatReadableSize(sum(bytes)) as size 
  FROM system.parts 
  WHERE active 
  GROUP BY table
  ```
- 查看 MergeTree 合并状态：`SELECT * FROM system.merges`
- 查看异步 Mutation 进度：`SELECT * FROM system.mutations WHERE is_done = 0`

## 高级查询函数

- 转换时间格式：`formatDateTime(EventTime, '%Y-%m-%d')`
- JSON 提取：`visitParamExtractString(json_str, 'key')`
- 数组操作：`arrayJoin([1, 2, 3])`
- 漏斗分析 (窗口模型)：`windowFunnel(3600)(timestamp, event='login', event='view', event='buy')`
