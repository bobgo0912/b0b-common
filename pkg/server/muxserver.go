package server

import (
	"context"
	"fmt"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/log"
	"github.com/bobgo0912/b0b-common/pkg/server/common"
	"net/http"
	"time"
)

type MuxServer struct {
	Server
	Router  *MuxRouter
	Options []Option
}

type Option func(*http.Server) error

func WriteTimeOut(duration time.Duration) Option {
	return func(server *http.Server) error {
		server.WriteTimeout = duration
		return nil
	}
}

func ReadTimeout(duration time.Duration) Option {
	return func(server *http.Server) error {
		server.ReadTimeout = duration
		return nil
	}
}
func IdleTimeout(duration time.Duration) Option {
	return func(server *http.Server) error {
		server.IdleTimeout = duration
		return nil
	}
}

func NewMuxServer(host string, port int, router *MuxRouter, options ...Option) *MuxServer {
	return &MuxServer{
		Server: Server{
			Ctx:      context.Background(),
			Type:     common.Http,
			Port:     port,
			Host:     host,
			HostName: config.Cfg.HostName,
		},
		Router:  router,
		Options: options,
	}
}

func (s *MuxServer) Start(ctx context.Context) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.Host, s.Port),
		Handler: s.Router.R,
	}

	for _, option := range s.Options {
		if option != nil {
			err := option(srv)
			if err != nil {
				log.Error("Start httpServer fail err=", err)
				return err
			}
		}
	}
	log.Infof("muxServer %s:%d start", s.Host, s.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Panicf("muxServer start fail err=", err)
		}
	}()

	select {
	case <-ctx.Done():
		break
	}

	timeCtx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	err := srv.Shutdown(timeCtx)
	if err != nil {
		log.Error("http Shutdown fail err=", err)
		return err
	}
	log.Info("muxServer stop")
	return nil
}

func (s *MuxServer) Ctx() context.Context {
	return s.Server.Ctx
}

func (s *MuxServer) GetInfo() Server {
	return s.Server
}
