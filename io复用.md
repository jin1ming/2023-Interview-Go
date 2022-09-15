目前Linux系统中提供了5种IO处理模型

1. **阻塞IO**：发起IO操作，等待成功或失败。

   优点：维护简单

   缺点：每个连接都需要一个单独的进程、线程进行处理。当并发量大时，进线程上下文切换的开销过大

2. **非阻塞IO**：阻塞时返回错误/轮询请求。

   优点：非阻塞，实时性比较好

   缺点：需要不断地轮询内核接口，这将占用大量的 CPU 时间

3. **IO多路复用**：select、poll、epoll。对一组文件描述符进行相关事件的注册，然后阻塞等待某些事件的发生或等待超时。

   - select：数组存储fd，默认1024个；每次调用需要将fd集合从用户态拷贝到内核态；内核轮询fd，效率低。

   - poll：与select不同的是链表存储fd，没有长度限制。

   - epoll：红黑树管理fd，并且会对所有添加到epoll中的fd设置事件回调。当相应的事件发生时，会将发生的事件添加到内核的rdlist双向链表中。

     工作模式

     - **LT**（水平触发，Level Triggered ）：默认的工作模式。只要这个 fd 还有数据可读，每次 epoll_wait 都会返回它的事件，提醒用户程序去操作。可能造成事件过多。
     - **ET**（边缘触发，Edge Triggered）：只会提示一次，直到下次再有数据流入之前都不会再提示了，无论fd中是否还有数据可读。所以在ET模式下，read一个 fd 的时候一定要把它的buffer读完，或者遇到EAGIN错误

4. **信号驱动IO**：利用信号机制，让内核告知应用程序文件描述符的相关事件。

5. **异步IO**：相比信号驱动IO需要在程序中完成数据从用户态到内核态(或反方向)的拷贝，异步IO可以把拷贝这一步也帮我们完成之后才通知应用程序。我们使用 aio_read 来读，aio_write 写。

#### Reactor

Reactor是一种事件驱动框架，**基于同步IO**，有以下5个关键的参与者：

- 文件描述符
- 同步事件多路分离器（event demultiplexer）：使用select/poll/epoll封装，运行在主线程的事件循环上。
- 事件处理器（event handler）：事件处理方法，内嵌具体的业务逻辑，运行在任务线程上
- Reactor 管理器（reactor）：定义一组针对reactor的fd curd接口

经典用法：**主线程采用epoll管理事件，任务线程池处理具体的任务。**

#### Proactor

Proactor同样是一种事件驱动框架，**基于异步IO**。成功从IO设备读取数据后，才会通知事件分离器，然后再由事件分离器分发并处理事件。

#### epoll

- 接口

  ```c++
  #include <sys/epoll.h>
  
  // 数据结构
  // 每一个epoll对象都有一个独立的eventpoll结构体
  // 用于存放通过epoll_ctl方法向epoll对象中添加进来的事件
  // epoll_wait检查是否有事件发生时，只需要检查eventpoll对象中的rdlist双链表中是否有epitem元素即可
  struct eventpoll {
      /*红黑树的根节点，这颗树中存储着所有添加到epoll中的需要监控的事件*/
      struct rb_root  rbr;
      /*双链表中则存放着将要通过epoll_wait返回给用户的满足条件的事件*/
      struct list_head rdlist;
  };
  
  // API
  // 内核中间加一个 ep 对象，把所有需要监听的 socket 都放到 ep 对象中
  int epoll_create(int size); 
  // epoll_ctl 负责把 socket 增加、删除到内核红黑树
  int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event); 
  // epoll_wait 负责检测可读队列，没有可读 socket 则阻塞进程
  int epoll_wait(int epfd, struct epoll_event * events, int maxevents, int timeout);
  ```

- 使用用例

  ```c++
  int main(int argc, char* argv[])
  {
     /*
     * 在这里进行一些初始化的操作，
     * 比如初始化数据和socket等。
     */
  
      // 内核中创建ep对象
      epfd=epoll_create(256);
      // 需要监听的socket放到ep中
      epoll_ctl(epfd,EPOLL_CTL_ADD,listenfd,&ev);
   
      while(1) {
        // 阻塞获取
        nfds = epoll_wait(epfd,events,20,0);
        for(i=0;i<nfds;++i) {
            if(events[i].data.fd==listenfd) {
                // 这里处理accept事件
                connfd = accept(listenfd);
                // 接收新连接写到内核对象中
                epoll_ctl(epfd,EPOLL_CTL_ADD,connfd,&ev);
            } else if (events[i].events&EPOLLIN) {
                // 这里处理read事件
                read(sockfd, BUF, MAXLINE);
                //读完后准备写
                epoll_ctl(epfd,EPOLL_CTL_MOD,sockfd,&ev);
            } else if(events[i].events&EPOLLOUT) {
                // 这里处理write事件
                write(sockfd, BUF, n);
                //写完后准备读
                epoll_ctl(epfd,EPOLL_CTL_MOD,sockfd,&ev);
            }
        }
      }
      return 0;
  }
  ```
