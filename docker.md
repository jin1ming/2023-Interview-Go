[TOC]

## Dockerfile

- 结构
  - 基础镜像信息、维护者信息、镜像操作指令、容器启动时执行指令

- 指令：

  FROM、MAINTAINER、RUN、CMD、EXPOSE、ENV、ADD、COPY、ENTRYPOINT、VOLUME、USER、WORKDIR、ONBUILD

  MAINTAINER：维护者信息

  RUN：在镜像基础上执行命令，提交为新的镜像

  CMD：容器启动时执行的命令，每个Dockerfile只能有一个，如果存在多个只执行最后一个

  EXPOSE：告诉服务器暴露的端口，供外部连接使用

  ENV：指定环境变量，后续可以被RUN使用，以及容器运行起来后使用

  ADD：拷贝来自远程URL或本地的文件、目录，支持通配符，自动解压压缩文件

  COPY：拷贝宿主机上的文件或目录

  ENTRYPOINT：用于配置容器启动后执行的命令，不能被docker run提供的参数覆盖，最后一个生效

  VOLUME：创建在本地主机或其他容易可以挂载的数据卷

  USER：指定容器运行的用户名或UID，在之前可以先用RUN创建需要的用户

  WORKDIR：为RUN CMD ENTRYPOINT指定工作目录，可以使用多个

  ONBUILD：命令出现的Dockerfile里不会执行指令，但被其他Dockerfile所引用的时候会执行

  

- 镜像分层构建

  - Dockerfile中每个指令都会创建一个新的镜像层
  - 镜像层将会被缓存和复用
  - 当Dockerfile的指令修改了，复制的文件变化了，或者构建镜像时指定的变量不同了，对应的缓存就会失效
  - 某一层的镜像缓存失效后，它之后的镜像层缓存就会都失效
  - 镜像层是不可变的，即便下一层删除某个文件，其镜像中仍然会包含该文件
  
  

## 原理

### Namespce

主要用来做**资源隔离**。Linux内核共实现了以下几种Namespce：

- **UTS**：表示不同的 namespace 可以配置不同的 hostname
- **User**：表示不同的 namespace 可以配置不同的用户和组
- **Mount**：表示不同的 namespace 的文件系统挂载点是隔离的
- **PID**：表示不同的 namespace 有完全独立的 pid
- **Network**：表示不同的 namespace 有独立的网络协议栈
- **IPC**：表示不同的namespace有不同的IPC对象

查询Docker容器的Namespace信息如下：

```bash
# 获取容器ID
docker ps

# 获取容器对应的进程PID
docker inspect f604f0e34bc2

# 查询PID对应的ns信息
ls -l /proc/58212/ns 
lrwxrwxrwx 1 root root 0 Jul 16 19:19 ipc -> ipc:[4026532278]
lrwxrwxrwx 1 root root 0 Jul 16 19:19 mnt -> mnt:[4026532276]
lrwxrwxrwx 1 root root 0 Jul 16 01:43 net -> net:[4026532281]
lrwxrwxrwx 1 root root 0 Jul 16 19:19 pid -> pid:[4026532279]
lrwxrwxrwx 1 root root 0 Jul 16 19:19 user -> user:[4026531837]
lrwxrwxrwx 1 root root 0 Jul 16 19:19 uts -> uts:[4026532277]
```

操作Namespace的方式如下：

```
# 进入指定的namespace
nsenter --target 58212 --mount --uts --ipc --net --pid -- env --ignore-environment -- /bin/bash

# 离开当前ns，并加入新的ns
unshare --mount --ipc --pid --net --mount-proc=/proc --fork /bin/bash

# 创建子进程，并将子进程放到新的ns中。父进程不变
# arg可选：CLONE_NEWUTS、CLONE_NEWUSER、CLONE_NEWNS、CLONE_NEWPID。CLONE_NEWNET
int clone(int (*fn)(void *), void *child_stack, int flags, void *arg);

# 将当前进程放到已有的ns中
# nstype 用来指定 namespace 的类型，可以设置为 CLONE_NEWUTS、CLONE_NEWUSER、CLONE_NEWNS、CLONE_NEWPID 和 CLONE_NEWNET
int setns(int fd, int nstype);

# 使当前进程退出当前的 namespace，并加入到新创建的namespace中
int unshare(int flags);
```

### CGroups

全称Control Group，主要用来**限制资源使用**。CGroups定义了下面一系列子系统，每个子系统用于控制某一类资源：

- cpu：限制进程的 cpu 使用率
- cpuacct：统计 cgroups 中的进程的 cpu 使用报告
- cpuset：为 cgroups 中的进程分配单独的 cpu 节点或者内存节点
- memory：限制进程的 memory 使用量
- blkio：限制进程的块设备 io
- devices：控制进程能够访问某些设备
- net_cls：标记 cgroups 中进程的网络数据包，然后可以使用 tc 模块（traffic control）对数据包进行控制。
- freezer：可以挂起或者恢复 cgroups 中的进程

在 Linux 上，为了操作 Cgroup，有一个专门的 Cgroup 文件系统，位于/sys/fs/cgroup/目录下。目录结构如下所示：

```
drwxr-xr-x 5 root root  0 May 30 17:00 blkio
lrwxrwxrwx 1 root root 11 May 30 17:00 cpu -> cpu,cpuacct
lrwxrwxrwx 1 root root 11 May 30 17:00 cpuacct -> cpu,cpuacct
drwxr-xr-x 5 root root  0 May 30 17:00 cpu,cpuacct
drwxr-xr-x 3 root root  0 May 30 17:00 cpuset
drwxr-xr-x 5 root root  0 May 30 17:00 devices
drwxr-xr-x 3 root root  0 May 30 17:00 freezer
drwxr-xr-x 3 root root  0 May 30 17:00 hugetlb
drwxr-xr-x 5 root root  0 May 30 17:00 memory
lrwxrwxrwx 1 root root 16 May 30 17:00 net_cls -> net_cls,net_prio
drwxr-xr-x 3 root root  0 May 30 17:00 net_cls,net_prio
lrwxrwxrwx 1 root root 16 May 30 17:00 net_prio -> net_cls,net_prio
drwxr-xr-x 3 root root  0 May 30 17:00 perf_event
drwxr-xr-x 5 root root  0 May 30 17:00 pids
drwxr-xr-x 5 root root  0 May 30 17:00 systemd
```

Docker操作CGroups的方式如下图所示：

![下载 (assets/下载 (4).png)](assets/下载 (4).png)

## Q&A

- 如何借助宿主机上工具排查问题？考察命名空间

  kubectl debug、nsenter

- docker镜像的构建方法

  基于已有镜像创建：对容器进行docker commit

  基于本地模板创建：导入操作系统模板

  基于Dockerfile

- 容器间通信：

  none（容器网络堆栈）、host（主机网络堆栈）、default bridge（IP地址链接）、自定义网桥

  容器IP、宿主机IP、link、User-defined networks（桥接网络）

- 容器与宿主机一个网段该怎么实现？

  macvlan的原理是在宿主机物理网卡上虚拟出多个子网卡，通过不同的MAC地址在数据链路层进行网络数据转发的，它是比较新的网络虚拟化技术，需要较新的内核支持（Linux kernel v3.9–3.19 and 4.0+）