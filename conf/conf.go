package conf

import (
	"encoding/json"
	"flag"
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

	ChanRPCLen         int
	TimerDispatcherLen int

	ConsolePort   int
	ConsolePrompt string
	ProfilePath   string
}

var EnvPath string

func init() {
	flag.StringVar(&EnvPath, "env", "", "path of env file")
	flag.Parse()

	if EnvPath == "" {
		fmt.Println("param -env no set")
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(EnvPath + "ggs.env")
	if err != nil {
		fmt.Println("env file not found, path: " + EnvPath)
		os.Exit(1)
	}

	err = json.Unmarshal(data, &Env)
	if err != nil {
		fmt.Printf("env file format error: %v\n", err)
		os.Exit(1)
	}

	if Env.WSAddr != "" {
		if Env.MaxMsgLen <= 0 {
			Env.MaxMsgLen = 4096
			fmt.Println("invalid MaxMsgLen, reset to %v", Env.MaxMsgLen)
		}
		if Env.HTTPTimeout <= 0 {
			Env.HTTPTimeout = 10
			fmt.Println("invalid HTTPTimeout, reset to %v", Env.HTTPTimeout)
		}
		Env.HTTPTimeout *= time.Second
	}
}
