package database

// 从redis获取数据并写入到db中
// 每隔1分钟，读出更新标识，进行玩家数据保存
// 从redis中获取所有更新过的玩家id，根据玩家id 查找玩家redis 数据，然后将redis数据读出并更新到数据库中

type IRedisToDB interface {
}

type RedisToDB struct {
	updateList []int // 玩家id

}
