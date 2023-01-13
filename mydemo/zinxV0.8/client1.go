package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"../../zinx/znet"
)

/*
	模拟客户端
*/
func main() {
	fmt.Println("client 1 start")

	time.Sleep(1 * time.Second)

	//1直接连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		//发送封包的message消息 - MsgID=0
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(1, []byte("Zinx client 1 test..")))
		if err != nil {
			fmt.Println("Pack error:", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error: ", err)
			return
		}

		//服务区应该回复一个消息，MsgID=1，data=ping..ping..

		//先读取流中的head部分， 得到ID和长度
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error: ", err)
			break
		}

		//将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error: ", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			//根据DataLen做二次读取
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error: ", err)
				return
			}
			fmt.Println("---> recv Server Msg: ID = ", msg.Id,
				", len = ", msg.DataLen, ", data = ", string(msg.Data))

		}

		//cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
