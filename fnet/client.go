package fnet

type IClient interface {
}

type Client struct {
	netService INetService
}
