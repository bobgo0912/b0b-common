package config

import (
	"context"
	"github.com/bobgo0912/b0b-common/pkg/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	log.InitLog()
	config := NewConfig(Json)
	config.InitConfig()

	log.Infof("conf=%+v", Cfg)
	etcdCfg := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	}
	client, err := clientv3.New(etcdCfg)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	err = config.EtcdMerge(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	//config = NewConfig(Yaml)
	//config.InitConfig()
}

func TestConfigYaml(t *testing.T) {
	log.InitLog()
	config := NewConfig(Yaml)
	config.InitConfig()

	log.Infof("conf=%+v", Cfg)
	etcdCfg := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	}
	client, err := clientv3.New(etcdCfg)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	err = config.EtcdMerge(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	//config = NewConfig(Yaml)
	//config.InitConfig()
}
