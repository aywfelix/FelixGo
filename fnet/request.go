package fnet

// 包装消息和连接
type IRequest interface {
	GetSession() ISession
	GetData() []byte
	GetMsgID() uint32
	GetMsgType() uint8
}

type Request struct {
	session ISession
	msg     IMessage
}

func NewRequest(session ISession, msg IMessage) *Request {
	return &Request{
		session: session,
		msg:     msg,
	}
}

func (r *Request) GetSession() ISession {
	return r.session
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}

func (r *Request) GetMsgType() uint8 {
	return r.msg.GetMsgType()
}
