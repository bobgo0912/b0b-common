package server

import (
	"b0b-common/internal/config"
	"b0b-common/internal/constant"
	"b0b-common/internal/etcd"
	"b0b-common/internal/log"
	bp "b0b-common/internal/server/proto"
	"bytes"
	"context"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"testing"
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

func TestServer(t *testing.T) {
	ctx := context.Background()
	log.InitLog()
	newConfig := config.NewConfig(config.Json)
	newConfig.Category = "../config"
	newConfig.InitConfig()
	server := NewMainServer()
	etcdClient := etcd.NewClientFromCnf()

	r := NewRouter()
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
	r.HandleProtoFunc1("/proto", func(req *proto.Message, w http.ResponseWriter) {
		log.Info(req)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("xxxxxxx"))
	}, &bp.HelloRequest{}).Methods("POST")
	httpServer := NewHttpServer(config.Cfg.Host, config.Cfg.Port, r)
	grpcServer := NewGrpcServer(config.Cfg.Host, config.Cfg.RpcPort)
	grpcServer.RegService(&bp.Greeter_ServiceDesc, &HelleServer{})
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
	//conn, err := grpc.Dial(address, grpc.WithInsecure())
	//if err != nil {
	//	panic(err)
	//}
	//defer conn.Close()
	//c := bp.NewGreeterClient(conn)
	//hello, err := c.SayHello(context.Background(), &bp.HelloRequest{Name: "bobby"})
	//if err != nil {
	//	panic(err)
	//}
	//t.Log(hello)

	select {}
}

func TestHttpGrpc(t *testing.T) {

	request := bp.HelloRequest{Name: "xxx"}
	marshal, err := proto.Marshal(&request)
	if err != nil {
		t.Fatal(err)
	}

	newRequest, err := http.NewRequest("POST", "http://localhost:8888/proto", bytes.NewReader(marshal))
	if err != nil {
		t.Fatal(err)
	}
	newRequest.Header.Add(constant.ProtoHeader, "HelloRequest")
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
	var resp bp.HelloReply
	err = proto.Unmarshal(all, &resp)
	if err != nil {
		t.Error(string(all))
		t.Fatal(err)
	}
	t.Log(resp)
}
