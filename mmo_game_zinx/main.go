package main

import (
	"fmt"

	"../zinx/ziface"
	"../zinx/znet"
	"./api"
	"./core"
)

func OnConnectionAdd(conn ziface.IConnection) {
	//创建一个Player对象
	player := core.NewPlayer(conn)

	//给客户端发送 msgID=1 消息
	player.SyncPid()

	//给客户端发送 msgID=200 消息
	player.BroadcastStartPosition()

	//将新上线的玩家添加到WorldManager中
	core.WorldMgrObj.AddPlayer(player)

	//将该链接绑定一个pid 玩家ID的属性
	conn.SetProperty("pid", player.Pid)

	//广播给周边玩家，当前玩家的位置
	player.SyncSurrounding()

	fmt.Println("-------> Player pid = ", player.Pid, " is arrived<------")
}

func OnConnectionLost(conn ziface.IConnection) {
	// 给server注册链接断开前hook
	pid, _ := conn.GetProperty("pid")
	// 通过链接属性得到pid
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	// 触发玩家下线业务
	player.Offline()

	fmt.Println("------>Player pid = ", pid, " offline<-----")
}

func main() {
	//创建zinx server handler
	s := znet.NewServer("Zinx MMO Game")

	//链接创建和销毁的hook
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	//注册路由业务
	s.AddRouter(2, &api.WorldChatApi{})
	s.AddRouter(3, &api.MoveApi{})

	//启动服务
	s.Serve()
}
