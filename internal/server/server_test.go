package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bobgo0912/b0b-common/internal/config"
	"github.com/bobgo0912/b0b-common/internal/constant"
	"github.com/bobgo0912/b0b-common/internal/etcd"
	"github.com/bobgo0912/b0b-common/internal/log"
	hello "github.com/bobgo0912/b0b-common/internal/server/proto"
	"github.com/bobgo0912/b0b-common/internal/util"
	"github.com/gorilla/handlers"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"os"
	"os/signal"
	"testing"
	"time"
)

type HelleServer struct {
	hello.UnimplementedGreeterServer
}

func (s *HelleServer) SayHello(ctx context.Context, request *hello.HelloRequest) (*hello.HelloReply,
	error) {
	return &hello.HelloReply{
		Message: "hello " + request.Name,
	}, nil
}

func TestServer(t *testing.T) {
	ctx, ca := context.WithCancel(context.Background())
	log.InitLog()
	newConfig := config.NewConfig(config.Json)
	newConfig.Category = "../config"
	newConfig.InitConfig()
	server := NewMainServer()
	etcdClient := etcd.NewClientFromCnf()

	r := NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"Proto"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	r.Use(handlers.CORS(headersOk, originsOk, methodsOk, handlers.AllowCredentials()))
	r.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		log.Info("test")
		ip := RemoteIp(request)
		log.Info(ip)
		writer.Write([]byte("ttt"))
	}).Methods("GET")
	//r.Use(GrpcMid)
	r.HandleProtoFunc("/proto", func(req proto.Message, w http.ResponseWriter) {
		log.Info(req)
		request := req.(*hello.HelloRequest)
		fmt.Println(request)
		w.WriteHeader(http.StatusOK)
		reply := hello.HelloReply{Message: "drrr"}
		marshal, _ := proto.Marshal(&reply)
		w.Write(marshal)
	}, &hello.HelloRequest{}).Methods("POST", "OPTIONS")
	r.HandleProtoFunc1("/proto1", func(req any, w http.ResponseWriter) {
		log.Info(req)
		request := req.(*hello.HelloRequest)
		fmt.Println(request)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("xxxxxxx"))
	}, &hello.HelloRequest{}).Methods("POST")

	httpServer := NewMuxServer(config.Cfg.Host, config.Cfg.Port, r)
	//irisServer := NewIrisServer(config.Cfg.Host, 9999)
	//irisServer.Iris.Get("test", func(c context2.Context) {
	//	log.Info("xxdd")
	//})
	grpcServer := NewGrpcServer(config.Cfg.Host, config.Cfg.RpcPort)
	grpcServer.RegService(&hello.Greeter_ServiceDesc, &HelleServer{})
	server.AddServer(httpServer)
	server.AddServer(grpcServer)
	//server.AddServer(irisServer)
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
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	ca()
	//conn, err := grpc.Dial(address, grpc.WithInsecure())
	//if err != nil {
	//	panic(err)
	//}
	//defer conn.Close()
	//c := hello.NewGreeterClient(conn)
	//hello, err := c.SayHello(context.Background(), &hello.HelloRequest{Name: "bobby"})
	//if err != nil {
	//	panic(err)
	//}
	//t.Log(hello)

	select {}
}

func TestHttpGrpc(t *testing.T) {

	request := hello.HelloRequest{Name: "xxx"}
	marshal, err := proto.Marshal(&request)
	if err != nil {
		t.Fatal(err)
	}

	newRequest, err := http.NewRequest("POST", "http://localhost:8888/proto", bytes.NewReader(marshal))
	if err != nil {
		t.Fatal(err)
	}
	newRequest.Header.Add(constant.ProtoHeader, "hello.HelloRequest")
	do, err := http.DefaultClient.Do(newRequest)
	if err != nil {
		t.Fatal(err)
	}
	defer do.Body.Close()
	all, err := io.ReadAll(do.Body)
	if err != nil {
		t.Fatal(err)
	}
	if do.StatusCode != http.StatusOK {
		t.Error(string(all))
		return
	}
	var resp hello.HelloReply
	err = proto.Unmarshal(all, &resp)
	if err != nil {
		t.Error(string(all))
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestA(t *testing.T) {
	util.GetIp()

}
