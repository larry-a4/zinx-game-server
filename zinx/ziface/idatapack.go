package ziface

/*
	封拆包模块
	直接面向TCP数据流，处理粘包问题
*/
type IDataPack interface {
	//获取包的头的长度
	GetHeadLen() uint32
	//封包方法
	Pack(msg IMessage) ([]byte, error)
	//拆包方法
	Unpack([]byte) (IMessage, error)
}
