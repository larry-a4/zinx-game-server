package utils

import (
	"encoding/json"
	"io/ioutil"

	"../ziface"
)

/*
	存储有关zinx框架的全局参数，供其他模块使用
	一些参数可以通过zinx.json由用户配置
*/
type GlobalObj struct {
	/*
		Server
	*/
	Name      string
	TcpServer ziface.IServer //当前zinx全局的Server对象
	Host      string
	TcpPort   int

	/*
		Zinx
	*/
	Version        string //当前zinx的版本好
	MaxConn        int    //当前服务器主机允许的最大链接数
	MaxPackageSize uint32 //当前zinx框架数据包的最大值
}

/*
	定义一个全局的对外globalobj
*/
var GlobalObject *GlobalObj

/*
	从zinx.json家在用户自定义的参数
*/
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json文件解析为struct
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
	提供一个init方法，初始化GlobalObj
*/
func init() {
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	GlobalObject.Reload()
}
