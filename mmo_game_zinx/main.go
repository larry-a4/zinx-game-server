package main

import (
	"fmt"

	"../zinx/ziface"
	"../zinx/znet"
	"./core"
)

func OnConnectionAdd(conn ziface.IConnection) {
	//创建一个Player对象
	player := core.NewPlayer(conn)

	//给客户端发送 msgID=1 消息
	player.SyncPid()

	//给客户端发送 msgID=200 消息
	player.BroadcastStartPosition()

	fmt.Println("-------> Player pid = ", player.Pid, " is arrived<------")
}

func main() {
	//创建zinx server handler
	s := znet.NewServer("Zinx MMO Game")

	//链接创建和销毁的hook

	//注册路由业务

	//启动服务
	s.Serve()
}
