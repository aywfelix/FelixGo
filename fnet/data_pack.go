package fnet

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// // write
// v := uint32(500)
// buf := make([]byte, 4)
// binary.BigEndian.PutUint32(buf, v)
//
// // read
// x := binary.BigEndian.Uint32(buf)

type IDataPack interface {
	UnPack(binaryData []byte) (IMessage, error)
	Pack(msg IMessage) ([]byte, error)
	GetHeadLen() uint32
}

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) UnPack(binaryData []byte) (IMessage, error) {
	buffer := bytes.NewReader(binaryData)

	msg := &Message{}
	if err := binary.Read(buffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(buffer, binary.LittleEndian, &msg.MsgType); err != nil {
		return nil, err
	}
	if err := binary.Read(buffer, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}
	// 判读是否超出限定长度
	if msg.DataLen > max_packet_size {
		return nil, errors.New("too large msg data received")
	}
	// 只需要把header解包，然后通过head的长度，直接读取消息即可
	return msg, nil
}

func (dp *DataPack) Pack(msg IMessage) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(buffer, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, msg.GetMsgType()); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (dp *DataPack) GetHeadLen() uint32 {
	return default_header_len
}
