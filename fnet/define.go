package fnet

type ConnState uint32

const (
	CS_NONE         ConnState = 0
	CS_CONNECTING   ConnState = 1
	CS_CONNECTED    ConnState = 2
	CS_DISCONNECTED ConnState = 3
	CS_RECONNECTING ConnState = 4
	CS_RECONNECTED  ConnState = 5
)

type ServState uint32

const (
	SS_NONE  ServState = 0
	SS_START ServState = 1
	SS_STOP  ServState = 2
	SS_FATAL ServState = 3
)

type NetState uint32

const (
	NS_NONE  NetState = 0
	NS_START NetState = 1
	NS_STOP  NetState = 2
)

const MSG_CHAN_BUFF_LEN = 10 * 1024

type MsgType uint8

const (
	MT_OTHER MsgType = 0
	MT_PROTO MsgType = 1
	MT_JSON  MsgType = 2
	MT_HTTP  MsgType = 3
)

const DEFAULT_HEADER_LEN uint32 = 9
const MAX_PACKAGE_SIZE = 32 * 1024
const MAX_MSG_WORKER_SIZE = 100
const MSG_WORKER_POOL_SIZE = 1000

type MsgErrCode int64

const (
	MEC_OK    MsgErrCode = 0
	MEC_ERR   MsgErrCode = 1
	MEC_MSGID MsgErrCode = 2
)

type NetEventType uint32

const (
	NET_EOF       NetEventType = 0
	NET_ERROR     NetEventType = 1
	NET_TIMEOUT   NetEventType = 2
	NET_CONNECTED NetEventType = 3
)
