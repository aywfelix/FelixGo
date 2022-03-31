package fnet

// netservice 接口
type INetService interface {
	GetSocket() ISocket
	SetNodeService(service INodeService)
	GetNodeService() INodeService
}

type NetService struct {
	socket      ISocket
	nodeService INodeService
}
