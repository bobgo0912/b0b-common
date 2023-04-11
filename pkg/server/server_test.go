package server

import (
	"context"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/etcd"
	"github.com/bobgo0912/b0b-common/pkg/log"
	"github.com/bobgo0912/b0b-common/pkg/server/common"
	"os"
	"os/signal"
	"testing"
	"time"
)

// import (
//
//	"bytes"
//	"context"
//	"fmt"
//	"github.com/bobgo0912/b0b-common/pkg/config"
//	"github.com/bobgo0912/b0b-common/pkg/constant"
//	"github.com/bobgo0912/b0b-common/pkg/etcd"
//	"github.com/bobgo0912/b0b-common/pkg/log"
//	hello "github.com/bobgo0912/b0b-common/pkg/server/proto"
//	"github.com/bobgo0912/b0b-common/pkg/util"
//	"github.com/gorilla/handlers"
//	"google.golang.org/protobuf/proto"
//	"io"
//	"net/http"
//	"os"
//	"os/signal"
//	"testing"
//	"time"
//
// )
//
//	type HelleServer struct {
//		hello.UnimplementedGreeterServer
//	}
//
// func (s *HelleServer) SayHello(ctx context.Context, request *hello.HelloRequest) (*hello.HelloReply,
//
//		error) {
//		return &hello.HelloReply{
//			Message: "hello " + request.Name,
//		}, nil
//	}
func TestServer(t *testing.T) {
	ctx, ca := context.WithCancel(context.Background())
	log.InitLog()
	newConfig := config.NewConfig(config.Json)
	newConfig.Category = "../config"
	newConfig.InitConfig()
	server := NewMainServer()
	etcdClient := etcd.NewClientFromCnf()

	irisServer := NewIrisServer(config.Cfg.Host, config.Cfg.Port)
	server.AddServer(irisServer)
	//server.AddServer(irisServer)
	err := server.Start(ctx)
	if err != nil {
		t.Fatal(err)
	}
	server.Discover(ctx, etcdClient)

	time.Sleep(time.Second * 5)
	n, err := GetNodeN("testServers", common.Http)
	if err != nil {
		t.Error("ddddd")
	}
	t.Log(n)
	//
	time.Sleep(time.Second * 15)

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

//
//func TestHttpGrpc(t *testing.T) {
//
//	request := hello.HelloRequest{Name: "xxx"}
//	marshal, err := proto.Marshal(&request)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	newRequest, err := http.NewRequest("POST", "http://localhost:8888/proto", bytes.NewReader(marshal))
//	if err != nil {
//		t.Fatal(err)
//	}
//	newRequest.Header.Add(constant.ProtoHeader, "hello.HelloRequest")
//	do, err := http.DefaultClient.Do(newRequest)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer do.Body.Close()
//	all, err := io.ReadAll(do.Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if do.StatusCode != http.StatusOK {
//		t.Error(string(all))
//		return
//	}
//	var resp hello.HelloReply
//	err = proto.Unmarshal(all, &resp)
//	if err != nil {
//		t.Error(string(all))
//		t.Fatal(err)
//	}
//	t.Log(resp)
//}
//
//func TestA(t *testing.T) {
//	util.GetIp()
//
//}
//
//func TestBackServer_Start(t *testing.T) {
//	ctx, ca := context.WithCancel(context.Background())
//	defer ca()
//	log.InitLog()
//	newConfig := config.NewConfig(config.Json)
//	newConfig.Category = "../config"
//	newConfig.InitConfig()
//	server := NewMainServer()
//
//	etcdClient := etcd.NewClientFromCnf()
//	backServer := NewBackServer(config.Cfg.Host)
//	server.AddServer(backServer)
//
//	server.Discover(ctx, etcdClient)
//	err := server.Start(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	c := make(chan os.Signal, 1)
//	signal.Notify(c, os.Interrupt, os.Kill)
//	<-c
//
//}
