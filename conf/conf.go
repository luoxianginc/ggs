package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

var Env struct {
	StackBufLen int
	LogLevel    string
	LogPath     string

	WSAddr      string
	MaxMsgLen   uint32
	HTTPTimeout time.Duration
}

func init() {
	data, err := ioutil.ReadFile("ggs.env")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Env)
	if err != nil {
		log.Fatal("%v", err)
	}
}
