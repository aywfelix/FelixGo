package configure

import "errors"

var (
	errPathNotExist = errors.New("config path is not exit")
)

var JsonConfig *jsonConfig
var IniConfig *iniConfig

// 定义需要的配置节点
var (
	MasterCfg *MasterNode
	GateCfg *GateNode
	GameCfg *GameNode
	WorldCfg *WorldNode
	LoginCfg *LoginNode
	ChatCfg *ChatNode
	GateUserCfg *GateUserNode
	MysqlCfg *MysqlConf
	RedisCfg *RedisConf
	LogCfg *LogConf
	GlobalCfg *GlobalConf
)

func LoadIniConfig(cfg string) error {
	if err := IniConfig.Load(cfg); err != nil {
		return err
	}
	
	MasterCfg = &MasterNode{
		NetNode: NetNode{
			NodeId: IniConfig.GetInt("MasterServer", "NodeId"),
			NodeName: IniConfig.GetString("MasterServer", "NodeName"),
			NodeIP: IniConfig.GetString("MasterServer", "NodeIP"),
			NodePort: IniConfig.GetInt("MasterServer", "NodePort"),
			MaxConnect: IniConfig.GetInt("MasterServer", "MaxConnect"),
		},
	}

	GateCfg = &GateNode {
		NetNode: NetNode{
			NodeId:IniConfig.GetInt("GateServer", "NodeId"),
			NodeName: IniConfig.GetString("GateServer", "NodeName"),
			NodeIP: IniConfig.GetString("GateServer", "NodeIP"),
			NodePort: IniConfig.GetInt("GateServer", "NodePort"),
			MaxConnect: IniConfig.GetInt("GateServer", "MaxConnect"),
		},
	}

	GameCfg = &GameNode {
		NetNode: NetNode{
			NodeId:IniConfig.GetInt("GameServer", "NodeId"),
			NodeName: IniConfig.GetString("GameServer", "NodeName"),
			NodeIP: IniConfig.GetString("GameServer", "NodeIP"),
			NodePort: IniConfig.GetInt("GameServer", "NodePort"),
			MaxConnect: IniConfig.GetInt("GameServer", "MaxConnect"),
		},
	}

	WorldCfg = &WorldNode {
		NetNode: NetNode{
			NodeId:IniConfig.GetInt("WorldServer", "NodeId"),
			NodeName: IniConfig.GetString("WorldServer", "NodeName"),
			NodeIP: IniConfig.GetString("WorldServer", "NodeIP"),
			NodePort: IniConfig.GetInt("WorldServer", "NodePort"),
			MaxConnect: IniConfig.GetInt("WorldServer", "MaxConnect"),
		},
	}

	LoginCfg = &LoginNode {
		NetNode: NetNode{
			NodeId:IniConfig.GetInt("LoginServer", "NodeId"),
			NodeName: IniConfig.GetString("LoginServer", "NodeName"),
			NodeIP: IniConfig.GetString("LoginServer", "NodeIP"),
			NodePort: IniConfig.GetInt("LoginServer", "NodePort"),
			MaxConnect: IniConfig.GetInt("LoginServer", "MaxConnect"),
		},
	}

	ChatCfg = &ChatNode {
		NetNode: NetNode{
			NodeId:IniConfig.GetInt("ChatServer", "NodeId"),
			NodeName: IniConfig.GetString("ChatServer", "NodeName"),
			NodeIP: IniConfig.GetString("ChatServer", "NodeIP"),
			NodePort: IniConfig.GetInt("ChatServer", "NodePort"),
			MaxConnect: IniConfig.GetInt("ChatServer", "MaxConnect"),
		},
	}

	GateUserCfg = &GateUserNode {
		NetNode: NetNode{
			NodeId:IniConfig.GetInt("GateUserServer", "NodeId"),
			NodeName: IniConfig.GetString("GateUserServer", "NodeName"),
			NodeIP: IniConfig.GetString("GateUserServer", "NodeIP"),
			NodePort: IniConfig.GetInt("GateUserServer", "NodePort"),
			MaxConnect: IniConfig.GetInt("GateUserServer", "MaxConnect"),
		},
	}
	
	MysqlCfg = &MysqlConf {
		IP: IniConfig.GetString("MysqlServer", "IP"),
		Port: IniConfig.GetInt("MysqlServer", "Port"),
		User: IniConfig.GetString("MysqlServer", "User"),
		Password: IniConfig.GetString("MysqlServer", "Password"),
		DataBase: IniConfig.GetString("MysqlServer", "DataBase"),
	}

	RedisCfg = &RedisConf {
		IP: IniConfig.GetString("RedisServer", "IP"),
		Port: IniConfig.GetInt("RedisServer", "Port"),
		Password: IniConfig.GetString("RedisServer", "Password"),
		DB: IniConfig.GetInt("RedisServer", "DB"),
	}

	LogCfg = &LogConf {
		Level: IniConfig.GetInt("LogServer", "Level"),
		Path: IniConfig.GetString("LogServer", "Path"),	
		RollType: IniConfig.GetInt("LogServer", "RollType"),
		RollTime: IniConfig.GetString("LogServer", "RollTime"),
		RollSize: IniConfig.GetInt("LogServer", "RollSize"),
	}
	return nil
}

func init() {
	IniConfig = &iniConfig{}
}