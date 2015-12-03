package service

import (
	"ggs/chanrpc"
	"ggs/conf"
	"ggs/log"
)

type Skeleton struct {
	chanRPCServer *chanrpc.Server
}

func NewSkeleton() *Skeleton {
	skeleton := &Skeleton{
		chanRPCServer: chanrpc.NewServer(conf.Env.ChanRPCLen),
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
		}
	}
}

func (s *Skeleton) RegisterChanRPC(id interface{}, f interface{}) {
	if s.chanRPCServer == nil {
		panic("invalid ChanRPCServer")
	}

	s.chanRPCServer.Register(id, f)
}