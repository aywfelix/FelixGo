package cluster

import (
	. "github.com/aywfelix/felixgo/cluster"
	. "github.com/aywfelix/felixgo/common"
	. "github.com/aywfelix/felixgo/configure"
)

type MasterServer struct {
	masterServer IMasterNodeServer
	IBaseServer
}

func NewMasterServer() *MasterServer {
	master := &MasterServer{
		masterServer: NewMasterNodeServer(),
		IBaseServer:  NewBaseServer(),
	}
	// 获取本服信息
	master.SetServerInfo(
		int32(MasterCfg.NodeId),
		MasterCfg.NodeName,
		MasterCfg.NodeIP,
		int32(MasterCfg.NodePort),
		int32(MasterCfg.MaxConnect),
		SS_NORMAL,
		1)
	master.masterServer.InitNetServer(MasterCfg.NodeIP, MasterCfg.NodePort)
	return master
}
