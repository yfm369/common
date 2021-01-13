/*!
 * <TCP_Client>
 *
 * Copyright (c) 2017.11.16 by <yfm/ Fermin Co.>
 */

package network

import (
	"net"
	pub "public_yfm"
	"time"
)

const (
	NET_CLIENT_LOG = "tcp_client"
)

type TcpClient struct {
	Addr       string
	Conn       *TcpConn
	boclose    bool
	msgparser  *TCPMsgParser
	recv_cb    func([]byte, *TcpConn)
	conn_cb    func(*TcpConn)
	disconn_cb func(*TcpConn)
}

func (this *TcpClient) Start(addr string, recv func([]byte, *TcpConn), conn func(*TcpConn),
	disconn func(*TcpConn)) {
	if this == nil {
		panic("TcpClient::Start this == nil")
	}

	this.Addr = addr
	this.boclose = false
	this.msgparser = new(TCPMsgParser)
	this.recv_cb = recv
	this.conn_cb = conn
	this.disconn_cb = disconn

	go this.Conncet()
}

func (this *TcpClient) Dial() net.Conn {
	if this == nil {
		pub.PrintFileLog(NET_CLIENT_LOG, "Dial this == nil")
		return nil
	}

	for {
		conn, err := net.Dial("tcp", this.Addr)
		if err == nil || this.boclose {
			return conn
		}

		pub.PrintFileLog(NET_CLIENT_LOG, "Dial Error :", err.Error())

		time.Sleep(5 * time.Second)
	}
}

func (this *TcpClient) Conncet() {
reconn:
	conn := this.Dial()
	if conn == nil {
		return
	}

	newconn := new(TcpConn)
	newconn.boclose = false
	newconn.conn = conn
	newconn.conn_cb = this.conn_cb
	newconn.disconn_cb = this.disconn_cb
	newconn.recv_cb = this.recv_cb
	newconn.writechan = make(chan []byte, 65535)
	newconn.msgparser = this.msgparser

	this.Conn = newconn
	this.conn_cb(newconn)

	go newconn.RunWrite()

	newconn.RunRead()

	this.disconn_cb(newconn)
	goto reconn
}
