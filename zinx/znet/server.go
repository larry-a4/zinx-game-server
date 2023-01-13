package znet

import (
	"fmt"
	"net"

	"../utils"
	"../ziface"
)

//iServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的IP版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	//当前server的消息管理模块
	MsgHandler ziface.IMsgHandler
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, listenner at IP: %s, Port: %d is starting\n",
		s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
		//0 开启消息队列以及worker pool
		s.MsgHandler.StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		//2 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, " err ", err)
			return
		}

		fmt.Println("start Zinx server succ, ", s.Name, " succ, Listening...")
		var cid uint32
		cid = 0

		//3 阻塞的等待客户端连接，处理客户端连接业务（读写）
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			//将处理新链接的业务方法和conn进行绑定，得到我们的链接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			//启动当前的业务处理
			go dealConn.Start()

		}

	}()
}

func (s *Server) Stop() {
	//todo：将一些服务器的资源/状态或已经开辟的连接信息，进行停止或回收
}

func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//todo：做一些启动服务器之后的额外业务

	//阻塞状态
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!!")
}

/*
	初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandler(),
	}

	return s
}
