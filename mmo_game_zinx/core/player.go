package core

import (
	"fmt"
	"math/rand"
	"sync"

	"../../zinx/ziface"
	"../pb"
	"google.golang.org/protobuf/proto"
)

type Player struct {
	Pid  int32              //玩家ID
	Conn ziface.IConnection //当前玩家的链接（用于和客户端的链接）
	X    float32            //平面的经度
	Y    float32            //高度
	Z    float32            //平面的纬度
	V    float32            //旋转的0-360角度
}

//全局计数器，用于生成玩家ID
var PidGen int32 = 1
var IdLock sync.Mutex

func NewPlayer(conn ziface.IConnection) *Player {
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	return &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)),
		Y:    0,
		Z:    float32(140 + rand.Intn(20)),
		V:    0,
	}
}

/*
	发送给客户端消息的方法
	主要将protobuf序列化之后，再调用SendMsg方法
*/
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	//将proto message序列化 转换成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err: ", err)
		return
	}

	//将二进制文件 通过zinx框架的sendMsg发送给客户端
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}
	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player SendMsg error: ", err)
		return
	}
}

//告知客户端玩家Pid，同步已经生成的玩家ID给客户端
func (p *Player) SyncPid() {
	//组建MsgID=0 的proto数据
	protoMsg := &pb.SyncPid{
		Pid: p.Pid,
	}

	//将消息发送给客户端
	p.SendMsg(1, protoMsg)
}

//广播玩家自己的出生地点
func (p *Player) BroadcastStartPosition() {
	//组建MsgID=200 的proto数据
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//发给客户端
	p.SendMsg(200, protoMsg)
}
