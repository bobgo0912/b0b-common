package config

import (
	"github.com/bobgo0912/b0b-common/pkg/constant"
	"github.com/bobgo0912/b0b-common/pkg/util"
	"reflect"
)

type ServerCfg struct {
	NodeId      string               `json:"nodeId" yaml:"nodeId"`
	ServiceName string               `json:"serviceName" yaml:"serviceName"`
	Host        string               `json:"host" yaml:"host"`
	HostName    string               `json:"hostName" yaml:"hostName"`
	Port        int                  `json:"port" yaml:"port"`
	RpcPort     int                  `json:"rpcPort" yaml:"rpcPort"`
	ENV         constant.ENV         `json:"-" yaml:"-"`
	MysqlCfg    map[string]*MysqlCfg `json:"mysql" yaml:"mysql"`
	RedisCfg    *RedisCfg            `json:"redis" yaml:"redis"`
	NatsCfg     *NatsCfg             `json:"nats" yaml:"nats"`
	EtcdCfg     *EtcdCfg             `json:"etcd" yaml:"etcd"`
	OtelCfg     *OtelCfg             `json:"otel" yaml:"otel"`
	Version     string               `json:"-" yaml:"-"`
}

type MysqlCfg struct {
	UserName string `json:"userName" yaml:"userName"`
	Password string `json:"password" yaml:"password" mask:"true"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Database string `json:"database" yaml:"database"`
}

func (c *MysqlCfg) String() string {
	return util.ConfigMask(reflect.ValueOf(*c))
}

type RedisCfg struct {
	Hosts    []string `json:"hosts" yaml:"hosts"`
	Password string   `json:"password" yaml:"password" mask:"true"`
	Port     int      `json:"port" yaml:"port"`
	Size     int      `json:"size" yaml:"size"`
	Db       int      `json:"db" yaml:"db"`
	Host     string   `json:"host" yaml:"host"`
}

func (c *RedisCfg) String() string {
	return util.ConfigMask(reflect.ValueOf(*c))
}

type NatsCfg struct {
	Host     string `json:"host"  yaml:"host"`
	UserName string `json:"userName" yaml:"userName"`
	Password string `json:"password" yaml:"password" mask:"true"`
	Port     int    `json:"port" yaml:"port"`
}

func (c *NatsCfg) String() string {
	return util.ConfigMask(reflect.ValueOf(*c))
}

type EtcdCfg struct {
	Hosts    []string `json:"hosts"  yaml:"hosts"`
	UserName string   `json:"userName" yaml:"userName"`
	Password string   `json:"password" yaml:"password" mask:"true"`
}

func (c *EtcdCfg) String() string {
	return util.ConfigMask(reflect.ValueOf(*c))
}

type OtelCfg struct {
	Host   string            `json:"host" yaml:"host"`
	Port   int               `json:"port" yaml:"port"`
	Type   constant.OtelType `json:"type" yaml:"type"`
	Secure bool              `json:"secure" yaml:"secure"`
}
