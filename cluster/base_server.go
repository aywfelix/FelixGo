package cluster

import (
	"context"
	"fmt"

	. "github.com/aywfelix/felixgo/fnet"
	pb "github.com/aywfelix/felixgo/fnet/proto"
	. "github.com/aywfelix/felixgo/logger"

	. "github.com/aywfelix/felixgo/common"
)

type IBaseServer interface {
	Start()
	Serve()
	Stop()

	Context() context.Context

	ReadConfig() // 读取配置：读取redis mysql配置
	StartLog()   // 启动日志管理

	SetOnStart(onStart func() bool)
	SetOnStop(onStop func() bool)

	// 设置服务器信息，用于广播给其他节点
	SetServerInfo(id int32, name string, ip string, port int32, maxConn int32, state ServerState, serverType int32)
	GetServerInfo() pb.ServerReport
	AddRouter(msgID int32, router IRouter)

	// 指令打印服务器信息
	PrintServerInfo()
}

// GFNodeServer 包含服务器端和客户端
type BaseServer struct {
	// 服务器启动关闭
	onStart func() bool
	onStop  func() bool
	// 本服务器详细信息
	ServerInfo pb.ServerReport
	// ctx
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewBaseServer() *BaseServer {
	return &BaseServer{
		onStart:    nil,
		onStop:     nil,
		ServerInfo: pb.ServerReport{},
	}
}

// 启动服务器
func (s *BaseServer) Start() {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	if s.onStart != nil {
		s.onStart()
	}
}

func (s *BaseServer) Context() context.Context{
	return s.ctx
}

// 关闭服务器
func (s *BaseServer) Stop() {
	s.cancel()
}

// 运行服务器执行
func (s *BaseServer) Serve() {
	for {
		select {
		case <-s.ctx.Done():
			// 服务器退出
			if s.onStop != nil {
				s.onStop()
			}
			return
		}
	}
}

// 读取配置
func (s *BaseServer) ReadConfig() {

}

//
func (s *BaseServer) SetOnStart(onStart func() bool) {
	if onStart != nil {
		s.onStart = onStart
	}
}

func (s *BaseServer) SetOnStop(onStop func() bool) {
	if onStop != nil {
		s.onStop = onStop
	}
}

func (s *BaseServer) StartLog() {

}

func (s *BaseServer) SetServerInfo(id int32, name string, ip string, port int32, maxConn int32, state ServerState, serverType int32) {
	s.ServerInfo.ServerId = id
	s.ServerInfo.ServerName = []byte(name)
	s.ServerInfo.ServerIp = []byte(ip)
	s.ServerInfo.ServerPort = port
	s.ServerInfo.MaxOnline = maxConn
	s.ServerInfo.CurOnline = 0
	s.ServerInfo.ServerState = int32(state)
	s.ServerInfo.ServerType = serverType
}

func (s *BaseServer) GetServerInfo() pb.ServerReport {
	return s.ServerInfo
}

func (s *BaseServer) AddRouter(msgID int32, router IRouter) {

}

func (s *BaseServer) PrintServerInfo() {
	serverInfo := fmt.Sprintf("ServerInfo:\n%#v", s.ServerInfo)
	LogInfo(serverInfo)
}
