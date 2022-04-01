package fnet

import (
	"fmt"
	"net"
	
	. "github.com/felix/felixgo/logger"
	. "github.com/felix/felixgo/utils"
)

type ISocket interface {
	Close()
	Recv(data []byte) (n int, err error)
	Send(data []byte) (int, error)
	GetConn() net.Conn
	LocalAddr() string
	RemoteAddr() string
}

type ISSocket interface {
	Listen() error
}
type ICSocket interface {
	Connect(ipVer, ip string, port int) (error, ISession)
}

type Socket struct {
	Closed    bool
	IP        string
	Port      int
	IPVersion string

	conn       net.Conn
	localAddr  string
	remoteAddr string

	// 所属的哪个服务器器端
	netService INetService
}

func NewSocket(ip string, port int, conn net.Conn) *Socket {
	return &Socket{
		IP:        ip,
		Port:      port,
		IPVersion: "tcp",
		Closed:    false,
		conn:      conn,
	}
}

func (s *Socket) Close() {
	s.Closed = true
	s.conn.Close()
	LogInfo("close connection, local:", s.conn.LocalAddr().String(), " remote:", s.conn.RemoteAddr().String())
}

func (s *Socket) Recv(data []byte) (n int, err error) {
	return s.conn.Read(data)
}

func (s *Socket) Send(data []byte) (n int, err error) {
	return s.conn.Write(data)
}

func (s *Socket) GetConn() net.Conn {
	return s.conn
}

func (s *Socket) LocalAddr() string {
	return s.conn.LocalAddr().String()
}

func (s *Socket) RemoteAddr() string {
	return s.conn.RemoteAddr().String()
}

func (s *Socket) setSocketOpts() {

}

type SSocket struct {
	listener net.Listener
	Socket
}

func NewSSocket(netService INetService, ip string, port int) *SSocket {
	ss := &SSocket{
		listener: nil,
		Socket: Socket{
			IP:        ip,
			Port:      port,
			IPVersion: "tcp",
			Closed:    false,
			conn:      nil,
		},
	}
	ss.netService = netService
	return ss
}

func (ss *SSocket) Listen() error {
	address := fmt.Sprintf("%s:%d", ss.IP, ss.Port)
	addr, err := net.ResolveTCPAddr(ss.IPVersion, address)
	if err != nil {
		return err
	}
	ss.listener, err = net.ListenTCP(ss.IPVersion, addr)
	if err != nil {
		return err
	}
	LogInfo("server socket listen port:%d\n", ss.Port)
	go func() {
		defer func() { // server quit
			err := recover()
			if err != nil {
				LogError("server some wrong, err:", err)
			}
			ss.listener.Close()
			LogError("server stopped, server addr:", ss.listener.Addr().String())
		}()
		for {
			if ss.Closed {
				break
			}
			conn, err := ss.listener.Accept()
			if err != nil {
				LogError("server accept err:", err.Error())
				continue
			}
			// TODO：判断最大连接数
			LogInfo("server accept client, remote address=%s", conn.RemoteAddr().String())

			socket := NewSocket(ss.IP, ss.Port, conn)
			// 将连接保存到会话中
			session := NewSession(SnowFlake.GenInt(), socket)
			session.SetNetService(ss.netService)
			// TODO:每次连接创建一个协程，可以优化使用协程池处理
			go func() {
				defer func() {
					err := recover()
					if err != nil {
						LogError("csocket some wrong, err:", err)
						// 出错，关闭连接
						ss.Close()
						// 关闭session
						session.Stop()
					} else {
						// 将session加入到管理器中
						netServer := ss.netService.(INetServer)
						sessionMgr := netServer.GetSessionManager()
						sessionMgr.Add(session)
					}
				}()
				session.Start()
			}()
		}
	}()
	return nil
}

type CSocket struct {
	Socket
}

func NewCSocket(netService INetService) *CSocket {
	cs := &CSocket{
		Socket: Socket{},
	}
	cs.netService = netService
	return cs
}

func (cs *CSocket) Connect(ipVer, ip string, port int) (error, ISession) {
	cs.IP = ip
	cs.Port = port
	cs.IPVersion = ipVer
	// 连接服务器端
	address := fmt.Sprintf("%s:%d", cs.IP, cs.Port)
	conn, err := net.Dial(cs.IPVersion, address)
	if err != nil {
		LogError("client connect server(%s) failed, %s", address, err.Error())
		return err, nil
	}
	cs.conn = conn
	// 将连接保存到会话中
	session := NewSession(SnowFlake.GenInt(), cs) // TODO: 优化session 复用处理
	session.SetNetService(cs.netService)
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				LogError("csocket some wrong, err: %v", err)
				// 出错，关闭连接
				cs.Close()
				// 关闭session
				session.Stop()
			}
		}()
		session.Start()
		LogInfo("client connect server ok... %s", cs.conn.RemoteAddr().String())
	}()
	return nil, session
}
