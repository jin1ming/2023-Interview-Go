  - P(processor): 处理器（调度器、非CPU），代表着运行Go代码的必要资源，以及调度goroutine的能力。

  - G(gooutine): Go协程，轻量级用户线程。

  - M(machine): 代表一个工作线程。

  - 模型：

    ```
    -M-M-M-M- M列表
    G0-M-G
       |
       P-G-G-G-G-G
       |
    gFree(G) 执行完的
```
    
**在M上有一个P和G，P是绑定到M上的，G是通过P的调度获取的，在某一时刻，一个M上只有一个G（g0除外）。在P上拥有一个G队列，里面是已经就绪的G，是可以被调度到线程栈上执行的协程，称为运行队列。**
    
g0: 特殊的 Goroutine 涉及许多其他操作（较大空间的对象分配，cgo 等），需要较大的栈来保证我们的程序进行更高效的管理操作，以保持程序的低内存打印效率。
    
    ```
    全局G队列：-G-G-G-G-
    sudoG: -G-G-G-
    游离P：
    	P-syscall：P-G-G-G-G
    	P-idle：P-P-P
```
    
    **每个进程都有一个全局的G队列，也拥有P的本地执行队列，同时也有不在运行队列中的G。如正处于channel的阻塞状态的G，还有脱离P绑定在M的(系统调用)G，还有执行结束后进入P的gFree列表中的G等等。**
    
    全局runqueue队列：由多个处理器共享，访问通过互斥锁来完成。
    处理器P中的协程G额外再创建的协程会加入到本地的runqueues中。
    两种情况下会放入全局队列中：1. 本地队列已满 2. 阻塞的协程被唤醒
    全局队列会被处理器P周期性的摘取来调度。
    
  - 关于M的数目：
    M的个数是根据实际情况自行创建的，一般稍大于P的个数，为了保证runtime包的内置任务的运行。在运行中不够用时，也会再重新创建一个。

  - 关于P的数目：
    P的个数默认为CPU的核数，在IO密集的场景下可以适当提高P的个数。设置方式有两种，例：
    设置环境变量:export GOMAXPROCS=80 或者 runtime.GOMAXPROCS(80)

### P 调度器的设计策略

- 线程复用、并行
- work stealing 工作窃取机制：当 P 本地队列无运行 G 时，会去其他线程绑定的 P 窃取 G ，若其他 P 本地队列也没有时会去 G 全局队列进行窃取
- hand off 分离机制：当 G 因为系统调用阻塞时，P 会和 M 解绑，将 G 和 M 绑定，P 会和空闲的线程进行绑定
- 主动让出机制：当 G 占用了 CPU 超过 10MS 会主动让出（sysmon 轮询） 

### GMP 模型中为什么需要 P

- 假设没有 P 会发生什么事情，（没有 P 本地队列，没有 P 全局队列）

  a. 出现资源竞争，由于没有 P 的本地以及全局队列的多级缓存，所以 G 都会放在一起，多个 M 去获取时会出现资源竞争

  b. 协程切换资源消耗大，由于没有 P ，也没有 g0 去负责切换协程堆栈，相当于协程堆栈以及一些运行现场全部都在内核态去维护 

  c. 系统调用阻塞时切换成本高，M 会经常被阻塞和解阻塞切换内核态和用户态（类似于进程间切换），消耗大


- 引入了 P 相当于解决了上面的大部分问题，甚至引入了新的特性去榨干 CPU 的性能

  a. 引入本地队列和全局队列做多级缓存，获取 G 时都会从本地队列获取，没有竞争，就算去全局队列获取也比较少的几率出现大量 P 去获取，降低了资源竞争的概率

  b. 切换协程堆栈效率提高，使用了 g0 的协程负责去管理协程切换的堆栈以及保护现场的工作，进入内核态去切换的时候少了很多切换指令以及寄存器的使用，甚至引入了 g0 去负责做垃圾回收部分工作的职能 

  c. 系统调用的曲线救国方案，当 G 和 M 发生了系统调用时，P 会解绑 M ，带着本地队列的 P 去找空闲的 M 或者新创建的 M 去继续剩下的工作 还引入了「工作窃取」的功能，让基于和 M 绑定的 P 更加灵活的让每个 M 都能够最大限度的运行 task，榨干 CPU 

### **G状态**

**_Gidle**：刚刚被分配并且还没有被初始化，值为0，为创建goroutine后的默认值
**_Grunnable**： 没有执行代码，没有栈的所有权，存储在运行队列中，可能在某个P的本地队列或全局队列中(如上图)。
**_Grunning**： 正在执行代码的goroutine。
**_Gsyscall**：正在执行系统调用，拥有栈的所有权，与P脱离，但是与某个M绑定，会在调用结束后被分配到运行队列(如上图)。
**_Gwaiting**：被阻塞的goroutine，阻塞在某个channel的发送或者接收队列(如上图)。
**_Gdead**： 当前goroutine未被使用，没有执行代码，可能有分配的栈，分布在空闲列表gFree，可能是一个刚刚初始化的goroutine，也可能是执行了goexit退出的goroutine(如上图)。
**_Gcopystac**：栈正在被拷贝，没有执行代码，不在运行队列上，执行权在
**_Gscan** ： GC 正在扫描栈空间，没有执行代码，可以与其他状态同时存在

### **P的状态**

**_Pidle** ：处理器没有运行用户代码或者调度器，被空闲队列或者改变其状态的结构持有，运行队列为空
**_Prunning** ：被线程 M 持有，并且正在执行用户代码或者调度器(如上图)
**_Psyscall**：没有执行用户代码，当前线程陷入系统调用(如上图)
**_Pgcstop** ：被线程 M 持有，当前处理器由于垃圾回收被停止
**_Pdead** ：当前处理器已经不被使用
### **M的状态**
**自旋线程**：处于运行状态但是没有可执行goroutine的线程，数量最多为GOMAXPROC，若是数量大于GOMAXPROC就会进入休眠。
**非自旋线程**：处于运行状态有可执行goroutine的线程。