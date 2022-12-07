package main

import (
	"context"
	"github.com/bobgo0912/b0b-common/internal/config"
	"github.com/bobgo0912/b0b-common/internal/etcd"
	"github.com/bobgo0912/b0b-common/internal/log"
	"github.com/bobgo0912/b0b-common/internal/server"
	bp "github.com/bobgo0912/b0b-common/internal/server/proto"
	"google.golang.org/protobuf/proto"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type HelleServer struct {
	bp.UnimplementedGreeterServer
}

func (s *HelleServer) SayHello(ctx context.Context, request *bp.HelloRequest) (*bp.HelloReply,
	error) {
	return &bp.HelloReply{
		Message: "hello " + request.Name,
	}, nil
}
func main() {

	ctx, ca := context.WithCancel(context.Background())
	log.InitLog()
	newConfig := config.NewConfig(config.Json)
	newConfig.Category = "../config"
	newConfig.InitConfig()
	mainServer := server.NewMainServer()
	etcdClient := etcd.NewClientFromCnf()

	r := server.NewRouter()
	r.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		log.Info("test")
		writer.Write([]byte("ttt"))
	}).Methods("GET")
	//r.Use(GrpcMid)
	r.HandleProtoFunc("/proto", func(req *proto.Message, w http.ResponseWriter) {
		log.Info(req)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("xxxxxxx"))
	}, &bp.HelloRequest{}).Methods("POST")
	httpServer := server.NewHttpServer(config.Cfg.Host, config.Cfg.Port, r)
	grpcServer := server.NewGrpcServer(config.Cfg.Host, config.Cfg.RpcPort)
	grpcServer.RegService(&bp.Greeter_ServiceDesc, &HelleServer{})
	mainServer.AddServer(httpServer)
	mainServer.AddServer(grpcServer)
	err := mainServer.Start(ctx)
	if err != nil {
		log.Panic(err)
	}
	mainServer.Discover(ctx, etcdClient)
	//
	time.Sleep(time.Second * 15)
	address := server.GetRpcNodeAddress("testServers")
	if address == "" {
		log.Info(" bad address")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	ca()
	log.Info("xxx")
}
