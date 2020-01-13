package znet

import (
	"errors"
	"fmt"
	"net"
	"time"
	"utils"
	"ziface"
)

//iserver的接口实现，定义一个Server类
type Server struct {
	//服务器名称
	Name string
	//ip类型
	IPVersion string
	//服务器绑定地址
	IP string
	//端口
	Port int

	//当前Server由用户绑定的回调router,也就是Server注册的链接对应的处理业务
	Router ziface.IRouter
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
}

//======================================定义当前连接客户端的HandleApi====================================
func CallBackToClient(conn *net.TCPConn,data []byte,cnt int) error {
	//回显业务
	fmt.Println("[Conn Handle] CallBackToClient ... ")
	//回显
	if _,err := conn.Write(data[:cnt]);err!=nil{
		fmt.Println("write back buf err ",err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

//======================================实现iserver接口中的全部方法======================================
//开启网络服务
func (s *Server) Start(){
	fmt.Printf("[START] Server listening at IP: %s, Port :%d, is starting\n",s.IP,s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize:%d\n",
	utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	//开启一个go 去做Listener服务
	go func() {
		//1. 获取一个 tcp 的Addr
		addr , err := net.ResolveTCPAddr(s.IPVersion,fmt.Sprintf("%s:%d",s.IP,s.Port))
		if err!=nil{
			fmt.Println("resolve tcp addr err: ",err)
			return
		}

		//2.监听服务器地址
		listenner , err := net.ListenTCP(s.IPVersion,addr)
		if err != nil{
			fmt.Println("listen tcp error :",err)
			return
		}

		//已经监听成功
		fmt.Println("start zinx server  ",s.Name," succ now listenning....")

		//TODO server.go 应该有一个自动生成ID的方法
		var cid uint32
		cid = 0
		//3.启动server网络连接服务
		for {
			//3.1阻塞等待客户端建立连接请求
			conn , err := listenner.AcceptTCP()
			if err!=nil{
				fmt.Println("accept err: ",err)
				continue
			}

			//3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			//3.3 TODO Server.Start() 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的

			dealConn:= NewConntion(conn,cid,s.Router)
			cid++
			//3.4启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}

//关闭网络服务
func (s *Server) Stop(){
	fmt.Println("[STOP] Zinx server , name " , s.Name)

	//TODO Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并 停止或者清理
}

//启动网络业务服务
func (s *Server) Serve(){
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加
	//阻塞，否则主go退出，listenner的go将会退出
	for{
		time.Sleep(10*time.Second)
	}
}

func NewServer() ziface.IServer{
	//先初始化全局配置文件
	utils.GlobalObject.Reload()
	s:=&Server{
		Name:      utils.GlobalObject.HostName,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:	   nil,
	}
	return s
}