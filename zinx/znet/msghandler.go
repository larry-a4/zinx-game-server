package znet

import (
	"fmt"

	"../ziface"
)

/*
消息处理模块的实现
*/
type MsgHandler struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	//1-从request中找到msgID
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgId(), " is NOT FOUND! Need register")
		return
	}
	//2-根据msgID调度对应router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	//1-判断当前msg绑定的API方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		//ID已经注册
		panic(fmt.Sprintf("repeat API, msgID = %d", msgID))
	}
	//2-添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " succ!")
}
