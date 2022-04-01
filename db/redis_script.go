package db

const (
	REDIS_SCRIPT_GET     = "redis_get"
	REDIS_SCRIPT_HGET    = "redis_hget"
	REDIS_SCRIPT_SET     = "redis_set"
	REDIS_SCRIPT_HSET    = "redis_hset"
	REDIS_SCRIPT_UNLOCK  = "redis_unlock"
	REDIS_SCRIPT_HUNLOCK = "redis_hunlock"
	REDIS_SCRIPT_DEL     = "redis_del"
	REDIS_SCRIPT_HDEL    = "redis_hdel"
	REDIS_SCRIPT_PUSH    = "redis_push"
	REDIS_SCRIPT_POP     = "redis_pop"
	REDIS_SCRIPT_INCRBY  = "redis_incrby"
	REDIS_SCRIPT_GLOCK   = "redis_glock"
	REDIS_SCRIPT_GUNLOCK = "redis_gunlock"
)

var (
	redis_script_get = `
		local key = KEYS[1]
		local lock = (ARGV[1]~="0")
		if lock then
			local locker = string.format("locker:%s", key)
			local ret = redis.call("SET", locker, 1, "EX", 5, "NX")
			if not ret then
				return "error"
			end
		end
		local tag = string.format("%s", key)
		return redis.call("GET", tag)
	`

	redis_script_hget = `
		local key = KEYS[1]
		local lock = (table.remove(ARGV, 1) ~= "0")
		local len = #ARGV
		local lockers = {}
		if lock then
			for i=1, len do
				local locker = string.format("locker:%s:%s", key, ARGV[i])
				local ret = redis.call("SET", locker, 1, "EX", 5, "NX")
				table.insert(lockers, locker)
				if not ret then
					redis.call("DEL", unpack(lockers))
					return "error"
				end
			end
		end
		local tag = string.format("%s", key)
		if len==0 then
			return redis.call("HGETALL", tag)
		end
		return redis.call("HMGET", tag, unpack(ARGV))
	`

	redis_script_set = `
		local key = KEYS[1]
		local unlock = (ARGV[1]~= "0")
		local data = ARGV[2]
		if unlock then
			local locker = string.format("locker:%s", key)
			redis.call("DEL", locker)
		end
		local tag = string.format("%s", key)
		if #ARGV > 2 then
			local expire = tonumber(ARGV[3])
			if expire >0 then
				redis.call("SET", tag, data, "EX", expire)
				return
			end
		end
		redis.call("SET", tag, data)
	`

	redis_script_hset = `
		local key = KEYS[1]
		local field = ARGV[1]
		local unlock = (ARGV[2] ~= "0")
		local data = ARGV[3]
		if unlock then
			local locker = string.format("locker:%s:%s", key, field)
			redis.call("DEL", locker)
		end
		reids.call("HSET", key, field, data)
	`
	redis_script_unlock = `
		local key = KEYS[1]
		local locker = string.format("locker:%s", key)
		redis.call("DEL", locker)
	`
	redis_script_hunlock = `
		local key = KEYS[1]
		local field = ARGV[1]
		local locker = string.format("locker:%s:%s", key, field)
		redis.call("DEL", locker)
	`
	redis_script_del = `
		local key = KEYS[1]
		local locker = string.format("locker:%s", key)
		redis.call("DEL", locker)
		redis.call("DEL", key)
	`
	redis_script_hdel = `
		local key = KEYS[1]
		local field = ARGV[1]
		local locker = string.format("locker:%s:%s", key, field)
		redis.call("DEL", locker)
		redis.call("HDEL", key, field)
	`
	redis_script_push = `
		local key = KEYS[1]
		local data = ARGV[1]
		local time = tonumber(ARGV[2])
		local tag = string.format("queue:%s", key)
		redis.call("ZADD", tag, "NX", time, data)
	`
	redis_script_pop = `
		local key = KEYS[1]
		local tag = string.format("queue:%s", key)
		local result = redis.call("ZRANGE", tag, 0, 0)
		if element then
			redis.call("ZREM", tag, element)
			return element
		end
	`
	redis_script_incrby = `
		local key = KEYS[1]
		local data = ARGV[1]
		local result = redis.call("GET", key)
		if not result then
			redis.call("SET", key, tonumber(data))
		else
			redis.call("SET", key, tonumber(data)+tonumber(tostring(result)))
		end
	`
	redis_script_glock = `
			local key = KEYS[1]
			local token = ARGV[1]
			local timeout = tonumber(ARGV[2])
			local ret = redis.call("GET", key)
			if ret then
				if ret == token then
					return "error"
				else
					return "different token"
				end
			end
			redis.call("SET", key, token, "EX", timeout)
			return "ok"
	`

	redis_script_gunlock = `
			local key = KEYS[1]
			local token = ARGV[1]
			local ret = redis.call("GET", key)
			if ret then
				redis.call("DEL", key)
				if ret ~= token then
					return "error"
				end
			end
			return "ok"
	`
)

func SetRedisScript() {
	RedisHelper.SetScript(REDIS_SCRIPT_GET, redis_script_get, 1)
	RedisHelper.SetScript(REDIS_SCRIPT_HGET, redis_script_hget, 1)
	RedisHelper.SetScript(REDIS_SCRIPT_SET, redis_script_set, 1)
	RedisHelper.SetScript(REDIS_SCRIPT_HSET, redis_script_hset, 1)
	RedisHelper.SetScript(REDIS_SCRIPT_UNLOCK, redis_script_unlock, 1)
	RedisHelper.SetScript(REDIS_SCRIPT_HUNLOCK, redis_script_hunlock, 1)
	RedisHelper.SetScript(REDIS_SCRIPT_DEL, redis_script_del, 1)
	RedisHelper.SetScript(REDIS_SCRIPT_HDEL, redis_script_hdel, 1)
	RedisHelper.SetScript(REDIS_SCRIPT_PUSH, redis_script_push, 1)
	RedisHelper.SetScript(REDIS_SCRIPT_POP, redis_script_pop, 1)
	RedisHelper.SetScript(REDIS_SCRIPT_INCRBY, redis_script_incrby, 1)
	RedisHelper.SetScript(REDIS_SCRIPT_GLOCK, redis_script_glock, 1)
	RedisHelper.SetScript(REDIS_SCRIPT_GUNLOCK, redis_script_gunlock, 1)
}
