package server

import (
	"context"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/etcd"
	"github.com/bobgo0912/b0b-common/pkg/log"
	"github.com/bobgo0912/b0b-common/pkg/meter"
	"github.com/bobgo0912/b0b-common/pkg/server/middleware"
	"github.com/bobgo0912/b0b-common/pkg/trac"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestMux(t *testing.T) {
	ctx, ca := context.WithCancel(context.Background())
	log.InitLog()
	newConfig := config.NewConfig(config.Json)
	newConfig.Category = "../config/d"

	newConfig.InitConfig()
	server := NewMainServer()
	etcdClient, err := etcd.NewClientFromCnf()
	if err != nil {
		log.Panicf("etcd init fail")
	}

	newRouter := NewRouter()
	muxServer := NewMuxServer(config.Cfg.Host, 1231, newRouter)
	server.AddServer(muxServer)
	prop := propagation.TraceContext{}
	newRouter.Use(otelmux.Middleware(config.Cfg.ServiceName, otelmux.WithPropagators(prop)))
	newRouter.Use(middleware.Meter)
	newRouter.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		//log.Otel(request.Context()).Infof("", "test")
		writer.Write([]byte("ttt"))
	}).Methods("GET")
	newRouter.HandleFunc("/test/{tt}", func(writer http.ResponseWriter, request *http.Request) {
		//s := trace.SpanContextFromContext(request.Context()).TraceID().String()
		//log.Info("tt ", s)
		field := zap.String("sd", "sds")
		//log.Info("tt ", s, field)
		log.Otel(request.Context()).Error("133311=", field)
		log.Otel(request.Context()).Info("111=", field)
		//log.Otel(request.Context()).Infof("%s", "sdsdsd=asdsa")
		writer.Write([]byte("zz"))
	}).Methods("GET")

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
	err = server.Start(ctx)
	if err != nil {
		t.Fatal(err)
	}
	server.Discover(ctx, etcdClient)

	time.Sleep(time.Second * 5)
	//n, err := GetNodeN("testServers", common.Http)
	//if err != nil {
	//	t.Error("ddddd")
	//}
	//t.Log(n)
	//
	time.Sleep(time.Second * 15)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	ca()
}
