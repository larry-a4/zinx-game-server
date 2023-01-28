package api

import (
	"fmt"

	"../../zinx/ziface"
	"../../zinx/znet"
	"../core"
	"../pb"
	"google.golang.org/protobuf/proto"
)

type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	// 	解析客户端protoMsg
	protoMsg := &pb.Position{}
	if err := proto.Unmarshal(request.GetData(), protoMsg); err != nil {
		fmt.Println("Move: Position unmarshal err: ", err)
		return
	}

	// 得到发送位置的玩家
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error: ", err)
		return
	}

	// fmt.Printf("Player pid=%d, move(%f,%f,%f,%f)",
	// 	pid, protoMsg.X, protoMsg.Y, protoMsg.Z, protoMsg.V)

	// 给其他玩家广播位置proto_msg
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	player.UpdatePos(protoMsg.X, protoMsg.Y, protoMsg.Z, protoMsg.V)
}
