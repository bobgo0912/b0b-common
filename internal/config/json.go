package config

import (
	"encoding/json"
	"github.com/bobgo0912/b0b-common/internal/log"
)

type JsonCfgHandle struct {
}

func (c *JsonCfgHandle) Read(data []byte) {
	err := json.Unmarshal(data, &Cfg)
	if err != nil {
		log.Panic("config json Unmarshal fail err=", err)
	}
}
