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
	newAgent   func(*WSConn) Agent
	upgrader   websocket.Upgrader
	conns      WebsocketConnSet
	mutexConns sync.Mutex
	wg         sync.WaitGroup
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

	handler.mutexConns.Lock()
	if handler.conns == nil {
		handler.mutexConns.Unlock()
		conn.Close()
		return
	}
	if len(handler.conns) >= conf.Env.MaxConnNum {
		handler.mutexConns.Unlock()
		conn.Close()
		log.Debug("too many connections")
		return
	}
	handler.conns[conn] = struct{}{}
	handler.mutexConns.Unlock()

	wsConn := NewWSConn(conn)
	agent := handler.newAgent(wsConn)
	agent.Run()

	wsConn.Close()
	handler.mutexConns.Lock()
	delete(handler.conns, conn)
	handler.mutexConns.Unlock()
	agent.OnClose()
}

func (server *WSServer) Start(newAgent func(*WSConn) Agent) {
	if newAgent == nil {
		log.Fatal("newAgent must not be nil")
	}

	ln, err := net.Listen("tcp", conf.Env.WSAddr)
	if err != nil {
		log.Fatal("%v", err)
	}

	server.ln = ln
	server.handler = &WSHandler{
		newAgent: newAgent,
		conns:    make(WebsocketConnSet),
		upgrader: websocket.Upgrader{
			HandshakeTimeout: conf.Env.HTTPTimeout,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}

	httpServer := &http.Server{
		Addr:           conf.Env.WSAddr,
		Handler:        server.handler,
		ReadTimeout:    conf.Env.HTTPTimeout,
		WriteTimeout:   conf.Env.HTTPTimeout,
		MaxHeaderBytes: 1024,
	}

	go httpServer.Serve(ln)
	log.Info("Websocket host: %v", conf.Env.WSAddr)
}

func (server *WSServer) Close() {
	server.ln.Close()

	server.handler.mutexConns.Lock()
	for conn := range server.handler.conns {
		conn.Close()
	}
	server.handler.conns = nil
	server.handler.mutexConns.Unlock()

	server.handler.wg.Wait()
}
