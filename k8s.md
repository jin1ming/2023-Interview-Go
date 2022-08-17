# Kubernetes笔记

目录：

[TOC]

## 核心组件

- etcd 保存了整个集群的状态，就是一个数据库；
- apiserver 提供了资源操作的唯一入口，并提供认证、授权、访问控制、API 注册和发现等机制；
- controller manager 负责维护集群的状态，比如故障检测、自动扩展、滚动更新等；
- scheduler 负责资源的调度，按照预定的调度策略将 Pod 调度到相应的机器上；
- kubelet 负责维护容器的生命周期，同时也负责 Volume（CSI）和网络（CNI）的管理；
- Container runtime 负责镜像管理以及 Pod 和容器的真正运行（CRI）；
- kube-proxy 负责为 Service 提供 cluster 内部的服务发现和负载均衡；

## k8s对象

- 资源对象

  Pod、ReplicaSet、ReplicationController、Deployment、StatefulSet、DaemonSet、Job、CronJob、HorizontalPodAutoscaling、Node、Namespace、Service、Ingress、Label、CustomResourceDefinition

- 存储对象

  Volume、PersistentVolume、Secret、ConfigMap

- 策略对象

  SecurityContext、ResourceQuota、LimitRange

- 身份对象

  ServiceAccount、Role、ClusterRole

### Pod

Pod在k8s中是创建和部署的最小单位，Pod代表集群中运行的进程，可能由一个或多个容器组合在一起共享资源。

Pod中可以共享两种资源：网络和存储。

很少直接在k8s中创建单个Pod。因为Pod生命周期短暂，用后即焚。

Pod不会自愈，Pod运行的Node故障或是调度器故障，Pod就会被删除。（Pod所在Node缺少资源或是Pod处于维护状态，Pod也会被驱逐）

一般使用Controller来管理Pod，提供副本管理、滚动升级和集群级别的自愈能力。

### ReplicaSet

ReplicaSet的目的是维护一组在任何时候都处于运行状态的Pod副本的稳定集合。

ReplicaSet确保任何时间都有指定数量的Pod副本在运行。

一般使用Deployment管理ReplicaSet，并向Pod提供声明式的更新及其他功能。

### Service

将运行在一组Pods上的应用程序公开为网络服务的抽象方法。

为Pods提供自己的IP地址，并为一组Pod提供相同的DNS名以及负载均衡。

解决了前端Pod无法跟踪后端Pod的IP变化问题。

## Controller

### 1. Deployment

管理Pod和ReplicaSet，应用场景有：

- 创建ReplicaSet -> ReplicaSet在后台创建Pods
- 滚动升级和回滚应用
- 扩容和缩容
- 暂停和继续Deployment
- 清理旧的ReplicaSet

### 2. StatefulSet

管理有状态Pod集合的部署和扩缩，并为这些Pod提供持久存储和持久标志符。

StatefulSet 为它们的每个 Pod 维护了一个有粘性的 ID。这些 Pod 是基于相同的规约来创建的， 但是不能相互替换：无论怎么调度，每个 Pod 都有一个永久不变的 ID。

对有以下需求的应用程序有很大价值：

- 稳定的、唯一的网络标识符。
- 稳定的、持久的存储。
- 有序的、优雅的部署和缩放。
- 有序的、自动的滚动更新。

限制：

- PersistentVolume 驱动
- 删除/收缩StatefulSet不会删除关联的存储卷
- 需要无头服务来负责Pod的网络标志
- 当删除 StatefulSets 时，StatefulSet 不提供任何终止 Pod 的保证。 为了实现 StatefulSet 中的 Pod 可以有序地且体面地终止，可以在删除之前将 StatefulSet 缩放为 0。
- 在默认 Pod 管理策略(OrderedReady) 时使用 滚动更新，可能进入需要**人工干预**才能修复的损坏状态。

### 3. DaemonSet

DaemonSet 确保全部（或者某些）节点上运行一个 Pod 的副本。

- 在每个节点上运行集群守护进程
- 在每个节点上运行日志收集守护进程
- 在每个节点上运行监控守护进程

### 4. Job

Job 会创建一个或者多个 Pods，并将继续重试 Pods 的执行，直到指定数量的 Pods 成功终止。 

随着 Pods 成功结束，Job 跟踪记录成功完成的 Pods 个数。 当数量达到指定的成功个数阈值时，任务（即 Job）结束。

删除 Job 的操作会清除所创建的全部 Pods。 

挂起 Job 的操作会删除 Job 的所有活跃 Pod，直到 Job 被再次恢复执行。

### 5. CronJob

CronJob 创建基于时隔重复调度的 Jobs。

**Cron时间表语法**

```
# ┌───────────── 分钟 (0 - 59)
# │ ┌───────────── 小时 (0 - 23)
# │ │ ┌───────────── 月的某天 (1 - 31)
# │ │ │ ┌───────────── 月份 (1 - 12)
# │ │ │ │ ┌───────────── 周的某天 (0 - 6)（周日到周一；在某些系统上，7 也是星期日）
# │ │ │ │ │                          或者是 sun，mon，tue，web，thu，fri，sat
# │ │ │ │ │
# │ │ │ │ │
# * * * * *
```



## yaml属性

​	只使用Maps、Lists结构。

- apiVersion：此处值是v1，这个版本号需要根据安装的Kubernetes版本和资源类型进行变化，记住不是写死的。
- kind：此处创建的是Pod，根据实际情况，此处资源类型可以是Deployment、Job、Ingress、Service等。
- metadata：包含Pod的一些meta信息，比如name、namespace、label等信息。
- spec：包括一些container，storage，volume以及其他Kubernetes需要的参数，以及诸如是否在容器失败时重新启动容器的属性。可在特定Kubernetes API找到完整的Kubernetes Pod的属性。

​	示例：

```yaml
#test-pod 
apiVersion: v1 #指定api版本，此值必须在kubectl apiversion中   
kind: Pod #指定创建资源的角色/类型   
metadata: #资源的元数据/属性   
  name: test-pod #资源的名字，在同一个namespace中必须唯一   
  labels: #设定资源的标签 
    k8s-app: apache   
    version: v1   
    kubernetes.io/cluster-service: "true"   
  annotations:            #自定义注解列表   
    - name: String        #自定义注解名字   
spec: #specification of the resource content 指定该资源的内容   
  restartPolicy: Always #表明该容器一直运行，默认k8s的策略，在此容器退出后，会立即创建一个相同的容器   
  nodeSelector:     #节点选择，先给主机打标签kubectl label nodes kube-node1 zone=node1   
    zone: node1   
  containers:   
  - name: test-pod #容器的名字   
    image: 10.192.21.18:5000/test/chat:latest #容器使用的镜像地址   
    imagePullPolicy: Never #三个选择Always、Never、IfNotPresent，每次启动时检查和更新（从registery）images的策略， 
                           # Always，每次都检查，默认策略
                           # Never，每次都不检查（不管本地是否有） 
                           # IfNotPresent，如果本地有就不检查，如果没有就拉取 
    command: ['sh'] #启动容器的运行命令，将覆盖容器中的Entrypoint,对应Dockefile中的ENTRYPOINT   
    args: ["$(str)"] #启动容器的命令参数，对应Dockerfile中CMD参数   
    env: #指定容器中的环境变量   
    - name: str #变量的名字   
      value: "/etc/run.sh" #变量的值   
    resources: #资源管理 
      requests: #容器运行时，最低资源需求，也就是说最少需要多少资源容器才能正常运行   
        cpu: 0.1 #CPU资源（核数），两种方式，浮点数或者是整数+m，0.1=100m，最少值为0.001核（1m） 
        memory: 32Mi #内存使用量   
      limits: #资源限制   
        cpu: 0.5   
        memory: 1000Mi   
    ports:   
    - containerPort: 80 #容器开发对外的端口 
      name: httpd  #名称 
      protocol: TCP   
    livenessProbe: #pod内容器健康检查的设置 
      httpGet: #通过httpget检查健康，返回200-399之间，则认为容器正常   
        path: / #URI地址   
        port: 80   
        #host: 127.0.0.1 #主机地址   
        scheme: HTTP   
      initialDelaySeconds: 180 #表明第一次检测在容器启动后多长时间后开始   
      timeoutSeconds: 5 #检测的超时时间   
      periodSeconds: 15  #检查间隔时间   
      #也可以用这种方法   
      #exec: 执行命令的方法进行监测，如果其退出码不为0，则认为容器正常   
      #  command:   
      #    - cat   
      #    - /tmp/health   
      #也可以用这种方法   
      #tcpSocket: //通过tcpSocket检查健康    
      #  port: number    
    lifecycle: #生命周期管理   
      postStart: #容器运行之前运行的任务   
        exec:   
          command:   
            - 'sh'   
            - 'yum upgrade -y'   
      preStop:#容器关闭之前运行的任务   
        exec:   
          command: ['service httpd stop']   
    volumeMounts:  #挂载持久存储卷 
    - name: volume #挂载设备的名字，与volumes[*].name 需要对应     
      mountPath: /data #挂载到容器的某个路径下   
      readOnly: True   
  volumes: #定义一组挂载设备   
  - name: volume #定义一个挂载设备的名字   
    #meptyDir: {}   
    hostPath:   
      path: /opt #挂载设备类型为hostPath，路径为宿主机下的/opt,这里设备类型支持很多种 
    #nfs
```

## Q&A

### Pod的创建流程

- 用户通过kubectl命名发起请求。

- apiserver通过对应的kubeconfig进行认证，认证通过后将yaml中的pod信息存到etcd。

- Controller-Manager通过apiserver的watch接口发现了pod信息的更新，执行该资源所依赖的拓扑结构整合，整合后将对应的信息交给apiserver，apiserver写到etcd，此时pod已经可以被调度了。

- Scheduler同样通过apiserver的watch接口更新到pod可以被调度，通过算法给pod分配节点，并将pod和对应节点绑定的信息交给apiserver，apiserver写到etcd，然后将pod交给kubelet。

- kubelet收到pod后，调用CRI接口去启动容器（先创建pause容器），调用CNI接口给pod创建pod网络（pause只共享命名空间不创建网络设备），调用CSI进行存储卷的挂载。

  pause容器两个作用，1.创建Linux命名空间方便之后共享，2.回收僵尸进程

- 网络，容器，存储创建完成后pod创建完成，等业务进程启动后，pod运行成功

### Pod出站流量

- Pod到Pod

  每个Pod有自己的IP地址，所有的Pod之间都可以保持三层网络的连通性。CNI就是用来实现这些网络功能的标准接口。

- Pod到Service

  Service就是Pod前面的4层负载均衡器。最常用的是ClusterIP，他会自动分配一个仅集群内可以访问的虚拟IP。

  通过kube-proxy组件实现这些功能，每台计算节点上都运行一个kube-proxy进程，通过复杂的iptables/IPVS规则在Pod和Service之间进行各种过滤和NAT。

- Pod到集群外

  通过SNAT来处理。SNAT做的工作就是将数据包的源从Pod内部的IP:Port替换为宿主机的IP:Port。当数据包返回时，再将目的地址从宿主机的IP:Port替换为Pod内的IP:Port，然后发送给Pod。中间过程对Pod来说时完全透明的，他们对地址转换不会有任何感知。

### Service流量转发过程

- 3（+1）种访问方式：

  - ClusterIP： 分配一个集群内的虚拟IP，默认类型。这种方式提供的ClusterIP只能在K8s集群内访问。

  - NodePort：在每个Node上提供一个相同静态端口，作为服务的端口映射。可以使用NodeIP:NodePort的方式从集群外进行服务的访问。

  - LoadBalancer：使用外部的负载均衡器来提供服务的访问功能。

  - 另外，还有一种ExternalName类型，通过将服务映射到某个域名来提供访问，这种方式需要使用1.7版本及以上的kube-dns组件，一般使用较少。

- kube-proxy三种运行模式：

  - userspace: 需要内核态用户态多次转换
  
  - iptables（默认）：依赖netfilter/iptable模块
  
    不同于userspace，iptables由kube-proxy动态的管理，kube-proxy不再负责转发，数据包的走向完全由iptables规则决定，这样的过程不存在内核态到用户态的切换，效率明显会高很多。但是随着service的增加，iptables规则会不断增加，导致内核十分繁忙（等于在读一张很大的没建索引的表）。
  
  - ipvs：也是基于netfilter实现的
  
    iptables是为防火墙设计的，IPVS则专门用于高性能负载均衡，并使用高效的数据结构Hash表，允许几乎无限的规模扩张。
  
    假设要禁止上万个IP访问我们的服务器，如果用iptables的话，就需要一条一条的添加规则，会生成大量的iptabels规则；但是用ipset的话，只需要将相关IP地址加入ipset集合中即可，这样只需要设置少量的iptables规则即可实现目标。
  
    由于ipvs无法提供包过滤、地址伪装、SNAT等功能，所以某些场景下（比如NodePort的实现）还要与iptables搭配使用。

### k8s接口与传统接口的区别

基于RESTful，不同在于list-watch

**Why list-watch?** 如果采用客户端轮询，将会加大`apiserver`压力，同时实时性很低；如果`apiserver`主动发起`HTTP`请求，无法保证消息的可靠性、大量端口占用问题；如果使用`TCP`~~负载高、引起互联网拥塞、大量端口占用~~。

使用`list`来描述返回资源集合的操作，以便于返回单个资源的`get`操作相区分。（实际使用的还是`GET`）

`watch`通过`http长链接`来实现快速检测变更。使用的`Chunked transfer encoding(分块传输编码)`，出现于HTTP/1.1，即response`的`HTTP Header`中设置`Transfer-Encoding`的值为`chunked。服务器会返回所提供的 `resourceVersion` 之后发生的所有变更（创建、删除和更新）。

> HTTP 分块传输编码允许服务器为动态生成的内容维持 HTTP 持久链接。通常，持久链接需要服务器在开始发送消息体前发送Content-Length消息头字段，但是对于动态生成的内容来说，在内容创建完之前是不可知的。使用分块传输编码，数据分解成一系列数据块，并以一个或多个块发送，这样服务器可以发送数据而不需要预先知道发送内容的总大小。

### k8s开发系统上线前夕要做些什么

（service loadbalancer负载均衡、k8s高可用、系统数据持久化、节点漂移备份）

### pod内容器是怎么共享网络和内存的

pause父容器来共享namespace

### 怎样让一个pod优先提供服务

设置pod优先级（priority）？

### list-watch用的http什么机制

chunked机制

### k8s高可用

- 堆控制平面+etcd节点
- 外部etcd节点

### kube-scheduler调度机制

维护一个`PodQueue`，按照一定算法将节点分配给Pods。

调度器先在集群中找到一个 Pod 的所有可调度节点【**过滤**】，然后根据一系列函数对这些可调度节点打分【**打分**】， 选出其中得分最高的 Node 来运行 Pod。之后，调度器将这个调度决定通知给 kube-apiserver，这个过程叫做 *绑定*。

在做调度决定时需要考虑的因素包括：单独和整体的资源请求、硬件/软件/策略限制、 亲和以及反亲和要求、数据局域性、负载间的干扰等等。

### iptables 五表五链工作原理，DNAT发生在哪条链？

iptables的底层实现是**netfilter**，其架构是在整个网络流程的若干位置放置一些钩子，并在每个钩子上挂载一些处理函数进行处理。它作为一个通用的、抽象的框架，提供一整套hook函数的管理机制，是的数据包过滤、包处理（设置标志位、修改TTL）、地址伪装、网络地址转换、透明代理、访问控制、基于协议类型的链接跟踪，甚至带宽限速等功能成为可能。

- **5 chain**

  IP层的5个钩子点的位置，对应iptables就是5条**内置链**，分别是**PREROUTING、POSTROUTING、INPUT、OUTPUT和FORWAR**。支持用户自定义链。

  - PREROUTING：可以在此处进行DNAT（destination NAT POSTROUTING，用于互联网访问局域网）
  - POSTROUTING：可以在此处进行SNAT（source NAT POSTROUTING，用于局域网访问互联网）
  - INPUT：处理输入本地进程的数据包
  - OUTPUT：处理本地进程的输出数据包
  - FORWARD：处理转发到其他机器/network namespace的数据包

- **5 table**

  优先级从高到低为：raw、mangle、nat、filter、security，不支持用户自定义表。

  - filter表：用于控制到达某条链上的数据包是否继续放行、直接丢弃（drop）或拒绝（reject）
  - nat标：用于修改数据包的源和目的地址
  - mangle表：用于修改数据包的IP头信息
  - raw表：iptables是有状态的，即iptables对数据包有链接追踪（connection tracking）机制，而raw是用来去除这种追踪机制的
  - security表：用于在数据包上应用SELinux

  不是每个链上都拥有这5个表：

  - raw存在于PREROUTING和OUTPUT。对应输入和输出经过的第一条链。
  - mangle存在于所有链中。
  - nat(SNAT)存在于POSTROUTING和INPUT。
  - nat(DNAT)存在于PREROUTING和OUTPUT。
  - filter、security存在于FORWARD、INPUT、OUTPUT。

- **rule**

  iptables的表示所有规则的5个逻辑集合，一跳iptables规则包含两部分信息：匹配条件和动作。

  匹配条件，即匹配数据包被这条iptables规则“捕获”的条件，例如协议类型、源IP、目的IP、源端口、目的端口、连接状态等。每条iptables规则允许多个匹配条件任意组合，从而实现多条件的匹配，多条件之间是逻辑与关系。

  常见动作有:

  - DROP：直接将数据包丢弃。应用场景：不让某个数据源意识到系统的存在，可以用来模拟宕机。
  - REJECT：给客户端返回connection refused或destination unreachable报文，应用场景：不让某个数据源访问系统，提示这里没有想要的服务内容。
  - QUEUE：将数据包放入用户空间的队列，供用户空间的程序处理。
  - RETURN：跳出当前链，该链后续规则不再执行。
  - ACCEPT：统一数据包通过，继续执行后续的规则。
  - JUMP：跳转到其他用户自定义的链继续执行。

### Deployment、DaemonSet、StatefulSet更新策略

- Deployment

  可以实时修改Deployment的内容并应用，k8s自动完成更新，更新发生错误可以用Rollback恢复版本

  支持两种更新策略：

  - **Recreate**：重建，spec.strategy.type=Recreate 表示Deployment会在更新Pod时先杀掉所有正在运行的Pod然后创建新Pod
  - **Rolling Update**：滚动更新 spec.strategy.type=RollingUpdate, 通过滚动更新的方式逐个更新Pod。
    - maxUnavailable 最大不可用：spec.strategy.type.RollingUpdate.MaxUnavailable 用于指定更新过程中不可用状态Pod数量的上限。 可整数可百分比
    - maxSurge最大超出： spec.strategy.type.RollingUpdate.MaxSurge 用于指定更新过程中Pod总数超过Pod期望副本数量部分的最大值。可整数可百分比

  使用 rollout history 检查 deploy的历史记录 --revision 制定版本
  使用 rollout undo 回滚到上个部署版本， 加 参数-to-revision 指定回滚版本

- StatefulSet

  3种升级策略：

  - **OnDelete** ： 默认升级策略，在创建好新的StatefulSetSet配置之后，新的Pod不会被自动创建，**用户需要手动删除旧版本的Pod**，才出发新建操作。
  - **RollingUpdate**： StatefulSet 控制器将删除并重新创建 StatefulSet 中的每个 Pod。它将按照与 Pod 终止相同的顺序进行（从最大的序数到最小的），依次更新每个 Pod。Kubernetes 控制平面会等到更新的 Pod 运行并准备好，然后再更新其前任。如果您已设置.spec.minReadySeconds（请参阅“最小就绪秒数”），则控制平面会在 Pod 准备就绪后额外等待该时间，然后再继续。
  - **partitioned**：分区滚动更新。如果指定了分区，则在更新 StatefulSet 时，将更新序数大于或等于该分区的所有 Pod .spec.template。所有序号小于分区的 Pod 都不会更新，即使删除了，也会在之前的版本中重新创建。如果 StatefulSet 的.spec.updateStrategy.rollingUpdate.partition大于其.spec.replicas，则对其的更新.spec.template将不会传播到其 Pod。在大多数情况下，您不需要使用分区，但如果您想要暂存更新、roll out Canary 或执行分阶段roll out，它们会很有用

- DaemonSet

  - **OnDelete** ： 默认升级策略，在创建好新的DaemonSet配置之后，新的Pod不会被自动创建，用户需要手动删除旧版本的Pod，才出发新建操作。
  - **RollingUpdate**： 旧版本的POD 将被自动杀掉，然后自动创建新版的DaemonSet Pod。与Deployment 不同为不支持查看和管理DaemonSet的更新记录；回滚操作是通过再次提交旧版本配置而不是 rollback命令实现