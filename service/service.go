package service

import (
	"ggs/conf"
	"ggs/log"
	"runtime"
	"sync"
)

type Service interface {
	OnInit()
	OnDestroy()
	Run(closeSig chan bool)
}

type service struct {
	si       Service
	closeSig chan bool
	wg       sync.WaitGroup
}

var (
	services   []*service
	serviceCnt int
)

func Register(si Service) {
	s := new(service)
	s.si = si
	s.closeSig = make(chan bool, 1)

	services = append(services, s)
}

func Init() {
	serviceCnt = len(services)

	for i := 0; i < serviceCnt; i++ {
		services[i].si.OnInit()
	}

	for i := 0; i < serviceCnt; i++ {
		go run(services[i])
	}
}

func Destroy() {
	for i := serviceCnt - 1; i >= 0; i-- {
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
			if conf.StackBufLen > 0 {
				buf := make([]byte, conf.StackBufLen)
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("%v", r)
			}
		}
	}()

	s.si.OnDestroy()
}
