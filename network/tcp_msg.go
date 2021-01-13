/*!
 * <消息解析>
 *
 * Copyright (c) 2017.11.16 by <yfm/ Fermin Co.>
 */

package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	pub "github.com/yfm369/common/public_yfm"
)

// 小数据包
// ------------------------------------
// | len(uint16) | cmd(uint16) | data |
// ------------------------------------
// 大数据包
// ----------------------------------------------------
// | 65535(uint16) | len(uint32) | cmd(uint16) | data |
// ----------------------------------------------------

type TCPMsgParser struct {
}

//接收数据
func (this *TCPMsgParser) Read(conn *TcpConn) ([]byte, error) {
	//获取数据包大小
	smallMsgLen := make([]byte, pub.SMALL_PACKET_LEN)
	if _, err := io.ReadFull(conn, smallMsgLen); err != nil {
		return nil, err
	}

	msgLen := uint32(binary.LittleEndian.Uint16(smallMsgLen))
	if msgLen == pub.BIG_PACKET_FLAG { //是大包 再读4个字节表示真正的数据包大小
		bigMsgLen := make([]byte, pub.BIG_PACKET_LEN)
		if _, err := io.ReadFull(conn, bigMsgLen); err != nil {
			return nil, err
		}
		msgLen = binary.LittleEndian.Uint32(bigMsgLen)
	}

	//检查数据包大小
	if msgLen < pub.MIN_PACKET_LEN {
		return nil, errors.New(fmt.Sprintf("message too short %d:", msgLen))
	}

	//数据包数据
	msgData := make([]byte, msgLen)
	if _, err := io.ReadFull(conn, msgData); err != nil {
		return nil, err
	}

	return msgData, nil
}

//发送数据
func (this *TCPMsgParser) Write(conn *TcpConn, data []byte) error {
	if conn == nil {
		return errors.New("Msgparser::Write conn is nil")
	}

	conn.Write(data)
	return nil
}
