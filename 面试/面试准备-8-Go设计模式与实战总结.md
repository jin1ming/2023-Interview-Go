# 面试准备 - Go设计模式与项目实战

本文档总结了 Go 语言中地道的设计模式实现，并结合 `Motion`, `Carina`, `Dipper` 等实际项目中的代码案例，展示如何在面试中体现架构设计能力。

## 一、面试回答策略

当面试官问到：“**你在项目中用到了哪些设计模式？**”时，建议按以下层次回答：

1.  **架构层 (解耦与扩展)**：
    > “在 `Motion` 运动引擎中，为了支持数十种运动项目的快速迭代，我设计了基于**注册表+反射工厂**的插件化架构，配合**状态机 (FSM)** 管理复杂的运动生命周期，这极大地提高了系统的可维护性。”

2.  **工程层 (标准与规范)**：
    > “在 `Carina` Web 服务中，我广泛使用**中间件模式**处理鉴权和审计，并严格利用 `sync.Once` 保证全局单例的安全性。”

3.  **跨语言/系统层**：
    > “在涉及硬件交互的 `Dipper` 项目中，我使用了**代理模式**，通过 gRPC 将底层的 C++ SDK 能力封装代理给上层 Go 业务，实现了软硬解耦。”

4.  **高并发治理 (秀肌肉)**：
    > “此外，在处理高并发热点数据时，我会引入 **Singleflight** 模式来防止缓存击穿，并利用 **Channel 信号量**机制控制并发水位，确保系统稳定性。”

## 二、核心设计模式 (Go Idiomatic Patterns)

Golang 的设计哲学是“组合优于继承”和“并发一等公民”，因此不要生搬硬套 Java/GoF 的模式。

### 1. 单例模式 (Singleton)
**场景**：全局唯一的对象（配置、数据库连接、HTTP Server）。
**实现**：使用 `sync.Once` 保证线程安全且只执行一次。
**项目案例**：`Carina` 的 HTTP Server 初始化。
```go
// carina/api/server/http.go
var (
    httpServer *HttpServer
    once       sync.Once
)
func GetInstance() *HttpServer {
    once.Do(func() {
        httpServer = &HttpServer{...}
    })
    return httpServer
}
```

### 2. 中间件/装饰器模式 (Middleware)
**场景**：Web 框架中横切通用逻辑（鉴权、日志、审计）。
**项目案例**：`Carina` 的 Router 中间件。
- `jwt.go`: 负责 Token 校验。
- `audit.go`: 负责操作审计。
- **亮点**：基于洋葱模型，业务 Controller 纯净，无需感知通用逻辑。

### 3. 注册表模式 + 工厂模式 (Registry & Factory)
**场景**：管理多种类型的插件或策略，实现开闭原则（OCP）。
**项目案例**：`Motion` (yavcd) 的运动项目管理。
- **注册**：利用 `init()` 和 `map[string]reflect.Type` 实现自注册。
- **工厂**：`CreateSport(name)` 根据字符串名称，利用反射动态创建具体的运动实例（如跑步、跳远）。
- **优势**：新增运动项目只需加新文件，无需修改核心工厂代码。

### 4. 状态机模式 (FSM / State)
**场景**：管理复杂对象的生命周期。
**项目案例**：`Motion` (yavcd) 的运动控制。
- 使用 FSM 管理运动会话的流转：`Init` -> `Start` -> `Stop` -> `CleanUp`。
- 相比大量的 `if-else`，状态机让状态流转更安全、清晰。

### 5. 代理/适配器模式 (Proxy & Adapter)
**场景**：屏蔽底层复杂性或不兼容接口。
**项目案例**：`Dipper` 项目。
- **Adapter**: C++ 层适配海康 SDK (`HCNetSDK`)。
- **Proxy**: Go 客户端通过 gRPC 代理调用底层硬件能力，上层业务无需关心底层是 C++ 还是硬件驱动。

---

## 三、高并发与稳定性模式 (Advanced Concurrency)

面试加分项，体现对高并发和资源治理的理解。

### 1. Singleflight (防缓存击穿)
**场景**：热点 Key 失效时，瞬间大量请求打库。
**实现**：合并并发请求，只让一个请求去执行，其余等待并共享结果。
**话术**：“在设计缓存层时，我常配合 `x/sync/singleflight` 使用，有效解决了缓存击穿问题，保护下游数据库。”

### 2. 信号量模式 (Semaphore / Bounded Concurrency)
**场景**：限制 Goroutine 数量，防止资源耗尽。
**实现**：使用带缓冲的 Channel (`make(chan struct{}, limit)`)。
**话术**：“虽然 Goroutine 很轻量，但我会使用 Buffered Channel 做信号量限制并发数，实现简单的背压机制。”
