- i++为什么不是线程安全的
  - 底层经历了读取数据、更新CPU缓存、存入内存等操作
  - 编译器编译和CPU处理时通过调整指令顺序进行优化
  - 锁
    - 原子锁/自旋锁:
      - `atomic.CompareAndSwapInt64(&flag), 0, 0)`
      - 缺点：1.CPU开销大；2.不能保证代码的原子性；3.ABA问题（版本号+时间戳解决）
    - 互斥锁：
      - sync.Mutex 更轻量级的锁，解决原子锁CPU开销大的问题。
      - 实现：
        ```go
        type Mutex struct {
           state int32 // 锁状态
           sema uint32 // 信号量
        }
        ```
        state锁状态为互斥锁位图，包含：
        - 1.正在等待被唤醒的协成数量；
        - 2.当前进入饥饿状态；
        - 3.协成准备从正常状态下被唤醒；
        - 4.锁定状态
        
        阶段：
        1. CAS快速抢占锁。成功就返回，失败调用lockSow。
           - lockSlow正常情况下会自旋尝试抢占锁一段时间，不立即进入休眠状态。
           - 存在4种情况，自旋终止：
             - 单核CPU
             - 逻辑处理器小于等于1
             - 当前协成所在逻辑处理器的本地队列上有其它协成待运行。
             - 自旋次数超过了设定的阈值
        2. 
    - 读写锁
- 数组和切片有什么区别？
    - go语言中数组是一种值类型，[2]int和[3]int是两种不同的类型；切片类型只和它的基础数据类型有关，如[]int和[]string.
    - 数组本身的赋值和传参都是以整体复制进行处理的；而切片复制的只有切片头部分信息，因为包含底层数据指针。
    - 切片基于数组：
      ```
      type SliceHeader struct {
        Data uintptr
        len int
        Cap int
      ```

- 字符串和数组有什么区别？
  - 字符串底层数据是对应的字节数组，但字符串的只读属性禁止了在程序中对字节数组的元素进行修改。
  - 字符串赋值只复制了数据地址和对应的长度，不会复制底层数据。（使用不当会造成内存泄漏）

- 内存泄漏场景？
  - 切片/字符串引用不当：先对需要引用的进行拷贝，再引用
  - 频繁的系统调用
  - for循环中使用defer：在for中构建一个局部函数，在函数内部执行defer
  - goroutine泄露
    如：Ticker使用忘记Stop，通常使用context来避免。

- MPG模型
  - M(machine): 工作线程，由操作系统调度。应该就是通常所说的内核线程。
  - P(processor): 处理器（非CPU），代表着运行Go代码的必要资源，以及调度goroutine的能力。个人觉得可以当作拥有自主调度权的算法模块，用于工作窃取（work stealing）。
  - G(gooutine): Go协程，轻量级用户线程。主要包含执行栈和调度管理器。这里的调度管理器指的是，统一并管理调度资源，等待被调度。
  - 关于M的数目：
    M的个数是根据实际情况自行创建的，一般稍大于P的个数，为了保证runtime包的内置任务的运行。在运行中不够用时，也会再重新创建一个。
  - 关于P的数目：
    P的个数默认为CPU的核数，在IO密集的场景下可以适当提高P的个数。设置方式有两种，例：
    设置环境变量:export GOMAXPROCS=80 或者 runtime.GOMAXPROCS(80)
  - 另外还有全局runqueue队列存在：全局队列由多个处理器共享，访问通过互斥锁来完成。
    处理器P中的协程G额外再创建的协程会加入到本地的runqueues中。
    两种情况下会放入全局队列中：1. 本地队列已满 2. 阻塞的协程被唤醒
    全局队列会被处理器P周期性的摘取来调度。
  - 调度策略
    - 队列轮转
    - 系统调用
    - 工作量窃取
    - 抢占式调度
- Context

  context包是Go 语言中用来设置截止日期、同步信号，传递请求相关值的结构体，是开发常用的并发控制技术。

  与WaitGroup的不同在于context可以控制多级的goroutine。

  Context是线程安全的。 

### Gin的启动过程

#### 项目的main函数

主函数位于项目根目录下的main.go中，代码如下：

```
package main

import (
	"github.com/LearnGin/handler"
	"github.com/LearnGin/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// init gin with default configs
	r := gin.Default()

	// append custom middle-wares
	middleware.RegisterMiddleware(r)
	// register custom routers
	handler.RegisterHandler(r)

	// run the engine
	r.Run()
}
```

主要步骤：

1. 初始化Gin

   ```
   gin.Default()
   ```

   执行Gin的初始化过程，默认的初始化包含两个中间件，

   1. **Logger**：日志中间件，将Gin的启动与响应日志输出到控制台；
   2. **Recovery**：恢复中间件，将Gin遇到的无法处理的请求按HTTP 500状态码返回。

2. **注册中间件**：本例的`middleware.RegisterMiddleware(r)`用于将项目中开发的中间件注册到Gin Engine上；

3. **注册事件处理**：本例的`handler.RegisterHandler(r)`用于将项目中开发的对应于指定URL的事件处理函数注册到Gin Engine上；

4. **启动Gin**：`r.Run()`负责启动Gin Engine，开始监听请求并提供HTTP服务。

