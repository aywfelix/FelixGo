package event

type EventType uint32

const (
	ET_LOGIN  EventType = 1 // 玩家登录进入游戏
	ET_LOGOUT EventType = 2 // 玩家退出游戏
)
