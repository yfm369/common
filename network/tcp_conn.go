/*!
 * <TCP_Conn>
 *
 * Copyright (c) 2017.11.16 by <yfm/ Fermin Co.>
 */

package network

import (
	"net"
	pub "public_yfm"
	"runtime/debug"
)

type TcpConn struct {
	conn       net.Conn
	msgparser  *TCPMsgParser
	writechan  chan []byte
	boclose    bool
	Sessionid  uint32
	recv_cb    func([]byte, *TcpConn)
	conn_cb    func(*TcpConn)
	disconn_cb func(*TcpConn)
}

//write线程
func (this *TcpConn) RunWrite() {
	defer func() {
		this.Close()
		if err := recover(); err != nil {
			pub.PrintFileLog(NET_SERVER_LOG, "TcpConn::RunWrite error:", err, " session:", this.Sessionid,
				", debug:", string(debug.Stack()))
		}
	}()

	for data := range this.writechan {
		if data == nil {
			pub.PrintFileLog(NET_SERVER_LOG, "TcpConn::RunWrite data == nil session:",
				this.Sessionid)
			break
		}

		_, err := this.conn.Write(data)
		if err != nil {
			pub.PrintFileLog(NET_SERVER_LOG, "TcpConn::RunWrite over :", err.Error(),
				" session:", this.Sessionid)
			break
		}
	}
}

//read线程
func (this *TcpConn) RunRead() {
	defer func() {
		this.Close()
		if err := recover(); err != nil {
			pub.PrintFileLog(NET_SERVER_LOG, "TcpConn::RunRead error:", err, " session:", this.Sessionid,
				", debug:", string(debug.Stack()))
		}
	}()

	for !this.boclose && this.conn != nil {
		data, err := this.ReadMsg()
		if err != nil {
			pub.PrintFileLog(NET_SERVER_LOG, "TcpConn::RunRead over :", err.Error())
			break
		}
		if data == nil {
			pub.PrintFileLog(NET_SERVER_LOG, "Tcp::RunRead over data == nil")
			break
		}

		this.recv_cb(data, this)
	}
}

func (this *TcpConn) Read(data []byte) (int, error) {
	return this.conn.Read(data)
}

func (this *TcpConn) ReadMsg() ([]byte, error) {
	return this.msgparser.Read(this)
}

//tcp连接
func (this *TcpConn) SendMsgByCmd(cmd uint16, data []byte) {
	datacmd, _ := pub.IntToByte(uint16(cmd))
	data = append(datacmd, data...)

	length := len(data)
	if length >= pub.MAX_PACKET_LEN {
		//增加大包的标识
		BigFlag, _ := pub.IntToByte(uint16(pub.BIG_PACKET_FLAG))
		datalen, _ := pub.IntToByte(uint32(length))
		data = append(datalen, data...)
		data = append(BigFlag, data...)
	} else {
		datalen, _ := pub.IntToByte(uint16(length))
		data = append(datalen, data...)
	}

	this.Write(data)
}

func (this *TcpConn) Close() {
	if this.boclose {
		return
	}

	this.boclose = true
	close(this.writechan)
	this.conn.Close()
}

func (this *TcpConn) Write(data []byte) {
	if this.boclose {
		return
	}

	select {
	case this.writechan <- data:
	default:
		pub.PrintFileLog(NET_SERVER_LOG, "write msg writechan full conn over ...")
		this.Close()
	}
}

func (this *TcpConn) Type() uint8 {
	return CONN_TYPE_TCP
}

func (this *TcpConn) GetSession() uint32 {
	return this.Sessionid
}
