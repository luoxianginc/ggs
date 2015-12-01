package network

import (
	"ggs/log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WSServer struct {
	Addr string
	//	MaxConnNum      int
	//	PendingWriteNum int
	MaxMsgLen uint32
	//	HTTPTimeout     time.Duration
	//	NewAgent    func(*WSConn) Agent
	HTTPTimeout time.Duration
	ln          net.Listener
	handler     *WSHandler
}

type WSHandler struct {
	//	maxConnNum      int
	//	pendingWriteNum int
	maxMsgLen uint32
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
	conn.SetReadLimit(int64(handler.maxMsgLen))

	handler.wg.Add(1)
	defer handler.wg.Done()

	log.Info("%v", handler.maxMsgLen)
}

func (server *WSServer) Start() {
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal("%v", err)
	}

	//	if server.MaxConnNum <= 0 {
	//		server.MaxConnNum = 100
	//		log.Info("invalid MaxConnNum, reset to %v", server.MaxConnNum)
	//	}
	//	if server.PendingWriteNum <= 0 {
	//		server.PendingWriteNum = 100
	//		log.Info("invalid PendingWriteNum, reset to %v", server.PendingWriteNum)
	//	}
	if server.MaxMsgLen <= 0 {
		server.MaxMsgLen = 4096
		log.Info("invalid MaxMsgLen, reset to %v", server.MaxMsgLen)
	}
	if server.HTTPTimeout <= 0 {
		server.HTTPTimeout = 10 * time.Second
		log.Info("invalid HTTPTimeout, reset to %v", server.HTTPTimeout)
	}
	//	if server.NewAgent == nil {
	//		log.Fatal("NewAgent must not be nil")
	//	}

	server.ln = ln
	server.handler = &WSHandler{
		maxMsgLen: server.MaxMsgLen,
		upgrader: websocket.Upgrader{
			HandshakeTimeout: server.HTTPTimeout,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}
	httpServer := &http.Server{
		Addr:           server.Addr,
		Handler:        server.handler,
		ReadTimeout:    server.HTTPTimeout,
		WriteTimeout:   server.HTTPTimeout,
		MaxHeaderBytes: 1024,
	}
	go httpServer.Serve(ln)
}

func (server *WSServer) Close() {
	server.ln.Close()

	server.handler.wg.Wait()
}
