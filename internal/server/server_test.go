package server

import (
	"b0b-common/internal/config"
	"b0b-common/internal/log"
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	ctx := context.Background()
	log.InitLog()
	newConfig := config.NewConfig(config.Json)
	newConfig.Category = "../config"
	newConfig.InitConfig()
	server := NewMainServer()
	r := mux.NewRouter()
	r.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		log.Info("test")
		writer.Write([]byte("ttt"))
	}).Methods("GET")
	httpServer := NewHttpServer(config.Cfg.Host, config.Cfg.Port, r)
	server.AddServer(httpServer)
	err := server.Start(ctx)
	if err != nil {
		t.Fatal(err)
	}
	select {}
}
