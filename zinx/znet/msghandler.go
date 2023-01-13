package znet

import (
	"fmt"

	"../utils"
	"../ziface"
)

/*
消息处理模块的实现
*/
type MsgHandler struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责worker的消息队列
	TaskQueue []chan ziface.IRequest
	//worker数量
	WorkerPoolSize uint32
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置获取
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

//启动一个worker工作池(只发生一次)
func (mh *MsgHandler) StartWorkerPool() {
	//根据poolSize开启worker，每个用一个go
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//1-当前的worker对应的消息队列 开辟空间 第0个worker用第0个channel...
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.startOneWorker(i, mh.TaskQueue[i])

	}
}

//启动一个worker工作流程
func (mh *MsgHandler) startOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started...")

	for {
		select {
		//如果有消息过来，出列的是一个client request，执行当前Request所绑定业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//将消息交给taskQueue，由worker进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	//1-将消息平均分配给不同的worker
	//根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		", request MsgID = ", request.GetMsgId(),
		" to WorkerID = ", workerID)

	//2-将消息发送给对应的worker的taskQueue
	mh.TaskQueue[workerID] <- request
}
