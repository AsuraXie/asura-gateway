package server

import (
	"asura-gateway/config"
	"asura-gateway/global"
	"asura-gateway/log"
	"github.com/gin-gonic/gin"
)

type GateWayServer struct{
	Server *gin.Engine
}

func (s *GateWayServer) Start(config *config.GateWayConfig){
	s.Server = gin.Default()
	s.Server.POST("metrics",global.MetricsHandler)
	s.Server.POST("health",global.HealthHandler)
	s.Server.POST("gw",GateWayHandler)
	s.Server.Run(config.GetListenAddr())
}

func (s *GateWayServer) Stop(){
	log.Info("stop")
}

func (s *GateWayServer) HeartBeat(){
	log.Info("heart pong pong pong")
}