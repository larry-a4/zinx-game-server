package api

import (
	"fmt"

	"../../zinx/ziface"
	"../../zinx/znet"
	"../core"
	"../pb"
	"google.golang.org/protobuf/proto"
)

//世界聊天 路由业务
type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	//1 解析客户端传递进来的proto协议
	proto_msg := &pb.Talk{}
	if err := proto.Unmarshal(request.GetData(), proto_msg); err != nil {
		fmt.Println("Talk Unmarshal error: ", err)
		return
	}

	//2 当前聊天数据属于哪个玩家发送的
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty(pid) error: ", err)
		return
	}

	//3 根据pid得到对应的player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//4 将这个消息广播给其他在线的玩家
	player.Talk(proto_msg.Content)
}
