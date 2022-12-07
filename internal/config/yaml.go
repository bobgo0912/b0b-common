package config

import (
	"github.com/bobgo0912/b0b-common/internal/log"
	"gopkg.in/yaml.v3"
)

type YamHandle struct {
}

func (y *YamHandle) Read(data []byte) {
	err := yaml.Unmarshal(data, &Cfg)
	if err != nil {
		log.Panic("config yaml Unmarshal fail err=", err)
	}
}
