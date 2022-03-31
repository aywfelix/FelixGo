package global

// 用于定义常量
const SIGHT float64 = 4.0        // 玩家视野宽度 默认玩家与怪物视野一致
const SCREEN_SIGHT float64 = 4.0 // 默认玩家视野看到与屏幕显示一致

type ServerType int32

const (
	SERVER_TYPE_NONE   ServerType = 0
	SERVER_TYPE_MASTER ServerType = 1
	SERVER_TYPE_GAME   ServerType = 2
	SERVER_TYPE_LOGIN  ServerType = 3
	SERVER_TYPE_WORLD  ServerType = 4
	SERVER_TYPE_GATE   ServerType = 5
	SERVER_TYPE_CHAT   ServerType = 6
	SERVER_TYPE_DB     ServerType = 7
	SERVER_TYPE_PLAYER ServerType = 8
)

type ServiceType uint32

const (
	ST_SERVER ServiceType = 0
	ST_CLIENT ServiceType = 1
)
