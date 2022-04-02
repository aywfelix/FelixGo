package configure

type MysqlConf struct {
	IP     string
	Port     int
	User     string
	Password string
	DataBase string
}

type RedisConf struct {
	IP     string
	Port     int
	Password string
	DB       int
}


type NetNode struct {
	NodeId     int
	NodeName   string
	NodeIP     string
	NodePort   int
	MaxConnect int
}

type GateNode struct {
	NetNode
}

type MasterNode struct {
	NetNode
}

type WorldNode struct {
	NetNode
}

type LoginNode struct {
	NetNode
}

type ChatNode struct {
	NetNode
}

type GameNode struct {
	NetNode
}

type GateUserNode struct {
	NetNode
}

type GlobalConf struct {
	MaxOnLine int
	Dev       int
}

type LogConf struct {
	Level int
	Path string
	RollType int
	RollTime string
	RollSize int
}