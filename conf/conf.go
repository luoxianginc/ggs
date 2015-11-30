package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var Env struct {
	StackBufLen int
	LogLevel    string
	LogPath     string
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
