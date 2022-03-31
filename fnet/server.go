package fnet

// import (
// 	"fmt"
// 	"net"
// )
//
// type IServer interface {
// 	Start()
// 	Serve()
// 	Stop()
//
// 	AddRouter(msgID uint32, router IRouter)
// 	GetConnMgr() IConnManager
// }
//
// type Server struct {
// 	servState  ServState
// 	Sign       string
// 	IPVersion  string
// 	IP         string
// 	Port       int
// 	msgHandler IMsgHandler
// 	ConnMgr    ISessionManager
// 	packet     IDataPack
//
// 	onConnStart func(conn ISession)
// 	onConnStop  func(conn ISession)
// 	netServer   INetServer
// }
//
// func NewServer(sign string, ip string, port int) IServer {
// 	s := &Server{
// 		servState:  SS_NONE,
// 		Sign:       sign,
// 		IPVersion:  "tcp4",
// 		IP:         ip,
// 		Port:       port,
// 		msgHandler: NewMsgHandler(),
// 		ConnMgr:    NewSessionManager(),
// 		packet:     NewDataPack(),
// 	}
// 	return s
// }
//
// func (s *Server) Start() {
// 	fmt.Println(s.Sign + " server start...")
//
// 	go func() {
// 		s.msgHandler.StartWorkPool()
// 		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
// 		if err != nil {
// 			fmt.Println("resolve tcp addr err:", err)
// 			return
// 		}
// 		listener, err := net.ListenTCP(s.IPVersion, addr)
// 		if err != nil {
// 			fmt.Println("listen", s.IPVersion, "err", err)
// 			return
// 		}
//
// 		fmt.Println("start server ", s.Sign, "succ, now listening...")
// 		var cID uint32
// 		cID = 0
//
// 		for {
// 			conn, err := listener.AcceptTCP()
// 			if err != nil {
// 				fmt.Println("Accept err ", err)
// 				continue
// 			}
// 			fmt.Println("Get conn remote addr=", conn.RemoteAddr().String())
// 			if s.ConnMgr.Len() >= 5000 {
// 				conn.Close()
// 				continue
// 			}
//
// 			session := NewSession(s, cID, s.msgHandler)
// 			cID++
// 			go session.Start()
// 		}
// 	}()
// }
//
// func (s *Server) Stop() {
//
// }
//
// func (s *Server) Serve() {
//
// }
//
// func (s *Server) AddRouter(msgID uint32, router IRouter) {
//
// }
//
// func (s *Server) GetConnMgr() IConnManager {
// 	return nil
// }
