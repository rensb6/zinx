package znet

import (
	"fmt"
	"net"
	"ziface"
)

type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID,也可称作为SessionID,ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool

	//该连接的处理方法API
	//handleAPI ziface.HandFunc

	//该连接的处理方法router
	Router ziface.IRouter

	//告知该连接已经退出/停止的channel
	ExitBuffChan chan bool
}

//创建连接的方法
func NewConntion(conn *net.TCPConn,connID uint32,router ziface.IRouter) *Connection{
	c:=&Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		Router:    router,
		ExitBuffChan: make(chan bool,1),
	}

	return c
}

//处理conn读数据的GoRoutine
func (c *Connection)StartReader() {
	fmt.Println("Reader GoRoutine is running")

	defer fmt.Println(c.RemoteAddr().String(),"conn reader exit!")
	defer c.Stop()

	for {
		//读取最大的数据到buf中
		buf := make([]byte,512)
		_,err := c.Conn.Read(buf)
		if err!=nil{
			fmt.Println("recv buf error  ",err)
			c.ExitBuffChan <-true
			continue
		}
		//if err:=c.handleAPI(c.Conn,buf,cnt);err!=nil{
		//	fmt.Println("connID ",c.ConnID," handle is error")
		//	c.ExitBuffChan <-true
		//	return
		//}
		//得到当前客户端请求的Request数据
		req:=Request{
			conn: c,
			data: buf,
		}
		//从路由Routers 中找到注册绑定Conn的对应Handle
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

//启动连接，让当前连接开始工作
func (c *Connection)Start()  {
	//开启处理该连接读取到客户端数据之后的请求业务
	go c.StartReader()

	for{
		select {
		case <-c.ExitBuffChan:
			//得到退出消息，不在阻塞
		}
	}
}

//停止连接，结束当前连接状态M
func (c *Connection)Stop(){
	//1.如果当前连接已经关闭
	if c.isClosed == true{
		return
	}
	c.isClosed = true
	//TODO Connection Stop() 如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用

	//关闭socket连接
	c.Conn.Close()

	//通知从缓冲队列里读数据的业务，该连接已经关闭
	c.ExitBuffChan <- true

	//关闭该连接全部管道
	close(c.ExitBuffChan)
}

//从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}
//获取当前连接ID
func (c *Connection) GetConnID() uint32{
	return c.ConnID
}
//获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}