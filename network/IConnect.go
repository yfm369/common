/*!
 * <网络连接接口>
 *
 * Copyright (c) 2019 by <yfm/ ALADING Co.>
 */

package network

const (
	CONN_TYPE_WS  = uint8(0) //websocket
	CONN_TYPE_TCP = uint8(1) //tcp
)

type IConnect interface {
	SendMsgByCmd(uint16, []byte)
	Close()
	Write([]byte)
	Type() uint8
	GetSession() uint32
}
