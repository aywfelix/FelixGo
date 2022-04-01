package common

// import (
// 	. "github.com/felix/felixgo/logger"
// 	. "github.com/felix/felixgo/utils"
// )

// type NodeConfig struct {
// 	NodeId     int    `json:"nodeId, omitempty"`
// 	NodeName   string `json:"nodeName, omitempty"`
// 	NodeIP     string `json:"nodeIP, omitempty"`
// 	NodePort   int    `json:"nodePort, omitempty"`
// 	MaxConnect int    `json:"max, omitempty"`
// }

// type GlobalConfig struct {
// 	Project   string `json:"project, omitempty"`
// 	MaxOnLine int    `json:"maxOnLine, omitempty"`
// 	Dev       int    `json:"dev, omitempty"`
// }

// type MysqlConfig struct {
// 	Host     string `json:"host, omitempty"`
// 	Port     int    `json:"port, omitempty"`
// 	User     string `json:"user, omitempty"`
// 	Password string `json:"password, omitempty"`
// 	DataBase string `json:"dataBase, omitempty"`
// }

// type RedisConfig struct {
// 	Host     string `json:"host, omitempty"`
// 	Port     int    `json:"port, omitempty"`
// 	Password string `json:"password, omitempty"`
// 	DB       int    `json:"db, omitempty"`
// }

// type ServerConfig struct {
// 	MasterServer   NodeConfig
// 	WorldServer    NodeConfig
// 	LoginServer    NodeConfig
// 	GateServer     NodeConfig
// 	GameServer     NodeConfig
// 	ChatServer     NodeConfig
// 	GateUserServer NodeConfig

// 	Mysql MysqlConfig
// 	Redis RedisConfig

// 	Global GlobalConfig
// }

// func (c *ServerConfig) GetNodeConfig(serverType ServerType) *NodeConfig {
// 	switch serverType {
// 	case SERVER_TYPE_MASTER:
// 		return &c.MasterServer
// 	case SERVER_TYPE_GAME:
// 		return &c.GameServer
// 	case SERVER_TYPE_LOGIN:
// 		return &c.LoginServer
// 	case SERVER_TYPE_WORLD:
// 		return &c.WorldServer
// 	case SERVER_TYPE_GATE:
// 		return &c.GateServer
// 	case SERVER_TYPE_CHAT:
// 		return &c.ChatServer
// 	}
// 	return nil
// }

// func (c *ServerConfig) GetMysql() MysqlConfig {
// 	return c.Mysql
// }

// func (c *ServerConfig) GetRedis() RedisConfig {
// 	return c.Redis
// }

// func (c *ServerConfig) GetGlobal() GlobalConfig {
// 	return c.Global
// }

// //===================================================
// var ServerConf *ServerConfig

// func LoadConfig(configPath string) error {
// 	if err := Configure.Load(configPath); err != nil {
// 		LogError("load config file error:", err)
// 		return err
// 	}
// 	ServerConf = &ServerConfig{}
// 	if err := Configure.AssignStruct(ServerConf); err != nil {
// 		LogError("parse config file failed:", err)
// 		return err
// 	}
// 	return nil
// }

// func RedisKeyPrefix() string {
// 	key := ""
// 	globalCfg := ServerConf.Global
// 	key += globalCfg.Project
// 	if globalCfg.Dev == 1 {
// 		key += ".dev."
// 	}
// 	return key
// }

// func init() {
// 	LoadConfig("./config/server.json")
// }
