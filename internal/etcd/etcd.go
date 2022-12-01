package etcd

import (
	"b0b-common/internal/config"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Client struct {
}

func NewClient(cfg clientv3.Config) *clientv3.Client {

	client, err := clientv3.New(cfg)
	if err != nil {
		fmt.Println("failed to connect etcd: ", err)
		return nil
	}
	return client
}
func NewClientFromCnf() *clientv3.Client {
	cfg := clientv3.Config{
		Endpoints:   config.Cfg.EtcdCfg.Hosts,
		DialTimeout: time.Second * 5,
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		fmt.Println("failed to connect etcd: ", err)
		return nil
	}
	return client
}
