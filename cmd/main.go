package main

import (
	"asura-gateway/config"
	"asura-gateway/server"
)

func main() {
	conf := config.LoadConfig()
	go func() {
		server.TestServer()
	}()

	server := &server.GateWayServer{}
	server.Start(conf)
}
