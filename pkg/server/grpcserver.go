package server

import (
	"context"
	"fmt"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/log"
	"github.com/bobgo0912/b0b-common/pkg/server/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

type GrpcServer struct {
	Server     Server
	GrpcServer *grpc.Server
}

func NewGrpcServer(host string, port int, options ...grpc.ServerOption) *GrpcServer {
	server := grpc.NewServer(options...)
	return &GrpcServer{
		Server: Server{
			Ctx:      context.Background(),
			Type:     common.RPC,
			Host:     host,
			Port:     port,
			HostName: config.Cfg.HostName,
		},
		GrpcServer: server,
	}
}

func (s *GrpcServer) RegService(sd *grpc.ServiceDesc, ss interface{}) {
	s.GrpcServer.RegisterService(sd, ss)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s.GrpcServer, healthServer)
}

func (s *GrpcServer) Start(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", s.Server.Host, s.Server.Port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Error("grpc Listen ", address, " fail err=", err)
		return err
	}
	log.Infof("GrpcServer %s start", address)
	go func() {
		err = s.GrpcServer.Serve(listen)
		if err != nil {
			log.Panic("grpc Serve start ", address, " fail err=", err)
			return
		}
	}()
	select {
	case <-ctx.Done():
		s.GrpcServer.GracefulStop()
		break
	}
	log.Info("rpcServer stop ", address)
	return nil
}

func (s *GrpcServer) Ctx() context.Context {
	return s.Server.Ctx
}

func (s *GrpcServer) GetInfo() Server {
	return s.Server
}
