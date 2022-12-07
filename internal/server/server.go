package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bobgo0912/b0b-common/internal/log"
	"github.com/bobgo0912/b0b-common/internal/util"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Type string

const (
	Http Type = "http"
	RPC  Type = "rpc"
)

var MainServers *MainServer

type Interface interface {
	Start(ctx context.Context) error
	Ctx() context.Context
	GetInfo() Server
}

type Server struct {
	Ctx      context.Context `json:"-"`
	Type     Type            `json:"type"`
	Port     int             `json:"port"`
	Host     string          `json:"host"`
	HostName string          `json:"hostName"`
}

type EtcdReg struct {
	Type     Type   `json:"type"`
	Port     int    `json:"port"`
	Host     string `json:"host"`
	HostName string `json:"hostName"`
}

func (s *Server) ToJson() (string, error) {
	s.Host = util.GetIp()
	marshal, err := json.Marshal(s)
	if err != nil {
		return "", nil
	}
	return string(marshal), nil
}

type MainServer struct {
	Servers         []Interface
	DiscoverServers *Servers
}

func NewMainServer() *MainServer {
	if MainServers != nil {
		return MainServers
	}
	m := &MainServer{
		Servers: make([]Interface, 0),
	}
	MainServers = m
	return m
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
				log.Errorf("server %+v  start fail err=%v", server, err)
			}
		}(server)
	}
	return nil
}

func (m *MainServer) Discover(ctx context.Context, etcd *clientv3.Client) {
	if m.DiscoverServers == nil {
		m.DiscoverServers = NewServices(ctx, etcd)
	}
	for _, server := range m.Servers {
		info := server.GetInfo()
		serverCtx := server.Ctx()
		discover := NewDiscover(serverCtx, etcd, info)
		go discover.Start()
	}
	m.DiscoverServers.Start()
}
