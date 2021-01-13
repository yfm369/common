/*!
 * <数据包解析>
 *
 * Copyright (c) 2019 by <yfm/ ALADING Co.>
 */

package network

import (
	"bytes"
	pub "common/public_yfm"
	"encoding/binary"
)

const (
	PACKET_LOG = "package_error"
)

type Packet struct {
	buffer  []byte
	readpos int
}

func (this *Packet) Init(data []byte) {
	this.Reset()
	this.buffer = append(this.buffer, data...)
}

func (this *Packet) Reset() {
	this.buffer = make([]byte, 0)
	this.readpos = 0
}

func (this *Packet) GetBuffer() []byte {
	return this.buffer
}

func (this *Packet) WriteInt(v interface{}) {
	by, err := pub.IntToByte(v)
	if err != nil {
		pub.PrintFileLog(PACKET_LOG, "WriteInt error :", err.Error())
	}

	this.buffer = append(this.buffer, by...)
}

func (this *Packet) ReadUint8() uint8 {
	var value uint8 = 0
	if this.readpos+1 > len(this.buffer) {
		return value
	}

	header := bytes.NewBuffer(this.buffer[this.readpos : this.readpos+1])
	err := binary.Read(header, binary.LittleEndian, &value)
	if err != nil {
		pub.PrintFileLog(PACKET_LOG, "ReadUint8 error :", err.Error())
	}
	this.readpos += 1

	return value
}

func (this *Packet) Readint8() int8 {
	var value int8 = 0
	if this.readpos+1 > len(this.buffer) {
		return value
	}

	header := bytes.NewBuffer(this.buffer[this.readpos : this.readpos+1])
	err := binary.Read(header, binary.LittleEndian, &value)
	if err != nil {
		pub.PrintFileLog(PACKET_LOG, "Readint8 error :", err.Error())
	}
	this.readpos += 1

	return value
}

func (this *Packet) ReadUint16() uint16 {
	var value uint16 = 0
	if this.readpos+2 > len(this.buffer) {
		return value
	}

	header := bytes.NewBuffer(this.buffer[this.readpos : this.readpos+2])
	err := binary.Read(header, binary.LittleEndian, &value)
	if err != nil {
		pub.PrintFileLog(PACKET_LOG, "ReadUint16 error :", err.Error())
	}
	this.readpos += 2

	return value
}

func (this *Packet) Readint16() int16 {
	var value int16 = 0
	if this.readpos+2 > len(this.buffer) {
		return value
	}

	header := bytes.NewBuffer(this.buffer[this.readpos : this.readpos+2])
	err := binary.Read(header, binary.LittleEndian, &value)
	if err != nil {
		pub.PrintFileLog(PACKET_LOG, "Readint16 error :", err.Error())
	}
	this.readpos += 2

	return value
}

func (this *Packet) ReadUint32() uint32 {
	var value uint32 = 0
	if this.readpos+4 > len(this.buffer) {
		return value
	}

	header := bytes.NewBuffer(this.buffer[this.readpos : this.readpos+4])
	err := binary.Read(header, binary.LittleEndian, &value)
	if err != nil {
		pub.PrintFileLog(PACKET_LOG, "ReadUint32 error :", err.Error())
	}
	this.readpos += 4

	return value
}

func (this *Packet) Readint32() int32 {
	var value int32 = 0
	if this.readpos+4 > len(this.buffer) {
		return value
	}

	header := bytes.NewBuffer(this.buffer[this.readpos : this.readpos+4])
	err := binary.Read(header, binary.LittleEndian, &value)
	if err != nil {
		pub.PrintFileLog(PACKET_LOG, "ReadUint32 error :", err.Error())
	}
	this.readpos += 4

	return value
}

func (this *Packet) WriteString(str string) {
	strby := []byte(str)
	strlen := len(strby)
	this.WriteInt(uint16(strlen))
	this.buffer = append(this.buffer, strby...)
}

func (this *Packet) ReadString() string {
	strlen := this.ReadUint16()
	if this.readpos+int(strlen) > len(this.buffer) {
		return ""
	}

	strbyte := this.buffer[this.readpos : this.readpos+int(strlen)]
	this.readpos += len(strbyte)

	return string(strbyte)
}
