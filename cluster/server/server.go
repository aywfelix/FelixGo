package main

import (
	"fmt"
	"time"

	. "github.com/aywfelix/felixgo/configure"
	. "github.com/aywfelix/felixgo/fnet"
)

func main() {
	// 读取配置
	err := LoadIniConfig("../felixgo.ini")
	if err != nil {
		fmt.Println("load config error")
		return
	}
	server := NewNodeServer(&MasterCfg.NetNode)
	server.Start()

	for {
		time.Sleep(time.Second * 1)
	}
}
