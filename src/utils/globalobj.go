package utils

import (
	"encoding/json"
	"io/ioutil"
	"ziface"
)

/*
存储一切有关Zinx框架的全局参数，供其他模块使用
一些参数也可以通过 用户根据 zinx.json来配置
*/
type GlobalObj struct {
	TcpServer ziface.IServer //当前zinx的全局Server对象
	Host string //当前服务器的主机IP
	TcpPort int //当前服务器的主机端口
	HostName string //当前服务器的主机名
	Version string //当前Zinx的版本号

	MaxPacketSize uint32 //数据包最大值
	MaxConn int //当前服务器支持的客户端最大连接数
}

/* 定义一个全局的对象 */
var GlobalObject *GlobalObj

//读取用户的配置文件
func (g *GlobalObj) Reload(){
	data ,err := ioutil.ReadFile("./conf/zinx.json")
	if err!=nil{
		panic(err) //异常时，终止下方代码运行
	}
	//将数据解析到struct中
	err = json.Unmarshal(data,GlobalObject)
	if err!=nil{
		panic(err)
	}
}

//提供init()方法，系统默认执行
func init(){
	//初始化全局配置对象，提供一定默认参数
	GlobalObject = &GlobalObj{
		HostName: "ZinxServerApp",
		Version: "V0.4",
		TcpPort: 7777,
		Host: "0.0.0.0",
		MaxConn: 12000,
		MaxPacketSize:4096,
	}
	//从配置文件中加载用户参数
	GlobalObject.Reload()
}
