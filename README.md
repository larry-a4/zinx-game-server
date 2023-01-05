# zinx框架

## ZinxV0.1-基础的server

#### 方法
````
启动服务器
基本的服务器开发 1创建addr 2创建listener 3处理客户端基本业务（回显）

停止服务器
做一些资源的回收和状态的回执

运行服务器
调用Start()方法，调用之后做阻塞处理，之后可以扩展

初始化
````
#### 属性
````
名称
监听的IP
监听的端口
````

## ZinxV0.2-简单的链接封装和业务绑定 

### 链接模块

#### 方法
````
启动链接Start()
停止链接Stop()
获取当前链接的conn对象GetTCPConnection()
得到链接ID - GetConnID()
得到客户端链接的地址和端口RemoteAddr()
发送数据Send()
链接所绑定的处理业务的函数类型
````
#### 属性
````
socket TCP套接字: *net.TCPConn
链接的ID: uint32
链接的状态（是否已经关闭): bool
与当前链接所绑定的处理业务方法: ziface.HandleFunc
等待链接被动退出的channel: chan bool
````

## ZinxV0.3-基础router模块 

### Request请求封装 - 绑定链接和数据

#### 方法
````
得到当前链接: ziface.IConnection
得到当前数据: []byte
新建一个request请求(过于简单，没必要实现)
````
#### 属性
````
链接Iconnection
请求数据
````

### Router模块

#### 抽象的IRouter
````
处理业务之前的方法 PreHandle(IRequest)
处理业务的主方法 Handle(IRequest)
处理业务之后的方法 PostHandle(IRequest)
````
#### 具体的BaseRouter
````
处理业务之前的方法
处理业务的主方法
处理业务之后的方法
````

### zinx集成router模块
````
IServer增添路由添加功能 - AddRouter(router IRouter)
Server增添Router成员
Connection类绑定一个Router成员
在Connection调用，已经注册的Router处理业务
````

## 测试
````
GO111MODULE=off go run server.go
GO111MODULE=off go run client.go
````
