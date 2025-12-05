# Motion Media 项目 - 完整架构设计

## 一、项目定位

**Motion Media** 是Dipper边缘端的**媒体流服务**，核心职责：

- 🎥 **多源支持**：海康、大华、HTTP/HTTPS视频源
- 📡 **实时流传输**：通过gRPC将视频流推送给客户端
- 💾 **视频录制**：支持指定时间段的精确录制
- 📸 **截图功能**：从视频源获取单帧截图
- ⚡ **高性能**：平均延迟<1ms，P99延迟<3ms

**技术栈**：Go 1.24 + gRPC + FFmpeg + Prometheus

---

## 二、系统全景图

```
┌─────────────────────────────────────────────────────────────┐
│                    gRPC Server                              │
│          (Unix Socket + TCP 38080)                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ OpenStream   │  │ StartRecording│  │CaptureSnapshot│   │
│  │ (流式推送)   │  │ (开始录制)    │  │ (截图)       │   │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘     │
│         │                 │                 │              │
│         └─────────────────┼─────────────────┘              │
│                           │                                │
│         ┌─────────────────▼─────────────────┐              │
│         │     Service Layer                 │              │
│         │  ├─ stream.go (流管理)            │              │
│         │  ├─ snapshot.go (截图)           │              │
│         │  └─ recording.go (录制)           │              │
│         └─────────────────┬─────────────────┘              │
│                           │                                │
│         ┌─────────────────▼─────────────────┐              │
│         │  Source Abstraction Layer         │              │
│         │  ├─ Hikvision (SDK)               │              │
│         │  ├─ Dahua (SDK)                   │              │
│         │  └─ Http (FFmpeg)                 │              │
│         └──────────────┬──────────────┬─────┘              │
│                        │              │                    │
│         ┌──────────────▼──┐  ┌───────▼──────┐             │
│         │ Stream Processing│  │ Circular Queue│            │
│         │ ├─ Frame Channel │  │ (I帧缓冲)    │            │
│         │ ├─ Video Writer  │  │ (1秒缓冲)    │            │
│         │ └─ Metrics       │  └──────────────┘            │
│         └──────────────────┘                              │
│                                                             │
│         ┌─────────────────────────────────────┐            │
│         │  Infrastructure                     │            │
│         │  ├─ Prometheus Metrics              │            │
│         │  ├─ Zerolog Logging                 │            │
│         │  ├─ Viper Configuration             │            │
│         │  └─ PProf Profiling                 │            │
│         └─────────────────────────────────────┘            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
        │                    │                    │
        ▼                    ▼                    ▼
    ┌────────┐          ┌────────┐          ┌────────┐
    │摄像头1 │          │摄像头2 │          │摄像头3 │
    │(Hik)   │          │(Dahua) │          │(HTTP)  │
    └────────┘          └────────┘          └────────┘
```

---

## 三、分层架构

### 3.1 表现层 - gRPC API

**核心服务**：
- `OpenStream`：打开视频流，流式推送帧
- `StartRecording`：开始录制
- `StopRecording`：停止录制
- `CaptureSnapshot`：获取截图

**服务器启动** (`server/media.go`)：
- 同时监听Unix Socket和TCP 38080
- 配置gRPC Keepalive参数
- 注册Health Check服务
- 支持MaxRecvMsgSize=10MB

### 3.2 业务逻辑层 - Service

**流管理** (`service/stream.go`)：
- 解析URI，创建Source
- 登录认证，启动流
- 流式推送帧给客户端
- 清理资源

**录制管理** (`service/stream.go`)：
- 检查循环队列中是否有指定时间的I帧
- 有则立即开始，无则等待下一个I帧
- 从最近的I帧开始写入所有帧
- 支持指定时长自动停止

**截图管理** (`service/snapshot.go`)：
- 海康/大华：调用SDK截图
- HTTP：使用FFmpeg解码第一帧并编码为JPEG
- 使用mutex防止并发截图

### 3.3 数据访问层 - Source抽象

**Source接口**：
```go
type Source interface {
  GetURI() *model.URI
  Login() error
  Logout() error
  StartStream(stream string) (*Stream, error)
  StopStream(stream string) error
  GetStream(stream string) *Stream
}
```

**三种实现**：
1. **Hikvision**：调用海康SDK，支持主码流和第三码流
2. **Dahua**：调用大华SDK
3. **Http**：使用go-astiav（FFmpeg Go绑定）

### 3.4 核心处理层 - Stream & Queue

**Stream处理**：
- 验证帧有效性
- 计算延迟指标
- 添加到循环队列
- 写入录制文件
- 发送到客户端

**循环队列**：
- 大小 = 帧率（通常25），保存最近1秒的帧
- 自动覆盖最旧的帧
- 查找最早的I帧
- 判断是否可以开始录制

---

## 四、核心业务流程

### 4.1 流式推送流程

```
1. 客户端请求 OpenStream(uri)
2. 解析URI (hik://ip:port/user:pass@main,sub)
3. 创建Source (Hikvision/Dahua/Http)
4. Source.Login() - 认证
5. StartStream("main") - 启动主码流
6. SDK回调 onFrame() → Stream.onFrame()
   ├─ 验证帧有效性
   ├─ 添加到循环队列
   ├─ 如果在录制，写入文件
   └─ 发送到FrameCh (100帧缓冲)
7. 通过gRPC发送给客户端
8. 客户端断开连接
9. Stream.Stop() - 清理资源
10. Source.Logout() - 登出
```

### 4.2 录制流程

```
1. StartRecording(recordingId, startTime, duration)
2. 检查循环队列中是否有指定时间的I帧
   ├─ 有 → 立即开始录制
   └─ 无 → 标记为scheduled，等待下一个I帧
3. onFrame() 收到帧
4. 如果是I帧 && scheduled
   ├─ 创建VideoWriter
   ├─ 从队列中找到最近的I帧，写入所有帧
   └─ 后续帧继续写入文件
5. 等待duration秒
6. StopRecording() - 关闭文件
```

### 4.3 截图流程

```
1. CaptureSnapshot(uri)
2. 根据scheme判断
   ├─ hik → 海康SDK截图
   │   ├─ HikLogin()
   │   ├─ HikCaptureImage()
   │   └─ HikLogout()
   └─ http → FFmpeg截图
       ├─ 打开输入流
       ├─ 找到视频流
       ├─ 创建解码器
       ├─ 读取第一帧
       ├─ 编码为JPEG
       └─ 返回JPEG数据
```

---

## 五、关键技术点

### 5.1 多源抽象设计

**优势**：
- 统一接口：所有源都实现Source接口
- 易于扩展：添加新源只需实现接口
- 代码复用：Stream、Queue等通用逻辑

**工厂函数**：
```go
func New(uri *model.URI) (Source, error) {
  switch uri.Scheme {
  case "hik":
    return NewHikvision(uri), nil
  case "dahua":
    return NewDahua(uri), nil
  case "http", "https":
    return NewHttp(uri), nil
  default:
    return nil, fmt.Errorf("unsupported scheme: %s", uri.Scheme)
  }
}
```

### 5.2 循环队列设计

**目的**：缓存最近1秒的视频帧，支持精确时间录制

**特点**：
- 固定大小：capacity = 帧率（通常25）
- 自动覆盖：队列满时覆盖最旧帧
- I帧对齐：录制时从最近的I帧开始

**时间判断逻辑**：
```
请求时间 t
    ↓
1. t < 最早I帧时间 → 无法录制（数据已过期）
2. 最早I帧 ≤ t < 下一个I帧 → 立即开始
3. t ≥ 下一个I帧时间 → 等待（标记为scheduled）
```

### 5.3 AVCC vs Annex-B格式转换

**问题**：HTTP源获取的H.264数据是AVCC格式，需要转换为Annex-B

**转换逻辑**：
```go
func avccToAnnexB(avcc []byte) []byte {
  var out []byte
  i := 0
  for i+4 <= len(avcc) {
    // 读取4字节长度
    nalLen := int(avcc[i])<<24 | int(avcc[i+1])<<16 | 
              int(avcc[i+2])<<8 | int(avcc[i+3])
    i += 4
    
    if i+nalLen > len(avcc) {
      break
    }
    
    // 添加Annex-B分隔符
    out = append(out, 0x00, 0x00, 0x00, 0x01)
    // 添加NAL单元
    out = append(out, avcc[i:i+nalLen]...)
    i += nalLen
  }
  return out
}
```

### 5.4 并发安全设计

**关键点**：
- **atomic操作**：`streamClosed`, `frameIndex`, `recodingScheduled`
- **读写锁**：`writerMu`保护writer的并发访问
- **互斥锁**：`mu`保护formatCtx的并发访问
- **Channel**：FrameCh用于线程间通信

---

## 六、性能指标

| 指标 | 数值 |
|------|------|
| 平均延迟 | 0.7ms |
| P99延迟 | 3ms |
| 最大延迟 | 10-20ms |
| 单路CPU占用 | 5-6% |
| 单路内存占用 | 120MB |
| 支持并发路数 | 3+ |

**性能优化**：
- 直接转发：不做任何处理，直接转发SDK回调的帧
- 缓冲最小化：FrameCh只缓冲100帧（4秒左右）
- 并发处理：读取、处理、发送并行进行
- 指标监控：实时监控延迟，发现问题及时告警

---

## 七、遇到的困难和解决方案

### 困难1：HTTP源内存泄漏

**问题**：长时间运行HTTP源，内存持续增长

**原因**：
- FFmpeg packet没有及时释放
- goroutine没有正确退出
- formatCtx没有正确关闭

**解决**：
- 及时释放packet：`packet.Unref()`
- defer确保资源释放：`formatCtx.CloseInput()`, `formatCtx.Free()`
- 强制GC：`runtime.GC()`
- 等待goroutine退出：使用done channel

### 困难2：录制时无法从指定时间开始

**问题**：用户请求录制某个时间的视频，但无法精确定位

**原因**：
- 视频只能从I帧开始解码
- 循环队列中可能没有指定时间的I帧

**解决**：
- 检查队列中是否有指定时间的I帧
- 有则立即开始，无则标记为scheduled
- 等待下一个I帧时触发录制

### 困难3：AVCC vs Annex-B格式转换

**问题**：HTTP源获取的H.264数据是AVCC格式，但需要Annex-B格式

**原因**：
- AVCC：MP4容器格式，NAL单元前有4字节长度
- Annex-B：标准H.264格式，NAL单元前有0x00000001

**解决**：实现转换函数，将AVCC格式转换为Annex-B格式

### 困难4：并发访问formatCtx导致崩溃

**问题**：StopStream和读取帧并发访问formatCtx，导致崩溃

**原因**：FFmpeg不是线程安全的

**解决**：
- 使用mutex保护formatCtx
- StopStream时只cancel，不释放
- 等待goroutine自己释放资源
- 使用done channel同步

---

## 八、设计亮点

### 8.1 多源统一抽象

通过Source接口统一了三种不同的视频源（海康、大华、HTTP），使得：
- 上层代码无需关心具体源类型
- 添加新源只需实现接口
- 代码逻辑清晰，易于维护

### 8.2 低延迟设计

- 直接转发SDK回调的帧，不做任何处理
- 缓冲最小化（100帧≈4秒）
- 并发处理，充分利用多核
- 实时监控延迟指标

### 8.3 精确时间录制

通过循环队列和I帧对齐，实现了：
- 精确的时间定位
- 从最近的I帧开始录制
- 支持等待下一个I帧的调度机制

### 8.4 内存管理

- 及时释放FFmpeg资源
- 固定大小循环队列，不会无限增长
- 对象池复用buffer
- 定期强制GC

### 8.5 并发安全

- 合理使用atomic、mutex、channel
- 保护共享资源的访问
- 支持多个客户端同时连接

---

## 九、总结

Motion Media是一个**高性能、高可靠**的媒体流服务，特点：

- **多源支持**：统一接口支持多种视频源
- **低延迟**：平均延迟<1ms
- **高并发**：支持多路并发流处理
- **精确录制**：支持指定时间段的精确录制
- **内存高效**：固定大小缓冲，不会无限增长
- **并发安全**：合理使用并发原语

这个项目展示了如何构建一个高性能的实时媒体流服务，涉及多个技术领域的深度理解。

