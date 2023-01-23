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

## ZinxV0.7-读写协程分离
````
1-添加Reader和Writer之间通信的channel
2-添加writer goroutine
3-Reader由直接发送给client，改成发送给通信channel
4-启动Reader和Writer一同工作
````

## ZinxV0.8-消息队列/多任务处理机制

### 1-创建消息队列（msgHandler模块）
#### 属性
````
TaskQueue []chan IRequest
WorkerPoolSize uint32 
````
### 2-创建多任务worker的工作池并且启动
#### 开启一个worker pool - StartWorkerPool()
````
1-创建pool size个worker，每个worker开一个go - StartOneWorker(workerID, taskQueue)
2-阻塞等待当前worker对应channel的消息
3-一旦有消息到来，对应的worker处理消息，调用DoMsgHandler()
````
### 3-将发送消息，改为 把消息发送给消息队列和worker池 处理
#### 定义一个方法，将消息发送给消息队列worker pool
````
1-保证每个worker收到的request均衡，也就是发给对应worker的taskQueue
2-将消息直接发送到对应的channel
````
### 4-集成到zinx框架
````
1-开启并调用worker pool（必须保证只有一个，最好在server模块init时开启）
2-将从client收到的消息，发送给worker pool来处理
````

## ZinxV0.9-链接管理模块
### 1-创建链接管理模块 ConnManager
#### 属性
````
已经创建的Connection - map[uint32]IConnection
针对map的互斥锁 - connLock sync.RWMutex
````
#### 方法
````
添加链接 - Add(IConnection)
删除链接 - Remove(IConnection)
查找链接 - Get(connID uint32) (IConnection, error)
总链接数 - Len() int
清理全部链接 - ClearConn()
````
### 2-集成到zinx框架中
````
将conn Manager加入Server模块中：添加ConnMgr属性，初始化ConnMgr，server停止时ClearConn
每次成功与client建立链接时，添加链接：NewConnection时将conn加入ConnMgr
判断当前的链接数量是否已经超出最大值MaxConn
每次与client断开链接时，删除链接：在Conn.Stop时，从ConnMgr中移除conn
````
### 3-提供业务所需要的hook
#### 属性
````
创建链接后hook - OnConnStart(hookFunc func (IConnection)
销毁链接前hook - OnConnStop(hookFunc func (IConnection))
````
#### 方法
````
注册OnConnStart 钩子的方法
注册OnConnStop 钩子的方法
调用OnConnStart 钩子的方法
调用OnConnStop 钩子的方法
````

## ZinxV1.0-链接属性配置
### 给Connection模块添加可配置属性
#### 属性
````
链接属性集合 - map[string]interface{}
保护链接属性的锁 - sync.RWMutex
````
#### 方法
````
设置链接属性 - SetProperty(key, value)
获取链接属性 - GetProperty(key)interface{}
移除链接属性 - RemoveProperty(key)
````

## Zinx应用-MMO多人在线网游

### AOI目标范围算法
#### 格子属性
````
格子ID - gID
格子左边界坐标
格子右边界坐标
格子上边界坐标
格子下边界坐标
格子给玩家/物体的ID集合 - map
保护当前集合的锁
````
#### 格子方法
````
初始化 - NewGrid(gID, minX, maxX, minY, maxY) *Grid
添加一个玩家/物体 - Add(playerID int)
删除一个玩家/物体 - Remove(playerID int)
得到当前格子中所有玩家/物体 - GetPlayerIDs() (playerIDs []int)
调试使用-打印出格子的基本信息 - String()
````
#### AOI地图管理属性
````
区域左边界坐标
区域右边界坐标
X方向的格子数量 - countX
区域上边界坐标
区域下边界坐标
Y方向的格子数量 - countY
当前区域中有哪些格子 - map[gID] *Grid
````
#### AOI地图管理方法
````
初始化 - NewAOIManager(minX, maxX, countX, minY, maxY, countY int) *AOIManager
调试使用-打印当前AOI模块 - String()
获取周边九宫格信息
添加playerID到格子
移除playerID从格子
获取一个格子中全部playerID
通过坐标将Player添加到一个格子
通过坐标将Player从一个格子中移除
通过坐标获取周边九宫格内全部playerID
通过坐标获取玩家所在的gID
````

### Protobuf传输协议

## 游戏业务
### 协议定义
````
msgID:1 - SyncPid { int32 Pid=1 } 同步玩家本次登录的ID，登录时由Server主动生成发送给Client
msgID:2 - Talk { string Content=1 } 由Client发起，聊天信息为Content
msgID:3 - Position { float X=1; float Y=2; float Z=3; float V=4 } 由Client发起，移动的坐标
msgID:200 - Broadcast {int32 Pid=1; int32 Tp=2; oneof Data} 由Server发起，Tp:1-世界聊天，2-坐标，3-动作，4-坐标信息更新
    Data { string Content=3; Position P=4; int32 ActionData=5 }
msgID:201 - SyncPid {int32 Pid=1 } 由Server发起，广播玩家掉线
msgID:202 - SyncPlayers { repeated Player ps=1 } 由Server发起，同步周围人的位置信息（包括自己）
    Player { int32 Pid=1; Position P2 }
````

### 项目构建
````
api - 存放用户自定义的路由业务，一个msgID对应一个业务
conf - zinx.json 存放zinx配置文件
pb - msg.proto 原始protobuf定义文件
     build.sh 编译 msg.proto 的脚本
     msg.pb.go 编译生成的go文件（只读）
core - 存放游戏核心功能
main.go - 服务器主入口
````

### 玩家上线
#### 先定义proto协议，生成对应的pb.go文件
#### 玩家模块属性
````
玩家ID
链接信息（用于和对应客户端通信的connection
玩家当前坐标
````
#### 玩家模块方法
````
创建玩家的方法 NewPlayer(conn) *Player
玩家和客户端通信的方法 SendMsg(msgId, proto.Message)
````

#### 实现上线业务功能
````
给server注册一个创建链接后的hook：给客户端发送 msgID=1 和 msgID=200
给player提供方法：1-同步PlayerID给客户端
给player提供方法：2-同步上线位置给客户端
````
#### 测试上线功能
````
GO111MODULE=off go run main.go
````

## 测试
````
GO111MODULE=off go run server.go
GO111MODULE=off go run client.go
````
