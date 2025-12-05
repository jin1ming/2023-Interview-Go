# MQTT 协议面试题

[TOC]

## 基础概念

### Q1: 什么是 MQTT？适用场景？
- **定义**：Message Queuing Telemetry Transport，基于发布/订阅 (Pub/Sub) 模式的轻量级通讯协议，构建于 TCP/IP 之上。
- **特点**：轻量（Header 最小仅 2 字节）、低带宽、低功耗、支持不可靠网络。
- **场景**：物联网 (IoT)、车联网、移动即时通讯、智能家居。

### Q2: MQTT 的核心组件有哪些？
1.  **Publisher (发布者)**：发送消息的客户端。
2.  **Subscriber (订阅者)**：接收消息的客户端。
3.  **Broker (代理/服务端)**：核心中转站，负责接收发布者的消息，并路由给匹配的订阅者（如 EMQX, Mosquitto）。
4.  **Topic (主题)**：消息的路由标签，支持层级结构（如 `home/livingroom/temp`）。

### Q3: MQTT 的 QoS (服务质量) 等级有哪些？
这是 MQTT 最重要的特性之一。
- **QoS 0 (At most once)**：**最多一次**。发完即忘，不保证到达。适用于传感器数据（丢一两条无所谓）。
- **QoS 1 (At least once)**：**至少一次**。保证到达，但可能重复。发送方需收到 `PUBACK`，否则重发。需要接收方做幂等处理。
- **QoS 2 (Exactly once)**：**恰好一次**。保证到达且不重复。机制复杂（四次握手），开销大。适用于支付、计费等关键数据。

## 进阶机制

### Q4: Topic 通配符有哪些？
- **`+` (单层通配)**：匹配一层。
  - `home/+/temp` 匹配 `home/livingroom/temp`, `home/bedroom/temp`。
  - 不匹配 `home/livingroom/temp/1`。
- **`#` (多层通配)**：匹配后续所有层，必须放在末尾。
  - `home/#` 匹配 `home/livingroom`, `home/livingroom/temp`, `home/a/b/c`。

### Q5: 什么是保留消息 (Retained Message)？
- 发布消息时设置 `Retain=true`，Broker 会存储该 Topic 下的**最后一条**消息。
- **作用**：新上线的订阅者能立即收到该 Topic 的最新状态，而不需要等待下一次发布。
- **场景**：设备上线查看开关的最新状态。

### Q6: 什么是遗嘱消息 (LWT - Last Will and Testament)？
- 客户端连接 (CONNECT) 时可以预设一条“遗嘱消息”。
- **触发条件**：当客户端**非正常断开**（网络异常、超时、断电）时，Broker 会自动发布这条消息。
- **作用**：通知其他设备“我掉线了”。
- **注意**：如果是客户端主动发送 `DISCONNECT` 断开，遗嘱不会触发。

### Q7: MQTT Keep Alive 与心跳机制？
- 客户端在 `CONNECT` 时设置 `Keep Alive` 时间（如 60s）。
- 如果在时间内没有数据传输，客户端必须发送 `PINGREQ` 包，Broker 回复 `PINGRESP`。
- 若 Broker 超过 1.5 倍 Keep Alive 时间未收到包，则断开连接（并触发遗嘱）。

### Q8: MQTT session (会话) 持久化？
- `Clean Session = false`：客户端断线重连后，Broker 会恢复之前的订阅关系，并补发离线期间 QoS > 0 的未读消息。
- `Clean Session = true`：每次连接都是全新的，不保留历史订阅和离线消息。
