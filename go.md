- Go优点
  - less is more. 减少无用特性带来的开发上的复杂性。
  - 性能与开发效率兼顾。
  - 与java性能类似的GC，但不像JAVA一样完全/繁琐的面向对象设计。与c/c++类似的性能，比c++更好的依赖管理以及有GC。和python一样的优雅，但是强类型语言，IDE自动补全以及编译期间检查错误。
  
- Go中有哪些锁

  ​	sync.Mutex 互斥锁
  ​	sync.RWMutex  读写锁
  
- golang的sync.atomic和C++11的atomic最显著的在golang doc⾥提到的差别在哪⾥，如何解决或 者说规避 

- 注意：Go sync/atomic包Load和Store并发不安全

- data race如何检测和解决？能不能不加锁解决这个问题？

- CSP并发模型？
  	CSP并发模型它并不关注发送消息的实体，而关注的是发送消息时使用的channel，go语言借用了process和channel这两个概念，process表现为go里面的goroutine，是实际并发执行的实体，每个实体之间是通过channel来进行匿名传递消息使之解藕，从而达到通讯来实现数据共享。

  ​	不要通过共享内存来通信，而要通过通信来实现内存共享。

  ​	1、sync.mutex 互斥锁（获取锁和解锁可以不在同一个协程，当获取到锁之后，未解锁，此时再次获取锁将会阻塞）
  ​	2、通过channel通信
  ​	3、sync.WaitGroup

- GPM模型指的是什么？goroutine的调度时机有哪些？如果syscall阻塞会发生什么？

  在go中是通过channel通信来共享内存的。

  G：指的是Goroutine，也就是协程，go中的协程做了优化处理，内存占用仅几kb
  且调度灵活，切换成本低。
  P：指的是processor,也就是处理器，感觉也可理解为协程调度器。
  M：指的是thread，内核线程。

  调度器的设计策略：
  	1、线程复用：当本线程无可运行的G时，M-P-G0会处于自旋状态，尝试从全局队列获取G，再从其他线程绑定的P队列中偷取G，而不是销毁线程；当本线程因为G进行系统调用阻塞时，线程会释放绑定的P队列，如果有空闲的线程可用就复用空闲的线程，不然就创建一个新的线程来接管释放出来的P队列。
  	2、利用并行：GOMAXPROCS设置P的数量，最多有这么多个线程分布在多个cpu上同时运行。
  	3、抢占：在coroutine中要等待一个协程主动让出CPU才执行下一个协程，在Go中，一个goroutine最多占用CPU 10ms，防止其他goroutine被饿死。

  go func的流程：
  	1、创建一个G，新建的G优先保存在P的本地队列中，如果满了则会保存到全局队列中。
  	2、G只能运行在M中，一个M必须持有一个P，M与P时1:1关系，M会从P的本地队列弹出一个可执行状态的G来执行。
  	3、一个M调度G执行的过程是一个循环机制。
  	4、如果G阻塞，则M也会被阻塞，runtime会把这个线程M从P摘除，再创建或者复用其他线程来接管P队列。
  	5、当G、M不在被阻塞，即系统调用结束，会先尝试找会之前的P队列，如果之前的P队列已经被其他线程接管，那么这个G会尝试获取一个空闲的P队列执行，并放入到这个P的本地队列。否则这个线程M会变成休眠状态，加入空闲线程队列，而G则会被放入全局队列中。

  M0：
  	M0是启动程序后的编号为0的主线程，这个M对应的实例会在全局变量runtime.m0中，不需要在heap上分配，M0负责执行初始化操作和启动第一个G，之后M0与其他的M一样。
  G0：
  	G0是每次启动一个M都会第一个创建的goroutine，G0仅负责调度，不指向任何可执行函数，每个M都会有一个自己的G0，在调度或者系统调用时会使用G0的栈空间，全局变量的G0是M0的。

  N:1-----出现阻塞的瓶颈，无法利用多个cpu
  1:1-----跟多线程/多进程模型无异，切换协程代价昂贵
  M:N-----能够利用多核，过于依赖协程调度器的优化和算法

  同步协作式调度
  异步抢占式调度

- 那些类型不能作map的为key？map的key为什么是无序的？
		map的key必须可以比较，func、map、slice这三种类型不可比较，只有在都是nil的情况下，才可与nil (== or !=)。因此这三种类型不能作为map的key。
	
- 数组或者结构体能够作为key？？？？有些能，有些不能，要看字段或者元素是否可比较
		1、map在扩容后，会发生key的搬迁，原来落在同一个bucket中的key可能分散，key的位置发生了变化。
		2、go中遍历map时，并不是固定从0号bucket开始遍历，每次都是从一个随机值序号的bucket开始遍历，并且是从这个bucket的一个随机序号的cell开始遍历。
		3、哈希查找表用一个哈希函数将key分配到不同的bucket(数组的下标index)。不同的哈希函数实现也会导致map无序。
	
	​	"迭代map的结果是无序的"这个特性是从go1.0开始加入的。
	
- 如何解决哈希查找表存在的"碰撞"问题（hash冲突）？

  ​	hash碰撞指的是：两个不同的原始值被哈希之后的结果相同，也就是不同的key被哈希分配到了同一个bucket。

  ​	链表法：将一个bucket实现成一个链表，落在同一个bucket中的key都会插入这个链表。

  ​	开放地址法：碰撞发生后，从冲突的下标处开始往后探测，到达数组末尾时，从数组开始处探测，直到找到一个空位置存储这个key，当找不到位置的情况下会触发扩容。

- map是线程安全的么？
		map不是线程安全的，sync.map是线程安全的。

	​	在查找、赋值、遍历、删除的过程中都会检测写标志，一旦发现写标志"置位"等于1，则直接panic,因为这表示有其他协程同时在进行写操作。赋值和删除函数在检测完写标志是"复位"之后，先将写标志位"置位"，才会进行之后的操作。

- 为什么sync.map为啥是线程安全？

    ```
    type Map struct {
    	// 互斥锁mu，操作dirty需先获取mu
    	mu Mutex 
    
    	// read是只读的数据结构，访问它无须加锁，sync.map的所有操作都优先读read
    	// read中存储结构体readOnly，readOnly中存着真实数据---entry（详见1.3），read是dirty的子集
    	// read中可能会存在脏数据：即entry被标记为已删除
    	read atomic.Value // readOnly
    
     	// dirty是可以同时读写的数据结构，访问它要加锁，新添加的key都会先放到dirty中
     	// dirty == nil的情况：1.被初始化 2.提升为read后，但它不能一直为nil，否则read和dirty会数据不一致。
    	// 当有新key来时，会用read中的数据 (不是read中的全部数据，而是未被标记为已删除的数据，详见3.2)填充dirty
    	// dirty != nil时它存着sync.map的全部数据（包括read中未被标记为已删除的数据和新来的数据）
    	dirty map[interface{}]*entry 
    
     	// 统计访问read没有未命中然后穿透访问dirty的次数
     	// 若miss等于dirty的长度，dirty会提升成read，提升后可以增加read的命中率，减少加锁访问dirty的次数   
     	misses int
    }
    ```

    

- map的底层实现原理是什么？
    ```
    type hmap struct {
        count      int   // len(map)元素个数
        flags      uint8 //写标志位
        B          uint8 // buckets数组的长度的对数，buckets数组的长度是2^B
        noverflow  uint16
        hash0      uint32
        buckets    unsafe.Pointer // 指向buckets数组
        oldbuckets unsafe.Pointer // 扩容的时候，buckets长度会是oldbuckets的两倍
        nevacuate  uintptr
        extra      *mapextra
    }
    
    // 编译期间动态创建的bmap
    type bmap struct {
        topbits  [8]uint8
        keys     [8]keytype
        values   [8]valuetype
        pad      uintptr
        overflow uintptr
    }
	```
    ​	在go中map是数组存储的，采用的是哈希查找表，通过哈希函数将key分配到不同的bucket，每个数组下标处存储的是一个bucket，每个bucket中可以存储8个kv键值对，当每个bucket存储的kv对到达8个之后，会通过overflow指针指向一个新的bucket，从而形成一个链表。

- map的扩容过程是怎样的？
	相同容量扩容
	2倍容量扩容
	扩容时机:
	1、当装载因子超过6.5时，表明很多桶都快满了，查找和插入效率都变低了，触发扩容。
		扩容策略：元素太多，bucket数量少，则将B加1，buctet最大数量(2^B)直接变为
		原来bucket数量的2倍，再渐进式的把key/value迁移到新的内存地址。
	2、无法触发条件1，overflow bucket数量太多，查找、插入效率低，触发扩容。
		(可以理解为：一座空城，房子很多，但是住户很少，都分散了，找起人来很困难)
		扩容策略：开辟一个新的bucket空间，将老bucket中的元素移动到新bucket，使得
		同一个bucket中的key排列更紧密，节省空间，提高bucket利用率。
	
- map的key的定位过程是怎样的？
		对key计算hash值，计算它落到那个桶时，只会用到最后B个bit位，再用哈希值的高8位找到key在bucket中的位置。桶内没有key会找第一个空位放入，冲突则从前往后找到第一个空位。
	
- iface和eface的区别是什么？值接收者和指针接收者的区别？
		iface和eface都是Go中描述接口的底层结构体，区别在于iface包含方法。而eface则是不包含任何方法的空接口：interface{}

	​	注意：编译器会为所有接收者为T的方法生成接收者为*T的包装方法，但是链接器会把程序中确定不会用到的方法都裁剪掉。因此*T和T不能定义同名方法。
	​	生成包装方法是为了接口，因为接口不能直接使用接收者为值类型的方法。
	​	如果方法的接收者是值类型，无论调用者是对象还是对象指针，修改的都是对象的副本，不影响调用者；如果方法的接收者是指针类型，则调用者修改的是指针指向的对象本身。
	​	如果类型具备"原始的本质"，如go中内置的原始类型，就定义值接收者就好。
	​	如果类型具备"非原始的本质"，不能被安全的复制，这种类型总是应该被共享，则可定义为指针接收者。

- context是什么？如何被取消？有什么作用？
    ```
    type Context interface {
        // 当context被取消或者到了deadline，返回一个被关闭的channel
        Done() <-chan struct{}
        // 在channel Done关闭后，返回context取消原因
        Err() error
        // 返回context是否会被取消以及自动取消时间(即deadline)
        Deadline() (deadline time.Time,ok boll)
        // 获取key对应的value
        Value(key interface{}) interface{}
    }
	
    type canceler interface {
        cancel(removeFromParent bool, err error)
        Done() <-chan struct{} 
    }
    ```
	context：goroutine的上下文，包含goroutine的运行状态、环境、现场等信息。实现了canceler接口的Context，就表明是可取消的。
	
	context用来解决goroutine之间退出通知、元数据传递的功能。比如并发控制和超时控制。
	
	注意事项：
		1、不要将Context塞到结构体里，直接将Context类型作为函数的第一参数，而且一般都命名为ctx。
		2、不要向函数传入一个nil的Context，如果你实在不知道传什么，标准库给你准备好了一个Context：todo
		3、不要把本应该作为函数参数的类型塞到Context中，Context存储的应该是一些共同的数据。例如：登陆的session、cookie等
		4、同一个Context可能会被传递到多个goroutine，Context是并发安全的。
	
- slice的底层数据结构是怎样的？
	```
	type slice struct {
    	array unsafe.Pointer // 底层数组的起始位置
    	len int
    	cap int
	}
	```

	​	slice的元素要存在一段连续的内存中，底层数据是数组，slice是对数组的封装，它描述一个数组的片段。
	​	slice可以向后扩展，不可以向前扩展。
	​	s[i]不可以超越len(s),向后扩展不可以超越底层数组cap(s)。
	​	make会为slice分配底层数组，而new不会为slice分配底层数组，所以array其实位置会是nil，可以通过append来分配底层数组。
	
	​	slice扩容方案计算：
	​		1、预估扩容后的容量：即假设扩容后的 cap 等于扩容后元素的个数
	​		if
	​			oldCap * 2 < cap，则newCap = cap
	​		else
	​			oldLen < 1024,则newCap = oldCap * 2 
	​			oldLen >= 1024,则newCap = oldCap * 1.25
	
	​		2、预估内存大小（int一个元素占8子节，string一个元素占16子节）
	​			假设元素类型是int，预估容量newCap = 5，那么预估内存 = 5 * 8 = 40 byte 
	​		3、匹配到合适的内存规格（内存分配规格为8、16、32、48、64、80....）
	​			实际申请的内存为 48 byte， cap = 48 / 8 = 6
	
- 你了解GC么？常见的GC实现方式有哪些？
	标记-清扫（三色标记）、标记-压缩、半空间复制、引用计数、分代GC、
	
- go的GC有那三个阶段？流程是什么？如果内存分配速度超过了标记清除速度怎么办？

	goV1.8三色+混合写屏障机制，栈不启动屏障，流程如下：
		1、GC开始将栈上的对象全部扫描并标记为黑色(之后不再进行重复扫描，无需STW)；
		2、GC期间，任何在栈上创建的新对象均标记为黑色；
		3、被删除的对象和被添加的对象均标记为灰色；
		4、回收白色集合中的所有对象。

	总结：
		v1.3普通标记清除法，整体过程需要STW，效率极低；
		v1.5三色标记法+屏障，堆空间启动写屏障，栈空间不启动，全部扫描之后，需要重新扫描一次栈(需要STW)，效率普通；
		v1.8三色标记法+混合写屏障，**消除栈的重扫过程，因为一旦栈被扫描变为黑色，则它会继续保持黑色， 并要求将对象分配为黑色**。堆空间启动，栈空间不启动屏障，整体过程几乎不需要STW,效率较高。
	
	​	如果申请内存的速度超过预期，运行时就会让申请内存的应用程序辅助完成垃圾收集的扫描阶段，在标记和标记终止阶段结束之后就会进入异步的清理阶段，将不用的内存增量回收。并发标记会设置一个标志，并在mallocgc调用时进行检查，当存在新的内存分配时，会暂停分配内存过快的哪些goroutine，并将其转去执行一些辅助标记的工作，从而达到放缓内存分配和加速GC工作的目的。

- 内存泄漏如何解决？
	1、通过pprof工具获取内存相差较大的两个时间点heap数据。htop可以查看内存增长情况。
	2、通过go tool pprof比较内存情况，分析多出来的内存。
	3、分析代码、修复代码。

- 内存逃逸分析？
	在函数中申请一个新的对象，如果分配在栈中，则函数执行结束可自动将内存回收；
	如果分配在堆中，则函数执行结束可交给GC处理。

	案例：
	函数返回局部变量指针；
	申请内存过大超过栈的存储能力。

- 你是如何实现单元测试的？有哪些框架？
	testing、GoMock、testify
