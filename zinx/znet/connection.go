package znet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"../ziface"
)

/*
	链接模块
*/
type Connection struct {
	//当前链接的socket TCP套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前的链接状态
	isClosed bool

	//告知当前链接已经退出的channel
	ExitChan chan bool

	//该链接处理的方法
	Router ziface.IRouter
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}

	return c
}

//链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader goroutine is running...")
	defer fmt.Println("ConnID = ", c.ConnID, "Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端数据到buf
		// buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		// _, err := c.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("recv buf err ", err)
		// 	continue
		// }

		//创建拆包解包对象
		dp := NewDataPack()

		//读取客户端的Msg Head（二进制流8字节）
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error: ", err)
			break
		}

		//拆包，得到 MsgId 和 MsgDataLen，放入msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error: ", err)
			break
		}

		//根据 DataLen 再次读取 Data，放在 msg.Data 中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error: ", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn的request数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
		//从路由中找到注册绑定的Conn对应的router调用

	}
}

//启动链接：让当前链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start().. ConnID = ", c.ConnID)

	//启动从当前链接读数据的业务
	go c.StartReader()

	//todo启动从当前链接写数据的业务
}

//停止链接：结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	if c.isClosed {
		return
	}
	c.isClosed = true

	//关闭socket链接
	c.Conn.Close()

	//回收资源
	close(c.ExitChan)
}

//获取当前链接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//提供 SendMsg 方法，将我们要发送客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}
	//将data封包, MsgDataLen|MsgId|Data
	dp := NewDataPack()

	msgInBytes, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg id =", msgId)
		return errors.New("pack error msg")
	}

	//将数据发送给客户端
	if _, err := c.Conn.Write(msgInBytes); err != nil {
		fmt.Println("write msg id ", msgId, " error: ", err)
		return errors.New("conn write error")
	}

	return nil
}
