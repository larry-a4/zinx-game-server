package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//只负责拆包/封包单元测试
func TestDataPack(t *testing.T) {
	/*	模拟服务器	*/
	//1-创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}

	//创建go 承载 负责从client处理业务
	go func() {
		//2-从client读取数据，拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error: ", err)
			}

			go func(conn net.Conn) {
				//处理client request
				//------》拆包《------
				//定义拆包的对象
				dp := NewDataPack()
				for {
					//1-第一次从conn读，把head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err: ", err)
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err: ", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//msg是有数据的，需要进行第二次读取
						//2-第二次从conn读，根据长度在读取data
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						//根据dataLen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data error: ", err)
							return
						}

						//完整的消息已经读取完毕
						fmt.Println("--->Recv MsgID: ", msg.Id, ", dataLen=", msg.DataLen)
						fmt.Println("Data=", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	/*  模拟客户端 */
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}

	//创建封包对象dp
	dp := NewDataPack()

	//模拟粘包过程，封装两个msg一起发送
	//封装第一个msg1包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte("zinx"),
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error:", err)
		return
	}
	//封装第一个msg2包
	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte("nihao!!"),
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error:", err)
		return
	}

	//将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)

	//一次性发送给服务端
	conn.Write(sendData1)

	//客户端阻塞
	select {}
}
