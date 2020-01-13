package main

import (
	"fmt"
	"ziface"
	"znet"
)

type PingRouter struct {
	znet.BaseRouter //一定要先继承BaseRouter
}

//test PreHandle
func (this *PingRouter) PreHandle(request ziface.IRequest){
	fmt.Println("call  Router PreHandle")
	_,err :=request.GetConnection().GetTCPConnection().Write([]byte("before ping ....\n"))
	if err!=nil{
		fmt.Println("call back pre ping ping ping err",err)
	}
}
//test Handle
func (this *PingRouter) Handle(request ziface.IRequest){
	fmt.Println("call  Router Handle")
	_,err :=request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err!=nil{
		fmt.Println("call back handle ping ping ping err",err)
	}
}
//test PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest){
	fmt.Println("call  Router PostHandle")
	_,err :=request.GetConnection().GetTCPConnection().Write([]byte("After ping .....\n\n"))
	if err!=nil{
		fmt.Println("call back handle ping ping ping err",err)
	}
}

func main(){
	//创建一个server句柄
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	//2 开启服务
	s.Serve()
}