# PostgreSQL面试题汇总

本目录包含PostgreSQL数据库的常见面试题及详细答案，内容涵盖基础概念、索引优化、高级特性、复制和高可用等方面。

## 目录结构

### 1. [基础篇](./PostgreSQL面试题-基础篇.md)
- PostgreSQL简介和特点
- PostgreSQL与MySQL的区别
- MVCC（多版本并发控制）
- VACUUM机制
- 事务隔离级别
- 锁机制
- 连接池

**核心知识点：**
- MVCC实现原理
- 事务隔离级别的区别
- 表级锁和行级锁
- 连接池的作用和配置

### 2. [索引和优化篇](./PostgreSQL面试题-索引和优化篇.md)
- 索引类型（B-tree、Hash、GiST、GIN、BRIN等）
- 索引优化策略
- 慢查询分析和优化
- JOIN类型和优化
- EXPLAIN使用

**核心知识点：**
- 各种索引类型的适用场景
- GIN vs GiST的选择
- 部分索引、表达式索引、覆盖索引
- EXPLAIN输出解读
- JOIN执行策略

### 3. [高级特性篇](./PostgreSQL面试题-高级特性篇.md)
- JSONB类型
- 数组类型
- 范围类型
- CTE（公共表表达式）
- 递归查询
- 窗口函数
- 全文搜索
- 表分区
- 外部数据包装器（FDW）
- 扩展（PostGIS、pg_trgm等）
- 性能监控

**核心知识点：**
- JSONB的操作和索引
- 数组的使用场景
- 范围类型的实际应用
- 递归CTE的使用
- 窗口函数的常见用法
- 表分区策略

### 4. [复制和高可用篇](./PostgreSQL面试题-复制和高可用篇.md)
- 流复制
- 逻辑复制
- 同步复制 vs 异步复制
- 读写分离
- 高可用方案
- Patroni自动故障转移
- 备份和恢复
- PITR（时间点恢复）
- 监控

**核心知识点：**
- 流复制配置
- 逻辑复制的优缺点
- 高可用架构设计
- 备份策略
- 监控指标

## 学习建议

### 初级（1-2年经验）
重点掌握：
- 基础概念（MVCC、事务、锁）
- 常用索引类型（B-tree、GIN）
- 基本查询优化
- 简单的备份恢复

### 中级（3-5年经验）
重点掌握：
- 各种索引类型的选择和优化
- 复杂查询优化（JOIN、子查询、CTE）
- JSONB、数组等高级数据类型
- 流复制配置
- 性能监控和调优

### 高级（5年以上经验）
重点掌握：
- 高可用架构设计
- 分区表设计
- 逻辑复制
- 自动故障转移（Patroni）
- 大规模数据库优化
- 源码级别的理解

## 常见面试问题

### 必问问题
1. PostgreSQL与MySQL的区别？
2. 什么是MVCC？如何实现？
3. PostgreSQL有哪些索引类型？
4. 如何优化慢查询？
5. 如何配置主从复制？

### 进阶问题
1. GIN索引和GiST索引的区别？
2. 如何实现读写分离？
3. 什么是VACUUM？为什么需要？
4. 如何实现高可用？
5. JSONB和JSON的区别？

### 高级问题
1. 如何设计分区表？
2. 逻辑复制的原理和应用场景？
3. 如何实现自动故障转移？
4. 如何进行PITR恢复？
5. 如何监控和调优PostgreSQL？

## 实践建议

1. **搭建测试环境**
   - 使用Docker快速搭建PostgreSQL
   - 配置主从复制环境
   - 尝试各种索引类型

2. **性能测试**
   - 使用pgbench进行压力测试
   - 分析EXPLAIN输出
   - 优化慢查询

3. **高可用实践**
   - 配置流复制
   - 使用Patroni实现自动故障转移
   - 模拟故障场景

4. **数据类型实践**
   - 使用JSONB存储非结构化数据
   - 使用数组类型
   - 使用范围类型实现预订系统

5. **监控和调优**
   - 配置pg_stat_statements
   - 使用Prometheus + Grafana监控
   - 分析慢查询日志

## 参考资源

### 官方文档
- [PostgreSQL官方文档](https://www.postgresql.org/docs/)
- [PostgreSQL Wiki](https://wiki.postgresql.org/)

### 推荐书籍
- 《PostgreSQL实战》
- 《PostgreSQL技术内幕：查询优化深度探索》
- 《PostgreSQL高可用实战》

### 在线资源
- [PostgreSQL Tutorial](https://www.postgresqltutorial.com/)
- [Postgres Weekly](https://postgresweekly.com/)
- [Planet PostgreSQL](https://planet.postgresql.org/)

### 工具
- pgAdmin：图形化管理工具
- DBeaver：通用数据库工具
- pgBadger：日志分析工具
- pgBouncer：连接池
- Patroni：高可用方案

## 贡献

欢迎提交Issue和Pull Request来完善这份面试题集。

## 许可

本项目采用MIT许可证。
