package ggs

import (
	"ggs/conf"
	"ggs/log"
)

func Run() {
	log.Init(conf.LogLevel, conf.LogPath)
	log.Debug("GGS is starting up...")
}
