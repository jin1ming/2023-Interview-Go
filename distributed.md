# 分布式系统

#### 1. 一致性

一般来讲，我们将一致性分为三类：

1. **强一致性**：系统写入了什么，读出来的就是什么。
2. **弱一致性**：写操作完成后，系统不能保证后续的访问都能读到更新后的值。
3. **最终一致性**：如果对某个对象没有新的写操作了，最终所有后续访问都能读到相同的最近更新的值。

**业界比较推崇的是最终一致性级别，但是某些对数据一致性要求十分严格的场景（比如银行转账）还是要保证强一致性。**



#### 2. 可用性

- 可用性说的是，任何来自客户端的请求，不管访问哪个非故障节点，都能得到响应数据。
- 但不保证是同一份最新数据，可用性强调的是服务可用。



#### 3. CAP理论

`CAP`理论对分布式系统的特性做了高度抽象，形成了三个指标：

1. **一致性（`Consistency`）**：
   - 一致性说的是，所有节点都能访问到同一份最新的数据副本。
   - 客户端的每次读操作，不管访问哪个节点，要么读到的都是同一份最新写入的数据，要么读取失败。
   - **一致性强调的是数据正确。**
   - 可以把一致性看作是分布式系统，对访问自己的客户端的一种承诺：
     - 不管你访问哪个节点，要么我给你返回的都是绝对一致的最新写入的数据，要么你读取失败。
     - 对客户端而言，每次读都能读取到最新写入的数据。
     - **一致性这个指标，描述的是分布式系统非常重要的一个特性，强调的是数据正确。**
2. **可用性（`Availability`）**：
   - 非故障的节点在合理的时间内返回合理的响应（不是错误或者超时的响应）。
   - 可用性说的是，任何来自客户端的请求，不管访问哪个非故障节点，都能得到响应数据，但不保证是同一份最新数据。
   - **可用性这个指标，强调的是服务可用，但不保证数据正确。**
   - 也可以把可用性看作是分布式系统对访问本系统的客户端的另外一种承诺：
     - 我尽力给你返回数据，不会不响应你，但是我不保证每个节点给你的数据都是最新的。
3. **分区容错性（`Partition Tolerance`）**：
   - 分布式系统出现网络分区的时候，仍然能够对外提供服务。
   - 当节点间出现任意数量的消息丢失或高延迟的时候，系统仍然在继续工作。
   - **分区容错性这个指标，强调的是集群对分区故障的容错能力。**
   - 分区容错性是分布式系统对访问本系统的客户端的一种承诺：
     - 不管我的内部出现什么样的数据同步问题，我会一直运行。
   - 因为分布式系统与单机系统不同，它涉及到多节点间的通讯和交互，节点间的分区故障是必然发生的。
     - 因此，**在分布式系统中分区容错性是必须要考虑的！**
   - **follow-up：什么是网络分区？**
     - 在分布式系统中，多个节点之间的网络本来是连通的。
     - 但是因为某些故障（比如部分节点网络出了问题），某些节点之间不连通了。
     - 整个网络就分成了几个区域，这就叫网络分区。

##### follow-up： CAP不可能三角

> **`CAP`不可能三角说的是，对于一个分布式系统而言，一致性（`Consistency`）、可用性（`Availability`）、分区容错性（`Partition Tolerance`）三个指标不可兼得，只能在`3`个指标中选择`2`个。**

1. **不是所谓的“3选2”**

   大部分人解释这一定律时，常常简单的表述为：“一致性、可用性、分区容错性三者你只能同时达到其中两个，不可能同时达到”。

   - 实际上这是一个非常具有误导性质的说法。

   > - **当发生网络分区的时候，如果我们要继续服务，那么强一致性和可用性只能2选1。**
   > - **也就是说当网络分区之后`P`是前提，决定了`P`之后才有`C`和`A`的选择。**
   > - **也就是说分区容错性（`Partition tolerance`）我们是必须要实现的。**
   >
   > ***简而言之就是：CAP理论中分区容错性`P`是一定要满足的，在此基础上，只能满足可用性`A`或者一致性`C`。***

   因此，**分布式系统理论上不可能选择`CA`架构，只能选择`CP`或者`AP`架构。**

   ![image-20210325113753731](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210325113753731.png)

   1. **`CA`模型**：
      - 在分布式系统中不存在。
      - 因为舍弃`P`，意味着舍弃分布式系统。
        - 比如单机版关系型数据库`MySQL`，如果`MySQL`要考虑主备或集群部署时，它必须考虑`P`。
   2. **`CP`模型**：
      - 采用`CP`模型的分布式系统，舍弃了可用性。
      - 一定会读到最新数据，不会读到旧数据。
      - 一旦因为消息丢失、延迟过高发生了网络分区，就影响用户的体验和业务的可用性。
        - 比如基于`Raft`的强一致性系统，此时可能无法执行读操作和写操作。
      - `CP`模型的典型应用有：`Etcd`、`Consul`和`HBase`。
   3. **`AP`模型**：
      - 采用`AP`模型的分布式系统，舍弃了一致性，实现了服务的高可用。
      - 用户访问系统的时候，都能得到响应数据。
      - 不会出现响应错误，但会读到旧数据。
      - `AP`模型的典型应用有：`Cassandra`和`DynamoDB`。

2. **为什么不能同时保证`CA`呢？**

   举个例子：

   - 若系统出现“分区”，系统中的某个节点在进行写操作。
   - 为了保证`C`（一致性），必须要禁止其他节点的读写操作，这就和`A`（可用性）发生冲突了。
   - 如果要保证`A`（可用性），其他节点的读写操作正常的话，那就和`C`（一致性）发生冲突了。

3. **选择`CP`模型还是`AP`模型？**

   - 选择关键在于当前的业务场景，没有定论。
     - 比如对于需要确保强一致性的场景（比如银行），一般会选择`CP`模型。



#### 4. ACID理论

- **A（Atomic) 原子性**：
  - 事务是不可分割的工作单位。
  - 一个事务要么成功，要么不成功。
  - 事务中任何一个`SQL`语句如果执行失败，已经成功执行的`SQL`语句必须撤销，数据库状态应该退回到执行事务前的状态。
  - 回滚可以用回滚日志（`Undo Log`）来实现。
  - 回滚日志记录着事务所执行的修改操作，在回滚时反向执行这些修改操作即可。
- **C (Consistency) 一致性**：
  - 数据库在事务执行前后都保持一致性状态。
  - 在一致性状态下，所有事务对同一个数据的读取结果都是相同的。
- **I  (Isolation) 隔离性**：
  - 一个事务所做的修改在最终提交前，对其他事务是不可见的。
- **D  (Durability）持久性**：
  - 一旦事务提交，则其所做的修改将会永远保存到数据库中。
  - 即使系统发生崩溃，事务执行的结果也不能丢失。
  - 系统发生崩溃可以用重做日志（`Redo Log`）进行恢复，从而实现持久性。
  - 与回滚日志记录数据的逻辑修改不同，重做日志记录的是数据页的物理修改。



#### 5. BASE理论

> **`BASE`是基本可用（`Basically Available`）、软状态（`Soft State`）和最终一致性（`Eventually Consistent`）三个短语的缩写。**
>
> `BASE`理论是对`CAP`中一致性和可用性权衡的结果：
>
> - **`BASE`理论的的核心思想是：即使无法做到强一致性，但每个应用都可以根据自身业务特点，采用适当的方式来使系统达到最终一致性。**
> - **也就是牺牲数据的一致性来满足系统的高可用性，系统中一部分数据不可用或者不一致时，仍然需要保持系统整体“主要可用”。**
> - **`BASE`理论来源于对大规模互联网系统实践的总结，是基于`CAP`理论逐步演化而来的，它大大降低了我们对系统的要求。**
> - **`BASE`理论本质上是对`CAP`的延伸和补充。更具体地说，是对`CAP`中`AP`方案的一个补充：**
>   - **`AP`方案只是在系统发生分区的时候放弃一致性，而不是永远放弃一致性。**
>   - **在分区故障恢复后，系统应该达到最终一致性。这一点就是`BASE`理论延伸的地方。**

##### 基本可用

- 基本可用是指分布式系统在出现不可预知的故障的时候，允许损失部分可用性。

- 但是，这绝不等价于系统不可用。

- 即，在分布式系统出现故障的时候，保证核心可用，允许损失部分可用性。

- 什么叫允许损失部分可用性呢？

  - **响应时间上的损失**：
    - 正常情况下，处理用户请求需要`0.5`返回结果。
    - 但是由于系统出现故障，处理用户请求的时间变为`3s`。
  - **系统功能上的损失**：
    - 正常情况下，用户可以使用系统的全部功能。
    - 但是由于系统访问量突然剧增，系统的部分非核心功能无法使用。
  - 例如：
    - 电商在做促销时，为了保证购物系统的稳定性，部分消费者可能会被引导到一个降级的页面。

- **实现基本可用的4板斧——流量削峰、延迟响应、体验降级、过载保护**

  - 基本可用是说，当分布式系统在出现不可预知的故障时，**允许损失部分功能的可用性，保障核心功能的可用性。**
  - 可以把基本可用理解成，当系统节点出现大规模故障的时候。
  - 比如专线的光纤被挖断、突发流量导致系统过载（出现了突发事件，服务被大量访问），这个时候可以通过服务降级，牺牲部分功能的可用性，保障系统的核心功能可用。

  就拿`12306`订票系统基本可用的设计为例，这个订票系统在春运期间，因为开始售票后先到先得的缘故，会出现极其海量的请求峰值，如何处理这个问题呢？

  1. 咱们可以在不同的时间，出售不同区域的票，将访问请求错开，削弱请求峰值。
     - 比如，在春运期间，深圳出发的火车票在`8`点开售，北京出发的火车票在`9`点开售。这就是我们常说的**流量削峰**。
  2. 另外，你可能已经发现了，在春运期间，自己提交的购票请求，往往会在队列中排队等待处理，可能几分钟或十几分钟后，系统才开始处理，然后响应处理结果，这就是你熟悉的**延迟响应**。 
  3. 再比如，你正负责一个互联网系统，突然出现了网络热点事件，好多用户涌进来，产生了海量的突发流量，系统过载了，大量图片因为网络超时无法显示。
     - 那么这个时候你可以通过哪些方法，保障系统的基本可用呢？
     - 相信你马上就能想到**体验降级**。
     - 比如用小图片来替代原始图片，通过降低图片的清晰度和大小，提升系统的处理能力。
  4. 然后你还能想到**过载保护**：
     - 比如把接收到的请求放在指定的队列中排队处理，如果请求等待时间超时了（假设是 `100ms`），这个时候直接拒绝超时请求。
     - 再比如队列满了之后，就清除队列中一定数量的排队请求，保护系统不过载，实现系统的基本可用。

##### 软状态

- **软状态**指允许系统中的数据存在中间状态（`CAP`理论中的数据不一致）。
- 并认为该中间状态不会影响系统整体可用性。
- 即允许系统不同节点的数据副本之间进行同步的过程存在时延。

##### 最终一致性

- 最终一致性强调的是，系统中所有的数据副本，在经过一段时间的同步后，最终能达到一致的状态。
- 也就是说，在数据一致性上，存在一个短暂的延迟。
- 因此，**最终一致性的本质是需要系统保证最终数据能够达到一致，而不需要实时保证系统数据的强一致性**。
- `ACID`要求强一致性，通常运用在传统的数据库系统上。
- 而`BASE`要求最终一致性，通过牺牲强一致性来达到可用性，通常应用在大型分布式系统中。
- **实现最终一致性的具体方式**：
  - **读时修复**：
    - 在读取数据时，检测数据的不一致，进行修复。
    - 比如`Cassandra`的`Read Repair`实现。
      - 具体来说，在向`Cassandra`系统查询数据的时候，如果检测到不同节点的副本数据不一致，系统就会自动修复数据。
  - **写时修复**：
    - 在写入数据，检测数据的不一致时，进行修复。
    - 比如`Cassandra`的`Hinted Handoff`实现。
      - 具体来说，`Cassandra`集群的节点之间远程写数据的时候，如果写失败就将数据缓存下来，然后定时重传，修复数据的不一致性。
  - **异步修复**：
    - 这个是最常用的方式，通过定时对账检测副本数据的一致性，并修复。

在实际的分布式场景中，不同业务单元和组件对一致性的要求是不同的，因此`ACID`和`BASE`往往会结合在一起使用。



#### 6. Raft算法

- `Raft`算法属于`Multi-Paxos`算法，它是在兰伯特`Multi-Paxos`思想的基础上，做了一些简化和限制。
- 比如增加了日志必须是连续的，只支持领导者、跟随者和候选人三种状态，在理解和算法实现上都相对容易许多。

> **Raft算法是现在分布式系统开发首选的共识算法。**
>
> **从本质上说，Raft 算法是通过一切以领导者为准的方式，实现一系列值的共识和各节点日志的一致。**

##### 有哪些成员身份？——领导者（Leader）、跟随者（Follower）、候选人（Follower）

- 成员身份，又叫做服务器节点状态。
- `Raft`算法支持领导者（`Leader`）、跟随者（`Follower`）和候选人（`Candidate`）`3`种状态。
- 在任何时候，每一个服务器节点都处于这`3`个状态中的`1`个。

- **跟随者**：

  - 就相当于普通群众，默默地接收和处理来自领导者的消息。
  - 当等待领导者心跳信息超时的时候，就主动站出来，推荐自己当候选人。

- **候选人**：

  - 候选人将向其他节点发送请求投票（`RequestVote`）`RPC` 消息，通知其他节点来投票。
  - 如果赢得了大多数选票，就晋升当领导者。

- **领导者**：

  - 蛮不讲理的霸道总裁，一切以我为准，平常的主要工作内容就是`3`部分：**处理写请求**、**管理日志复制**和**不断地发送心跳信息**。
  - 心跳信息：通知其他节点“我是领导者，我还活着，你们现在不要发起新的选举，找个新领导者来替代我。”

  注意，**Raft是强领导者模型，集群中只能有一个领导者。**

##### 选举领导者

- 首先，在初始状态下，集群中所有的节点都是跟随者的状态。

  ![image-20210328140607503](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328140607503.png)

- `Raft`算法实现了随机超时时间的特性。

- 也就是说，每个节点等待领导者节点心跳信息的超时时间间隔是随机的。

  - 通过上面的图片你可以看到，集群中没有领导者，而节点`A`的等待超时时间最小（`150ms`），它会最先因为没有等到领导者的心跳信息，发生超时。

- 这个时候，节点`A`就增加自己的任期编号，并推举自己为候选人，先给自己投上一张选票。

- 然后向其他节点发送请求投票`RPC`消息，请它们选举自己为领导者。

  ![image-20210328140721477](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328140721477.png)

- 如果其他节点接收到候选人`A`的请求投票`RPC`消息，在编号为`1`的这个任期内，也还没有进行过投票，那么它就把选票投给节点`A`，并增加自己的任期编号。

  ![image-20210328140741036](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328140741036.png)

- 如果候选人在选举超时时间内赢得了大多数的选票，那么它就会成为本届任期内新的领导者。

  ![image-20210328140757617](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328140757617.png)

- 节点`A`当选领导者后，它会周期性地发送心跳信息，通知其他服务器我是领导者，阻止跟随者发起新的选举、篡权。

  ![image-20210328140816537](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328140816537.png)

##### 选举过程的四个问题

###### 1. 节点间如何通讯？

- 在`Raft`算法中，服务器节点间的沟通联络采用的是远程过程调用（`RPC`）。
- 在领导者选举过程中，需要用到这样两类的`RPC`：
  1. **请求投票（`Request Vote`）`RPC`**：
     - 是由候选人在选举期间发起的，用于通知各节点进行投票。
  2. **日志服务（`Append Entries`）`RPC`**：
     - 是由领导者发起的，用于复制日志和提供心跳信息。
- 注意：**日志复制`RPC`只能由领导者发起，这是实现强领导者模型的关键之一**。

###### 2. 什么是任期？

- 我们知道，议会选举中的领导者是有任期的，领导者任命到期后，要重新开会再次选举。
- `Raft`算法中的领导者也是有任期的，每个任期由单调递增的数字（任期编号）标识。
  - 比如节点`A` 的任期编号是`1`。
- 任期编号是随着选举的举行而变化的，这是在说下面几点：
  1. 跟随者在等待领导者心跳信息超时后，推举自己为候选人时，会增加自己的任期号。
     - 比如节点`A`的当前任期编号为`0`，那么在推举自己为候选人时，会将自己的任期编号增加为`1`。
  2. 如果一个服务器节点，发现自己的任期编号比其他节点小，那么它会更新自己的编号到较大的编号值。
     - 比如节点`B`的任期编号是`0`，当收到来自节点`A`的请求投票`RPC`消息时，因为消息中包含了节点`A` 的任期编号，且编号为`1`，那么节点`B` 将把自己的任期编号更新为`1`。
- **与现实议会选举中的领导者的任期不同，Raft 算法中的任期不只是时间段，而且任期编号的大小，会影响领导者选举和请求的处理。**
  1. 在`Raft`算法中约定，如果一个候选人或者领导者，发现自己的任期编号比其他节点小，那么它会立即恢复成跟随者状态。
     - 比如分区错误恢复后，任期编号为`3`的领导者节点`B`，收到来自新领导者的，包含任期编号为`4`的心跳消息，那么节点`B`将立即恢复成跟随者状态。
  2. `Raft`算法还约定，如果一个节点接收到一个包含较小的任期编号值的请求，那么它会直接拒绝这个请求。
     - 比如节点`C`的任期编号为`4`，收到包含任期编号为`3`的请求投票`RPC`消息，那么它将拒绝这个消息。

###### 3. 选举有哪些规则？

- 在议会选举中，比成员的身份、领导者的任期还要重要的就是选举的规则。

- 比如一人一票、弹劾制度等。

- “无规矩不成方圆”，在`Raft`算法中，也约定了选举规则，主要有这样几点：

  1. 领导者周期性地向所有跟随者发送心跳消息（即不包含日志项的日志复制`RPC`消息）。通知大家我是领导者，阻止跟随者发起新的选举。

  2. 如果在指定时间内，跟随者没有接收到来自领导者的消息，那么它就认为当前没有领导者，推举自己为候选人，发起领导者选举。

  3. 在一次选举中，赢得大多数选票的候选人，将晋升为领导者。

  4. 在一个任期内，领导者一直都会是领导者，直到它自身出现问题（比如宕机），或者因为网络延迟，其他节点发起一轮新的选举。

  5. 在一次选举中，每一个服务器节点最多会对一个任期编号投出一张选票，并且按照“先来先服务”的原则进行投票。

     - 比如节点`C`的任期编号为`3`，先收到了`1`个包含任期编号为`4`的投票请求（来自节点`A`）。
     - 然后又收到了`1`个包含任期编号为`4`的投票请求（来自节点`B`）。
     - 那么节点`C`将会把唯一一张选票投给节点`A`。
     - 当再收到节点`B`的投票请求`RPC`消息时，对于编号为`4`的任期，已没有选票可投了。

     ![image-20210328141440466](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328141440466.png)

  6. 日志完整性高的跟随者（也就是最后一条日志项对应的任期编号值更大，索引号更大），拒绝投票给日志完整性低的候选人。

     - 比如节点`B`的任期编号为`3`，节点`C`的任期编号是`4`。
     - 节点`B`的最后一条日志项对应的任期编号为`3`，而节点`C`为`2`，那么当节点 C 请求节点 `B` 投票给自己时，节点`B`将拒绝投票。

     ![image-20210328141622273](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328141622273.png)

注意：

- **选举是跟随者发起的，推举自己为候选人。**
- **大多数选票是指集群成员半数以上的选票。**
- **大多数选票规则的目标，是为了保证在一个给定的任期内最多只有一个领导者。**
- 在选举中，除了选举规则外，我们还需要避免一些会导致选举失败的情况，**比如同一任期内，多个候选人同时发起选举，导致选票被瓜分，选举失败。那么在`Raft`算法中，如何避免这个问题呢？答案就是随机超时时间。**


###### 4. 如何理解随机超时时间？

- 在议会选举中，常出现未达到指定票数，选举无效，需要重新选举的情况。
- 在`Raft`算法的选举中，也存在类似的问题，那它是如何处理选举无效的问题呢？
- **`Raft`算法巧妙地使用随机选举超时时间的方法，把超时时间都分散开来。**
  - **在大多数情况下只有一个服务器节点先发起选举，而不是同时发起选举，这样就能减少因选票瓜分导致选举失败的情况。**
- 在`Raft`算法中，随机超时时间是有 2 种含义的：
  1. 跟随者等待领导者心跳信息超时的时间间隔，是随机的。
  2. 如果候选人在一个随机时间间隔内，没有赢得过半票数，那么选举就无效。然后候选人发起新一轮的选举，也就是说，等待选举超时的时间间隔，是随机的。

###### Leader选举总结

- `Raft`算法和兰伯特的`Multi-Paxos`不同之处，主要有 2 点：
  1. 首先，在`Raft`中，不是所有节点都能当选领导者，只有日志较完整的节点（也就是日志完整度不比半数节点低的节点），才能当选领导者。
  2. 其次，在`Raft`中，日志必须是连续的。
- `Raft`算法通过任期、领导者心跳消息、随机选举超时时间、先来先服务的投票原则、大多数选票原则等，保证了一个任期只有一位领导，也极大地减少了选举失败的情况。
- 本质上，`Raft`算法以领导者为中心，选举出的领导者，以“一切以我为准”的方式，达成值的共识，和实现各节点日志的一致。

##### 日志复制

> **选完领导者之后，领导者需要处理来自客户端的写请求，并通过日志复制实现各节点日志的一致。**
>
> **在`Raft`算法中，副本数据是以日志的形式存在的，领导者接收到来自客户端写请求后，处理写请求的过程就是一个复制和应用（`Apply`）日志项到状态机的过程。**

###### 1. 如何理解日志

- 副本数据是以日志的形式存在的，日志是由日志项组成的。

- 日志项是一种数据格式，它主要包含用户指定的数据，也就是指令（`Command`），还包含一些附加信息，比如索引值（`Log index`）、任期编号（`Term`）。

  ![image-20210328142136861](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328142136861.png)

  - 指令：一条由客户端请求指定的、状态机需要执行的指令。可以将指令理解成客户端指定的数据。
  - 索引值：日志项对应的整数索引值。它其实就是用来标识日志项的，是一个连续的、单调递增的整数号码。
  - 任期编号：创建这条日志项的领导者的任期编号。

- 一届领导者任期，往往有多条日志项。而且日志项的索引值是连续的。

###### 2. 如何复制日志

- 可以把`Raft`的日志复制理解成一个优化后的二阶段提交（将二阶段优化成了一阶段），减少了一半的往返消息，也就是降低了一半的消息延迟。

- 具体过程：

  1. 首先，领导者进入第一阶段，通过日志复制（`AppendEntries`）`RPC`消息，将日志项复制到集群其他节点上。

  2. 接着，如果领导者接收到大多数的“复制成功”响应后，它将日志项应用到它的状态机，并返回成功给客户端。

     - 如果领导者没有接收到大多数的“复制成功”响应，那么就返回错误给客户端。

     - 领导者将日志项应用到它的状态机，怎么没通知跟随者应用日志项呢？

       - 这是`Raft`中的一个优化，领导者不直接发送消息通知其他节点应用指定日志项。
       - 因为领导者的日志复制`RPC`消息或心跳消息，包含了当前最大的，将会被提交（`Commit`）的日志项索引值。
       - 所以通过日志复制`RPC`消息或心跳消息，跟随者就可以知道领导者的日志提交位置信息。
       - 因此，当其他节点接受领导者的心跳消息，或者新的日志复制 RPC 消息后，就会将这条日志项应用到它的状态机。
       - 而这个优化，降低了处理客户端请求的延迟，将二阶段提交优化为了一段提交，降低了一半的消息延迟。

       ![image-20210328142307011](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328142307011.png)

- **日志复制过程总结**：

  1. 接收到客户端请求后，领导者基于客户端请求中的指令，创建一个新日志项，并附加到本地日志中。
  2. 领导者通过日志复制`RPC`，将新的日志项复制到其他的服务器。
  3. 当领导者将日志项，成功复制到大多数的服务器上的时候，领导者会将这条日志项应用到它的状态机中。
  4. 领导者将执行的结果返回给客户端。
  5. 当跟随者接收到心跳信息，或者新的日志复制`RPC`消息后，如果跟随者发现领导者已经提交了某条日志项，而它还没应用，那么跟随者就将这条日志项应用到本地的状态机中。

  不过，这是一个理想状态下的日志复制过程。在实际环境中，复制日志的时候，你可能会遇到进程崩溃、服务器宕机等问题，这些问题会导致日志不一致。那么在这种情况下，Raft 算法是如何处理不一致日志，实现日志的一致的呢？

###### 3. 如何实现日志的一致

- 在`Raft`算法中，领导者通过强制跟随者直接复制自己的日志项，处理不一致日志。

- 也就是说，`Raft`是通过以领导者的日志为准，来实现各节点日志的一致的。

- 具体有`2`个步骤：

  1. 首先，领导者通过日志复制`RPC`的一致性检查，找到跟随者节点上，与自己相同日志项的最大索引值。
     - 也就是说，这个索引值之前的日志，领导者和跟随者是一致的，之后的日志是不一致的了。
  2. 然后，领导者强制跟随者更新覆盖的不一致日志项，实现日志的一致。

- **详细过程**：

  - 为了方便演示，引入`2`个新变量：

    - `PrevLogEntry`：表示当前要复制的日志项，前面一条日志项的索引值。

      - 比如在图中，如果领导者将索引值为`8`的日志项发送给跟随者，那么此时 `PrevLogEntry`值为`7`。

    - `PrevLogTerm`：表示当前要复制的日志项，前面一条日志项的任期编号。

      - 比如在图中，如果领导者将索引值为`8`的日志项发送给跟随者，那么此时 `PrevLogTerm`值为`4`。

      ![image-20210328142447298](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328142447298.png)

    1. 领导者通过日志复制`RPC`消息，发送当前最新日志项到跟随者。
       - 为了演示方便，假设当前需要复制的日志项是最新的。
       - 这个消息的`PrevLogEntry`值为`7`，`PrevLogTerm`值为`4`。
    2. 如果跟随者在它的日志中，找不到与`PrevLogEntry`值为`7`、`PrevLogTerm`值为`4`的日志项，也就是说它的日志和领导者的不一致了，那么跟随者就会拒绝接收新的日志项，并返回失败信息给领导者。
    3. 这时，领导者会递减要复制的日志项的索引值，并发送新的日志项到跟随者。
       - 这个消息的`PrevLogEntry`值为`6`，`PrevLogTerm`值为`3`。
    4. 如果跟随者在它的日志中，找到了`PrevLogEntry`值为`6`、`PrevLogTerm`值为`3`的日志项，那么日志复制`RPC`返回成功。
       - 这样，领导者就知道在`PrevLogEntry`值为`6`、`PrevLogTerm`值为`3`的位置，跟随者的日志项与自己相同。
    5. 领导者通过日志复制`RPC`，复制并更新覆盖该索引值之后的日志项（也就是不一致的日志项），最终实现了集群各节点日志的一致。

  > - **从上面步骤中可以看到，领导者通过日志复制`RPC`一致性检查，找到跟随者节点上与自己相同日志项的最大索引值，然后复制并更新覆盖该索引值之后的日志项，实现了各节点日志的一致。**
  > - **需要注意的是，跟随者中的不一致日志项会被领导者的日志覆盖，而且领导者从来不会覆盖或者删除自己的日志。**

###### 4. 日志复制总结

- 在`Raft`中，副本数据是以日志的形式存在的，其中日志项中的指令表示用户指定的数据。
- 兰伯特的`Muliti-Paxos`不要求日志是连续的，但在`Raft`中日志必须是连续的。
- 而在`Raft`中，日志不仅是数据的载体，日志的完整性还影响领导者选举的结果。
- 也就是说，日志完整性最高的节点才能当选领导者。
- `Raft`是通过以领导者的日志为准，来实现日志的一致的。

可以发现，值的共识和日志的一致都是由领导者决定的，领导者的唯一性很重要。

那如果需要对集群进行扩容或者缩容，比如将`3`节点集群扩容为`5`节点集群，这时候是可能同时出现两个领导者的。

这时就需要解决这个问题。（成员变更）

##### 成员变更

- 在日常工作中，你可能会遇到服务器故障的情况，这时你就需要替换集群中的服务器。
- 如果遇到需要改变数据副本数的情况，则需要增加或移除集群中的服务器。
- 总的来说，在日常工作中，集群中的服务器数量是会发生变化的。

讲到这儿，也许你会问：“`Raft`是共识算法，对集群成员进行变更时（比如增加`2`台服务器），会不会因为集群分裂，出现`2`个领导者呢？”

- 的确会出现这个问题，因为`Raft`的领导者选举，建立在“大多数”的基础之上，那么当成员变更时，集群成员发生了变化，就可能同时存在新旧配置的`2`个“大多数”，出现`2`个领导者，破坏了 `Raft`集群的领导者唯一性，影响了集群的运行。
- 关于成员变更，不仅是`Raft`算法中比较难理解的一部分，非常重要，也是`Raft`算法中唯一被优化和改进的部分。
  - 比如，最初实现成员变更的是联合共识（`Joint Consensus`），但这个方法实现起来难。
  - 后来`Raft`的作者就提出了一种改进后的方法，单节点变更（`single-server changes`）。

先介绍一下“配置”（`Configuration`）这个词儿：

- 它就是在说集群是哪些节点组成的，是集群各节点地址信息的集合。
  - 比如节点`A`、`B`、`C`组成的集群，那么集群的配置就是`[A, B, C]`集合。

###### 成员变更的问题

- 在集群中进行成员变更的最大风险是，可能会同时出现`2`个领导者。

  - 比如在进行成员变更时，节点`A`、`B`和`C`之间发生了分区错误，节点`A`、`B`组成旧配置中的“大多数”，也就是变更前的`3`节点集群中的“大多数”，那么这时的领导者（节点`A`）依旧是领导者。

  - 另一方面，节点`C`和新节点`D`、`E`组成了新配置的“大多数”，也就是变更后的`5`节点集群中的“大多数”，它们可能会选举出新的领导者（比如节点`C`）。

  - 那么这时，就出现了同时存在`2`个领导者的情况。

    ![image-20210328143303112](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328143303112.png)

- 如果出现了`2`个领导者，那么就违背了“领导者的唯一性”的原则，进而影响到集群的稳定运行。

- 你要如何解决这个问题呢？也许有的同学想到了一个解决方法——**固定配置**。

  - 因为我们在启动集群时，配置是固定的，不存在成员变更，在这种情况下，`Raft`的领导者选举能保证只有一个领导者。
  - 也就是说，这时不会出现多个领导者的问题，那我可以先将集群关闭再启动新集群啊。
  - 也就是先把节点`A`、`B`、`C`组成的集群关闭，然后再启动节点`A`、`B`、`C`、`D`、`E`组成的新集群。
  - **这个方法不可行。 **
    - 为什么呢？因为你每次变更都要重启集群，意味着在集群变更期间服务不可用。
    - 肯定不行啊，太影响用户体验了。想象一下，你正在玩王者荣耀，时不时弹出一个对话框通知你：系统升级，游戏暂停`3`分钟。这体验糟糕不糟糕？
  - 既然这种方法影响用户体验，根本行不通，那**到底怎样解决成员变更的问题呢？**
    - **最常用的方法就是单节点变更。**

###### 如何通过单节点变更解决成员变更的问题？

> **单节点变更，就是通过一次变更一个节点实现成员变更。如果需要变更多个节点，那你需要执行多次单节点变更。**

- 比如将`3`节点集群扩容为`5`节点集群，这时你需要执行`2`次单节点变更，先将`3`节点集群变更为 `4`节点集群，然后再将`4`节点集群变更为`5`节点集群，就像下图的样子：

  ![image-20210328143546517](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328143546517.png)

- **这样一来，我们就通过一次变更一个节点的方式，完成了成员变更，保证了集群中始终只有一个领导者，而且集群也在稳定运行，持续提供服务。**

- **大多数情况下，不存在新旧配置两个“大多数”：**

  - 我想说的是，**在正常情况下，不管旧的集群配置是怎么组成的，旧配置的“大多数”和新配置的“大多数”都会有一个节点是重叠的。 **
    - 也就是说，不会同时存在旧配置和新配置`2`个“大多数”：

  ![image-20210328143838487](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328143838487.png)

  ![image-20210328143850157](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328143850157.png)

  - 从上图中你可以看到，不管集群是偶数节点，还是奇数节点，不管是增加节点，还是移除节点，新旧配置的“大多数”都会存在重叠（图中的橙色节点）。
  - 需要你注意的是，在分区错误、节点故障等情况下，如果我们并发执行单节点变更，那么就可能出现一次单节点变更尚未完成，新的单节点变更又在执行，导致集群出现`2`个领导者的情况。
  - 如果你遇到这种情况，可以在领导者启动时，创建一个`NO_OP`日志项（也就是空日志项），只有当领导者将`NO_OP`日志项应用后，再执行成员变更请求。

###### 成员变更总结

1. 成员变更的问题，主要在于进行成员变更时，可能存在新旧配置的`2`个“大多数”，导致集群中同时出现两个领导者，破坏了`Raft`的领导者的唯一性原则，影响了集群的稳定运行。
2. 单节点变更是利用“一次变更一个节点，不会同时存在旧配置和新配置`2`个‘大多数’”的特性，实现成员变更。
3. 因为联合共识实现起来复杂，不好实现，所以绝大多数 Raft 算法的实现，采用的都是单节点变更的方法（比如`Etcd`、`Hashicorp Raft`）。
   - 其中，`Hashicorp Raft`单节点变更的实现，是由`Raft`算法的作者迭戈·安加罗（`Diego Ongaro`）设计的，很有参考价值。



#### 7. 一致性哈希

- 虽然领导者模型简化了算法实现和共识协商，但写请求只能限制在领导者节点上处理，导致了集群的接入性能约等于单机，那么随着业务发展，集群的性能可能就扛不住了，会造成系统过载和服务不可用，这时该怎么办呢？
- 其实这是一个非常常见的问题。在我看来，这时我们就要通过分集群，突破单集群的性能限制了。
- 说到这儿，有同学可能会说了，分集群还不简单吗？加个`Proxy`层，由`Proxy`层处理来自客户端的读写请求，接收到读写请求后，通过对`Key`做哈希找到对应的集群就可以了啊。

##### Q: What's wrong with 哈希算法？A: 数据迁移成本高

- 是的，哈希算法的确是个办法，但它有个明显的缺点：当需要变更集群数时（比如从`2`个集群扩展为`3`个集群），这时大部分的数据都需要迁移，重新映射，数据的迁移成本是非常高的。

  - 假设我们有一个由`A`、`B`、`C`三个节点组成（为了方便演示，我使用节点来替代集群）的`KV` 服务，每个节点存放不同的`KV`数据：

    ![image-20210328145931077](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328145931077.png)

  - 那么，使用哈希算法实现哈希寻址时，到底有哪些问题呢？

    - 通过哈希算法，每个`key`都可以寻址到对应的服务器。

      - 比如，查询`key`是`key-01`，计算公式为`hash(key-01) % 3` ，经过计算寻址到了编号为`1`的服务器节点`A`（就像下图 的样子）：

        ![image-20210328145955731](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328145955731.png)

      - 但如果服务器数量发生变化，基于新的服务器数量来执行哈希算法的时候，就会出现路由寻址失败的情况，`Proxy`无法找到之前寻址到的那个服务器节点，这是为什么呢？

        - 想象一下，假如`3`个节点不能满足业务需要了。

        - 这时我们增加了一个节点，节点的数量从`3`变化为`4`，那么之前的`hash(key-01) % 3 = 1`，就变成了`hash(key-01) % 4 = X`，因为取模运算发生了变化，所以这个`X`大概率不是`1`（可能`X`为`2`）。

        - 这时你再查询，就会找不到数据了，因为`key-01`对应的数据，存储在节点`A` 上，而不是节点`B`：

          ![image-20210328150031456](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328150031456.png)

        - 同样的道理，如果我们需要下线`1`个服务器节点（也就是缩容），也会存在类似的可能查询不到数据的问题。

        - 而解决这个问题的办法，在于我们要迁移数据，基于新的计算公式`hash(key-01) % 4`，来重新对数据和节点做映射。

        - 需要你注意的是，数据的迁移成本是非常高的。

        - **迁移成本是非常高昂的，这在实际生产环境中也是无法想象的。**

- 那么如何解决哈希算法，数据迁移成本高的痛点呢？

  - 答案就是一致哈希（`Consistent Hashing`）。

##### 如何使用一致性哈希实现哈希寻址

- 一致哈希算法也用了取模运算，但与哈希算法不同的是，哈希算法是对节点的数量进行取模运算，而一致哈希算法是对`2^32`进行取模运算。

- 你可以想象下，一致哈希算法，将整个哈希值空间组织成一个虚拟的圆环，也就是哈希环：

  ![image-20210328150157077](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328150157077.png)

- 哈希环的空间是按顺时针方向组织的，圆环的正上方的点代表`0`，`0`点右侧的第一个点代表`1`，以此类推，`2`、`3`、`4`、`5`、`6`……直到`2^32-1`，也就是说`0`点左侧的第一个点代表`2^32-1`。

- 在一致哈希中，你可以通过执行哈希算法（为了演示方便，假设哈希算法函数为`c-hash()`），将节点映射到哈希环上。

  - 比如选择节点的主机名作为参数执行`c-hash()`，那么每个节点就能确定其在哈希环上的位置了：

    ![image-20210328150222793](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328150222793.png)

- 当需要对指定`key`的值进行读写的时候，你可以通过下面`2`步进行寻址：

  - 首先，将`key`作为参数执行`c-hash()`计算哈希值，并确定此`key`在环上的位置。
  - 然后，从这个位置沿着哈希环顺时针“行走”，遇到的第一节点就是`key`对应的节点。

- 例如，假设`key-01`、`key-02`、`key-03`三个`key`，经过哈希算法`c-hash()`计算后，在哈希环上的位置就像下图的样子：

  ![image-20210328150307415](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328150307415.png)

- 那么根据一致哈希算法，`key-01`将寻址到节点`A`，`key-02`将寻址到节点`B`，`key-03`将寻址到节点`C`。讲到这儿，你可能会问：“那一致哈希是如何避免哈希算法的问题呢？”

- 接下来分别以增加节点和移除节点为例，具体说一说一致性哈希是如何避免上面的问题的。

- 假设，现在有一个节点故障了（比如节点`C`）：

  ![image-20210328150335298](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328150335298.png)

- 你可以看到，`key-01`和`key-02`不会受到影响，只有`key-03`的寻址被重定位到`A`。

- 一般来说，在一致性哈希算法中，如果某个节点宕机不可用了，那么受影响的数据仅仅是，会寻址到此节点和前一节点之间的数据。

  - 比如当节点`C`宕机了，受影响的数据是会寻址到节点`B`和节点`C`之间的数据（例如`key-03`），寻址到其他哈希环空间的数据（例如`key-01`），不会受到影响。

- 那如果此时集群不能满足业务的需求，需要扩容一个节点（也就是增加一个节点，比如`D`）：

  ![image-20210328150409524](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328150409524.png)

- 你可以看到，`key-01`、`key-02`不受影响，只有`key-03`的寻址被重定位到新节点`D`。

  > **一般而言，在一致哈希算法中，如果增加一个节点，受影响的数据仅仅是，会寻址到新节点和前一节点之间的数据，其它数据也不会受到影响。**

- 总的来说，使用了一致哈希算法后，扩容或缩容的时候，都只需要重定位环空间中的一小部分数据。

- 也就是说，**一致哈希算法具有较好的容错性和可扩展性。**

需要注意的是，在哈希寻址中常出现这样的问题：

- 客户端访问请求集中在少数的节点上，出现了有些机器高负载，有些机器低负载的情况。
- 那么在一致哈希中，有什么办法能让数据访问分布的比较均匀呢？
- 答案就是虚拟节点。

##### 虚拟节点

> **虚拟节点是解决一致性哈希数据倾斜问题（数据访问分布不均问题）的。**

- 在一致性哈希中，如果节点太少，容易因为节点分布不均匀造成数据访问的冷热不均，也就是说大多数访问请求都会集中少量几个节点上：

  ![image-20210328150639004](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328150639004.png)

- 你能从图中看到，虽然有`3`个节点，但访问请求主要集中的节点`A`上。那如何通过虚拟节点解决冷热不均的问题呢？

- 其实，就是对每一个服务器节点计算多个哈希值，在每个计算结果位置上，都放置一个虚拟节点，并将虚拟节点映射到实际节点。

  - 比如，可以在主机名的后面增加编号，分别计算`Node-A-01`、`Node-A-02`、`Node-B-01`、`Node-B-02`、`Node-C-01`、`Node-C-02`的哈希值，于是形成`6`个虚拟节点：

    ![image-20210328150700167](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328150700167.png)

  - 你可以从图中看到，增加了节点后，节点在哈希环上的分布就相对均匀了。

  - 这时，如果有访问请求寻址到`Node-A-01`这个虚拟节点，将被重定位到节点`A`。

  - 你看，这样我们就解决了冷热不均的问题。

- **当节点数越多的时候，使用哈希算法时，需要迁移的数据就越多，使用一致性哈希时，需要迁移的数据就越少。**

  > - **使用一致哈希实现哈希寻址时，可以通过增加节点数降低节点宕机对整个集群的影响，以及故障恢复时需要迁移的数据量。**
  > - **后续在需要时，你可以通过增加节点数来提升系统的容灾能力和故障恢复效率。**

##### 一致性哈希总结

- 一致哈希是一种特殊的哈希算法，在使用一致哈希算法后，节点增减变化时只影响到部分数据的路由寻址，也就是说我们只要迁移部分数据，就能实现集群的稳定了。
- 当节点数较少时，可能会出现节点在哈希环上分布不均匀的情况。这样每个节点实际占据环上的区间大小不一，最终导致业务对节点的访问冷热不均。这个问题可以通过引入更多的虚拟节点来解决。
- 一致哈希本质上是一种路由寻址算法，适合简单的路由寻址场景。
  - 比如在`KV`存储系统内部，它的特点是简单，不需要维护路由信息。



#### 8. 拜占庭将军问题

在我看来，拜占庭将军问题（`The Byzantine Generals Problem`），它其实是借拜占庭将军的故事展现了分布式共识问题，还探讨和论证了解决的办法。而大多数人觉得它难理解，除了因为分布式共识问题比较复杂之外，还与莱斯利·兰伯特（`Leslie Lamport`）的讲述方式有关，他在一些细节上（比如，口信消息型拜占庭问题之解的算法过程上）没有说清楚。

实际上，它是分布式领域最复杂的一个容错模型，一旦搞懂它，你就能掌握分布式共识问题的解决思路，还能更深刻地理解常用的共识算法，在设计分布式系统的时候，也能根据场景特点选择适合的算法，或者设计适合的算法了。

那么接下来，我就以战国时期六国抗秦的故事为主线串联起整篇文章，让你读懂、学透。

**苏秦的困境**

战国时期，齐、楚、燕、韩、赵、魏、秦七雄并立，后来秦国的势力不断强大起来，成了东方六国的共同威胁。于是，这六个国家决定联合，全力抗秦，免得被秦国各个击破。一天，苏秦作为合纵长，挂六国相印，带着六国的军队叩关函谷，驻军在了秦国边境，为围攻秦国作准备。但是，因为各国军队分别驻扎在秦国边境的不同地方，所以军队之间只能通过信使互相联系，这时，苏秦面临了一个很严峻的问题：如何统一大家的作战计划？

万一一些诸侯国在暗通秦国，发送误导性的作战信息，怎么办？如果信使被敌人截杀，甚至被敌人间谍替换，又该怎么办？这些都会导致自己的作战计划被扰乱，然后出现有的诸侯国在进攻，有的诸侯国在撤退的情况，而这时，秦国一定会趁机出兵，把他们逐一击破的。

所以，**如何达成共识，制定统一的作战计划呢？**苏秦他很愁。

**这个故事，是拜占庭将军问题的一个简化表述，苏秦面临的就是典型的共识难题，也就是如何在可能有误导信息的情况下，采用合适的通讯机制，让多个将军达成共识，制定一致性的作战计划？**你可以先停下来想想，这个问题难在哪儿？我们又是否有办法，帮助诸侯国们达成共识呢？

**二忠一叛的难题**

为了便于你理解和层层深入，我先假设只有`3`个国家要攻打秦国，这三个国家的三位将军，咱们简单点儿，分别叫齐、楚、燕。同时，又因为秦国很强大，所以只有半数以上的将军参与进攻，才能击败敌人（注意，这里是假设哈，你别较真），在这个期间，将军们彼此之间需要通过信使传递消息，然后协商一致之后，才能在同一时间点发动进攻。

举个例子，有一天，这三位将军各自一脸严肃地讨论明天是进攻还是撤退，并让信使传递信息，按照“少数服从多数”的原则投票表决，两个人意见一致就可以了，比如：

1. 齐根据侦查情况决定撤退；
2. 楚和燕根据侦查信息，决定进攻。

那么按照原则，齐也会进攻。最终，`3`支军队同时进攻，大败秦军。

![image-20210328111908151](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328111908151.png)

可是，问题来了： 一旦有人在暗通秦国，就会出现作战计划不一致的情况。

比如齐向楚、燕分别发送了“撤退”的消息，燕向齐和楚发送了“进攻”的消息。撤退：进攻 =1:1，无论楚投进攻还是撤退，都会成为 2:1，这个时候还是会形成一个一致性的作战方案。

但是，楚这个叛徒在暗中配合秦国，让信使向齐发送了“撤退”，向燕发送了“进攻”，那么：燕看到的是，撤退：进攻 =1:2；齐看到的是，撤退：进攻 =2:1。

按照“少数服从多数”的原则，就会出现燕单独进攻秦军，当然，最后肯定是因为寡不敌众，被秦军给灭了。

![image-20210328112011723](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328112011723.png)

在这里，你可以看到，叛将楚通过发送误导信息，非常轻松地干扰了齐和燕的作战计划，导致这两位忠诚将军被秦军逐一击败。这就是所说的**二忠一叛难题**。 那么苏秦应该怎么解决这个问题呢？我们来帮苏秦出出主意。

**苏秦该怎么办？**

##### 解决办法一：口信消息型拜占庭问题之解

> **增加忠将数目，如果叛将人数为`m`，那么将军总人数必须不小于`3m + 1`。**

先来说说第一个解决办法。首先，三位将军都分拨一部分军队，由苏秦率领，苏秦参与作战计划讨论并执行作战指令。这样，`3`位将军的作战讨论，就变为了`4`位将军的作战讨论，这能够增加讨论中忠诚将军的数量。

然后呢，`4`位将军还约定了，如果没有收到命令，就执行预设的默认命令，比如“撤退”。除此之外，还约定一些流程来发送作战信息、执行作战指令，比如，进行两轮作战信息协商。为什么要执行两轮呢？先卖个关子，你一会儿就知道了。

**第一轮：**

- 先发送作战信息的将军作为指挥官，其他的将军作为副官。
- 指挥官将他的作战信息发送给每位副官。
- 每位副官，将从指挥官处收到的作战信息，作为他的作战指令；如果没有收到作战信息，将把默认的“撤退”作为作战指令。

**第二轮：**

- 除了第一轮的指挥官外，剩余的`3`位将军将分别作为指挥官，向另外`2`位将军发送作战信息。
- 然后，这`3`位将军按照“少数服从多数”，执行收到的作战指令。

为了帮助你直观地理解苏秦的整个解决方案，我来演示一下作战信息协商过程。而且，我会分别以忠诚将军和叛将先发送作战信息为例来演示， 这样可以完整地演示叛将对作战计划干扰破坏的可能性。

首先是`3`位忠诚的将军先发送作战信息的情况。

为了演示方便，假设苏秦先发起作战信息，作战指令是“进攻”。那么在第一轮作战信息协商中，苏秦向齐、楚、燕发送作战指令“进攻”。

![image-20210328112307109](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328112307109.png)

在第二轮作战信息协商中，齐、楚、燕分别作为指挥官，向另外`2`位发送作战信息“进攻”，因为楚已经叛变了，所以，为了干扰作战计划，他就对着干，发送“撤退”作战指令。

![image-20210328112322231](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328112322231.png)

最终，齐和燕收到的作战信息都是“进攻、进攻、撤退”，按照原则，齐和燕与苏秦一起执行作战指令“进攻”，实现了作战计划的一致性，保证了作战的胜利。

那么，如果是叛徒楚先发送作战信息，干扰作战计划，结果会有所不同么？我们来具体看一看。在第一轮作战信息协商中，楚向苏秦发送作战指令“进攻”，向齐、燕发送作战指令“撤退”。

![image-20210328112411308](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328112411308.png)

然后，在第二轮作战信息协商中，苏秦、齐、燕分别作为指挥官，向另外两位发送作战信息。

![image-20210328112531810](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328112531810.png)

最终，苏秦、齐和燕收到的作战信息都是“撤退、撤退、进攻”，按照原则，苏秦、齐和燕一起执行作战指令“撤退”，实现了作战计划的一致性。也就是说，无论叛将楚如何捣乱，苏秦、齐和燕，都执行一致的作战计划，保证作战的胜利。

这个解决办法，其实是兰伯特在论文`《The Byzantine Generals Problem》`中提到的***口信消息型拜占庭问题之解：如果叛将人数为`m`，将军人数不能少于`3m + 1` ，那么拜占庭将军问题就能解决了。 ***不过，作者在论文中没有讲清楚一些细节，为了帮助你阅读和理解论文，在这里我补充一点：

**这个算法有个前提，也就是叛将人数`m`，或者说能容忍的叛将数`m`，是已知的。**在这个算法中，叛将数`m`决定递归循环的次数（也就是说，叛将数`m`决定将军们要进行多少轮作战信息协商），即`m+1`轮（所以，你看，只有楚是叛变的，那么就进行了两轮）。你也可以从另外一个角度理解：**`n`位将军，最多能容忍`(n - 1) / 3`位叛将。**关于这个公式，你只需要记住就好了，推导过程你可以参考论文。

不过，**这个算法虽然能解决拜占庭将军问题，但它有一个限制：如果叛将人数为`m`，那么将军总人数必须不小于`3m + 1`。**

在二忠一叛的问题中，在存在`1`位叛将的情况下，必须增加`1`位将军，将`3`位将军协商共识，转换为`4`位将军协商共识，这样才能实现忠诚将军的一致性作战计划。那么有没有办法，在不增加将军人数的时候，直接解决二忠一叛的难题呢？

##### 解决办法二：签名消息型拜占庭问题之解

> **能够容忍任意数量的叛徒，通过消息的签名，约束了叛徒的作恶行为，也就是说，任何篡改和伪造忠将的消息的行为，都会被发现。**

其实，苏秦还可以通过签名的方式，在不增加将军人数的情况下，解决二忠一叛的难题。首先，苏秦要通过印章、虎符等信物，实现这样几个特性：

- 忠诚将军的签名无法伪造，而且对他签名消息的内容进行任何更改都会被发现；
- 任何人都能验证将军签名的真伪。

这时，如果忠诚的将军，比如齐先发起作战信息协商，一旦叛将小楚修改或伪造收到的作战信息，那么燕在接收到楚的作战信息的时候，会发现齐的作战信息被修改，楚已叛变，这时他将忽略来自楚的作战信息，最终执行齐发送的作战信息。

![image-20210328112837332](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328112837332.png)

**如果叛变将军楚先发送误导的作战信息，那么，齐和燕将按照一定规则（比如取中间的指令）在排序后的所有已接收到的指令中（比如撤退、进攻）中选取一个指令，进行执行，最终执行一致的作战计划。**

![image-20210328112918083](https://raw.githubusercontent.com/MachineGunLin/markdown_pics/master/img/image-20210328112918083.png)

这个解决办法，是兰伯特在论文中提到的签名消息型拜占庭问题之解。而通过签名机制约束叛将的叛变行为，任何叛变行为都会被发现，也就会实现无论有多少忠诚的将军和多少叛将，忠诚的将军们总能达成一致的作战计划。

##### 拜占庭将军问题类比计算机分布式场景

- 故事里的各位将军，你可以理解为计算机节点。
- 忠诚的将军，你可以理解为正常运行的计算机节点。
- 叛变的将军，你可以理解为出现故障并会发送误导信息的计算机节点。
- 信使被杀，可以理解为通讯故障、信息丢失。
- 信使被间谍替换，可以理解为通讯被中间人攻击，攻击者在恶意伪造信息和劫持通讯。

那么我想强调的是，拜占庭将军问题描述的是最困难的，也是最复杂的一种分布式故障场景，除了存在故障行为，还存在恶意行为的一个场景。你要注意，在存在恶意节点行为的场景中（比如在数字货币的区块链技术中），必须使用拜占庭容错算法（`Byzantine Fault Tolerance`，`BFT`）。除了故事中提到两种算法，常用的拜占庭容错算法还有：`PBFT`算法，`PoW`算法。

而在计算机分布式系统中，最常用的是非拜占庭容错算法，即故障容错算法（`Crash Fault Tolerance`，`CFT`）。`CFT`解决的是分布式的系统中存在故障，但不存在恶意节点的场景下的共识问题。 也就是说，这个场景可能会丢失消息，或者有消息重复，但不存在错误消息，或者伪造消息的情况。常见的算法有`Paxos`算法、`Raft`算法、`ZAB`协议。

那么，**如何在实际场景选择合适的算法类型呢？答案是：如果能确定该环境中各节点是可信赖的，不存在篡改消息或者伪造消息等恶意行为（例如`DevOps`环境中的分布式路由寻址系统），推荐使用非拜占庭容错算法。反之，推荐使用拜占庭容错算法，例如在区块链中使用`PoW`算法。**


