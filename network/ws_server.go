/*!
 * <WS_Server>
 *
 * Copyright (c) 2019 by <yfm/ ALADING Co.>
 */
package network

import (
	"net"
	"net/http"
	"sync"
	"time"

	pub "github.com/yfm369/common/public_yfm"

	"github.com/gorilla/websocket"
)

type WSServer struct {
	listen  net.Listener
	handler *WSHandler
}

type WSHandler struct {
	upgrader   websocket.Upgrader
	mutexConns sync.Mutex

	clients    map[uint32]*WebSokConn
	msgparser  *WSMsgParser
	recv_cb    func([]byte, *WebSokConn)
	conn_cb    func(*WebSokConn)
	disconn_cb func(*WebSokConn)
	ns         uint32 //session计数
}

func (handler *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pub.PrintLog("ServeHTTP new client")
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	conn, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		pub.PrintFileLog(NET_SERVER_LOG, "upgrade error: %v", err.Error())
		return
	}

	handler.mutexConns.Lock()
	if handler.clients == nil {
		handler.mutexConns.Unlock()
		conn.Close()
		return
	}

	handler.ns = pub.GetSessionID()

	newconn := new(WebSokConn)
	newconn.conn = conn
	conn.SetReadDeadline(time.Now().Add(35 * time.Second))
	newconn.Sessionid = handler.ns
	newconn.conn_cb = handler.conn_cb
	newconn.disconn_cb = handler.disconn_cb
	newconn.recv_cb = handler.recv_cb
	newconn.msgparser = new(WSMsgParser)
	newconn.writechan = make(chan []byte, 65535) //暂定
	//newconn.alivetime = time.Now().Unix()

	handler.clients[newconn.Sessionid] = newconn
	handler.mutexConns.Unlock()

	handler.conn_cb(newconn)

	go newconn.RunWrite()

	newconn.RunRead()

	handler.disconn_cb(newconn)
	handler.mutexConns.Lock()
	delete(handler.clients, newconn.Sessionid)
	handler.mutexConns.Unlock()
}

//初始化tcp服务器
func (this *WSServer) Start(addr string, recv func([]byte, *WebSokConn),
	conn func(*WebSokConn), dis func(*WebSokConn)) bool {

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		pub.PrintFileLog(NET_SERVER_LOG, "listen error :", err.Error())
		return false
	}

	pub.PrintLog("初始化tcp服务器")

	// https的处理
	// config := &tls.Config{}
	// config.NextProtos = []string{"http/1.1"}
	// var errtls error
	// config.Certificates = make([]tls.Certificate, 1)
	// config.Certificates[0], errtls = tls.LoadX509KeyPair("1.pem",
	// 	"1.key")
	// if errtls != nil {
	// 	pub.PrintLog("tls error:", errtls.Error())
	// 	return false
	// }
	// ln = tls.NewListener(ln, config)

	this.listen = ln
	this.handler = &WSHandler{
		clients:    make(map[uint32]*WebSokConn),
		recv_cb:    recv,
		conn_cb:    conn,
		disconn_cb: dis,
		upgrader: websocket.Upgrader{
			HandshakeTimeout: 5 * time.Second,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}

	httpServer := &http.Server{
		Addr:           addr,
		Handler:        this.handler,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1024,
	}

	go httpServer.Serve(ln)

	return true
}
