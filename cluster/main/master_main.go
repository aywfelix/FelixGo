package main

import (
	"fmt"

	. "github.com/aywfelix/felixgo/cluster/master"
	. "github.com/aywfelix/felixgo/configure"
)

// 服务器开始启动执行入口
func main() {
	// 加载服务器配置
	err := LoadIniConfig("../config/felixgo.ini")
	if err != nil {
		fmt.Println("load config error")
		return
	}
	// 初始化节点服务
	server := NewMasterServer()
	// 启动服务器
	server.Start()
	// 服务器开始执行
	server.Serve()
	// 服务器关闭
	//server.Stop()
}
