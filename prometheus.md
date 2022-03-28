# Prometheus笔记

目录：

[TOC]

介绍：

- Prometheus 属于一站式监控告警平台，依赖少，功能齐全。
- Prometheus 支持对云或容器的监控，其他系统主要对主机监控。
- Prometheus 数据查询语句表现力更强大，内置更强大的统计函数。
- Prometheus 在数据存储扩展性以及持久性上没有 InfluxDB，OpenTSDB，Sensu 好。

## 组件

### Prometheus Server

功能：

1. 对监控数据的获取、存储以及查询
   - 时序数据库
   - 提供PromQL
   - 联邦集群
2. 静态/动态管理监控目标

### Exporters

Exporter将监控数据采集的端点通过HTTP服务的形式暴露给Prometheus Server，供其采集监控数据。

两类：

1. 直接采集：内置了对Prometheus监控支持的Endpoint，如cAdvisor，Kubernetes，Etcd，Gokit
2. 间接采集：通过Prometheus提供的Client Library编写该监控目标的监控采集程序。例如： Mysql Exporter，JMX Exporter，Consul Exporter等。

### AlertManager

基于PromQL创建告警规则

支持企业微信、邮件、Slack、钉钉（Webhook）

### PushGateway

PushGateway作为中转站，解决Prometheus无法与Exporter通信来Pull数据的问题。（不在一个子网或者防火墙）

### Client Library

对接 Prometheus Server, 可以查询和上报数据。

## PromQL

### Metrics类型



### 4个黄金指标

- 延迟：服务请求所需时间。
- 通讯量：监控当前系统的流量，用于衡量服务的容量需求。
- 错误：监控当前系统所有发生的错误请求，衡量当前系统错误发生的速率。
- 饱和度：衡量当前服务的饱和度。