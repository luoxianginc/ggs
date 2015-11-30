package main

import (
	"ggs/conf"
	"ggs/log"
	"ggs/service"
	"ggs/service/gate"
	"os"
	"os/signal"
)

func Run(services ...service.Service) {
	log.Init(conf.LogLevel, conf.LogPath)
	log.Info("GGS is starting up...")

	for i := 0; i < len(services); i++ {
		service.Register(services[i])
	}
	service.Init()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Info("GGS is closing down (signal: %v)", sig)
	service.Destroy()
}

func main() {
	Run(gate.Service)
}