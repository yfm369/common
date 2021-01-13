/*!
 * <消息解析>
 *
 * Copyright (c) 2019 by <yfm/ ALADING Co.>
 */

package network

import (
	"errors"
)

type WSMsgParser struct {
}

//接收数据
func (this *WSMsgParser) Read(conn *WebSokConn) ([]byte, error) {
	data, _, err := conn.Read()
	if err != nil {
		return nil, err
	}

	return data, nil
}

//发送数据
func (this *WSMsgParser) Write(conn *WebSokConn, data []byte) error {
	if conn == nil {
		return errors.New("WSMsgParser::Write conn is nil")
	}

	conn.Write(data)
	return nil
}
