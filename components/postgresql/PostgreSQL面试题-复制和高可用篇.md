# PostgreSQL面试题 - 复制和高可用篇

[TOC]

## 复制

### 1. PostgreSQL支持哪些复制方式？

**答案：**

**1. 流复制（Streaming Replication）**
- 基于WAL日志的物理复制
- 主从架构，异步或同步
- 最常用的复制方式

**2. 逻辑复制（Logical Replication）**
- 基于发布/订阅模式
- 可以复制部分表
- 支持不同版本间复制
- PostgreSQL 10+

**3. 级联复制（Cascading Replication）**
- 从库可以作为其他从库的主库
- 减轻主库压力

### 2. 如何配置流复制？

**答案：**

**主库配置：**

```bash
# postgresql.conf
wal_level = replica
max_wal_senders = 10
wal_keep_size = 1GB  # PostgreSQL 13+
# 或 wal_keep_segments = 64  # PostgreSQL 12及以下

# 同步复制（可选）
synchronous_standby_names = 'standby1,standby2'

# pg_hba.conf
# 允许从库连接
host replication replicator 192.168.1.0/24 md5
```

```sql
-- 创建复制用户
CREATE ROLE replicator WITH REPLICATION LOGIN PASSWORD 'password';
```

**从库配置：**

```bash
# 1. 使用pg_basebackup创建从库
pg_basebackup -h master_host -D /var/lib/postgresql/data -U replicator -P -v -R

# -R 参数会自动创建standby.signal文件和配置

# 2. 手动配置（如果没用-R）
# 创建standby.signal文件
touch /var/lib/postgresql/data/standby.signal

# postgresql.conf或postgresql.auto.conf
primary_conninfo = 'host=master_host port=5432 user=replicator password=password'
restore_command = 'cp /archive/%f %p'  # 可选，用于归档恢复
```

**验证复制：**

```sql
-- 主库查看复制状态
SELECT * FROM pg_stat_replication;

-- 从库查看复制延迟
SELECT now() - pg_last_xact_replay_timestamp() AS replication_delay;

-- 从库查看是否在恢复模式
SELECT pg_is_in_recovery();
```

### 3. 同步复制和异步复制的区别？

**答案：**

**异步复制（Asynchronous）：**
- 主库提交事务后立即返回，不等待从库确认
- 性能好，延迟低
- 可能丢失数据（主库宕机时）
- 默认模式

**同步复制（Synchronous）：**
- 主库提交事务后等待从库确认
- 数据安全性高，不会丢失数据
- 性能较差，延迟高
- 需要配置synchronous_standby_names

**配置同步复制：**

```sql
-- postgresql.conf
synchronous_commit = on  -- on, remote_apply, remote_write, local, off
synchronous_standby_names = 'FIRST 1 (standby1, standby2)'

-- FIRST 1：至少1个从库确认
-- ANY 1：任意1个从库确认
-- standby1, standby2：从库名称（application_name）
```

**从库设置application_name：**

```bash
# postgresql.conf或primary_conninfo
primary_conninfo = 'host=master_host port=5432 user=replicator password=password application_name=standby1'
```

**synchronous_commit级别：**
- `off`：不等待WAL写入（最快，可能丢数据）
- `local`：等待本地WAL写入
- `remote_write`：等待从库接收WAL（未刷盘）
- `on`：等待从库WAL刷盘（默认）
- `remote_apply`：等待从库应用WAL（最安全，最慢）

### 4. 如何配置逻辑复制？

**答案：**

**发布端（Publisher）配置：**

```sql
-- 1. 修改配置
-- postgresql.conf
wal_level = logical
max_replication_slots = 10
max_wal_senders = 10

-- 重启数据库
-- 2. 创建发布
CREATE PUBLICATION my_publication FOR ALL TABLES;

-- 或发布特定表
CREATE PUBLICATION my_publication FOR TABLE users, orders;

-- 或发布特定操作
CREATE PUBLICATION my_publication FOR TABLE users 
WITH (publish = 'insert,update');

-- 查看发布
SELECT * FROM pg_publication;
SELECT * FROM pg_publication_tables;
```

**订阅端（Subscriber）配置：**

```sql
-- 1. 创建相同结构的表
CREATE TABLE users (...);
CREATE TABLE orders (...);

-- 2. 创建订阅
CREATE SUBSCRIPTION my_subscription
CONNECTION 'host=publisher_host port=5432 dbname=mydb user=replicator password=password'
PUBLICATION my_publication;

-- 查看订阅
SELECT * FROM pg_subscription;
SELECT * FROM pg_subscription_rel;

-- 查看复制状态
SELECT * FROM pg_stat_subscription;
```

**管理逻辑复制：**

```sql
-- 禁用订阅
ALTER SUBSCRIPTION my_subscription DISABLE;

-- 启用订阅
ALTER SUBSCRIPTION my_subscription ENABLE;

-- 刷新订阅（重新同步表列表）
ALTER SUBSCRIPTION my_subscription REFRESH PUBLICATION;

-- 删除订阅
DROP SUBSCRIPTION my_subscription;

-- 删除发布
DROP PUBLICATION my_publication;
```

**逻辑复制优点：**
- 可以复制部分表
- 支持不同版本间复制
- 支持双向复制
- 可以在从库写入数据

**逻辑复制缺点：**
- 不复制DDL（需要手动同步）
- 不复制序列
- 性能不如流复制

### 5. 如何实现读写分离？

**答案：**

**方案1：应用层实现**

```python
# Python示例
import psycopg2

# 主库连接（写）
master_conn = psycopg2.connect(
    host='master_host',
    database='mydb',
    user='user',
    password='password'
)

# 从库连接（读）
slave_conn = psycopg2.connect(
    host='slave_host',
    database='mydb',
    user='user',
    password='password'
)

# 写操作用主库
cursor = master_conn.cursor()
cursor.execute("INSERT INTO users (name) VALUES ('Alice')")
master_conn.commit()

# 读操作用从库
cursor = slave_conn.cursor()
cursor.execute("SELECT * FROM users")
results = cursor.fetchall()
```

**方案2：使用连接池（PgBouncer）**

```ini
# pgbouncer.ini
[databases]
mydb_master = host=master_host port=5432 dbname=mydb
mydb_slave = host=slave_host port=5432 dbname=mydb

[pgbouncer]
listen_addr = *
listen_port = 6432
pool_mode = transaction
```

**方案3：使用Pgpool-II**

```bash
# pgpool.conf
backend_hostname0 = 'master_host'
backend_port0 = 5432
backend_weight0 = 0  # 不用于负载均衡
backend_flag0 = 'ALWAYS_MASTER'

backend_hostname1 = 'slave1_host'
backend_port1 = 5432
backend_weight1 = 1
backend_flag1 = 'DISALLOW_TO_FAILOVER'

backend_hostname2 = 'slave2_host'
backend_port2 = 5432
backend_weight2 = 1
backend_flag2 = 'DISALLOW_TO_FAILOVER'

load_balance_mode = on
master_slave_mode = on
master_slave_sub_mode = 'stream'
```

**方案4：使用HAProxy**

```bash
# haproxy.cfg
frontend pgsql_front
    bind *:5432
    default_backend pgsql_back

backend pgsql_back
    option pgsql-check user haproxy
    server master master_host:5432 check
    server slave1 slave1_host:5432 check backup
    server slave2 slave2_host:5432 check backup
```

## 高可用

### 6. 如何实现PostgreSQL高可用？

**答案：**

**常见高可用方案：**

**1. 主从 + 手动故障转移**
- 最简单的方案
- 需要手动提升从库为主库
- 停机时间较长

**2. 主从 + 自动故障转移（推荐）**
- 使用工具自动检测和切换
- 常用工具：Patroni、repmgr、Pacemaker

**3. 共享存储**
- 使用共享存储（SAN、NAS）
- 快速切换，但成本高

**4. 多主复制**
- BDR（Bi-Directional Replication）
- 复杂度高，适合特殊场景

### 7. 如何使用Patroni实现自动故障转移？

**答案：**

**Patroni架构：**
- Patroni：PostgreSQL管理工具
- etcd/Consul/ZooKeeper：分布式配置存储
- HAProxy：负载均衡

**安装Patroni：**

```bash
# 安装依赖
pip install patroni[etcd]
pip install psycopg2-binary

# 配置Patroni
# /etc/patroni/patroni.yml
scope: postgres-cluster
namespace: /db/
name: node1

restapi:
  listen: 0.0.0.0:8008
  connect_address: node1:8008

etcd:
  host: etcd_host:2379

bootstrap:
  dcs:
    ttl: 30
    loop_wait: 10
    retry_timeout: 10
    maximum_lag_on_failover: 1048576
    postgresql:
      use_pg_rewind: true
      parameters:
        wal_level: replica
        hot_standby: on
        max_wal_senders: 10
        max_replication_slots: 10

  initdb:
    - encoding: UTF8
    - data-checksums

  pg_hba:
    - host replication replicator 0.0.0.0/0 md5
    - host all all 0.0.0.0/0 md5

postgresql:
  listen: 0.0.0.0:5432
  connect_address: node1:5432
  data_dir: /var/lib/postgresql/data
  bin_dir: /usr/lib/postgresql/14/bin
  authentication:
    replication:
      username: replicator
      password: password
    superuser:
      username: postgres
      password: password

tags:
  nofailover: false
  noloadbalance: false
  clonefrom: false
  nosync: false
```

**启动Patroni：**

```bash
# 启动Patroni
patroni /etc/patroni/patroni.yml

# 查看集群状态
patronictl -c /etc/patroni/patroni.yml list

# 手动切换主库
patronictl -c /etc/patroni/patroni.yml switchover

# 手动故障转移
patronictl -c /etc/patroni/patroni.yml failover
```

**配置HAProxy：**

```bash
# haproxy.cfg
global
    maxconn 100

defaults
    log global
    mode tcp
    retries 2
    timeout client 30m
    timeout connect 4s
    timeout server 30m
    timeout check 5s

listen stats
    mode http
    bind *:7000
    stats enable
    stats uri /

listen postgres
    bind *:5432
    option httpchk
    http-check expect status 200
    default-server inter 3s fall 3 rise 2 on-marked-down shutdown-sessions
    server node1 node1:5432 maxconn 100 check port 8008
    server node2 node2:5432 maxconn 100 check port 8008
    server node3 node3:5432 maxconn 100 check port 8008

listen postgres_replica
    bind *:5433
    option httpchk GET /replica
    http-check expect status 200
    default-server inter 3s fall 3 rise 2 on-marked-down shutdown-sessions
    server node1 node1:5432 maxconn 100 check port 8008
    server node2 node2:5432 maxconn 100 check port 8008
    server node3 node3:5432 maxconn 100 check port 8008
```

### 8. 如何进行备份和恢复？

**答案：**

**物理备份（推荐）：**

**1. pg_basebackup（在线备份）**

```bash
# 完整备份
pg_basebackup -h localhost -D /backup/base -U postgres -P -v -R -X stream

# -D：备份目录
# -P：显示进度
# -v：详细输出
# -R：创建恢复配置
# -X stream：流式传输WAL

# 压缩备份
pg_basebackup -h localhost -D /backup/base -U postgres -P -v -Z 9 -F tar

# 恢复
# 1. 停止PostgreSQL
# 2. 清空数据目录
# 3. 解压备份到数据目录
# 4. 启动PostgreSQL
```

**2. 文件系统快照**

```bash
# 使用LVM快照
# 1. 执行检查点
psql -c "SELECT pg_start_backup('snapshot');"

# 2. 创建快照
lvcreate -L 10G -s -n pg_snapshot /dev/vg/pg_data

# 3. 结束备份
psql -c "SELECT pg_stop_backup();"

# 4. 挂载快照并复制
mount /dev/vg/pg_snapshot /mnt/snapshot
cp -a /mnt/snapshot /backup/
```

**逻辑备份：**

**1. pg_dump（单个数据库）**

```bash
# 备份数据库
pg_dump -h localhost -U postgres -d mydb -F c -f mydb.dump

# -F c：自定义格式（推荐）
# -F p：纯文本SQL
# -F t：tar格式

# 备份特定表
pg_dump -h localhost -U postgres -d mydb -t users -F c -f users.dump

# 备份schema
pg_dump -h localhost -U postgres -d mydb -n public -F c -f schema.dump

# 只备份schema结构
pg_dump -h localhost -U postgres -d mydb -s -f schema.sql

# 只备份数据
pg_dump -h localhost -U postgres -d mydb -a -f data.sql

# 恢复
pg_restore -h localhost -U postgres -d mydb mydb.dump

# 恢复到新数据库
createdb newdb
pg_restore -h localhost -U postgres -d newdb mydb.dump

# 并行恢复
pg_restore -h localhost -U postgres -d mydb -j 4 mydb.dump
```

**2. pg_dumpall（所有数据库）**

```bash
# 备份所有数据库
pg_dumpall -h localhost -U postgres -f all.sql

# 只备份全局对象（角色、表空间）
pg_dumpall -h localhost -U postgres -g -f globals.sql

# 恢复
psql -h localhost -U postgres -f all.sql
```

**持续归档和PITR（时间点恢复）：**

```bash
# 1. 配置WAL归档
# postgresql.conf
wal_level = replica
archive_mode = on
archive_command = 'cp %p /archive/%f'

# 2. 基础备份
pg_basebackup -h localhost -D /backup/base -U postgres -P -v -X fetch

# 3. 恢复到特定时间点
# 停止PostgreSQL
# 恢复基础备份
cp -a /backup/base/* /var/lib/postgresql/data/

# 创建recovery.signal
touch /var/lib/postgresql/data/recovery.signal

# 配置恢复参数
# postgresql.conf或postgresql.auto.conf
restore_command = 'cp /archive/%f %p'
recovery_target_time = '2024-01-15 10:30:00'
# 或
recovery_target_xid = '12345'
# 或
recovery_target_name = 'before_drop_table'

# 启动PostgreSQL，会自动恢复到指定时间点
```

**备份策略建议：**

```bash
# 每日全量备份 + WAL归档
# 备份脚本示例
#!/bin/bash
DATE=$(date +%Y%m%d)
BACKUP_DIR=/backup/$DATE

# 全量备份
pg_basebackup -h localhost -D $BACKUP_DIR -U postgres -P -v -X stream

# 压缩
tar -czf $BACKUP_DIR.tar.gz $BACKUP_DIR
rm -rf $BACKUP_DIR

# 删除7天前的备份
find /backup -name "*.tar.gz" -mtime +7 -delete

# 定时任务
# crontab -e
# 0 2 * * * /path/to/backup.sh
```

### 9. 如何监控PostgreSQL？

**答案：**

**系统视图：**

```sql
-- 查看活动连接
SELECT * FROM pg_stat_activity;

-- 查看数据库统计
SELECT * FROM pg_stat_database;

-- 查看表统计
SELECT * FROM pg_stat_user_tables;

-- 查看索引统计
SELECT * FROM pg_stat_user_indexes;

-- 查看复制状态
SELECT * FROM pg_stat_replication;

-- 查看锁信息
SELECT * FROM pg_locks;

-- 查看慢查询
SELECT * FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10;
```

**监控工具：**

**1. pgAdmin**
- 官方图形化管理工具
- 支持监控、查询、备份等

**2. pg_stat_statements**
```sql
CREATE EXTENSION pg_stat_statements;
-- 查看最慢的查询
SELECT query, calls, total_time, mean_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 20;
```

**3. Prometheus + Grafana**
```bash
# 安装postgres_exporter
docker run -d \
  -p 9187:9187 \
  -e DATA_SOURCE_NAME="postgresql://user:password@localhost:5432/postgres?sslmode=disable" \
  prometheuscommunity/postgres-exporter

# Prometheus配置
scrape_configs:
  - job_name: 'postgresql'
    static_configs:
      - targets: ['localhost:9187']
```

**4. pgBadger**
```bash
# 分析PostgreSQL日志
pgbadger /var/log/postgresql/postgresql.log -o report.html
```

**关键监控指标：**
- 连接数
- QPS/TPS
- 缓存命中率
- 复制延迟
- 锁等待
- 慢查询
- 表和索引膨胀
- 磁盘使用率
