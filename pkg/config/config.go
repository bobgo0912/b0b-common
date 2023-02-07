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
	MysqlCfg    map[string]*MysqlCfg `json:"mysql" yaml:"mysqlCfg"`
	RedisCfg    RedisCfg             `json:"redis" yaml:"redisCfg"`
	NatsCfg     NatsCfg              `json:"nats" yaml:"natsCfg"`
	EtcdCfg     EtcdCfg              `json:"etcd" yaml:"etcdCfg"`
	OtelCfg     OtelCfg              `json:"otel" yaml:"otel"`
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

func (c *MysqlCfg) UnmarshalYAML(unmarshal func(interface{}) error) error {
	cn := c
	cn.Password = "***"
	return unmarshal(&cn)
}

type RedisCfg struct {
	Hosts    []string `json:"hosts"`
	Password string   `json:"password" mask:"true"`
	Port     int      `json:"port"`
	Size     int      `json:"size"`
	Db       int      `json:"db"`
	Host     string   `json:"host"`
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
