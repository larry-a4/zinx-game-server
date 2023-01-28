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

### 游戏业务协议定义
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

### 世界聊天
#### proto3聊天协议：message Talk
#### 聊天业务的实现
````
注册msgID=2的路由业务：在main中添加router
解析聊天的proto协议：当前发送消息是谁（创建链接时，绑定pid，记录当前链接的玩家）
根据pid得到当前玩家的player对象
将聊天的数据广播给全部在线玩家：Talk(string)
玩家上线时，将玩家加入世界中
````
#### 测试聊天功能

### 世界管理模块：得到全部玩家信息和AOI地图信息
#### 属性
````
AOIManager 当前世界地图AOI的管理模块
当前全部在线的玩家集合
保护玩家集合的锁
````
#### 方法
````
初始化
添加一个玩家 - AddPlayer(*Player)
删除一个玩家 - RemovePlayer(pid int32)
通过玩家ID查询 - GetPlayerByPid(pid int32)
获取全部的在线玩家 - GetAllPlayers() []*Player
````

### 玩家上线广播
````
定义proto: SyncPlayers, Player
获取玩家周围九宫格的玩家
当前位置通过MsgID=200广播给周围（别人看到自己）
将周围玩家的位置发送给当前玩家（自己看到别人）
当前玩家上线后，触发功能
````

### 移动位置广播
````
注册路由MsgID=3
解析客户端protoMsg
得到发送位置的玩家
给其他玩家广播位置
    更新当前player坐标
    组建广播MsgID=200，Tp4
    获取周围九宫格的玩家
    向周围玩家广播当前玩家位置
````

### 玩家下限
#### 断开链接前处理下线业务
````
给server注册链接断开前hook
通过链接属性得到pid
获取周围九宫格玩家
广播MsgID=201
将当前玩家 从世界管理器 删除
将当前玩家 从AOI管理器 删除
````

## 测试
````
GO111MODULE=off go run main.go
````
