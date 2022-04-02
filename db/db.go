package db

import (
	"time"

	. "github.com/aywfelix/felixgo/configure"
	. "github.com/aywfelix/felixgo/logger"
	. "github.com/aywfelix/felixgo/thread"
)

var MysqlHelper *Mysql
var RedisHelper *Redis

func StartMysql(cfg MysqlConf) error {
	if MysqlHelper != nil {
		return nil
	}

	MysqlHelper = NewMysql()
	return MysqlHelper.Connect(cfg.User, cfg.Password, cfg.IP, cfg.Port, cfg.DataBase)
}

func StopMysql() {
	MysqlHelper.Close()
}

func StartRedis(cfg RedisConf) {
	if RedisHelper != nil {
		return
	}
	RedisHelper = NewRedis()
	RedisHelper.InitConnect(cfg.IP, cfg.Port, cfg.Password, cfg.DB)
}

func StopRedis() {
	RedisHelper.Close()
}

func DebugDB() {
	TimeLoop(time.Second*5, func() {
		if RedisHelper != nil {
			redisIdle := RedisHelper.GetIdleCount()
			redisActive := RedisHelper.GetActiveCount()
			LogDebug("db connection stats, redis: %d  %d", redisIdle, redisActive)
		}

		if MysqlHelper != nil {
			mysqlIdle, mysqlActive := -1, -1
			myStats := MysqlHelper.Stats()
			if myStats != nil {
				mysqlIdle = myStats.Idle
				mysqlActive = myStats.InUse
			}
			LogDebug("db connection stats, mysql: %d  %d", mysqlIdle, mysqlActive)
		}
	})
}
