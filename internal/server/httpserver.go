package server

import (
	"context"
	"fmt"
	"github.com/bobgo0912/b0b-common/internal/config"
	"github.com/bobgo0912/b0b-common/internal/log"
	"net/http"
	"time"
)

type HttpServer struct {
	Server
	Router  *Router
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

func NewHttpServer(host string, port int, router *Router, options ...Option) *HttpServer {
	return &HttpServer{
		Server: Server{
			Ctx:      context.Background(),
			Type:     Http,
			Port:     port,
			Host:     host,
			HostName: config.Cfg.HostName,
		},
		Router:  router,
		Options: options,
	}
}

func (s *HttpServer) Start(ctx context.Context) error {
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
	log.Infof("httpServer %s:%d start", s.Host, s.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("http ListenAndServe fail err=", err)
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
	log.Info("httpServer stop")
	return nil
}

func (s *HttpServer) Ctx() context.Context {
	return s.Server.Ctx
}

func (s *HttpServer) GetInfo() Server {
	return s.Server
}
