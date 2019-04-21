package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"strconv"
	"strings"
)

var confPath = flag.String("conf", "../conf/gateway.conf", "gateway config file")

type GateWayConfig struct {
	ListenAddress string `json:"listen_address"`
	ListenPort    int    `json:"listen_port"`
}

func LoadConfig() *GateWayConfig {
	flag.Parse()
	globalConfig := &GateWayConfig{}

	confValue, err := ioutil.ReadFile(*confPath)
	if err != nil {
		panic("config read error:" + err.Error())
	}

	if err = json.Unmarshal(confValue, globalConfig); err != nil {
		panic("config unmarshal error:" + err.Error())
	}
	return globalConfig
}

func (self *GateWayConfig) GetListenAddr() string {
	port := strconv.Itoa(self.ListenPort)

	var listenAddr string
	if strings.Contains(self.ListenAddress, "*") {
		listenAddr = ":" + port
	} else {
		listenAddr = self.ListenAddress + ":" + port
	}
	return listenAddr
}
