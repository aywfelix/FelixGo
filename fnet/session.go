package fnet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"

	. "github.com/aywfelix/felixgo/logger"
	. "github.com/aywfelix/felixgo/thread"
	"github.com/golang/protobuf/proto"
)

type ISession interface {
	Start()
	Stop()
	GetID() int64
	Context() context.Context

	GetConn() net.Conn  // 获取原始连接
	RemoteAddr() string // 获取远程客户端地址
	LocalAddr() string  //

	SendProto(msgID uint32, msg *proto.Message) error
	SendJson(msgID uint32, msg interface{}) error

	SendProtoBuffer(msgID uint32, data []byte) error
	SendJsonBuffer(msgID uint32, data []byte) error

	SetProperty(key string, value interface{})
	GetProperty(key string) interface{}
	RemoveProperty(key string)

	SetNetService(service INetService)
	GetNetService() INetService

	GetNodeService() INodeService
}

type Session struct {
	connID       int64
	msgChan      chan []byte
	msgBuffChan  chan []byte
	property     map[string]interface{}
	propertyLock sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	msgHandler   IMsgHandler
	dataPack     IDataPack

	sync.RWMutex
	socket         ISocket
	onConnected    FuncArgs
	onDisconnected FuncArgs
	// session 属于哪个服务器
	netService INetService
}

func NewSession(connID int64, socket ISocket) *Session {
	session := &Session{
		connID:         connID,
		msgChan:        make(chan []byte),
		msgBuffChan:    make(chan []byte, MSG_CHAN_BUFF_LEN),
		property:       nil,
		socket:         socket,
		onConnected:    nil,
		onDisconnected: nil,
		msgHandler:     nil,
	}

	return session
}

func (s *Session) Start() {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	go s.startReader()
	go s.startWriter()
	// 执行钩子方法
	if s.onDisconnected != nil {
		s.onConnected(s)
	}
}

func (s *Session) Stop() {
	if s.onDisconnected != nil {
		s.onDisconnected(s)
	}
	s.cancel()
}

func (s *Session) GetID() int64 {
	return s.connID
}

func (s *Session) Context() context.Context {
	return s.ctx
}

func (s *Session) GetConn() net.Conn {
	return s.socket.GetConn()
}

func (s *Session) RemoteAddr() string {
	return s.socket.RemoteAddr()
}

func (s *Session) LocalAddr() string {
	return s.socket.LocalAddr()
}

func (s *Session) SendProto(msgID uint32, msg *proto.Message) error {
	data, err := proto.Marshal(*msg)
	if err != nil {
		return err
	}
	return s.sendMsg(msgID, uint8(MT_PROTO), data)
}

func (s *Session) SendJson(msgID uint32, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return s.sendMsg(msgID, uint8(MT_JSON), data)
}

func (s *Session) sendMsg(msgID uint32, msgType uint8, data []byte) error {
	s.RLock()
	defer s.RUnlock()
	msg, err := s.dataPack.Pack(NewMessage(msgID, msgType, data))
	if err != nil {
		return fmt.Errorf("pack msg err, msgid=%d", msgID)
	}
	s.msgChan <- msg
	return nil
}

func (s *Session) SendProtoBuffer(msgID uint32, data []byte) error {
	return s.sendBufMsg(msgID, uint8(MT_PROTO), data)
}
func (s *Session) SendJsonBuffer(msgID uint32, data []byte) error {
	return s.sendBufMsg(msgID, uint8(MT_JSON), data)
}

func (s *Session) sendBufMsg(msgID uint32, msgType uint8, data []byte) error {
	s.RLock()
	defer s.RUnlock()
	msg, err := s.dataPack.Pack(NewMessage(msgID, msgType, data))
	if err != nil {
		return errors.New(fmt.Sprintf("pack msg err, msgid=%d", msgID))
	}
	s.msgBuffChan <- msg
	return nil
}

func (s *Session) SetProperty(key string, value interface{}) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.property[key] = value
}

func (s *Session) GetProperty(key string) interface{} {
	if value, ok := s.property[key]; ok {
		return value
	}
	return nil
}

func (s *Session) RemoveProperty(key string) {
	if _, ok := s.property[key]; ok {
		delete(s.property, key)
	}
}

// 读消息Goroutine，用于从客户端中读取数据
func (s *Session) startReader() {
	LogInfo("start conn reader goroutine...")
	defer LogInfo("reader goroutine exit! remote %s, local %s", s.RemoteAddr(), s.LocalAddr())
	defer s.Stop()

	for {
		select {
		case <-s.ctx.Done():
			LogInfo("session closed ...")
			return
		default:
			// 读取消息头部
			headerData := make([]byte, s.dataPack.GetHeadLen())
			if _, err := s.socket.Recv(headerData); err != nil {
				LogError("read msg err, %v", err.Error())
				return
			}

			msg, err := s.dataPack.UnPack(headerData)
			if err != nil {
				LogError("unpack msg err, %v", err.Error())
				return
			}
			//得到当前客户端请求的数据
			request := NewRequest(s, msg)
			if s.msgHandler.IsUseWorkPool() {
				s.msgHandler.DispatchByMsgID(request)
			} else {
				s.msgHandler.DoMsg(request)
			}
		}
	}
}

// 写消息Goroutine, 用户将数据发送给客户端
func (s *Session) startWriter() {
	LogInfo("start conn writer goroutine...")
	defer LogInfo("writer goroutine exit! remote %s, local %s", s.RemoteAddr(), s.LocalAddr())
	defer s.Stop()
	for {
		select {
		case data, ok := <-s.msgChan:
			if ok {
				if _, err := s.socket.Send(data); err != nil {
					LogError("send data to client err, %v", err.Error())
					return
				}
				LogInfo("send data to client ok, %v", data)
			} else {
				LogInfo("msgChan closed")
				break
			}
		case data, ok := <-s.msgBuffChan:
			if ok {
				if _, err := s.socket.Send(data); err != nil {
					LogError("send data to client err, %v", err.Error())
					return
				}
				LogInfo("send data to client ok, %v", data)
			} else {
				LogInfo("msgBuffChan closed")
				break
			}
		case <-s.ctx.Done():
			LogInfo("session closed ...")
			return
		}
	}
}

func (s *Session) SetNetService(service INetService) {
	s.netService = service
	nodeService := service.GetNodeService()
	s.dataPack = service.GetNodeService().GetDataPack()
	s.msgHandler = service.GetNodeService().GetMsgHandler()
	s.onConnected = nodeService.GetOnConnected()
	s.onDisconnected = nodeService.GetOnDisconnected()
}

func (s *Session) GetNetService() INetService {
	return s.netService
}

func (s *Session) GetNodeService() INodeService {
	return s.netService.GetNodeService()
}
