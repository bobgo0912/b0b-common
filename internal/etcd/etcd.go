package etcd

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
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
