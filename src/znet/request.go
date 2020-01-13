package znet

import "ziface"

type Request struct {
	conn ziface.Iconnection //已经和客户端建立好的连接
	data []byte	//客户端请求的数据
}

//获取请求连接信息
func (request *Request)GetConnection() ziface.Iconnection{
	return request.conn
}
//获取客户端请求数据
func (request Request)GetData() []byte{
	return request.data
}
