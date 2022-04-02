package db

import (
	"errors"
	"fmt"
	"time"

	. "github.com/aywfelix/felixgo/logger"
	redigo "github.com/gomodule/redigo/redis"
)

const (
	REDIS_MAX_IDLE         = 30
	REDIS_MAX_ACTIVE       = 100
	REDIS_MAX_DILE_TIMEOUT = 180 * time.Second
)

var (
	errScriptNotExist = errors.New("script not exit")
	errAssertReturn   = errors.New("assert return value error")
	errReturnNilValue = errors.New("return nil value")
)

type IRedis interface {
	InitConnect(ip string, port int, password string, db int)
	SetScript(key string, script string, keyCount int)
	DoScript(key string, args ...interface{}) (reply interface{}, err error)
	GetIdleCount() int
	GetActiveCount() int
	DoCmd(cmd string, args ...interface{}) (reply interface{}, err error)
	Close()
}

type Redis struct {
	scripts map[string]*redigo.Script
	pool    *redigo.Pool
}

func NewRedis() *Redis {
	return &Redis{
		scripts: make(map[string]*redigo.Script, 0),
		pool:    nil,
	}
}

func (r *Redis) InitConnect(ip string, port int, password string, db int) {
	dsn := fmt.Sprintf("redis://%s:%d", ip, port)
	dial := func() (redigo.Conn, error) {
		conn, err := redigo.DialURL(dsn)
		if err != nil {
			return nil, err
		}
		_, err = conn.Do("AUTH", password)
		if err != nil {
			return nil, err
		}
		_, err = conn.Do("SELECT", db)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}

	onBorrow := func(c redigo.Conn, t time.Time) error {
		_, err := c.Do("PING")
		if err != nil {
			return err
		}
		return nil
	}
	r.pool = &redigo.Pool{
		MaxIdle:      REDIS_MAX_IDLE,
		MaxActive:    REDIS_MAX_ACTIVE,
		IdleTimeout:  REDIS_MAX_DILE_TIMEOUT,
		Dial:         dial,
		TestOnBorrow: onBorrow,
		Wait:         true,
	}
}

func (r *Redis) SetScript(key string, script string, keyCount int) {
	r.scripts[key] = redigo.NewScript(keyCount, script)
}

func (r *Redis) DoScript(key string, args ...interface{}) (reply interface{}, err error) {
	script, ok := r.scripts[key]
	if !ok {
		return nil, errScriptNotExist
	}

	conn := r.tryGetConn()
	defer conn.Close()
	return script.Do(conn, args...)
}

func (r *Redis) DoCmd(cmd string, args ...interface{}) (reply interface{}, err error) {
	conn := r.tryGetConn()
	defer conn.Close()
	return conn.Do(cmd, args...)
}

func (r *Redis) getConn() (redigo.Conn, error) {
	conn := r.pool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}
	return conn, nil
}

func (r *Redis) tryGetConn() redigo.Conn {
	for {
		conn := r.pool.Get()
		if err := conn.Err(); err != nil {
			LogError("redis conn failed, error:", err.Error())
			time.Sleep(time.Second * 1)
			continue
		}
		return conn
	}
}

func (r *Redis) GetIdleCount() int {
	if r.pool == nil {
		return -1
	}
	return r.pool.IdleCount()
}

func (r *Redis) GetActiveCount() int {
	if r.pool == nil {
		return -1
	}
	return r.pool.ActiveCount()
}

func (r *Redis) Close() {
	r.scripts = nil
	r.pool.Close()
	r.pool = nil
}

// region redis 指令操作接口

// region key
func (r *Redis) DelKey(key string) error {
	if _, err := r.DoCmd("DEL", key); err != nil {
		LogError("DEL: key=%v, err=%v", key, err)
		return err
	}
	return nil
}

func (r *Redis) SetKeyExpire(key string, duration int) error {
	if _, err := r.DoCmd("EXPIRE", key, duration); err != nil {
		LogError("EXPIRE: key=%v, err=%v", key, err)
		return err
	}
	return nil
}

func (r *Redis) SetNX(key string, value string) error {
	if _, err := r.DoCmd("SETNX", key, value); err != nil {
		LogError("SETNX: key=%v, err=%v", key, err)
		return err
	}
	return nil
}
func (r *Redis) SetEX(key string, duration int, value string) error {
	if _, err := r.DoCmd("SETEX", key, duration, value); err != nil {
		LogError("SETEX: key=%v, err=%v", key, err)
		return err
	}
	return nil
}

func (r *Redis) Incrby(key string, inc int) error {
	if _, err := r.DoCmd("INCRBY", key, inc); err != nil {
		LogError("INCRBY: key=%v, err=%v", key, err)
		return err
	}
	return nil
}

func (r *Redis) Decrby(key string, dec int) error {
	if _, err := r.DoCmd("DECRBY", key, dec); err != nil {
		LogError("DECRBY: key=%v, err=%v", key, err)
		return err
	}
	return nil
}

// endregion

// region string
func (r *Redis) SetString(key string, value string) error {
	if _, err := r.DoCmd("SET", key, value); err != nil {
		LogError("SET: key=%v, err=%v", key, err)
		return err
	}
	return nil
}

func (r *Redis) GetString(key string) (string, error) {
	if res, err := r.DoCmd("GET", key); err != nil {
		LogError("GET: key=%v, err=%v", key, err)
		return "", err
	} else {
		if res == nil {
			return "", errReturnNilValue
		}
		return string(res.([]byte)), nil
	}
}

func (r *Redis) StrLen(key string) (int64, error) {
	if len, err := r.DoCmd("STRLEN", key); err != nil {
		LogError("STRLEN: key=%v, err=%v", key, err)
		return 0, err
	} else {
		return len.(int64), nil
	}
}

// endregion

// region hash
func (r *Redis) Hdel(key string, fields ...interface{}) error {
	if _, err := r.DoCmd("HDEL", key, fields); err != nil {
		LogError("HDEL: key=%v, err=%v", key, err)
		return err
	}
	return nil
}

func (r *Redis) HExists(key string, field interface{}) (bool, error) {
	if res, err := r.DoCmd("HEXISTS", key, field); err != nil {
		LogError("HEXISTS: key=%v, err=%v", key, err)
		return false, err
	} else {
		if 1 == res.(int64) {
			return true, nil
		}
	}
	return false, nil
}

func (r *Redis) HSet(key string, field interface{}, value interface{}) error {
	if _, err := r.DoCmd("HSET", key, field, value); err != nil {
		LogError("HSET: key=%v, err=%v", key, err)
		return err
	}
	return nil
}

func (r *Redis) HGet(key string, field interface{}) (string, error) {
	if res, err := r.DoCmd("HGET", key, field); err != nil {
		LogError("HGET: key=%v, err=%v", key, err)
		return "", err
	} else {
		if res == nil {
			return "", errReturnNilValue
		}
		return string(res.([]byte)), nil
	}
}

func (r *Redis) HGetAll(key string) (map[string]string, error) {
	if res, err := r.DoCmd("HGETALL", key); err != nil {
		LogError("HGETALL: key=%v, err=%v", key, err)
		return nil, err
	} else {
		array := res.([]interface{})
		if len(array) == 0 {
			return nil, errReturnNilValue
		}
		result := make(map[string]string, 0)
		for i := 0; i < len(array); {
			key := string(array[i].([]byte))
			i++
			value := string(array[i].([]byte))
			i++
			result[key] = value
		}
		return result, nil
	}
}

func (r *Redis) HKeys(key string) ([]string, error) {
	if res, err := r.DoCmd("HKEYS", key); err != nil {
		LogError("HKEYS: key=%v, err=%v", key, err)
		return nil, err
	} else {
		array := res.([]interface{})
		if len(array) == 0 {
			return nil, errReturnNilValue
		}
		result := make([]string, 0)
		for _, value := range array {
			result = append(result, string(value.([]byte)))
		}
		return result, nil
	}
}

func (r *Redis) HLen(key string) (int64, error) {
	if len, err := r.DoCmd("HLEN", key); err != nil {
		LogError("HLEN: key=%v, err=%v", key, err)
		return 0, err
	} else {
		return len.(int64), nil
	}
}

func (r *Redis) HMget(key string, fields ...interface{}) ([]string, error) {
	if res, err := r.DoCmd("HMGET", key, fields); err != nil {
		LogError("HMGET: key=%v, err=%v", key, err)
		return nil, err
	} else {
		array := res.([]interface{})
		if len(array) == 0 {
			return nil, errReturnNilValue
		}
		result := make([]string, 0)
		for _, value := range array {
			result = append(result, string(value.([]byte)))
		}
		return result, nil
	}
}

func (r *Redis) HMset(key string, fields ...interface{}) error {
	if _, err := r.DoCmd("HMSET", key, fields); err != nil {
		LogError("HMSET: key=%v, err=%v", key, err)
		return err
	}
	return nil
}

// endregion

// region list
func (r *Redis) LIndex(key string, index int) (string, error) {
	if res, err := r.DoCmd("LINDEX", key, index); err != nil {
		LogError("LINDEX: key=%v, err=%v", key, err)
		return "", err
	} else {
		return string(res.([]byte)), nil
	}
}

func (r *Redis) LPush(key string, values ...interface{}) error {
	if _, err := r.DoCmd("LPUSH", key, values); err != nil {
		LogError("LPUSH: key=%v, err=%v", key, err)
		return err
	}
	return nil
}

func (r *Redis) RPush(key string, values ...interface{}) error {
	if _, err := r.DoCmd("RPUSH", key, values); err != nil {
		LogError("RPUSH: key=%v, err=%v", key, err)
		return err
	}
	return nil
}

func (r *Redis) LPop(key string) (string, error) {
	if res, err := r.DoCmd("LPOP", key); err != nil {
		LogError("LPOP: key=%v, err=%v", key, err)
		return "", err
	} else {
		return string(res.([]byte)), nil
	}
}

func (r *Redis) RPop(key string) (string, error) {
	if res, err := r.DoCmd("RPOP", key); err != nil {
		LogError("RPOP: key=%v, err=%v", key, err)
		return "", err
	} else {
		return string(res.([]byte)), nil
	}
}

func (r *Redis) LRange(key string, start, end int) ([]string, error) {
	if res, err := r.DoCmd("LRANGE", key, start, end); err != nil {
		LogError("LRANGE: key=%v, err=%v", key, err)
		return nil, err
	} else {
		array := res.([]interface{})
		if len(array) == 0 {
			return nil, errReturnNilValue
		}
		result := make([]string, 0)
		for _, value := range array {
			result = append(result, string(value.([]byte)))
		}
		return result, nil
	}
}

func (r *Redis) LLen(key string) (int64, error) {
	if len, err := r.DoCmd("LLEN", key); err != nil {
		LogError("LLEN: key=%v, err=%v", key, err)
		return 0, err
	} else {
		return len.(int64), nil
	}
}

// endregion

// region set
func (r *Redis) SAdd(key string, member interface{}) error {
	if _, err := r.DoCmd("SADD", key, member); err != nil {
		LogError("SAdd: key=%v, err=%v", key, err)
		return err
	}
	return nil
}
func (r *Redis) SCard(key string) (int64, error) {
	if count, err := r.DoCmd("SCARD", key); err != nil {
		LogError("SCARD: key=%v, err=%v", key, err)
		return 0, err
	} else {
		return count.(int64), nil
	}
}
func (r *Redis) SIsMember(key string, member interface{}) (bool, error) {
	if res, err := r.DoCmd("SISMEMBER", key, member); err != nil {
		LogError("SISMEMBER: key=%v, err=%v", key, err)
		return false, nil
	} else {
		if 1 == res.(int64) {
			return true, nil
		}
	}
	return false, nil
}

func (r *Redis) SMembers(key string) ([]string, error) {
	if res, err := r.DoCmd("SMEMBERS", key); err != nil {
		LogError("SMEMBERS: key=%v, err=%v", key, err)
		return nil, err
	} else {
		array := res.([]interface{})
		if len(array) == 0 {
			return nil, errReturnNilValue
		}
		result := make([]string, 0)
		for _, value := range array {
			result = append(result, string(value.([]byte)))
		}
		return result, nil
	}
}
func (r *Redis) SPop(key string) (string, error) {
	if res, err := r.DoCmd("SPOP", key); err != nil {
		LogError("SPOP: key=%v, err=%v", key, err)
		return "", err
	} else {
		return string(res.([]byte)), nil
	}
}
func (r *Redis) SRem(key string, members ...interface{}) error {
	if _, err := r.DoCmd("SREM", key, members); err != nil {
		LogError("SREM: key=%v, err=%v", key, err)
		return err
	}
	return nil
}

// endregion

// region sorted set
func (r *Redis) ZAdd(key string, scoreAndmembers ...interface{}) error {
	if _, err := r.DoCmd("ZADD", key, scoreAndmembers); err != nil {
		LogError("ZADD: key=%v, scoreAndmembers=%v", key, scoreAndmembers)
		return err
	}
	return nil
}
func (r *Redis) ZCard(key string) (int64, error) {
	if count, err := r.DoCmd("LLEN", key); err != nil {
		LogError("LLEN: key=%v, err=%v", key, err)
		return 0, err
	} else {
		return count.(int64), nil
	}
}
func (r *Redis) ZCount(key string, min, max int64) (int64, error) {
	if count, err := r.DoCmd("LLEN", key); err != nil {
		LogError("LLEN: key=%v, err=%v", key, err)
		return 0, err
	} else {
		return count.(int64), nil
	}
}
func (r *Redis) ZRevRange(key string, start, end int64, withScores bool) ([][]string, error) {
	var res interface{}
	var err error
	if withScores {
		res, err = r.DoCmd("ZREVRANGE", key, start, end, "withscores")
	} else {
		res, err = r.DoCmd("ZREVRANGE", key, start, end)
	}
	if err != nil {
		LogError("ZREVRANGE: key=%v, err=%v", key, err)
		return nil, err
	}
	array := res.([]interface{})
	if len(array) == 0 {
		return nil, errReturnNilValue
	}
	result := make([][]string, 0)
	if withScores {
		for i := 0; i < len(array); {
			result[i] = append(result[i], array[i].(string))
			i++
			result[i] = append(result[i], array[i].(string))
			i++
		}
	} else {
		for i := 0; i < len(array); i++ {
			result[i] = append(result[i], array[i].(string))
		}
	}
	return result, nil
}
func (r *Redis) ZRank(key string, member interface{}) ([]string, error) {
	if res, err := r.DoCmd("ZRANK", key, member); err != nil {
		LogError("ZRANK: key=%v, err=%v", key, err)
		return nil, err
	} else {
		array := res.([]interface{})
		if len(array) == 0 {
			return nil, errReturnNilValue
		}
		result := make([]string, 0)
		for _, value := range array {
			result = append(result, string(value.([]byte)))
		}
		return result, nil
	}
}

func (r *Redis) ZRem(key string, members ...interface{}) error {
	if _, err := r.DoCmd("ZREM", key, members); err != nil {
		LogError("ZREM: key=%v, err=%v", key, err)
		return err
	}
	return nil
}
func (r *Redis) ZRevRank(key string) (int64, error) {
	if order, err := r.DoCmd("ZREVRANK", key); err != nil {
		LogError("ZREVRANK: key=%v, err=%", key, err)
		return 0, err
	} else {
		return order.(int64), nil
	}
}
func (r *Redis) ZScore(key string, member interface{}) (int64, error) {
	if score, err := r.DoCmd("ZSCORE", key); err != nil {
		LogError("ZSCORE: key=%v, err=%", key, err)
		return 0, err
	} else {
		return score.(int64), nil
	}
}

// endregion

//endregion
