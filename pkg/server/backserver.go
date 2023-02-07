package server

import (
	"context"
	"fmt"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/log"
)

type BackServer struct {
	Server Server
}

func NewBackServer(host string) *BackServer {
	return &BackServer{
		Server: Server{
			Ctx:      context.Background(),
			Type:     Back,
			Host:     host,
			Port:     0,
			HostName: config.Cfg.HostName,
		},
	}
}
func (s *BackServer) Start(ctx context.Context) error {
	fmt.Println(" bstart")
	select {
	case <-ctx.Done():
		log.Info("backServer stop ")
		break
	}
	fmt.Println("xxxxxx")
	return nil
}

func (s *BackServer) Ctx() context.Context {
	return s.Server.Ctx
}

func (s *BackServer) GetInfo() Server {
	return s.Server
}
