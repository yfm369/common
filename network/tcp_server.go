/*!
 * <TCP_Server>
 *
 * Copyright (c) 2017.11.16 by <yfm/ Fermin Co.>
 */
package network

import (
	"log"
	"net"
	"runtime/debug"
	"sync"
	"time"

	pub "github.com/yfm369/common/public_yfm"
)

const (
	NET_SERVER_LOG = "net_server" //记录net异常的日志
)

type TcpServer struct {
	sync.RWMutex
	listen     net.Listener
	clients    map[uint32]*TcpConn
	Addr       string
	msgparser  *TCPMsgParser
	recv_cb    func([]byte, *TcpConn)
	conn_cb    func(*TcpConn)
	disconn_cb func(*TcpConn)
	ns         uint32 //session计数
}

//初始化tcp服务器
func (this *TcpServer) Start(addr string, recv func([]byte, *TcpConn),
	conn func(*TcpConn), dis func(*TcpConn)) bool {
	if this == nil {
		pub.PrintFileLog(NET_SERVER_LOG, "TcpServer Init Error nil")
		return false
	}

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		pub.PrintFileLog(NET_SERVER_LOG, "TcpServer Init Listen Error :", err.Error())
		return false
	}

	this.listen = listen
	this.clients = make(map[uint32]*TcpConn)
	this.recv_cb = recv
	this.conn_cb = conn
	this.disconn_cb = dis
	this.msgparser = new(TCPMsgParser)

	//接收连接
	go this.Accept_Run()

	return true
}

//接收连接
func (this *TcpServer) Accept_Run() {
	defer func() {
		this.listen.Close()
		if err := recover(); err != nil {
			pub.PrintFileLog(NET_SERVER_LOG, "TcpServer::Accept_Run Error:", err, ", debug:",
				string(debug.Stack()))
		}
	}()

	var delay time.Duration
	for true {
		conn, err := this.listen.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if delay == 0 {
					delay = 5 * time.Millisecond
				} else {
					delay *= 2
				}
				if max := 1 * time.Second; delay > max {
					delay = max
				}

				pub.PrintFileLog(NET_SERVER_LOG, "Accept_Run Error :", err.Error(), " ReTrying in ", delay)

				time.Sleep(delay)
				continue
			}
			log.Println("server accept error :", err.Error())
			return
		}

		this.ns = pub.GetSessionID()

		newconn := new(TcpConn)
		newconn.conn = conn
		newconn.msgparser = this.msgparser
		newconn.writechan = make(chan []byte, 65535) //暂定
		newconn.Sessionid = this.ns
		newconn.recv_cb = this.recv_cb

		this.Lock()
		this.clients[newconn.Sessionid] = newconn
		this.Unlock()
		this.conn_cb(newconn)

		go newconn.RunWrite()

		go func() {
			newconn.RunRead()

			this.disconn_cb(newconn)
			this.Lock()
			delete(this.clients, newconn.Sessionid)
			this.Unlock()
		}()
	}
}
