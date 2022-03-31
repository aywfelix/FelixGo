package fnet

import (
	"context"
	. "github.com/felix/felixgo/global"
	. "github.com/felix/felixgo/thread"
)

// node 服务接口
type INodeService interface {
	GetMsgHandler() IMsgHandler
	GetDataPack() IDataPack
	GetNodeConfig() *NodeConfig
	RegisterRouter(msgID uint32, router IRouter)

	SetOnStart(onStart FuncArgs)
	SetOnStop(onStop FuncArgs)

	SetOnConnected(onConnected FuncArgs)
	SetOnDisconnected(onDisconnected FuncArgs)
	GetOnConnected() FuncArgs
	GetOnDisconnected() FuncArgs

	IsStop() bool
	GetNodeID() int32
}

type NodeService struct {
	msgHandler IMsgHandler
	dataPack   IDataPack

	onStarted FuncArgs
	onStopped FuncArgs

	onConnected    FuncArgs
	onDisconnected FuncArgs

	isStopped  bool
	nodeConfig *NodeConfig

	ctx    context.Context
	cancel func()
}

func (ns *NodeService) GetMsgHandler() IMsgHandler {
	return ns.msgHandler
}
func (ns *NodeService) GetDataPack() IDataPack {
	return ns.dataPack
}
func (ns *NodeService) GetNodeConfig() *NodeConfig {
	return ns.nodeConfig
}
func (ns *NodeService) RegisterRouter(msgID uint32, router IRouter) {
	ns.msgHandler.AddRouter(msgID, router)
}

func (ns *NodeService) SetOnStart(onStart FuncArgs) {
	ns.onStarted = onStart
}

func (ns *NodeService) SetOnStop(onStop FuncArgs) {
	ns.onStopped = onStop
}

func (ns *NodeService) SetOnConnected(onConnected FuncArgs) {
	ns.onConnected = onConnected
}
func (ns *NodeService) SetOnDisconnected(onDisconnected FuncArgs) {
	ns.onDisconnected = onDisconnected
}
func (ns *NodeService) GetOnConnected() FuncArgs {
	return ns.onConnected
}
func (ns *NodeService) GetOnDisconnected() FuncArgs {
	return ns.onDisconnected
}
func (ns *NodeService) IsStop() bool {
	return ns.isStopped
}
func (ns *NodeService) GetNodeID() int32 {
	return int32(ns.nodeConfig.NodeId)
}
