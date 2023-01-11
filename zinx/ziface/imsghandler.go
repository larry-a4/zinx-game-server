package ziface

type IMsgHandler interface {
	//调度/直冲对应的router消息处理方法
	DoMsgHandler(request IRequest)
	//为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
}