package config

import (
	"encoding/json"
	"github.com/bobgo0912/b0b-common/pkg/constant"
)

type ServerCfg struct {
	NodeId      string       `json:"nodeId" yaml:"nodeId"`
	ServiceName string       `json:"serviceName" yaml:"serviceName"`
	Host        string       `json:"host" yaml:"host"`
	HostName    string       `json:"hostName" yaml:"hostName"`
	Port        int          `json:"port" yaml:"port"`
	RpcPort     int          `json:"rpcPort" yaml:"rpcPort"`
	ENV         constant.ENV `json:"-" yaml:"-"`
	MysqlCfg    MysqlCfg     `json:"mysql" yaml:"mysqlCfg"`
	RedisCfg    RedisCfg     `json:"redis" yaml:"redisCfg"`
	NatsCfg     NatsCfg      `json:"nats" yaml:"natsCfg"`
	EtcdCfg     EtcdCfg      `json:"etcd" yaml:"etcdCfg"`
	Version     string       `json:"-" yaml:"-"`
}

type MysqlCfg struct {
	UserName string `json:"userName" yaml:"userName"`
	Password string `json:"password" yaml:"password"`
	password string
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Database string `json:"database" yaml:"database"`
}

func (c *MysqlCfg) MarshalJSON() ([]byte, error) {
	cn := c
	cn.Password = "***"
	return json.Marshal(&cn)
}
func (c *MysqlCfg) UnmarshalYAML(unmarshal func(interface{}) error) error {
	cn := c
	cn.Password = "***"
	return unmarshal(&cn)
}

type RedisCfg struct {
	Hosts    []string `json:"hosts"`
	Password string   `json:"password"`
	password string
	Port     int `json:"port"`
}

func (c *RedisCfg) MarshalJSON() ([]byte, error) {
	cn := c
	cn.Password = "***"
	return json.Marshal(&cn)
}

func (c *RedisCfg) UnmarshalYAML(unmarshal func(interface{}) error) error {
	cn := c
	cn.Password = "***"
	return unmarshal(&cn)
}

type NatsCfg struct {
	Host     string `json:"host"  yaml:"host"`
	UserName string `json:"userName" yaml:"userName"`
	Password string `json:"password" yaml:"password"`
	password string
	Port     int `json:"port" yaml:"port"`
}

func (c *NatsCfg) MarshalJSON() ([]byte, error) {
	cn := c
	cn.Password = "***"
	return json.Marshal(&cn)
}
func (c *NatsCfg) UnmarshalYAML(unmarshal func(interface{}) error) error {
	cn := c
	cn.Password = "***"
	return unmarshal(&cn)
}

type EtcdCfg struct {
	Hosts    []string `json:"hosts"  yaml:"hosts"`
	UserName string   `json:"userName" yaml:"userName"`
	Password string   `json:"password" yaml:"password"`
	password string
}

func (c *EtcdCfg) MarshalJSON() ([]byte, error) {
	cn := c
	cn.Password = "***"
	return json.Marshal(&cn)
}
func (c *EtcdCfg) UnmarshalYAML(unmarshal func(interface{}) error) error {
	cn := c
	cn.Password = "***"
	return unmarshal(&cn)
}
