package service

import (
	"ggs/chanrpc"
	"ggs/conf"
	"ggs/log"
	"ggs/timer"
	"time"
)

type Skeleton struct {
	chanRPCServer *chanrpc.Server
	dispatcher    *timer.Dispatcher
}

func NewSkeleton() *Skeleton {
	skeleton := &Skeleton{
		chanRPCServer: chanrpc.NewServer(conf.Env.ChanRPCLen),
		dispatcher:    timer.NewDispatcher(conf.Env.TimerDispatcherLen),
	}
	return skeleton
}

func (s *Skeleton) ChanRPCServer() *chanrpc.Server {
	return s.chanRPCServer
}

func (s *Skeleton) Run(closeSig chan bool) {
	for {
		select {
		case ci := <-s.chanRPCServer.ChanCall:
			err := s.chanRPCServer.Exec(ci)
			if err != nil {
				log.Error("%v", err)
			}
		case t := <-s.dispatcher.ChanTimer:
			t.Cb()
		}
	}
}

func (s *Skeleton) RegisterChanRPC(id interface{}, f interface{}) {
	if s.chanRPCServer == nil {
		panic("invalid ChanRPCServer")
	}

	s.chanRPCServer.Register(id, f)
}

func (s *Skeleton) AfterFunc(d time.Duration, cb func()) *timer.Timer {
	return s.dispatcher.AfterFunc(d, cb)
}

func (s *Skeleton) CronFunc(cronExpr *timer.CronExpr, cb func()) *timer.Cron {
	return s.dispatcher.CronFunc(cronExpr, cb)
}
