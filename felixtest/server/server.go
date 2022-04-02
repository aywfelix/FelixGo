package main

import (
	. "github.com/aywfelix/felixgo/fnet"
)

func main() {
	nodeServer := NewNodeServer()
	nodeServer.Start()
}
