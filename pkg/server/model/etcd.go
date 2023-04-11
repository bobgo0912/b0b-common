package model

import (
	"github.com/bobgo0912/b0b-common/pkg/server/common"
)

type EtcdReg struct {
	Type     common.Type `json:"type"`
	Port     int         `json:"port"`
	Host     string      `json:"host"`
	HostName string      `json:"hostName"`
	Key      string      `json:"-"`
}
