package fnet

import (
	"time"

	. "github.com/aywfelix/felixgo/logger"
)

type INetClients interface {
	// 添加将要连接的服务
	AddConnServer(connData *ConnData)
	// 进行连接服务
	ProcessExecute()
	// 重置连接状态
	ResetConnState(nodeID int32, state ConnState)
}

type NetClients struct {
	// 保存所有需要连接的服务器信息
	connMap ConnDataMap
	// 临时连接集合
	tempConnMap ConnDataMap
	nodeService INodeService
}

func NewNetClients(service INodeService) *NetClients {
	cs := new(NetClients)
	cs.connMap = ConnDataMap{}
	cs.tempConnMap = ConnDataMap{}
	cs.nodeService = service
	return cs
}

func (cs *NetClients) ProcessExecute() {
	go func() {
		for {
			if cs.nodeService != nil && !cs.nodeService.IsStop() {
				cs.processTempConnect()
				cs.processConnect()
			}
			time.Sleep(time.Second * 1)
		}
	}()
}

func (cs *NetClients) AddConnServer(connData *ConnData) {
	nodeID := connData.NodeID
	if _, ok := cs.tempConnMap[nodeID]; ok {
		return
	}
	cs.tempConnMap[nodeID] = connData
}

func (cs *NetClients) processConnect() {
	// 根据连接状态处理现有的连接情况
	for _, connData := range cs.connMap {
		if connData.netService == nil {
			continue
		}
		netClient := connData.netService.(*NetClient)
		switch connData.ConnState {
		case CS_CONNECTING:
			if ret := netClient.Connect(connData.IP, int(connData.Port)); ret {
				connData.ConnState = CS_CONNECTED
				netClient.SetNodeService(cs.nodeService)
			} else {
				connData.ConnState = CS_DISCONNECTED
			}
		case CS_DISCONNECTED:
			connData.ConnState = CS_RECONNECTING
			if ret := netClient.Connect(connData.IP, int(connData.Port)); ret {
				connData.ConnState = CS_RECONNECTED
				netClient.SetNodeService(cs.nodeService)
			}
		case CS_RECONNECTING:
			connData.ConnState = CS_DISCONNECTED
			if ret := netClient.Connect(connData.IP, int(connData.Port)); ret {
				connData.ConnState = CS_RECONNECTED
				netClient.SetNodeService(cs.nodeService)
			} else {
				connData.ConnState = CS_RECONNECTING
			}
		default:
		}
	}
}

func (cs *NetClients) processTempConnect() {
	if cs.tempConnMap == nil || len(cs.tempConnMap) == 0 {
		return
	}
	for _, connData := range cs.tempConnMap {
		// 开启创建客户端连接，并连接相应的服务器
		connData.ConnState = CS_CONNECTING
		netClient := NewNetClient(cs.nodeService)
		connData.netService = netClient
		if ret := netClient.Connect(connData.IP, int(connData.Port)); ret {
			connData.ConnState = CS_CONNECTED
			netClient.SetNodeService(cs.nodeService)
		} else {
			connData.ConnState = CS_DISCONNECTED
		}
		// 将创建的连接挂到已连接列表中
		cs.connMap.Add(connData)
	}
	// 清空临时临时连接表
	cs.tempConnMap = nil
}

func (cs *NetClients) ResetConnState(nodeID int32, state ConnState) {
	LogInfo("net client reset conn state....")
	if connData, ok := cs.connMap[nodeID]; ok {
		connData.ConnState = state
	}
}
