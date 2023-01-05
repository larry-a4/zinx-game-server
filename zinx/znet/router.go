package znet

import "../ziface"

//实现router时，先嵌入这个BaseRouter基类，然后根据需要重写
type BaseRouter struct{}

/*
	以下方法故意留空，因为有的Router不需要PreHandle/PostHandle
	Router全部继承Base Router，不需要实现PreHanle/PostHandle
*/

//在处理conn业务之前的hook方法
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

//在处理业务的主方法
func (br *BaseRouter) Hanle(request ziface.IRequest) {}

//在处理conn业务之后的hook方法
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
