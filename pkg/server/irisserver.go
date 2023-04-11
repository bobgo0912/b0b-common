package server

import (
	"context"
	"fmt"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/log"
	"github.com/bobgo0912/b0b-common/pkg/server/common"
	"github.com/kataras/iris/v12"
	"time"
)

type IrisServer struct {
	Server
	//Router  *MuxRouter
	//Options []Option
	Iris *iris.Application
}

func NewIrisServer(host string, port int) *IrisServer {
	return &IrisServer{
		Server: Server{
			Ctx:      context.Background(),
			Type:     common.Http,
			Port:     port,
			Host:     host,
			HostName: config.Cfg.HostName,
		},
		Iris: iris.New(),
	}
}
func (s *IrisServer) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.Server.Host, s.Server.Port)
	log.Infof("irisServer %s start", addr)
	go func() {

		err := s.Iris.Run(
			iris.Addr(addr), // 启动服务，监控地址和端口
			iris.WithoutServerError(iris.ErrServerClosed), // 忽略服务器错误
			iris.WithOptimizations,                        // 让程序自身尽可能的优化
			iris.WithCharset("UTF-8"))
		if err != nil {
			log.Panic("irisServer start fail err=", err)
		}
	}()

	select {
	case <-ctx.Done():
		break
	}
	log.Info("irisServer stop")
	timeCtx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	err := s.Iris.Shutdown(timeCtx)
	if err != nil {
		log.Error("irisServer Shutdown fail err=", err)
		return err
	}
	return nil
}
func (s *IrisServer) Ctx() context.Context {
	return s.Server.Ctx
}

func (s *IrisServer) GetInfo() Server {
	return s.Server
}
