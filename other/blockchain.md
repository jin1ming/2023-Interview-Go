## 软分叉和硬分叉

- 硬分叉是指当区块链代码发生改变后 **旧节点拒绝接受由新节点创建的区块**，不符合原规则的区块将会被忽略，矿工会按照原规则在他们最后验证的区块之后创建新的区块，区块链领域最有名的硬分叉案例，便是“以太坊”分叉。
- 软分叉是指区块链代码发生改变后，**旧的节点并不会感知到区块链代码发生改变，并继续接受由新节点创建的区块**，矿工们可能会在他们完全没有理解或验证过的区块上进行工作，软分叉新旧节点双方始终都工作在同一条链上。

## 介绍下DID，和传统方式有什么区别？

- 传统标识 vs 分布式标识

  - 数据集中托管-数据自主存储
  - 标识统一分配-标识自生成、自分配、自管理
  - 权威机构背书-加密算法背书
  - 数据泄露风险-无数据泄露风险

- 解决方案及流程
  1、在身份证明机构、数据持有机构、数据使用机构间搭建区块链网络，机构作为节点接入并注册DID

  2、由身份证明机构为用户生成DID

  3、用户授权数据使用机构使用自己的数据

  4、数据使用机构生成授权凭证（Credential），标明授权对象、数据属主、有效期、授权内容等属性，并使用机构私钥进行签名；数据使用机构可选择将授权凭证生成摘要并写入区块链，达到增信目的

  5、数据使用机构出示授权凭证给数据持有机构

  6、数据持有机构通过凭证验证（Verify）接口，验证授权凭证

  7、验证通过，数据持有机构将数据发送给数据使用机构

- 默克尔树中怎么查找数据是否存在？

  问默克尔树可以扯git

  0. 从网络上获取并保存最长链的所有block header至本地；
  1. 计算该交易的hash值tx_hash；
  2. 定位到包含该tx_hash所在的区块，验证block header是否包含在已知的最长链中；
  3. 从区块中获取构建merkle tree所需的hash值；
  4. 根据这些hash值计算merkle_root_hash；
  5. 若计算结果与block header中的merkle_root_hash相等，则交易真实存在。
  6. 根据该block header所处的位置，确定该交易已经得到多少个确认。

- 在信通院实习时，他们用的什么共识算法

  pos权益证明

## 怎么完成共识？

- PoW（Proof of Work）工作证明

  PoW是最早出现的共识机制，被比特币系统和以太坊（2.0版本之下）所采用，其本质是牺牲运算时间、通过竞争计算来选举出记账节点，最终达成共识、解决一致性问题的过程。在每一次达成共识的过程中，所有节点都需要参与进来，完成一个计算困难但是验证容易的算法，第一个计算出结果且经过验证的节点将会获得记账权，并获得奖励。如比特币系统，每个节点中发起了一笔交易时，会将这笔交易广播到整个网络中，网络中所有节点将会把收到的所有交易进行打包，然后进行哈希计算存储到Merkle树中，然后封装成区块头。此时，每个节点将使用一个不断变化的nonce值，与区块头的80字节数据不断进行哈希计算。倘若哈希值小于目标值，即运算结束，PoW完成。

- PoS（Proof of Stake）权益证明

  PoS共识机制用投票来代替竞争获得记账权，每个节点的影响力来自于它的持股比例，当收到51%以上的节点的投票时才会获得记账权。为了让投票更加高效、及时，也可以让获得票数最多的一部分节点按照顺序轮流获得记账权，产生区块。与PoW相比，PoS不需要消耗大量的运算时间，随之也减少了电能的消耗。但是正是因为不需要这些运算时间的消耗，没有了挖矿带来的奖励的激励，参与区块链网络的用户将会变少，随之而来的是每个节点承受更大的压力，代币的发行会收到阻碍。

- PBFT 实用拜占庭容错机制

  PBFT算法是原始的拜占庭算法所改良而来的，在原基础上大大增加了效率，得以真正应用在实际的系统中。对于公有链来说，任何人都可以自由加入这个区块链网络，倘若有的人不怀好意，发出错误信息，就会破坏整个系统的一致性。PBFT就是针对这类问题出现的，这类具有破坏能力的节点我们暂且称作作恶节点。定义N为节点总数，F为作恶节点数量，PBFT的核心方案便是：N≥3F+1下可以保证一致性。即在保证安全性、可用性的同时，只要作恶节点不超过(N-1)/3，系统便能够正常运行。

- Raft

  Raft类似于PBFT通过选举leader来达成共识，不同的是Raft中有三个角色：leader（领袖）、follower（追随者）、candidate（候选人）。正常情况下，只会有一个leader节点，剩余的都是follower节点。当Raft初始化或者leader节点出现故障时，便会发动选举（基于心跳机制检测）。当其它节点检测到leader的心跳超时，follower节点便会将term counter（任期编号）加一，给自己投票，再向其它节点拉票（此时成为了candidate）。每个节点每个任期只有一票，且投给最早拉票的节点。当candidate节点收到其它candidate节点的拉票，且term counter不小于自己的，就会退出竞选，重新成为follower身份，选举这个candidate为领袖。当某个candidate节点收到了超过一半的票便晋升为leader，若一定时间没选举出来，便终止这一轮重新开始。
