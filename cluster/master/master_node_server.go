package cluster

// master 用于服务器发现功能
// 1、所有服务器启动时先连接master，如果连接不上，则终止程序。
// 2、其他服务器连接上master，优先上报个人服务器信息，master保留连接上服务器信息
// 3、master将保存的服务器信息，同步给所有已经连接上服务器

import (
	"fmt"

	. "github.com/aywfelix/felixgo/configure"
	. "github.com/aywfelix/felixgo/fnet"
	pb "github.com/aywfelix/felixgo/fnet/proto"
	"github.com/golang/protobuf/proto"
)

type ServerInfoMap map[int32]*pb.ServerReport

type IMasterNodeServer interface {
	// 保存服务器信息
	SaveServer(serverInfo *pb.ServerReport)
	SyncServer(session ISession)

	InitNetServer(ip string, port int)
}

type MasterNodeServer struct {
	serverMap ServerInfoMap
	// // 网络服务
	// netService INetService
	// 消息处理
	msgHandler IMsgHandler
	// 继承自 nodeservice
	INodeService
}

func NewMasterNodeServer() *MasterNodeServer {
	s := &MasterNodeServer{
		serverMap:  make(ServerInfoMap),
		msgHandler: NewMsgHandler(),
	}
	s.INodeService = NewNodeServer(&MasterCfg.NetNode)
	s.AddRouter()
	return s
}

func (s *MasterNodeServer) InitNetServer(ip string, port int) {
	s.netService = NewNetServer(s, ip, port)
}

func (s *MasterNodeServer) AddRouter() {
	s.msgHandler.AddRouter(uint32(pb.ServerNodeMsgID_REPORT_CLIENT_INFO_TO_SERVER), &ServerRouter{})
}

// define master router
type ServerRouter struct {
	BaseRouter
}

func (router *ServerRouter) Handle(request IRequest) MsgErrCode {
	fmt.Println("master receive msg from server")
	msgID := request.GetMsgID()
	if msgID != uint32(pb.ServerNodeMsgID_REPORT_CLIENT_INFO_TO_SERVER) {
		fmt.Println("receive err msg, msgID=", msgID)
		return MEC_MSGID
	}
	// 保存信息
	iSession := request.GetSession()
	iNodeService := iSession.GetNodeService()
	master := iNodeService.(*MasterNodeServer)
	dataBytes := request.GetData()
	reportInfo := &pb.ServerReport{}
	if err := proto.Unmarshal(dataBytes, reportInfo); err != nil {
		fmt.Println("proto unmarshal msg error, err: %+v", err)
		return MEC_MSG_UNMARSHAL
	}
	master.SaveServer(reportInfo)
	master.SyncServer(iSession)
	return MEC_OK
}

func (s *MasterNodeServer) SaveServer(serverInfo *pb.ServerReport) {
	if _, ok := s.serverMap[serverInfo.ServerId]; ok {
		return
	}
	// 保存服务器信息
	s.serverMap[serverInfo.ServerId] = serverInfo
}

func (s *MasterNodeServer) SyncServer(session ISession) {
	// 1、将其他服务器信息同步给此会话客户端
	// 2、将此会话服务器信息同步给其他服务器
	msgID := uint32(pb.ServerNodeMsgID_MASTER_REPORT_SERVER_INFO_TO_SERVER)
	for _, serverInfo := range s.serverMap {
		if dataBytes, err := proto.Marshal(serverInfo); err == nil {
			session.SendProtoBuffer(msgID, dataBytes)
		}
	}
}
