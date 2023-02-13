package sql

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/etcd"
	"github.com/bobgo0912/b0b-common/pkg/log"
	"github.com/bobgo0912/b0b-common/pkg/trac"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"testing"
)

func TestCon(t *testing.T) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		"root", "123456", "123456", 3306, "edu",
	)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(db)
}

// func TestMain(m *testing.M) {
//
//		m.Run()
//	}
func TestSDS(t *testing.T) {
	cfg := config.ServerCfg{
		MysqlCfg: map[string]*config.MysqlCfg{"edu": {
			UserName: "root",
			Password: "123456",
			Host:     "127.0.0.1",
			Port:     3306,
			Database: "edu",
		}},
	}
	config.Cfg = &cfg
	store, err := GetStuStore()
	if err != nil {
		t.Fatal(err)
	}
	//id, err := store.QueryById(context.Background(), 1)
	////id, err := store.Store.QueryById(context.Background(), 1)
	////id, err := store.QueryById(context.Background(), 1)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(id)

	page, err := store.QueryPage(context.Background(), squirrel.Select("*"), 2, 10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(page)
}
func TestWithOtel(t *testing.T) {
	ctx, can := context.WithCancel(context.Background())
	defer can()
	log.InitLog()
	newConfig := config.NewConfig(config.Json)
	newConfig.Category = "../config"
	newConfig.InitConfig()
	etcdClient := etcd.NewClientFromCnf()
	err := newConfig.EtcdMerge(ctx, etcdClient)
	if err != nil {
		t.Fatal(err)
	}
	otelGrpc, err := trac.NewOtelGrpc(ctx, otlptracegrpc.WithEndpoint("localhost:4317"), otlptracegrpc.WithInsecure())
	defer otelGrpc.ShutDown(ctx)
	if err != nil {
		t.Fatal(err)
	}
	store, err := GetStuStore()
	if err != nil {
		t.Fatal(err)
	}
	//page, err := store.QueryPage(context.Background(), squirrel.Select("*").Where(squirrel.Eq{"age": 10}), 2, 10)
	//if err != nil {
	//	t.Fatal(err)
	//}

	list, err := store.QueryList(ctx, squirrel.Select("*").Where(squirrel.Eq{"id": 1, "name": "23"}))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)
}
