# gRPC & RPC 面试题

[TOC]

## 基础原理

### Q1: 什么是 RPC？与 RESTful API 有什么区别？
- **RPC (Remote Procedure Call)**：远程过程调用，像调用本地函数一样调用远程服务。
- **区别**：
  - **侧重点**：RPC 侧重于“动作/方法” (Function)，REST 侧重于“资源” (Resource)。
  - **协议**：RPC 通常基于 TCP 或 HTTP/2 (如 gRPC)；REST 基于 HTTP/1.1。
  - **传输效率**：RPC 常用二进制序列化 (Protobuf, Thrift)，体积小、解析快；REST 常用 JSON，文本冗余大。
  - **场景**：RPC 适合内部微服务高性能通信；REST 适合对外开放接口。

### Q2: gRPC 的核心特性是什么？
1.  **HTTP/2 协议**：支持多路复用、流式传输 (Streaming)、头部压缩。
2.  **Protobuf 序列化**：高效的二进制序列化协议，强类型，跨语言。
3.  **代码生成**：通过 `.proto` 文件自动生成客户端和服务端代码 (Stub)。
4.  **多语言支持**：Go, Java, Python, C++, Node.js 等。

### Q3: gRPC 有哪四种通信模式？
1.  **一元 RPC (Unary)**：客户端发一个请求，服务端回一个响应（类似普通 HTTP）。
2.  **服务端流式 (Server Streaming)**：客户端发一个请求，服务端回一串消息流（如股票行情）。
3.  **客户端流式 (Client Streaming)**：客户端发一串消息流，服务端回一个响应（如大文件上传）。
4.  **双向流式 (Bidirectional Streaming)**：双方建立长连接，随时互发消息（如聊天室、游戏）。

## 进阶机制

### Q4: Protobuf 为什么比 JSON 快/小？
1.  **二进制存储**：没有 `{}` `""` 等分隔符，紧凑。
2.  **Tag-Value 结构**：字段名被映射为整数 Tag，不再传输冗长的字段名字符串。
3.  **Varint 编码**：变长整数编码，小的整数占用字节更少（如 int32 存 1 只占 1 字节）。
4.  **ZigZag 编码**：优化负数存储。

### Q5: gRPC 如何做负载均衡？
gRPC 连接是长连接，普通的 L4 负载均衡（如 LVS）无法感知请求。
1.  **客户端负载均衡 (Client-side)**：客户端定期从注册中心（如 ETCD, Consul）拉取服务列表，内置轮询/随机算法选择节点。
2.  **代理负载均衡 (Proxy/Sidecar)**：使用 Envoy、Nginx 等支持 HTTP/2 的网关或 Service Mesh 进行转发。

### Q6: gRPC 的 Deadlines 和 Cancellation 是什么？
- **Deadlines (超时)**：客户端设置超时时间，会透传到服务端。如果超时，双方都会收到 `DEADLINE_EXCEEDED` 错误，服务端应自动取消计算以节省资源。
- **Cancellation (取消)**：客户端或服务端可以主动取消 RPC 调用，另一方会收到通知并中断操作。

### Q7: gRPC 如何处理 metadata？
类似于 HTTP Header。
- 客户端通过 `metadata.NewOutgoingContext` 发送。
- 服务端通过 `metadata.FromIncomingContext` 接收。
- 常用于传递 Token、TraceID 等认证或链路追踪信息。
