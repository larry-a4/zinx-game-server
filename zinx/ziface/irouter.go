package ziface

/*
	router路由：提供指令和对应的处理方式
	路由抽象接口，
	路由里的数据都是IRequest
*/
type IRouter interface {
	//在处理conn业务之前的hook方法
	PreHandle(request IRequest)

	//在处理业务的主方法
	Handle(request IRequest)

	//在处理conn业务之后的hook方法
	PostHandle(request IRequest)
}
