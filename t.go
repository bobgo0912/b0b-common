package main

import (
	"context"
	"fmt"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/etcd"
	"github.com/bobgo0912/b0b-common/pkg/log"
	"github.com/bobgo0912/b0b-common/pkg/server"
	"os"
	"os/signal"
	"time"
)

func main() {
	fmt.Println(config.Version)
	ctx, ca := context.WithCancel(context.Background())
	log.InitLog()
	newConfig := config.NewConfig(config.Json)
	newConfig.Category = "pkg/config"
	newConfig.InitConfig()
	mainServer := server.NewMainServer()

	etcdClient := etcd.NewClientFromCnf()
	backServer := server.NewBackServer(config.Cfg.Host)
	mainServer.AddServer(backServer)

	mainServer.Discover(ctx, etcdClient)
	err := mainServer.Start(ctx)
	if err != nil {
		log.Panic(err)
	}

	c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	ca()
	time.Sleep(5 * time.Second)
	fmt.Println("dddd")
	os.Exit(1)
}
