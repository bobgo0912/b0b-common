package server

import (
	"b0b-common/internal/log"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
)

type GrpcServer struct {
	Server     Server
	GrpcServer *grpc.Server
}

func NewGrpcServer(host string, port int, options ...grpc.ServerOption) *GrpcServer {
	server := grpc.NewServer(options...)
	return &GrpcServer{
		Server: Server{
			Ctx:  context.Background(),
			Type: RPC,
			Host: host,
			Port: port,
		},
		GrpcServer: server,
	}
}

func (s *GrpcServer) RegService(sd *grpc.ServiceDesc, ss interface{}) {
	s.GrpcServer.RegisterService(sd, ss)
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
			log.Error("grpc Serve ", address, " fail err=", err)
			return
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	log.Info("rpcServer stop ", address)
	return nil
}

func (s *GrpcServer) Ctx() context.Context {
	return s.Server.Ctx
}

func (s *GrpcServer) GetInfo() Server {
	return s.Server
}
