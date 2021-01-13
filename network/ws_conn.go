/*!
 * <WebSocket_Conn>
 *
 * Copyright (c) 2019 by <yfm/ ALADING Co.>
 */

package network

import (
	"runtime/debug"
	"time"

	pub "github.com/yfm369/common/public_yfm"

	"github.com/gorilla/websocket"
)

type WebSokConn struct {
	conn       *websocket.Conn
	msgparser  *WSMsgParser
	writechan  chan []byte
	boclose    bool
	Sessionid  uint32
	recv_cb    func([]byte, *WebSokConn)
	conn_cb    func(*WebSokConn)
	disconn_cb func(*WebSokConn)
	alivetime  int64
	FinalMsg   []byte
}

//write线程
func (this *WebSokConn) RunWrite() {
	defer func() {
		if err := recover(); err != nil {
			pub.PrintFileLog(NET_SERVER_LOG, "WebSokConn::RunWrite error:", err, " session:", this.Sessionid,
				", debug:", string(debug.Stack()))
		}
	}()

	for data := range this.writechan {
		if data == nil {
			// pub.PrintFileLog(NET_SERVER_LOG, "WebSokConn::RunWrite data == nil session:",
			// 		this.Sessionid)
			break
		}

		err := this.conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			pub.PrintFileLog(NET_SERVER_LOG, "WebSokConn::RunWrite over :", err.Error(),
				" session:", this.Sessionid)
			break
		}
	}

	if len(this.FinalMsg) != 0 && this.conn != nil {
		this.conn.WriteMessage(websocket.BinaryMessage, this.FinalMsg)
	}

	close(this.writechan)
	this.Close(true)
}

//read线程
func (this *WebSokConn) RunRead() {
	defer func() {
		this.Close(false)
		if err := recover(); err != nil {
			pub.PrintFileLog(NET_SERVER_LOG, "WebSokConn::RunRead error:", err, " session:", this.Sessionid,
				", debug:", string(debug.Stack()))
		}
	}()

	for !this.boclose && this.conn != nil {
		data, err := this.ReadMsg()
		if err != nil {
			pub.PrintFileLog(NET_SERVER_LOG, "WebSokConn::RunRead over :", err.Error())
			break
		}
		if data == nil {
			pub.PrintFileLog(NET_SERVER_LOG, "WebSokConn::RunRead over data == nil")
			break
		}

		this.conn.SetReadDeadline(time.Now().Add(35 * time.Second))
		this.recv_cb(data, this)
	}

	select {
	case this.writechan <- nil:
	default:
		pub.PrintFileLog(NET_SERVER_LOG, "RunRead write msg writechan full conn over ...")
		this.Close(false)
	}
}

func (this *WebSokConn) Read() (data []byte, length int, err error) {
	_, data, err = this.conn.ReadMessage()
	length = len(data)
	return
}

func (this *WebSokConn) ReadMsg() ([]byte, error) {
	return this.msgparser.Read(this)
}

func (this *WebSokConn) SendMsgByCmd(cmd uint16, data []byte) {
	datacmd, _ := pub.IntToByte(uint16(cmd))
	data = append(datacmd, data...)

	length := len(data)
	if length >= pub.MAX_PACKET_LEN {
		//增加大包的标识
		length += 6
		BigFlag, _ := pub.IntToByte(uint16(pub.BIG_PACKET_FLAG))
		datalen, _ := pub.IntToByte(uint32(length))
		data = append(datalen, data...)
		data = append(BigFlag, data...)
	} else {
		length += 2
		datalen, _ := pub.IntToByte(uint16(length))
		data = append(datalen, data...)
	}

	this.Write(data)
}

func (this *WebSokConn) Close(bflag bool) {
	if this.boclose {
		return
	}

	this.boclose = true

	if bflag {
		this.conn.Close()
	}
}

func (this *WebSokConn) GetRemoteAddr() string {
	return this.conn.RemoteAddr().String()
}

func (this *WebSokConn) Write(data []byte) {
	if this == nil {
		return
	}
	if this.boclose {
		return
	}

	select {
	case this.writechan <- data:
	default:
		pub.PrintFileLog(NET_SERVER_LOG, "write msg writechan full conn over ...")
		this.Close(false)
	}
}

func (this *WebSokConn) Type() uint8 {
	return CONN_TYPE_WS
}

func (this *WebSokConn) GetSession() uint32 {
	return this.Sessionid
}

func (this *WebSokConn) SendDirect(data []byte) {
	if this.conn == nil || this.boclose {
		pub.PrintFileLog(NET_SERVER_LOG, " session:", this.Sessionid, " data:", string(data))
		return
	}

	this.FinalMsg = data
	pub.PrintLog("直接发送消息 ！")
}
