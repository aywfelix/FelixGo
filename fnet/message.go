package fnet

// 4字节 消息长度
// 1字节 消息分类 1 proto 2 json 3 http 0 其他
// 4字节 消息id
// n字节 消息内容

type IMessage interface {
	GetDataLen() uint32
	GetMsgID() uint32
	GetMsgType() uint8
	GetData() []byte

	SetMsgID(msgID uint32)
	SetMsgType(msgType uint8)
	SetData(data []byte)
	SetDataLen(dataLen uint32)
}

type Message struct {
	ID      uint32
	DataLen uint32
	MsgType uint8
	Data    []byte
}

func NewMessage(msgID uint32, msgType uint8, data []byte) *Message {
	m := &Message{
		ID:      msgID,
		DataLen: uint32(len(data)),
		MsgType: msgType,
		Data:    data,
	}
	return m
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) GetMsgID() uint32 {
	return m.ID
}

func (m *Message) GetMsgType() uint8 {
	return m.MsgType
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgID(msgID uint32) {
	m.ID = msgID
}

func (m *Message) SetMsgType(msgType uint8) {
	m.MsgType = msgType
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetDataLen(dataLen uint32) {
	m.DataLen = dataLen
}
