// sample.go
package main

import (
	"ggs/conf"
	"ggs/log"
)

func main() {
	log.Init(conf.LogLevel, conf.LogPath)
	defer log.Close()

	log.Info("GGS Starting up...")
}
