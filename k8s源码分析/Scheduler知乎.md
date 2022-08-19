![img](https://pic3.zhimg.com/80/v2-623640db58796c066027e2c8f3112096_1440w.jpg)

## 源码分析

> Based on Kubernetes v1.19.11.

K8s scheduler 主要的数据结构是：

1. Scheduler。
2. SchedulingQueue。

相关的代码流程主要分为两个部分：

1. `cmd/kube-scheduler`，这里是我们调度器的起始处，主要是读取配置，初始化并启动调度器。
2. `pkg/scheduler`，这里是调度器的核心代码。

## 数据结构

### Scheduler

```go
// pkg/scheduler/scheduler.go

// Scheduler watches for new unscheduled pods. It attempts to find
// nodes that they fit on and writes bindings back to the api server.
type Scheduler struct {
    // It is expected that changes made via SchedulerCache will be observed
    // by NodeLister and Algorithm.
    SchedulerCache internalcache.Cache
    Algorithm core.ScheduleAlgorithm

    // NextPod should be a function that blocks until the next pod
    // is available. We don't use a channel for this, because scheduling
    // a pod may take some amount of time and we don't want pods to get
    // stale while they sit in a channel.
    NextPod func() *framework.QueuedPodInfo

    // Error is called if there is an error. It is passed the pod in
    // question, and the error
    Error func(*framework.QueuedPodInfo, error)

    // Close this to shut down the scheduler.
    StopEverything <-chan struct{}

    // SchedulingQueue holds pods to be scheduled
    SchedulingQueue internalqueue.SchedulingQueue

    // Profiles are the scheduling profiles.
    Profiles profile.Map

    scheduledPodsHasSynced func() bool
    client clientset.Interface
}
```

1. `SchedulerCache` ，保存了调度所需的 podStates 和 nodeInfos。
2. `Algorithm` ，会使用该对象的 `Schedule` 方法来运行调度逻辑。
3. `SchedulingQueue` ，调度队列。
4. `Profiles` ，调度器配置。

### SchedulingQueue

> Implementation

```go
type PriorityQueue struct {
    // PodNominator abstracts the operations to maintain nominated Pods.
    framework.PodNominator
    stop  chan struct{}
    clock util.Clock
    // pod initial backoff duration.
    podInitialBackoffDuration time.Duration
    // pod maximum backoff duration.
    podMaxBackoffDuration time.Duration

    lock sync.RWMutex
    cond sync.Cond
    // activeQ is heap structure that scheduler actively looks at to find pods to
    // schedule. Head of heap is the highest priority pod.
    activeQ *heap.Heap

    // podBackoffQ is a heap ordered by backoff expiry. Pods which have completed backoff
    // are popped from this heap before the scheduler looks at activeQ
    podBackoffQ *heap.Heap

    // unschedulableQ holds pods that have been tried and determined unschedulable.
    unschedulableQ *UnschedulablePodsMap

    // schedulingCycle represents sequence number of scheduling cycle and is incremented
    // when a pod is popped.
    schedulingCycle int64
    // moveRequestCycle caches the sequence number of scheduling cycle when we
    // received a move request. Unscheduable pods in and before this scheduling
    // cycle will be put back to activeQueue if we were trying to schedule them
    // when we received move request.
    moveRequestCycle int64

    // closed indicates that the queue is closed.
    // It is mainly used to let Pop() exit its control loop while waiting for an item.
    closed bool
}
```

1. PodNominator：调度算法调度的结果，保存了 Pod 和 Node 的关系。
2. cond：用来控制调度队列的 Pop 操作。
3. activeQ：用堆维护的优先队列，保存着待调度的 pod，其中优先级默认是根据 Pod 的优先级和创建时间来排序。
4. podBackoffQ：同样是用堆维护的优先队列，保存着运行失败的 Pod，优先级是根据 `backOffTime` 来排序，`backOffTime` 受 `podInitialBackoffDuration` 以及 `podMaxBackoffDuration` 两个参数影响。
5. unschedulableQ：是一个 Map 结构，保存着暂时无法调度（可能是资源不满足等情况）的 Pod。

## pkg/scheduler

### 运行调度器主流程

`Run` 会启动 scheduling queue，并不断调用 `sched.scheduleOne()` 进行调度。

```go
// Run begins watching and scheduling. It waits for cache to be synced, then starts scheduling and blocked until the context is done.
func (sched *Scheduler) Run(ctx context.Context) {
    if !cache.WaitForCacheSync(ctx.Done(), sched.scheduledPodsHasSynced) {
        return
    }
    sched.SchedulingQueue.Run()
    wait.UntilWithContext(ctx, sched.scheduleOne, 0)
    sched.SchedulingQueue.Close()
}
```

### 运行调度队列

```go
// Run starts the goroutine to pump from podBackoffQ to activeQ
func (p *PriorityQueue) Run() {
    go wait.Until(p.flushBackoffQCompleted, 1.0*time.Second, p.stop)
    go wait.Until(p.flushUnschedulableQLeftover, 30*time.Second, p.stop)
}
```

调度队列的运行逻辑：

1. 每隔 1s 检查 `podBackoffQ` 是否有 pod 可以放入 `activeQ` 中。检查的逻辑是判断 `backOffTime` 是否已经到期。
2. 每隔 30s 检查 `unschedulableQ` 是否有 pod 可以放入 `activeQ` 中。

### 单个 Pod 的调度 scheduleOne

在介绍 `scheduleOne` 之前，看这张 pod 调度流程图能有助于我们理清整个过程。同时这也是 k8s v1.15 开始支持的 Scheduling Framework 的 Plugin 扩展点。

![img](https://pic3.zhimg.com/80/v2-9e23ef8e9ba82e5b8e981a01d5cacad2_1440w.jpg)

`ScheduleOne` 是调度器的主流程，主要包括以下几步：

1. 调用 `sched.NextPod()` 拿到下一个需要调度的 pod。后面会对这个过程进行更详细的介绍。
2. 调用 `sched.profileForPod(pod)` ，根据 pod 中的 schedulerName 拿到针对该 pod 调度的 Profiles。这些 Profiles 就包括了调度插件的配置等。
3. 进行上图中的 Scheduling Cycle 部分，这部分是单线程运行的。
4. 调用 `sched.Algorithm.Schedule()`。此处包括好几个步骤，其中 `PreFilter`, `Filter` 被称为 **Predicate**，是对节点进行过滤，这里面考虑了节点资源，Pod Affinity，以及 Node Volumn 等情况。而 `PreScore` , `Score` , `Nomalize Score` 又被称为 **Priorities**，是对节点进行优选打分，这里会得到一个适合当前 Pod 分配上去的 Node。
5. 进行 `Reserve` 操作，将调度结果缓存。当后面的调度流程执行失败，会进行 `Unreserve` 进行数据回滚。
6. 进行 `Permit` 操作，这里是用户自定义的插件，可以使 Pod 进行 allow（允许 Pod 通过 Permit 阶段）、reject（Pod 调度失败）和 wait（可设置超时时间）这三种操作。对于 Gang Scheduling （一批 pod 同时创建成功或同时创建失败），可以在 `Permit` 对 Pod 进行控制。
7. 进行图中的 Binding Cycle 部分，这部分是起了一个 Goroutine 去完成工作的，不会阻塞调度主流程。
8. 最开始会进行 `WaitOnPermit` 操作，这里会阻塞判断 Pod 是否 Permit，直到 Pod Permit 状态为 allow 或者 reject 再往下继续运行。
9. 进行 `PreBind` , `Bind` , `PostBind` 操作。这里会调用 k8s apiserver 提供的接口，将待调度的 Pod 与选中的节点进行绑定，但是可能会绑定失败，此时会做 `Unreserve` 操作，将节点上面 Pod 的资源解除预留，然后重新放置到失败队列中。

这里还会涉及 scheduling framework 等知识，在这里不多做赘述。

## SchedulingQueue 细节

### 获取下一个运行的 Pod

调度的时候，需要获取一个调度的 pod，即 `sched.NextPod()` ，其中调用了 SchedulingQueue 的 `Pop()` 方法。

当 `activeQ` 中没有元素，会通过 `p.cond.Wait()` 阻塞，直到 `podBackoffQ` 或者 `unschedulableQ` 将元素加入 `activeQ` 并通过 `cond.Broadcast()` 来唤醒。

### 将 Pod 加入 activeQ

当 pod 加入 `activeQ` 后，还会从 `unschedulableQ` 以及 `podBackoffQ` 中删除对应 pod 的信息，并使用 `cond.Broadcast()` 来唤醒阻塞的 Pop。

```go
// Add adds a pod to the active queue. It should be called only when a new pod
// is added so there is no chance the pod is already in active/unschedulable/backoff queues
func (p *PriorityQueue) Add(pod *v1.Pod) error {
    p.lock.Lock()
    defer p.lock.Unlock()
    pInfo := p.newQueuedPodInfo(pod)
    if err := p.activeQ.Add(pInfo); err != nil {
        klog.Errorf("Error adding pod %v to the scheduling queue: %v", nsNameForPod(pod), err)
        return err
    }
    if p.unschedulableQ.get(pod) != nil {
        klog.Errorf("Error: pod %v is already in the unschedulable queue.", nsNameForPod(pod))
        p.unschedulableQ.delete(pod)
    }
    // Delete pod from backoffQ if it is backing off
    if err := p.podBackoffQ.Delete(pInfo); err == nil {
        klog.Errorf("Error: pod %v is already in the podBackoff queue.", nsNameForPod(pod))
    }
    metrics.SchedulerQueueIncomingPods.WithLabelValues("active", PodAdd).Inc()
    p.PodNominator.AddNominatedPod(pod, "")
    p.cond.Broadcast()

    return nil
}
```

### 当 Pod 调度失败时进入失败队列

当 pod 调度失败时，会调用 `sched.Error()` ，其中调用了 `p.AddUnschedulableIfNotPresent()` .

决定 pod 调度失败时进入 `podBackoffQ` 还是 `unschedulableQ` ：如果 `moveRequestCycle` 大于 `podSchedulingCycle` ，则进入 `podBackoffQ` ，否则进入 `unschedulableQ` .



何时 `moveRequestCycle >= podSchedulingCycle` ：

1. 我们在集群资源变更的时候（例如添加 Node 或者删除 Pod），会有回调函数尝试将 `unschedulableQ` 中之前因为资源不满足需求的 pod 放入 `activeQ` 或者 `podBackoffQ` ，及时进行调度。
2. 调度队列会每隔 30s 定时运行 `flushUnschedulableQLeftover` ，尝试调度 `unschedulableQ` 中的 pod。

这两者都会调用 `movePodsToActiveOrBackoffQueue` 函数，并将 `moveRequestCycle` 设为 `p.schedulingCycle`.

```go
func (p *PriorityQueue) movePodsToActiveOrBackoffQueue(podInfoList []*framework.QueuedPodInfo, event string) {
    ...
    p.moveRequestCycle = p.schedulingCycle
    p.cond.Broadcast()
}
```

### podBackoffQ 中 pod 的生命周期

**加入 podBackoffQ**

有两种情况会让 pod 加入 podBackoffQ：

1. 调度失败。如果调度失败，并且集群资源发生变更，即 `moveRequestCycle >= podSchedulingCycle` ，pod 就会加入到 podBackoffQ 中。
2. 从 unschedulableQ 中转移。当集群资源发生变化的时候，最终会调用 `movePodsToActiveOrBackoffQueue` 将 unschedulableQ 的 pod 转移到 podBackoffQ 或者 activeQ 中。转移到 podBackoffQ 的条件是 `p.isPodBackingoff(pInfo)` ，即 pod 仍然处于 backoff 状态。

**退出 podBackoffQ**

调度器会定时让 pod 从 podBackoffQ 转移到 activeQ 中。

在 `sched.SchedulingQueue.Run` 中运行的 `flushBackoffQCompleted` cronjob 会每隔 1s 按照优先级（优先级是按照 backoffTime 排序）依次将满足 backoffTime 条件的 pod 从 podBackoffQ 转移到 activeQ 中，直到遇到一个不满足 backoffTime 条件的 pod。

### unschedulableQ 中 pod 的生命周期

**加入 unschedulableQ**

只有一种情况会让 pod 加入 unschedulableQ，那就是调度失败。如果调度失败，并且集群资源没有发生变更，即 `moveRequestCycle < podSchedulingCycle` ，那么 pod 就会加入到 unschedulableQ 中。

**退出 unschedulableQ**

调度器会同样定时让 pod 从 unschedulableQ 转移到 podBackoffQ 或者 activeQ 中。

在 `sched.SchedulingQueue.Run` 中运行的 `flushUnschedulableQLeftover` 最终会调用 `movePodsToActiveOrBackoffQueue` 将 pod 分别加入到 podBackoffQ 或者 activeQ 中。