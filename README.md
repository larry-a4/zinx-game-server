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

## ZinxV0.4-全局配置

### 服务器应用/conf/zinx.json(用户填写)

#### 步骤
````
创建zinx全局配置模块 utils/globalobj.go
提供一个全局的GlobalObject对象
init时读取用户配置文件，写入globalobject对象
将zinx框架中的硬代码，用globalobject中的参数替换
````

## ZinxV0.5-消息封装

### 定义消息结构Message

#### 属性
````
消息ID
消息长度
消息内容
````

### 解决TCP粘包问题的封拆包模块

#### 针对Message进行TLV格式封装 Pack(IMessage) ([]byte, error)
````
写message的长度
写message的ID
写message的内容
````
#### 针对Message进行TLV格式拆解 Unpack([]byte) (IMessage, error)
````
先读取固定长度
再根据长度，读取内容
````
#### 将消息封装机制集成到Zinx框架中
````
将Message添加到Request属性中
修改链接读取数据的机制：拆包并按照TLV形式读取
提供发包机制：将数据打包，再发送
````

## ZinxV0.6-多路由模式

### 消息管理模块（支持多路由）

#### 属性
````
消息ID与对应router的关系-map
````
#### 方法
````
根据msgID来索引调度路由方法--DoMsgHandler(IRequest)
添加路由方法到map集合中--AddRouter(uint32, IRequest)
````

### 将消息管理模块集成到zinx框架中
````
1-将Server模块的router替换成MsgHandler
2-修改AddRouter
3-将Connection中的router替换成MsgHandler
4-将Connection的之前调度router的业务替换成MsgHandler调度，修改StartReader
````

## 测试
````
GO111MODULE=off go run server.go
GO111MODULE=off go run client.go
````
