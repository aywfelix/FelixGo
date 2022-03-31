package utils

// import (
// 	"encoding/json"
// 	"io/ioutil"
// 	"os"
// )
//
// type ServerCfg struct {
// 	Name       string
// 	ServerPort int
// }
//
// type MysqlCfg struct {
// 	Host     string
// 	Port     int
// 	User     string
// 	Password string
// 	DataBase string
// }
//
// type RedisCfg struct {
// 	Host     string
// 	Port     int
// 	Password string
// }
//
// type Config struct {
// 	// server
// 	Server ServerCfg
// 	// mysql
// 	Mysql MysqlCfg
// 	// redis
// 	Redis RedisCfg
// 	// config
// 	ConfigPath string
// }
//
// func (g *Config) LoadConfig() error {
// 	if isExists, err := IsFileExist(g.ConfigPath); !isExists || err != nil {
// 		return err
// 	}
//
// 	data, err := ioutil.ReadFile(g.ConfigPath)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = json.Unmarshal(data, g)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// var GConfig *Config
//
// func init() {
// 	pwd, err := os.Getwd()
// 	if err != nil {
// 		pwd = "."
// 	}
// 	GConfig = &Config{
// 		Server: ServerCfg{Name: "localhost"},
// 		Mysql: MysqlCfg{
// 			Host:     "localhost",
// 			Port:     3306,
// 			User:     "root",
// 			Password: "123456",
// 			DataBase: "mysql",
// 		},
// 		Redis: RedisCfg{
// 			Host:     "localhost",
// 			Port:     6379,
// 			Password: "123456",
// 		},
// 		ConfigPath: pwd + "/src/utils/server.json",
// 	}
// }
