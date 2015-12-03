package conf

import (
	"encoding/json"
	"fmt"
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

	ChanRPCLen int
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

	if Env.WSAddr != "" {
		if Env.MaxMsgLen <= 0 {
			Env.MaxMsgLen = 4096
			fmt.Println("invalid MaxMsgLen, reset to %v", Env.MaxMsgLen)
		}
		if Env.HTTPTimeout <= 0 {
			Env.HTTPTimeout = 10 * time.Second
			fmt.Println("invalid HTTPTimeout, reset to %v", Env.HTTPTimeout)
		}
	}
}
