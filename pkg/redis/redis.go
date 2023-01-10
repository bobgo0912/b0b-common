package redis

import (
	"context"
	"fmt"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/log"
	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"time"
)

type Client struct {
}

func NewClient() (*redis.Client, error) {
	host := config.Cfg.RedisCfg.Host
	port := config.Cfg.RedisCfg.Port
	password := config.Cfg.RedisCfg.Password
	size := config.Cfg.RedisCfg.Size
	db := config.Cfg.RedisCfg.Db
	addr := fmt.Sprintf("%s:%d", host, port)
	redis := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    password, // no password set
		DB:          db,       // use default DB
		PoolSize:    size,     // 连接池大小
		DialTimeout: time.Duration(30) * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := redis.Ping(ctx).Result()
	if err != nil {
		log.Infof("InitRedis: %s", err.Error())
		return nil, errors.Wrap(err, "InitRedis fail")
	}
	log.Infof("InitRedis: %s", "success")
	return redis, nil
}
func NewClusterClient() (*redis.ClusterClient, error) {
	hosts := config.Cfg.RedisCfg.Hosts
	password := config.Cfg.RedisCfg.Password
	size := config.Cfg.RedisCfg.Size
	redis := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:       hosts,
		Password:    password, // no password set
		PoolSize:    size,     // 连接池大小
		DialTimeout: time.Duration(30) * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := redis.Ping(ctx).Result()
	if err != nil {
		log.Infof("InitRedis: %s", err.Error())
		return nil, errors.Wrap(err, "InitRedis fail")
	}
	log.Infof("InitRedis: %s", "success")
	return redis, nil
}
