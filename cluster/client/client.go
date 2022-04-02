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
		fmt.Println("load config error, ", err.Error())
		return
	}
	client := NewNodeClient(&GateCfg.NetNode)
	client.Start()

	for {
		time.Sleep(time.Second * 1)
	}
}
