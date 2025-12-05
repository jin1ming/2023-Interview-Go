# ClickHouse原理篇

[TOC]

## 概述

ClickHouse 是一个用于联机分析 (OLAP) 的列式数据库管理系统 (DBMS)。它由俄罗斯的 Yandex 开发，并于 2016 年开源。ClickHouse 以其**极高的查询性能**而闻名，特别适合处理海量数据的分析查询。

### 核心特性

1.  **列式存储 (Columnar Storage)**：数据按列存储，分析查询时只读取必要的列，IO 效率极高。
2.  **向量化执行 (Vectorized Execution)**：利用 SIMD 指令集，一次处理多行数据，大幅提升 CPU 利用率。
3.  **数据压缩**：列式存储使得数据压缩率极高（通常 10:1 以上），节省磁盘空间并减少 IO。
4.  **线性扩展**：支持分片和副本，可水平扩展至数百台节点，处理 PB 级数据。
5.  **实时写入**：支持高吞吐量的实时数据插入（建议批量写入）。

## 数据类型

ClickHouse 提供了丰富的数据类型，选择合适的类型对性能至关重要。

### 1. 基础类型
*   **数值类型**：`UInt8`, `UInt16`, `UInt32`, `UInt64`, `Int8`, `Int16`... 以及 `Float32`, `Float64`, `Decimal`。
*   **字符串类型**：
    *   `String`: 任意长度的字符串，不限长。
    *   `FixedString(N)`: 固定长度 N 的字符串，性能略高。
*   **时间类型**：`Date` (精确到天), `DateTime` (精确到秒), `DateTime64` (精确到亚秒)。

### 2. 复杂与特殊类型
*   **Array(T)**: 数组类型，ClickHouse 对数组操作支持极好（Array Join）。
*   **Nullable(T)**: 允许存储 Null 值。**注意**：使用 Nullable 会创建额外的 Null 掩码文件，严重影响性能，**尽量避免使用**，建议用默认值（如 0, -1, ''）代替 Null。
*   **LowCardinality(String)**: 低基数字符串。类似于字典编码，适合枚举值很少（如性别、状态）的场景，能极大减少存储空间并加速查询。
*   **Map(Key, Value)**: 键值对映射。
*   **Nested**: 嵌套数据结构，相当于行内的“子表”。

## 索引体系

ClickHouse 的索引设计非常独特，主要分为主键索引和跳数索引。

### 1. 主键索引 (Primary Key)
*   **稀疏索引**：如前所述，`ORDER BY` 字段即为主键。它不存储每行位置，只存储每个 Granule（8192行）的起止值。
*   **常驻内存**：因为稀疏，所以索引非常小，启动时加载到内存。
*   **作用**：用于范围查询（Range Query）时的快速剪枝（Data Pruning）。

### 2. 跳数索引 (Skipping Indexes)
二级索引，用于进一步过滤 Granule。
*   **minmax**：记录 Granule 内某一列的最大值和最小值。适合单调递增/递减的列（如时间）。
*   **set(k)**：记录 Granule 内列的所有唯一值（最多 k 个）。适合基数较小的列。
*   **bloom_filter**: 布隆过滤器索引。适合高基数（Unique 值很多）的列，如 UserID, UUID，用于快速判断“不存在”。
*   **tokenbf_v1**: 对字符串分词后的布隆过滤器。适合文本搜索（LIKE '%text%'）。

## 核心架构与原理

### 1. 存储架构

#### MergeTree 引擎家族
MergeTree 是 ClickHouse 最核心的存储引擎系列，类似于 HBase 的 LSM-Tree 思想，但针对 OLAP 做了优化。

*   **数据分区 (Partitioning)**：数据按规则（如按月）物理拆分成不同目录，查询时可剪枝。
*   **主键索引 (Primary Key)**：稀疏索引。不同于传统数据库的 B+ 树，ClickHouse 的主键索引存储的是每隔几千行（Granularity，默认 8192）的第一个值。这使得索引非常小，常驻内存。
*   **数据排序 (Sorting)**：数据在磁盘上严格按主键排序存储。
*   **合并 (Merge)**：后台不断将小的 Part 文件合并成大的 Part，类似 LSM 的 Compaction。

#### 分布式架构
*   **分布式表 (Distributed Table)**：这是一个逻辑视图，不存储数据，只负责将查询路由到具体的本地表（Local Table）。
*   **复制表 (Replicated Table)**：基于 ZooKeeper 实现数据的高可用复制。

## 高级特性

### 1. 投影 (Projections)
类似于物化视图，但更轻量。可以为同一张表创建多种不同的物理排序存储，查询时自动选择最优的投影。

### 2. TTL (Time To Live)
原生支持数据生命周期管理。可以按列或按行设置 TTL，到期自动删除或移动到冷存储（Tiered Storage）。

### 3. 字典 (Dictionaries)
将维度表加载到内存中，用于加速 JOIN 查询。

## JSON 支持

ClickHouse 对 JSON 的支持经历了多个版本的迭代，目前主要有三种方式：

### 1. String + JSON 函数 (传统方式)
将 JSON 作为普通字符串存储，使用 `visitParamExtractString` 或 `JSONExtract` 系列函数进行查询。
*   **优点**：兼容性好，无需预定义 Schema。
*   **缺点**：查询时需要解析字符串，CPU 开销大。

### 2. Object('JSON') 类型 (新特性)
在 22.6+ 版本引入。它会自动推断 JSON 结构，并将每个键值对拆分成独立的**物理列**存储。
*   **优点**：查询性能极高（只读取需要的子列），支持动态 Schema。
*   **示例**：
    ```sql
    CREATE TABLE logs (data Object('JSON')) ENGINE = MergeTree ORDER BY tuple();
    INSERT INTO logs VALUES ('{"id": 1, "status": "ok"}');
    ```

### 3. Map(String, String)
适合 Key-Value 结构扁平且类型统一的 JSON。Schema 固定，性能稳定。

## 常见面试题

### Q1: ClickHouse 为什么这么快？
1.  **算法级优化**：大量使用特定场景下的高效算法（如 HyperLogLog, uniqCombined）。
2.  **底层优化**：列式存储 + 向量化执行 + SIMD 指令集。
3.  **存储引擎**：MergeTree 的稀疏索引和数据排序机制，极致减少磁盘 IO。
4.  **并行处理**：单条查询就能利用单机所有 CPU 核心（多线程）以及集群所有节点（分布式）。

### Q2: ClickHouse 适合什么场景？不适合什么场景？
*   **适合**：
    *   海量数据（亿级~PB级）的宽表分析。
    *   写少读多，或者只追加写（Append-only）。
    *   要求亚秒级响应的聚合查询。
*   **不适合**：
    *   高频的点对点更新和删除（Mutation 操作很重）。
    *   高并发的 OLTP 事务（不支持完整的 ACID 事务）。
    *   作为 Key-Value 存储使用（稀疏索引不适合点查）。

### Q3: ClickHouse 如何处理 JSON 数据？
*   **高性能场景**：在 ETL 阶段将 JSON 摊平（Flatten）成独立的列。
*   **灵活性场景**：使用 `Object('JSON')` 类型（推荐）或 `String` 类型配合 JSON 函数。

### Q4: MergeTree 的合并机制是怎样的？
数据写入时生成临时的 Part 文件，后台线程不定时将重叠的 Parts 进行合并。合并过程中会执行数据的物理排序和去重（如果是 ReplacingMergeTree）。这也是 ClickHouse "最终一致性" 的来源。

### Q4: ClickHouse 的稀疏索引原理是什么？为什么不使用 B+ 树？
*   **原理**：ClickHouse 的主键索引（Primary Key）不存储每一行的位置，而是存储**每一个颗粒（Granule，默认 8192 行）**的第一个值（标记值）。
*   **查询过程**：查询时，通过二分查找定位到可能的 Granule 范围，然后加载整个 Granule 的数据进行扫描过滤。
*   **为什么不用 B+ 树**：B+ 树适合点查，但对于海量分析查询，索引太大无法常驻内存。稀疏索引极小（亿级数据索引仅几 MB），可常驻内存，且适合范围扫描。

### Q5: ReplacingMergeTree 如何实现数据去重？有什么局限性？
*   **原理**：在后台合并（Merge）阶段，对于具有相同排序键（ORDER BY）的行，只保留版本号最大的一行。
*   **局限性**：
    1.  **最终一致性**：去重只在合并时发生，而合并时机不确定。查询时可能会看到重复数据（除非使用 `FINAL` 关键字，但性能差）。
    2.  **分片限制**：只能在同一分片内去重，无法跨分片去重。

### Q6: ClickHouse 的 JOIN 为什么容易 OOM？如何优化？
*   **原因**：ClickHouse 默认使用 Hash Join。它会将**右表**全量加载到内存构建 Hash 表。如果右表过大，内存就会爆掉。
*   **优化策略**：
    1.  **大表在左，小表在右**：始终把小表放在右边（Right Table）。
    2.  **使用 GLOBAL JOIN**：在分布式查询中，避免将右表分发到所有节点，而是预先在发起节点聚合好右表（适用于右表较小的情况）。
    3.  **使用字典 (Dictionary)**：如果是维度表 JOIN，直接用字典代替 JOIN，性能提升巨大。
    4.  **Nested Loop Join**：新版本支持，牺牲性能换取不 OOM。

### Q7: ClickHouse 如何实现高可用 (Replication)？
*   **依赖组件**：ZooKeeper / ClickHouse Keeper。
*   **机制**：使用 `ReplicatedMergeTree` 引擎。数据写入任意一个副本后，通过 ZooKeeper 分发日志（Log），其他副本根据日志异步拉取数据块。
*   **特点**：多主架构（Multi-Master），任意节点都可读写（但在不同分片上通常会做读写分离）。

### Q8: Update 和 Delete 操作为什么慢？
ClickHouse 是为追加写（Append-only）设计的。`ALTER TABLE ... UPDATE/DELETE` 是 Mutation 操作。
*   **重写机制**：Mutation 会强制重写整个数据分区（Part）的所有文件，代价极高。
*   **建议**：尽量避免单条更新，改用 `ReplacingMergeTree` 覆盖写，或 `CollapsingMergeTree` 逻辑删除。
