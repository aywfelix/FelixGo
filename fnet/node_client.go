package fnet

import (
	. "github.com/felix/felixgo/global"
	. "github.com/felix/felixgo/logger"
	nodeproto "github.com/felix/felixgo/fnet/proto"
)

type INodeClient interface {
	// 作为客户端上报服务器信息
	SetReportInfo(serverType ServerType)
	UpdateOnline(count int)
	UpdateServerState(state nodeproto.ServerState)
	OnMasterRouter()
	AddConnServer()
}

type NodeClient struct {
	NodeService
	netClients INetClients

	serverInfo nodeproto.ServerReport
}

func NewNodeClient(config *NodeConfig) *NodeClient {
	nc := &NodeClient{}
	nc.msgHandler = NewMsgHandler()
	nc.dataPack = NewDataPack()
	nc.nodeConfig = config
	nc.netClients = NewNetClients(nc)
	nc.isStopped = false

	// test default function
	nc.onStarted = nc.onStart
	nc.onStopped = nc.onStop
	nc.onConnected = nc.onSessionConnect
	nc.onDisconnected = nc.onSessionDisconnect
	return nc
}

func (nc *NodeClient) SetReportInfo(serverType ServerType) {
	nc.serverInfo.ServerId = int32(nc.nodeConfig.NodeId)
	nc.serverInfo.ServerName = []byte(nc.nodeConfig.NodeName)
	nc.serverInfo.ServerIp = []byte(nc.nodeConfig.NodeIP)
	nc.serverInfo.ServerPort = int32(nc.nodeConfig.NodePort)
	nc.serverInfo.MaxOnline = int32(nc.nodeConfig.MaxConnect)
	nc.serverInfo.CurOnline = 0
	nc.serverInfo.ServerState = nodeproto.ServerState_SS_NORMAL
	nc.serverInfo.ServerType = int32(serverType)
}

func (nc *NodeClient) UpdateOnline(count int) {
	nc.serverInfo.CurOnline = int32(count)
}

func (nc *NodeClient) UpdateServerState(state nodeproto.ServerState) {
	nc.serverInfo.ServerState = state
}

func (nc *NodeClient) OnMasterRouter() {
	// 所有node客户端 需要上报给master服务器信息
}

func (nc *NodeClient) Start() {
	LogInfo("node client start ...")
	if nc.onStarted != nil {
		nc.onStarted()
	}
	if !nc.isStopped {
		nc.netClients.ProcessExecute()
	}
	LogInfo("node client started ...")
}

func (nc *NodeClient) Serve() {

}

func (nc *NodeClient) Stop() {
	LogInfo("node client stop ...")
	nc.isStopped = true
	if nc.onStopped != nil {
		nc.onStopped()
	}
	LogInfo("node client stopped ...")
}

func (nc *NodeClient) AddConnServer() {
	// TODO: 临时测试
	connData := new(ConnData)
	connData.IP = nc.nodeConfig.NodeIP
	connData.Port = int32(nc.nodeConfig.NodePort)
	connData.NodeID = int32(nc.nodeConfig.NodeId)
	nc.netClients.AddConnServer(connData)
}

func (ns *NodeClient) onStart(args ...interface{}) {
	LogInfo("node client on start")
}

func (ns *NodeClient) onStop(args ...interface{}) {
	LogInfo("node client on stop")

}

func (ns *NodeClient) onSessionConnect(args ...interface{}) {
	LogInfo("node client on session connect")

}

func (ns *NodeClient) onSessionDisconnect(args ...interface{}) {
	LogInfo("node client on session disconnect")
	// 重置连接状态，重新连接服务器
	ns.netClients.ResetConnState(ns.GetNodeID(), CS_CONNECTING)
}
