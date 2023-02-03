package redis

import (
	"context"
	"fmt"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/log"
	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type Client struct {
	Client *redis.Client
}

const otelName = "b0b-common/redis"

func NewClient() (*Client, error) {
	host := config.Cfg.RedisCfg.Host
	port := config.Cfg.RedisCfg.Port
	password := config.Cfg.RedisCfg.Password
	size := config.Cfg.RedisCfg.Size
	db := config.Cfg.RedisCfg.Db
	addr := fmt.Sprintf("%s:%d", host, port)
	redisClient := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    password, // no password set
		DB:          db,       // use default DB
		PoolSize:    size,     // 连接池大小
		DialTimeout: time.Duration(30) * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Infof("InitRedis: %s", err.Error())
		return nil, errors.Wrap(err, "InitRedis fail")
	}
	log.Infof("InitRedis: %s", "success")
	return &Client{Client: redisClient}, nil
}

func newOTELSpan(ctx context.Context, name string) trace.Span {
	_, span := otel.Tracer(otelName).Start(ctx, name)
	span.SetAttributes(semconv.DBSystemRedis)
	return span
}
func NewClusterClient() (*redis.ClusterClient, error) {
	hosts := config.Cfg.RedisCfg.Hosts
	password := config.Cfg.RedisCfg.Password
	size := config.Cfg.RedisCfg.Size
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:       hosts,
		Password:    password, // no password set
		PoolSize:    size,     // 连接池大小
		DialTimeout: time.Duration(30) * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := clusterClient.Ping(ctx).Result()
	if err != nil {
		log.Infof("InitRedis: %s", err.Error())
		return nil, errors.Wrap(err, "InitRedis fail")
	}
	log.Infof("InitRedis: %s", "success")
	return clusterClient, nil
}
