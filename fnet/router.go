package fnet

type IRouter interface {
	PreHandle(request IRequest) MsgErrCode
	Handle(request IRequest) MsgErrCode
	PostHandle(request IRequest) MsgErrCode
}

type BaseRouter struct{}

// 消息处理前调用
func (br *BaseRouter) PreHandle(request IRequest) MsgErrCode { return MEC_OK }

// 消息处理
func (br *BaseRouter) Handle(request IRequest) MsgErrCode { return MEC_OK }

// 消息处理后调用
func (br *BaseRouter) PostHandle(request IRequest) MsgErrCode { return MEC_OK }
