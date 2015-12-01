package gate

import (
	"ggs/conf"
	"ggs/network"
)

var Service = new(Gate)

type Gate struct {
}

func (gate *Gate) OnInit() {
}

func (gate *Gate) Run(closeSig chan bool) {
	var wsServer *network.WSServer
	if conf.Env.WSAddr != "" {
		wsServer = new(network.WSServer)
	}
	
	if wsServer != nil {
		wsServer.Start(nil)
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
}

func (gate *Gate) OnDestroy() {

}
