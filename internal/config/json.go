package config

import (
	"b0b-common/internal/log"
	"encoding/json"
)

type JsonCfgHandle struct {
}

func (c *JsonCfgHandle) Read(data []byte) {
	err := json.Unmarshal(data, &Cfg)
	if err != nil {
		log.Panic("config json Unmarshal fail err=", err)
	}
}
