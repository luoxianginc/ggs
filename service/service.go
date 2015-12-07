package service

import (
	"ggs/conf"
	"ggs/log"
	"runtime"
	"sync"
)

type Service interface {
	OnInit()
	Run(closeSig chan bool)
	OnDestroy()
}

type service struct {
	si       Service
	closeSig chan bool
	wg       sync.WaitGroup
}

var services []*service

func Register(si Service) {
	s := new(service)
	s.si = si
	s.closeSig = make(chan bool, 1)

	services = append(services, s)
}

func Init() {
	for i := 0; i < len(services); i++ {
		services[i].si.OnInit()
	}

	for i := 0; i < len(services); i++ {
		go run(services[i])
	}
}

func Destroy() {
	for i := len(services) - 1; i >= 0; i-- {
		s := services[i]
		s.closeSig <- true
		s.wg.Wait()
		destroy(s)
	}
}

func run(s *service) {
	s.wg.Add(1)
	s.si.Run(s.closeSig)
	s.wg.Done()
}

func destroy(s *service) {
	defer func() {
		if r := recover(); r != nil {
			if conf.Env.StackBufLen > 0 {
				buf := make([]byte, conf.Env.StackBufLen)
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("%v", r)
			}
		}
	}()

	s.si.OnDestroy()
}
