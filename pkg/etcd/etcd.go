package etcd

import (
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/log"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

type Client struct {
}

func NewClient(cfg clientv3.Config) (*clientv3.Client, error) {
	client, err := clientv3.New(cfg)
	if err != nil {
		log.Error("failed to connect etcd: ", err)
		return nil, errors.Wrap(err, "failed to connect etcd")
	}
	return client, nil
}
func NewClientFromCnf() (*clientv3.Client, error) {
	split := strings.Split(config.Cfg.EtcdCfg.Host, ",")
	cfg := clientv3.Config{
		Endpoints:   split,
		DialTimeout: time.Second * 5,
		Username:    config.Cfg.EtcdCfg.Username,
		Password:    config.Cfg.EtcdCfg.Password,
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		log.Error("failed to connect etcd: ", err)
		return nil, errors.Wrap(err, "failed to connect etcd")
	}
	return client, nil
}
