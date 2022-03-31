package fnet

// 抽象网络管理器-客户端
type INetClient interface {
	GetSession() ISession
	Connect(ip string, port int) bool
}

type NetClient struct {
	NetService
	// 一个客户端对应一个会话
	session ISession
}

func NewNetClient(service INodeService) *NetClient {
	c := &NetClient{}
	c.socket = NewCSocket(c)
	c.nodeService = service
	return c
}

func (c *NetClient) Connect(ip string, port int) bool {
	csocket := c.socket.(ICSocket)
	err, _ := csocket.Connect("tcp", ip, port)
	if err == nil {
		return true
	}
	return false
}

func (c *NetClient) GetSocket() ISocket {
	return c.socket
}

func (c *NetClient) GetNodeService() INodeService {
	return c.nodeService
}

func (c *NetClient) SetNodeService(service INodeService) {
	if c.nodeService != nil || service == nil {
		return
	}
	c.nodeService = service
}

func (c *NetClient) GetSession() ISession {
	return c.session
}
