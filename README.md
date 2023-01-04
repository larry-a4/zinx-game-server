# zinx框架

## ZinxV0.1-基础的server

### 方法

#### 启动服务器
基本的服务器开发 1创建addr 2创建listener 3处理客户端基本业务（回显）

#### 停止服务器
做一些资源的回收和状态的回执

#### 运行服务器
调用Start()方法，调用之后做阻塞处理，之后可以扩展

#### 初始化

### 属性

#### 名称
#### 监听的IP
#### 监听的端口

## ZinxV0.2-简单的链接封装和业务绑定

### 链接的模块

### 方法

#### 启动链接Start()
#### 停止链接Stop()
#### 获取当前链接的conn对象GetTCPConnection()
#### 得到链接ID - GetConnID()
#### 得到客户端链接的地址和端口RemoteAddr()
#### 发送数据Send()
#### 链接所绑定的处理业务的函数类型

### 属性

#### socket TCP套接字: *net.TCPConn
#### 链接的ID: uint32
#### 链接的状态（是否已经关闭): bool
#### 与当前链接所绑定的处理业务方法: ziface.HandleFunc
#### 等待链接被动退出的channel: chan bool

## 测试
```GO111MODULE=off go run server.go```

```GO111MODULE=off go run client.go```
