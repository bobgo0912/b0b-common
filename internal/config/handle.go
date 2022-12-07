package config

import (
	"context"
	"fmt"
	"github.com/bobgo0912/b0b-common/internal/constant"
	"github.com/bobgo0912/b0b-common/internal/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Type string

const (
	Json Type = "json"
	Yaml Type = "yaml"
)

var Cfg *ServerCfg

type HandleInterface interface {
	Read([]byte)
}
type Config struct {
	Type     Type
	Handle   HandleInterface
	Category string
}

func NewConfig(readType Type) *Config {
	switch readType {
	case Yaml:
		return &Config{Handle: &YamHandle{}, Type: Yaml}
	case Json:
		return &Config{Handle: &JsonCfgHandle{}, Type: Json}
	default:
		return &Config{Handle: &JsonCfgHandle{}, Type: Json}
	}
}

func (c *Config) InitConfig() {
	env := os.Getenv(constant.EnvName)
	var cfgFile = fmt.Sprintf("config-dev.%s", c.Type)
	if env != "" {
		cfgFile = fmt.Sprintf("config-%s.%s", env, c.Type)
	}
	if c.Category != "" {
		cfgFile = fmt.Sprintf("%s/%s", c.Category, cfgFile)
	}
	abs, err := filepath.Abs(cfgFile)
	if err != nil {
		log.Panic("Abs cfgFile fail err=", err)
	}
	file, err := os.Open(abs)
	if err != nil {
		log.Panic("open cfgFile fail err=", err)
	}
	defer file.Close()
	all, err := io.ReadAll(file)
	if err != nil {
		log.Panic("open cfgFile ReadAll fail err=", err)
	}
	c.Handle.Read(all)
	Cfg.ENV = constant.ENV(env)
	if env == "" {
		Cfg.ENV = "dev"
	}
	log.Infof("config=%+v", Cfg)
}

func (c *Config) EtcdMerge(ctx context.Context, client *clientv3.Client) error {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*30)
	defer cancelFunc()
	key := fmt.Sprintf(constant.EtcdConfig, c.Type, Cfg.ENV)
	res, err := client.Get(timeout, key)
	if err != nil {
		log.Error("EtcdMerge get fail err=", err)
		return err
	}
	for _, kv := range res.Kvs {
		k := string(kv.Key)
		if key == k {
			c.Handle.Read(kv.Value)
		}
	}
	log.Infof("EtcdConfig=%+v", Cfg)
	return nil
}
