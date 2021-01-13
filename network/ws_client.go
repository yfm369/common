/*!
 * <ws_Client>
 *
 * Copyright (c) 2017.11.16 by <yfm/ Fermin Co.>
 */

package network

import (
	pub "common/public_yfm"
	"time"

	"github.com/gorilla/websocket"
)

const (
	WS_CLIENT_LOG = "ws_client"
)

type WSClient struct {
	Addr       string
	Conn       *WebSokConn
	boclose    bool
	msgparser  *WSMsgParser
	recv_cb    func([]byte, *WebSokConn)
	conn_cb    func(*WebSokConn)
	disconn_cb func(*WebSokConn)
}

func (this *WSClient) Start(addr string, recv func([]byte, *WebSokConn), conn func(*WebSokConn),
	disconn func(*WebSokConn)) {
	if this == nil {
		panic("WSClient::Start this == nil")
	}

	this.Addr = addr
	this.boclose = false
	this.msgparser = new(WSMsgParser)
	this.recv_cb = recv
	this.conn_cb = conn
	this.disconn_cb = disconn

	go this.Conncet()
}

func (this *WSClient) Dial() *websocket.Conn {
	if this == nil {
		pub.PrintFileLog(NET_CLIENT_LOG, "Dial this == nil")
		return nil
	}

	for {
		conn, _, err := websocket.DefaultDialer.Dial(this.Addr, nil)
		if err == nil || this.boclose {
			return conn
		}

		pub.PrintFileLog(NET_CLIENT_LOG, "Dial Error :", err.Error())

		time.Sleep(10 * time.Second)
	}
}

func (this *WSClient) Conncet() {
reconn:
	conn := this.Dial()
	if conn == nil {
		return
	}

	newconn := new(WebSokConn)
	newconn.boclose = false
	newconn.conn = conn
	newconn.conn_cb = this.conn_cb
	newconn.disconn_cb = this.disconn_cb
	newconn.recv_cb = this.recv_cb
	newconn.writechan = make(chan []byte, 65535)
	newconn.msgparser = this.msgparser
	newconn.conn.SetReadDeadline(time.Now().Add(35 * time.Second))

	this.Conn = newconn
	this.conn_cb(newconn)

	go newconn.RunWrite()

	newconn.RunRead()
	this.disconn_cb(newconn)

	goto reconn
}
