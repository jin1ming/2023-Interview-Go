# RPC

[TOC]

## 基本

### 什么是RPC

RPC即远程过程调用，它帮助我们屏蔽网络编程细节，实现调用远程方法就跟调用本地一样的体验，不需要编写太多与业务无关的代码。

其**作用**体现在两个方面：

- 屏蔽远程调用和本地调用的区别
- 隐藏底层网络通信的复杂性

一个完整的RPC**包含**：

网络传输（TCP/HTTP2）+序列化/反序列化+方法调用+动态代理

**动态代理**：RPC根据调用的服务接口提前生成动态代理实现类，并通过依赖注入等技术注入到声明了该接口的相关业务逻辑里面。该代理实现类会拦截所有的方法调用，在提供的方法处理逻辑里面完成一整套的远程调用，并发远程调用结果返回给调用方。

### 为什么微服务需要RPC

我们使用微服务化的一个好处就是，不限定服务的提供方使用什么技术选型，能够实现公司跨团队的技术解耦。

这样的话，如果没有统一的服务框架，RPC框架，各个团队的服务提供方就需要各自实现一套序列化、反序列化、网络框架、连接池、收发线程、超时处理、状态机等“业务之外”的重复技术劳动，造成整体的低效。

所以，统一RPC框架把上述“业务之外”的技术劳动统一处理，是服务化首要解决的问题。

## net/rpc

- 成员方法：只能两个参数（传入参数、&传出参数）

  成员方法必须公开（大写），必须返回一个error类型；

  当调用远程函数之后，如果返回的错误不为空，那么传出参数为空。

- 编码：默认gob，go独有，故不支持跨语言

  可选json(net/rpc/jsonrpc)来支持跨语言（不支持HTTP）

  - 服务端：

    `rpc.ServeCodec`或`jsonrpc.ServeConn`函数替代`rpc.ServeConn`

    ``rpc.ServeCodec(jsonrpc.NewServerCodec(conn))`

    `jsonrpc.ServeConn(conn)`

  - 客户端：

    `rpc.NewClientWithCodec`代替`rpc.NewClient`

    ``rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))`

    或

    `conn, err := jsonrpc.Dial("tcp", "127.0.0.1:8888")`

- 支持tcp/http，通过字符串的方式进行调用

  ```golang
  // TCP服务端
  listener, err := net.Listen("tcp", ":8888")
  defer listener.Close()
  conn, err := listener.Accept()
  rpc.RegisterName("zabbix", new(Zabbix))
  rpc.ServeConn(conn)
  // 客户端
  conn, err := net.Dial("tcp", "127.0.0.1:8888")
  defer conn.Close()
  client := rpc.NewClient(conn)
  err = client.Call("zabbix.MonitorHosts", "Nginx", &data)
  ```

  ```golang
  // HTTP服务端
  rpc.Register(new(Rect))
  rpc.HandleHTTP()
  err := http.ListenAndServe(":8080", nil)
  // 客户端
  client, err := rpc.DialHTTP("tcp", "127.0.0.1:8080")
  err = client.Call("Rect.Area", Params{30, 40}, &ret)
  ```

### grpc

GRPC是Google公司基于Protobuf开发的跨语言的开源RPC框架。

GRPC是一个高性能、开源和通用的 RPC 框架，面向移动和 HTTP/2 设计。目前提供C,Java和Go语言版本.

GRPC基于HTTP/2标准设计，带来诸如双向流、流控、头部压缩、单 TCP 连接上的多复用请求等特。这些特性使得其在移动设备上表现更好，更省电和节省空间占用。

- 成员方法：第一个参数是上下文，默认必填

  上下文可用来控制超时时间

  ```golang
  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  ```

- protobuf定义接口

  ```protobuf
  ////protobuf默认支持的版本是2.0,现在一般使用3.0版本,所以需要手动指定版本号
  //c语言的编程风格
  syntax = "proto3";
  //指定包名
  package pb;
  //定义传输数据的格式
  message People{
      string name = 1; //1表示表示字段是1   数据库中表的主键id等于1,主键不能重复,标示位数据不能重复
      //标示位不能使用19000 -19999  系统预留位
      int32 age = 2;
      //结构体嵌套
      student s = 3;
      //使用数组/切片
      repeated string phone = 4;
      //oneof的作用是多选一,c语言中的联合体
      oneof data{
          int32 score = 5;
          string city = 6;
          bool good = 7;
      }
  }
  
  message student{
      string name = 1;
      int32 age = 6;
  }
  //通过定义服务,然后借助框架,帮助实现部分的rpc代码
  service Hello{
      rpc World(student)returns(student);
  }
  ```

- 仅支持tcp，通过方法的方式进行调用

  ```golang
  // 服务端
  // 先定义一个结构体，继承protobuf生成的接口
  //先获取grpc对象
  grpcServer := grpc.NewServer()
  //注册服务
  pb.RegisterHelloServer(grpcServer,new(HelloService))
  //开启监听
  listener,err := net.Listen("tcp",":8888")
  defer listener.Close()
  //先获取grpc服务端对象
  grpcServer.Serve(listener)
  
  //客户端
  grpcCnn ,err := grpc.Dial("127.0.0.1:8888",grpc.WithInsecure())
  defer grpcCnn.Close()
  client := pb.NewHelloClient(grpcCnn)
  resp,err := client.World(context.TODO(),&s)
  ```

- [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)提供HTTP支持

- [go-grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware)提供中间件支持

- 4种调用方式

## Q&A

- 什么场景下适合使用RPC

  系统内部使用rpc，外部使用RESTful。需要网络环境稳定，以及设置合理的超时时间以及重试次数。