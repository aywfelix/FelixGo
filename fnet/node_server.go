package fnet

import (
	"github.com/aywfelix/felixgo/configure"
	. "github.com/aywfelix/felixgo/logger"
)

type INodeServer interface {
}

type NodeServer struct {
	netService INetService
	NodeService
}

func NewNodeServer(config *configure.NetNode) *NodeServer {
	s := new(NodeServer)
	s.msgHandler = NewMsgHandler()
	s.dataPack = NewDataPack()
	s.nodeConfig = config // 本服务器节点配置信息
	ip := config.NodeIP
	port := config.NodePort
	s.netService = NewNetServer(s, ip, port)
	s.isStopped = false
	// set default function
	s.onStarted = s.onStart
	s.onStopped = s.onStop
	s.onConnected = s.onSessionConnect
	s.onDisconnected = s.onSessionDisconnect
	return s
}

func (s *NodeServer) Start() {
	LogInfo("node server start...")
	if s.onStarted != nil {
		s.onStarted()
	}
	// 开启网络监听
	netServer := s.netService.(INetServer)
	if err := netServer.Listen(); err != nil {
		LogError("Error listening, %v", err)
		s.Stop()
		return
	}
	// 注册路由服务
	LogInfo("node server started...")
}

func (s *NodeServer) Serve() {}

func (s *NodeServer) Stop() {
	LogInfo("node server stop ...")
	s.isStopped = true
	if s.onStopped != nil {
		s.onStopped()
	}
	LogInfo("node server stopped...")
}

func (s *NodeServer) onStart(args ...interface{}) {
	LogInfo("node server on start")
}

func (s *NodeServer) onStop(args ...interface{}) {
	LogInfo("node server on stop")
}

func (s *NodeServer) onSessionConnect(args ...interface{}) {
	LogInfo("node server on session connect")
}

func (s *NodeServer) onSessionDisconnect(args ...interface{}) {
	LogInfo("node server on session disconnect")
}
