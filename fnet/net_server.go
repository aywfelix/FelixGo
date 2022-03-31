package fnet

// 抽象网络管理器-服务器端
type INetServer interface {
	GetSessionManager() ISessionManager
	GetSession(sessionID int64) (ISession, error)
	Listen() error
}

type NetServer struct {
	NetService
	sessionMgr ISessionManager
}

func NewNetServer(nodeService INodeService, servIP string, servPort int) *NetServer {
	s := &NetServer{}
	s.nodeService = nodeService
	s.sessionMgr = NewSessionManager()
	s.socket = NewSSocket(s, servIP, servPort)
	return s
}

func (s *NetServer) Listen() error {
	ssocket := s.socket.(ISSocket)
	return ssocket.Listen()
}

func (s *NetServer) GetSocket() ISocket {
	return s.socket
}

func (s *NetServer) GetNodeService() INodeService {
	return s.nodeService
}

func (s *NetServer) SetNodeService(service INodeService) {
	s.nodeService = service
}

func (s *NetServer) GetSessionManager() ISessionManager {
	return s.sessionMgr
}

func (s *NetServer) GetSession(sessionID int64) (ISession, error) {
	return s.sessionMgr.GetSession(sessionID)
}
