package main

import (
	. "github.com/felix/felixgo/fnet"
)

func main() {
	nodeServer := NewNodeServer()
	nodeServer.Start()
}