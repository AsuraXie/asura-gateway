package main

import (
	"asura-gateway/config"
	"asura-gateway/server"
	"flag"
)

var confPath = flag.String("conf", "../conf/gateway.conf", "gateway config file")

func main() {
	flag.Parse()
	config.ServerConfig.LoadConfig(*confPath)
	server.StartGateWay()
}
