package sql

import (
	"context"
	"fmt"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/jmoiron/sqlx"
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
		MysqlCfg: config.MysqlCfg{
			UserName: "root",
			Password: "123456",
			Host:     "127.0.0.1",
			Port:     3306,
			Database: "edu",
		},
	}
	config.Cfg = &cfg
	store := GetStuStore()
	id, err := store.QueryById(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)
}
