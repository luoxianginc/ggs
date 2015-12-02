package gate

import (
	"ggs/chanrpc"
	"ggs/conf"
	"ggs/log"
	"ggs/network"
)

type Gate struct {
	Processor    network.Processor
	AgentChanRPC *chanrpc.Server
}

func (gate *Gate) OnInit() {}

func (gate *Gate) Run(closeSig chan bool) {
	var wsServer *network.WSServer
	if conf.Env.WSAddr != "" {
		wsServer = new(network.WSServer)
	}

	if wsServer != nil {
		wsServer.Start(func(conn *network.WSConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go("NewAgent", a)
			}
			return a
		})
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
}

func (gate *Gate) OnDestroy() {}

type agent struct {
	conn     network.Conn
	gate     *Gate
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		if a.gate.Processor != nil {
			msg, err := a.gate.Processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = a.gate.Processor.Route(msg, a)
			if err != nil {
				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (a *agent) OnClose() {
	if a.gate.AgentChanRPC != nil {
		err := a.gate.AgentChanRPC.Open(0).Call0("CloseAgent", a)
		if err != nil {
			log.Error("chanrpc error: %v", err)
		}
	}
}
