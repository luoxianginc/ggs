package gate

import (
	"ggs/conf"
	"ggs/network"
	"time"
)

var Service = new(Gate)

type Gate struct {
}

func (gate *Gate) OnInit() {
}

func (gate *Gate) Run(closeSig chan bool) {
	var wsServer *network.WSServer = new(network.WSServer)
	if conf.Env.WSAddr != "" {
		wsServer.Addr = conf.Env.WSAddr
		wsServer.MaxMsgLen = conf.Env.MaxMsgLen
		wsServer.HTTPTimeout = conf.Env.HTTPTimeout * time.Microsecond
	}

	if wsServer != nil {
		wsServer.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
}

func (gate *Gate) OnDestroy() {

}
