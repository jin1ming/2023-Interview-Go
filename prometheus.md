# Prometheus笔记

目录：

[TOC]

介绍：

- Prometheus 属于一站式监控告警平台，依赖少，功能齐全。
- Prometheus 支持对云或容器的监控，其他系统主要对主机监控。
- Prometheus 数据查询语句表现力更强大，内置更强大的统计函数。
- Prometheus 在数据存储扩展性以及持久性上没有 InfluxDB，OpenTSDB，Sensu 好。

## 一、组件

### 1. Prometheus Server

功能：

1. 对监控数据的获取、存储以及查询
   - 时序数据库
   - 提供PromQL
   - 联邦集群
2. 静态/动态管理监控目标

### 2. Exporters

Exporter将监控数据采集的端点通过HTTP服务的形式暴露给Prometheus Server，供其采集监控数据。

两类：

1. 直接采集：内置了对Prometheus监控支持的Endpoint，如cAdvisor，Kubernetes，Etcd，Gokit
2. 间接采集：通过Prometheus提供的Client Library编写该监控目标的监控采集程序。例如： Mysql Exporter，JMX Exporter，Consul Exporter等。

### 3. AlertManager

基于PromQL创建告警规则

支持企业微信、邮件、Slack、钉钉（Webhook）

### 4. PushGateway

PushGateway作为中转站，解决Prometheus无法与Exporter通信来Pull数据的问题。（不在一个子网或者防火墙）

### 5. Client Library

对接 Prometheus Server, 可以查询和上报数据。

## 二、PromQL

### 1. Metrics类型

- Counter（计数器）

  只增不减，推荐使用_total作为后缀

- Gauge（仪表盘）

  Gauge类型的指标侧重于反应系统的当前状态，这类数据可增可减。

- Histogram（直方图）

  Histogram和Summary主用用于统计和分析样本的分布情况。

- Summary（摘要）

### 2. PromQL操作符

> 我们通过promQL语句查询得到的值主要有以下两种：
>
> 1. "瞬时向量"  # 查询得到最新的值，(实时数据)通常用于报警、实时监控
>
> 2. "区间向量"  # 查询某一段时间范围内所有的样本值，多用于数据分析、预测

- 数字运算符与集合运算符

  1-4为数学运算符，5-6为集合运算符，按优先级从高到低：

  1. `^`
  2. `*, /, %`
  3. `+, -`
  4. `==, !=, <=, <, >=, >`
  5. `and, unless`
  6. `or`

- bool修饰符

  布尔运算符的默认行为是对时序数据进行过滤。

  或者需要真正的布尔结果时使用。

  例如：

  ```PromQL
  http_requests_total > bool 1000 
  # 返回的是0或1
  ```

  ```PromQL
  http_requests_total > 1000 
  # 返回的是大于1000的具体值，如1001
  ```

- 匹配模式

  - 一对一

    语法：

    ```PromQL
    vector1 <operator> vector2
    ```

    如果两边标签不一致，需要使用`on(label list)`或者`ignoring(label list）`来修改便签的匹配行为。

    ```PromQL
    <vector expr> <bin-op> ignoring(<label list>) <vector expr>
    <vector expr> <bin-op> on(<label list>) <vector expr>
    ```

    例：

    ```PromQL
    method_code:http_errors:rate5m{code="500"} / ignoring(code) method:http_requests:rate5m
    ```

  - 多对一 和 一对多

    指的是“一”侧的每一个向量元素可以与"多"侧的多个元素匹配的情况。

    在这种情况下，必须使用group修饰符：group_left或者group_right来确定哪一个向量具有更高的基数（充当“多”的角色）。

    ```PromQL
    <vector expr> <bin-op> ignoring(<label list>) group_left(<label list>) <vector expr>
    <vector expr> <bin-op> ignoring(<label list>) group_right(<label list>) <vector expr>
    <vector expr> <bin-op> on(<label list>) group_left(<label list>) <vector expr>
    <vector expr> <bin-op> on(<label list>) group_right(<label list>) <vector expr>
    ```

    例：

    ```PromQL
    # 一对多模式
    method_code:http_errors:rate5m / ignoring(code) group_left method:http_requests:rate5m
    ```

### 3. PromQL聚合操作

语法：

```PromQL
<aggr-op>([parameter,] <vector expression>) [without|by (<label list>)]
# 其中只有count_values, quantile, topk, bottomk支持参数(parameter)。
```

支持的聚合操作符：

- `sum` (求和)
- `min` (最小值)
- `max` (最大值)
- `avg` (平均值)
- `stddev` (标准差)
- `stdvar` (标准差异)
- `count` (计数)
- `count_values` (对value进行计数)
- `bottomk` (后n条时序)
- `topk` (前n条时序)
- `quantile` (分布统计)

### 4. PromQL内置函数

#### 4.1 计算Counter指标增长率

- increase

  例，获取时间序列最近两分钟的所有样本，计算出最近两分钟的增长量，最后除以时间120秒得到node_cpu样本在最近两分钟的平均增长率：

  ```PromQL
  increase(node_cpu[2m]) / 120
  ```

- rate

  rate函数可以直接计算区间向量v在时间窗口内平均增长速率：

  ```PromQL
  rate(node_cpu[2m])
  ```

  与increase例子效果等同。

- irate

  irate反应出的是瞬时增长率，解决rate或increase容易陷入“长尾问题”。

  ```PromQL
  irate(node_cpu[2m])
  ```

#### 4.2 预测Gauge指标变化趋势

- predict_linear

   对时间序列变化趋势做出预测，用预测的阈值来代替固定阈值。

  例，基于2小时的样本数据，来预测主机可用磁盘空间的是否在4个小时后被占满：

  ```PromQL
  predict_linear(node_filesystem_free{job="node"}[2h], 4 * 3600) < 0
  ```

#### 4.3 统计Histogram指标的分位数

- histogram_quantile

  TODO：这里之后的未完成

#### 4.4 动态标签替换

- label_replace

### 4个黄金指标

- 延迟：服务请求所需时间。
- 通讯量：监控当前系统的流量，用于衡量服务的容量需求。
- 错误：监控当前系统所有发生的错误请求，衡量当前系统错误发生的速率。
- 饱和度：衡量当前服务的饱和度。