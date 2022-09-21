context主要用于父子任务之间的同步取消信号，本质上是一种协程调度的方式。另外在使用context时有两点值得注意：上游任务仅仅使用context通知下游任务不再需要，但不会直接干涉和中断下游任务的执行，由下游任务自行决定后续的处理操作，也就是说context的取消操作是无侵入的；context是线程安全的，因为context本身是不可变的（immutable），因此可以放心地在多个协程中传递使用。

  context包是Go 语言中用来设置截止日期、同步信号，传递请求相关值的结构体，是开发常用的并发控制技术。

  与WaitGroup的不同在于context可以控制多级的goroutine。

  Context是线程安全的。 

作用：

- 传递上下文。不推荐传递业务参数，常用来传递整个链路的trace id来作为链路追踪
- 协程间同步信息。比如在协程树中传递取消信号。

**Context主要用于多个协程间的信号同步**。定义如下：

```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}
```

共提供了6种实例：

```go
// 由context.emptyCtx初始化而来，没有任何功能
func Background() Context {...}
func TODO() Context {...}

// 功能context
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {...}
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {...}
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {...}
func WithValue(parent Context, key, val interface{}) Context {...}
```

案例：

1.创建过期时间为1s的上下文。 2. 将`context`传入`handle`函数中，函数使用500ms的时间处理请求。

```go
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go handle(ctx, 500*time.Millisecond)
	select {
	case <-ctx.Done():
		fmt.Println("main", ctx.Err())
	}
}

func handle(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("handle", ctx.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}
}
```

分析: `context`过期时间为1s，处理时间为0.5秒(`select`中的过期时间)，函数有足够的时间完成处理，也就是`<-time.After(duration):`会在`<-ctx.Done()`之前完成，故输出`process request with 500ms`。再过0.5s，`<-ctx.Done()`完成，这时候输出`main context deadline exceeded`。

**当父Context取消时，相关的子Context也会相应的结束**：

```go
func main() {
	wg := sync.WaitGroup{}
	parent, cancel := context.WithCancel(context.Background())
	child1, cancel1 := context.WithCancel(parent)
	defer cancel1()
	child2, cancel2 := context.WithCancel(parent)
	defer cancel2()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				fmt.Println("child1 ticker")
			case <-child1.Done():
				fmt.Println("child1 done")
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				fmt.Println("child2 ticker")
			case <-child2.Done():
				fmt.Println("child2 done")
				return
			}
		}
	}()

	time.Sleep(3 * time.Second)

	cancel()
	wg.Wait()
}
```

父Context与子Context间信号传递的原理如下：

```go
// 调用WithXXX时，会使用这个函数构建父、子Context之间的联系
func propagateCancel(parent Context, child canceler) {
	done := parent.Done()
	if done == nil {
		return // 说明parent context是一个非cancel类型的context，例如WithValue、TODO、Background
	}
	select {
	case <-done:
		child.cancel(false, parent.Err()) // 检测父上下文是否已经被取消
		return
	default:
	}
	
  // 将child context注册至parent context中。并启用select监听parent.Done和child.Done
	if p, ok := parentCancelCtx(parent); ok {	
		p.mu.Lock()
		if p.err != nil {
			child.cancel(false, p.err)
		} else {
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
	} else {
		go func() {
			select {
			case <-parent.Done():
				child.cancel(false, parent.Err())
			case <-child.Done():
			}
		}()
	}
}

// 取消核心实现。包括WithCancel、WithTimeout、WithDeadline
func (c *cancelCtx) cancel(removeFromParent bool, err error) {
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return
	}
	c.err = err
	if c.done == nil {
		c.done = closedchan
	} else {
		close(c.done)	// close done channel
	}
  
  // 遍历child context列表，并逐个执行child的cancel函数
	for child := range c.children {
		child.cancel(false, err)
	}
	c.children = nil
	c.mu.Unlock()
	
	if removeFromParent {
		removeChild(c.Context, c)
	}
}
```