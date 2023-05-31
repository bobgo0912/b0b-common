package server

import (
	"context"
	"fmt"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/etcd"
	"github.com/bobgo0912/b0b-common/pkg/log"
	"github.com/bobgo0912/b0b-common/pkg/meter"
	"github.com/bobgo0912/b0b-common/pkg/server/common"
	bp "github.com/bobgo0912/b0b-common/pkg/server/proto"
	h "github.com/bobgo0912/b0b-common/pkg/server/proto"
	"github.com/bobgo0912/b0b-common/pkg/trac"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"net/http"
	"os"
	"os/signal"
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
func TestOtelGrpc(t *testing.T) {
	ctx, ca := context.WithCancel(context.Background())
	log.InitLog()
	newConfig := config.NewConfig(config.Json)
	newConfig.Category = "../config/d"

	newConfig.InitConfig()
	mainServer := NewMainServer()
	etcdClient, err := etcd.NewClientFromCnf()
	if err != nil {
		log.Panicf("etcd init fail")
	}

	r := NewRouter()
	r.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		log.Info("test")
		writer.Write([]byte("ttt"))
	}).Methods("GET")
	//r.Use(GrpcMid)
	r.HandleProtoFunc("/proto", func(req proto.Message, w http.ResponseWriter) {
		log.Info(req)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("xxxxxxx"))
	}, &bp.HelloRequest{}).Methods("POST")
	httpServer := NewMuxServer(config.Cfg.Host, 2211, r)
	grpcServer := NewGrpcServer(config.Cfg.Host, 2212,
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	grpcServer.RegService(&bp.Greeter_ServiceDesc, &HelleServer{})
	mainServer.AddServer(httpServer)
	mainServer.AddServer(grpcServer)

	otelGrpc, err := trac.NewOtelGrpc(ctx, otlptracegrpc.WithEndpoint("localhost:4317"), otlptracegrpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer otelGrpc.ShutDown(ctx)
	otelMetricGrpc, err := meter.NewOtelMetricGrpc(ctx, otlpmetricgrpc.WithEndpoint("localhost:4317"), otlpmetricgrpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer otelMetricGrpc.ShutDown(ctx)

	err = mainServer.Start(ctx)
	if err != nil {
		log.Panic(err)
	}
	mainServer.Discover(ctx, etcdClient)
	//
	time.Sleep(time.Second * 15)
	address := GetRpcNodeAddress("testServers", common.Http)
	if address == "" {
		log.Info(" bad address")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	ca()
	log.Info("xxx")
}

func TestHelloServer(t *testing.T) {
	ctx, ca := context.WithCancel(context.Background())
	defer ca()
	log.InitLog()
	newConfig := config.NewConfig(config.Json)
	newConfig.Category = "../config/d"

	newConfig.InitConfig()
	otelGrpc, err := trac.NewOtelGrpc(ctx, otlptracegrpc.WithEndpoint("localhost:4317"), otlptracegrpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer otelGrpc.ShutDown(ctx)
	//stream
	conn, err := grpc.Dial("127.0.0.1:2212", grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := h.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &h.HelloRequest{Name: "bobby"})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)
}
