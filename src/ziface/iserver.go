package ziface

//定义服务器接口

type IServer interface {
	//服务器启动
	Start()
	//服务器停止
	Stop()
	//启动业务服务器
	Serve()

	//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	AddRouter(router IRouter)
}
