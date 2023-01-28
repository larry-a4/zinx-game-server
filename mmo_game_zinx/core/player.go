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

func (p *Player) Talk(content string) {
	//1 组建MsgID=200的proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1, //tp-1 聊天广播
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}
	//2 得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllPlayers()

	//3 向所有玩家（包括自己）发送MsgID=200消息
	for _, player := range players {
		//向每个player对应的客户端发送消息
		player.SendMsg(200, proto_msg)
	}
}

func (p *Player) SyncSurrounding() {
	// 1 获取玩家周围九宫格的玩家
	neighborIDs := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	neighbors := make([]*Player, len(neighborIDs))
	for i, pid := range neighborIDs {
		neighbors[i] = WorldMgrObj.GetPlayerByPid(int32(pid))
	}

	// 2 当前位置通过MsgID=200广播给周围（别人看到自己）
	// 2.1 组建广播消息200
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
	// 2.2 向每个周围玩家发送信息
	for _, n := range neighbors {
		n.SendMsg(200, protoMsg)
	}

	// 3 将周围玩家的位置发送给当前玩家 MsgID=202
	// 3.1 组建广播消息202
	neighborPositions := make([]*pb.Player, len(neighbors))
	for i, n := range neighbors {
		neighborPositions[i] = &pb.Player{
			Pid: n.Pid,
			P: &pb.Position{
				X: n.X,
				Y: n.Y,
				Z: n.Z,
				V: n.V,
			},
		}
	}
	protoMsg202 := &pb.SyncPlayers{
		Ps: neighborPositions[:],
	}
	// 3.2 发送给当前玩家
	p.SendMsg(202, protoMsg202)
}

//广播当前玩家位置移动信息
func (p *Player) UpdatePos(x, y, z, v float32) {
	//更新当前player坐标
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	//组建广播MsgID=200，Tp4
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//获取周围九宫格的玩家
	neighbors := p.GetSurroundingPlayers()

	for _, n := range neighbors {
		n.SendMsg(200, protoMsg)
	}
}

func (p *Player) Offline() {
	// 获取周围九宫格玩家
	neighbors := p.GetSurroundingPlayers()
	// 广播MsgID=201
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
	}
	for _, n := range neighbors {
		n.SendMsg(201, protoMsg)
	}
	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
	WorldMgrObj.RemovePlayer(p.Pid)
}

func (p *Player) GetSurroundingPlayers() []*Player {
	neighborIDs := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	neighbors := make([]*Player, len(neighborIDs))
	for i, pid := range neighborIDs {
		neighbors[i] = WorldMgrObj.GetPlayerByPid(int32(pid))
	}
	return neighbors
}
