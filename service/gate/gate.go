package gate

var Service = new(Gate)

type Gate struct {
	//wsServer *network.WSServer
}

func (gate *Gate) OnInit() {
	
}

func (gate *Gate) Run(closeSig chan bool) {

}

func (gate *Gate) OnDestroy() {
	
}