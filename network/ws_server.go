package network

import (
	"ggs/conf"
	"ggs/log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WSServer struct {
	ln      net.Listener
	handler *WSHandler
}

type WSHandler struct {
	//	newAgent func(*WSConn) Agent
	upgrader websocket.Upgrader
	//	conns WebsocketConnSet
	//	mutexConns      sync.Mutex
	wg sync.WaitGroup
}

func (handler *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Debug("upgrade error: %v", err)
		return
	}
	conn.SetReadLimit(int64(conf.Env.MaxMsgLen))

	handler.wg.Add(1)
	defer handler.wg.Done()
}

func (server *WSServer) Start(func(*WSConn) Agent) {
	ln, err := net.Listen("tcp", conf.Env.WSAddr)
	if err != nil {
		log.Fatal("%v", err)
	}

	//	if server.NewAgent == nil {
	//		log.Fatal("NewAgent must not be nil")
	//	}

	server.ln = ln
	server.handler = &WSHandler{
		upgrader: websocket.Upgrader{
			HandshakeTimeout: conf.Env.HTTPTimeout,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}
	httpServer := &http.Server{
		Handler:        server.handler,
		ReadTimeout:    conf.Env.HTTPTimeout,
		WriteTimeout:   conf.Env.HTTPTimeout,
		MaxHeaderBytes: 1024,
	}
	go httpServer.Serve(ln)
}

func (server *WSServer) Close() {
	server.ln.Close()

	server.handler.wg.Wait()
}
