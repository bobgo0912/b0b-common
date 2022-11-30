package server

import (
	"b0b-common/internal/log"
	"context"
	"errors"
)

type Type string

const (
	Http Type = "http"
	RPC  Type = "rpc"
)

type Interface interface {
	Start(ctx context.Context) error
}
type Server struct {
	Type Type
	Port int
	Host string
}

type MainServer struct {
	Servers []Interface
}

func NewMainServer() *MainServer {
	return &MainServer{
		Servers: make([]Interface, 0),
	}
}

func (m *MainServer) AddServer(server Interface) {
	m.Servers = append(m.Servers, server)
}

func (m *MainServer) Start(ctx context.Context) error {
	if len(m.Servers) < 1 {
		return errors.New("no server to start")
	}
	for _, server := range m.Servers {
		go func(server Interface) {
			err := server.Start(ctx)
			if err != nil {
				log.Errorf("server %+v", server)
			}
		}(server)
	}
	return nil
}
