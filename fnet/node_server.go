package fnet

import (
	. "github.com/felix/felixgo/global"
	. "github.com/felix/felixgo/logger"
)

type INodeServer interface {
}

type NodeServer struct {
	NodeService
	netService INetService
}

func NewNodeServer(config *NodeConfig) *NodeServer {
	ns := new(NodeServer)
	ns.msgHandler = NewMsgHandler()
	ns.dataPack = NewDataPack()
	ns.nodeConfig = config
	ip := config.NodeIP
	port := config.NodePort
	ns.netService = NewNetServer(ns, ip, port)
	ns.isStopped = false
	// set default function
	ns.onStarted = ns.onStart
	ns.onStopped = ns.onStop
	ns.onConnected = ns.onSessionConnect
	ns.onDisconnected = ns.onSessionDisconnect
	return ns
}

func (ns *NodeServer) Start() {
	LogInfo("node server start...")
	if ns.onStarted != nil {
		ns.onStarted()
	}
	// 开启网络监听
	netServer := ns.netService.(INetServer)
	if err := netServer.Listen(); err != nil {
		LogError("Error listening, %v", err)
		ns.Stop()
		return
	}
	// 注册路由服务
	LogInfo("node server started...")
}

func (ns *NodeServer) Serve() {}

func (ns *NodeServer) Stop() {
	LogInfo("node server stop ...")
	ns.isStopped = true
	if ns.onStopped != nil {
		ns.onStopped()
	}
	LogInfo("node server stopped...")
}

func (ns *NodeServer) onStart(args ...interface{}) {
	LogInfo("node server on start")
}

func (ns *NodeServer) onStop(args ...interface{}) {
	LogInfo("node server on stop")
}

func (ns *NodeServer) onSessionConnect(args ...interface{}) {
	LogInfo("node server on session connect")
}

func (ns *NodeServer) onSessionDisconnect(args ...interface{}) {
	LogInfo("node server on session disconnect")
}
