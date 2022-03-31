package db

import (
	"time"
	
	. "github.com/felix/felixgo/global"
	. "github.com/felix/felixgo/thread"
	. "github.com/felix/felixgo/logger"
)

var DbMysql *Mysql
var DbRedis *Redis

func StartMysql() error {
	if DbMysql != nil || ServerConf == nil {
		return nil
	}
	DbMysql = NewMysql()
	return DbMysql.Connect(ServerConf.Mysql.User,
		ServerConf.Mysql.Password,
		ServerConf.Mysql.Host,
		ServerConf.Mysql.Port,
		ServerConf.Mysql.DataBase)
}

func StopMysql() {
	DbMysql.Close()
}

func StartRedis() {
	if DbRedis != nil {
		return
	}
	DbRedis = NewRedis()
	DbRedis.InitConnect(ServerConf.Redis.Host,
		ServerConf.Redis.Port,
		ServerConf.Redis.Password,
		ServerConf.Redis.DB)
}

func StopRedis() {
	DbRedis.Close()
}

func DebugDB() {
	TimeLoop(time.Second*5, func() {
		if DbRedis != nil {
			redisIdle := DbRedis.GetIdleCount()
			redisActive := DbRedis.GetActiveCount()
			LogDebug("db connection stats, redis: %d  %d", redisIdle, redisActive)
		}

		if DbMysql != nil {
			mysqlIdle, mysqlActive := -1, -1
			myStats := DbMysql.Stats()
			if myStats != nil {
				mysqlIdle = myStats.Idle
				mysqlActive = myStats.InUse
			}
			LogDebug("db connection stats, mysql: %d  %d", mysqlIdle, mysqlActive)
		}
	})
}
