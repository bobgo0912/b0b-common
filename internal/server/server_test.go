package server

import (
	"b0b-common/internal/config"
	"b0b-common/internal/etcd"
	"b0b-common/internal/log"
	"b0b-common/internal/server/proto"
	"context"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"net/http"
	"testing"
	"time"
)

type HelleServer struct {
	proto.UnimplementedGreeterServer
}

func (s *HelleServer) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply,
	error) {
	return &proto.HelloReply{
		Message: "hello " + request.Name,
	}, nil
}

func TestServer(t *testing.T) {
	ctx := context.Background()
	log.InitLog()
	newConfig := config.NewConfig(config.Json)
	newConfig.Category = "../config"
	newConfig.InitConfig()
	server := NewMainServer()
	etcdClient := etcd.NewClientFromCnf()

	r := mux.NewRouter()
	r.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		log.Info("test")
		writer.Write([]byte("ttt"))
	}).Methods("GET")
	httpServer := NewHttpServer(config.Cfg.Host, config.Cfg.Port, r)
	grpcServer := NewGrpcServer(config.Cfg.Host, config.Cfg.RpcPort)
	grpcServer.RegService(&proto.Greeter_ServiceDesc, &HelleServer{})
	server.AddServer(httpServer)
	server.AddServer(grpcServer)
	err := server.Start(ctx)
	if err != nil {
		t.Fatal(err)
	}
	server.Discover(ctx, etcdClient)

	//
	time.Sleep(time.Second * 15)
	address := GetRpcNodeAddress("testServers")
	if address == "" {
		t.Log(" bad address")
		select {}

	}
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := proto.NewGreeterClient(conn)
	hello, err := c.SayHello(context.Background(), &proto.HelloRequest{Name: "bobby"})
	if err != nil {
		panic(err)
	}
	t.Log(hello)

	select {}
}
