package main

import (
	. "github.com/aywfelix/felixgo/cluster/master"
)

// 服务器开始启动执行入口
func main() {
	// 初始化节点服务
	server := MasterServer{}
	// 启动服务器
	server.Start()
	// 服务器开始执行
	server.Serve()
	// 服务器关闭
	server.Stop()
}
