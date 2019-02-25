package config

import (
	"encoding/json"
	"io/ioutil"
)

type GateWayConfig struct {
	ListenAddress string `json:"listen_address"`
	ListenPort    int `json:"listen_port"`
}

var ServerConfig GateWayConfig

func (myconfig *GateWayConfig) LoadConfig(path string) {
	confValue, err := ioutil.ReadFile(path)
	if err != nil {
		panic("config read error:" + err.Error())
	}

	if err = json.Unmarshal(confValue, &ServerConfig); err != nil {
		panic("config unmarshal error:" + err.Error())
	}
}
