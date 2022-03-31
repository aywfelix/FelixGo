package fnet

import (
	"fmt"
)

type IMsgHandler interface {
	DoMsg(request IRequest)
	AddRouter(msgID uint32, router IRouter)

	StartWorkPool()
	DispatchByMsgID(request IRequest)

	IsUseWorkPool() bool
}

type IServerMsgHandler interface {
	DispatchByServerID(request IRequest)
}

type IRoleMsgHandler interface {
	DispatchByRoleID(request IRequest)
}

type ISceneMsgHandler interface {
	DispatchBySceneID(request IRequest)
}

type MsgHandler struct {
	// 根据消息id指定路由
	Apis map[uint32]IRouter
	// 利用协程池分发处理消息
	isUsePool    bool
	workPool     []chan IRequest
	workPoolSize uint32
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:         make(map[uint32]IRouter),
		workPoolSize: msg_worker_pool_size,
		workPool:     make([]chan IRequest, msg_worker_pool_size),
	}
}

func (mh *MsgHandler) DoMsg(request IRequest) {
	router, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("msg not found, id=", request.GetMsgID())
		return
	}
	// 相应的路由处理消息
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		fmt.Println("duplicate add msg to api")
		return
	}
	mh.Apis[msgID] = router
}

func (mh *MsgHandler) StartWorkPool() {
	for i := 0; i < int(mh.workPoolSize); i++ {
		mh.workPool[i] = make(chan IRequest, max_msg_worker_size)
		go mh.startWorker(i)
	}
	mh.isUsePool = true
}

func (mh *MsgHandler) startWorker(workID int) {
	// TODO: 消息处理停止
	for {
		select {
		case request := <-mh.workPool[workID]:
			mh.DoMsg(request)
		}
	}
}

func (mh *MsgHandler) DispatchByMsgID(request IRequest) {
	// 根据消息ID分发到不同线程去执行
	workID := request.GetMsgID() % mh.workPoolSize
	mh.workPool[workID] <- request
}

func (mh *MsgHandler) DispatchByRoleID(request IRequest) {
	// 根据消息ID分发到不同线程去执行
	workID := request.GetMsgID() % mh.workPoolSize
	mh.workPool[workID] <- request
}

func (mh *MsgHandler) DispatchBySceneID(request IRequest) {
	// 根据消息ID分发到不同线程去执行
	workID := request.GetMsgID() % mh.workPoolSize
	mh.workPool[workID] <- request
}

func (mh *MsgHandler) DispatchByServerID(request IRequest) {
	// 根据消息ID分发到不同线程去执行
	workID := request.GetMsgID() % mh.workPoolSize
	mh.workPool[workID] <- request
}

func (mh *MsgHandler) IsUseWorkPool() bool {
	return mh.isUsePool
}
