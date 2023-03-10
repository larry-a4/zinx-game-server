package main

import (
	"fmt"

	"../../zinx/ziface"
	"../../zinx/znet"
)

/*
	基于Zinx框架开发的 服务器端应用程序
*/

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Test Handle
func (this *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	//先读取客户端数据，再回写ping..ping..ping

	fmt.Println("recv from client: msgID = ", req.GetMsgId(),
	  ", data = ", string(req.GetData()))

	err := req.GetConnection().SendMsg(1,[]byte("ping..ping..ping"))
	if err != nil{
		fmt.Println(err)
	}
}

func main() {
	//1 创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.5]")
	//2 给zinx添加一个自定义router
	s.AddRouter(&PingRouter{})
	//3 启动server
	s.Serve()
}
