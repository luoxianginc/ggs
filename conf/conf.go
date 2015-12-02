package conf

import (
	"encoding/json"
	"fmt"
	"ggs/log"
	"io/ioutil"
	"os"
	"time"
)

var Env struct {
	StackBufLen int
	LogLevel    string
	LogPath     string

	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32

	WSAddr      string
	HTTPTimeout time.Duration
}

func init() {
	data, err := ioutil.ReadFile("ggs.env")
	if err != nil {
		fmt.Println("file not found: ggs.env")
		os.Exit(1)
	}

	err = json.Unmarshal(data, &Env)
	if err != nil {
		fmt.Println("invalid format: ggs.env")
		os.Exit(1)
	}
}

func Init() {
	if Env.WSAddr != "" {
		log.Info("Loading websocket environment...")
		if Env.MaxMsgLen <= 0 {
			Env.MaxMsgLen = 4096
			log.Info("invalid MaxMsgLen, reset to %v", Env.MaxMsgLen)
		}
		if Env.HTTPTimeout <= 0 {
			Env.HTTPTimeout = 10 * time.Second
			log.Info("invalid HTTPTimeout, reset to %v", Env.HTTPTimeout)
		}
	}
}
