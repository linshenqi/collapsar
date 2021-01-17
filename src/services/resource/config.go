package resource

import (
	"github.com/linshenqi/collapsar/src/services/base"
	"github.com/linshenqi/sptty"
)

type Config struct {
	sptty.BaseConfig
	ResourceUrl string `yaml:"resource_url"`
}

func (s *Config) ConfigName() string {
	return base.ServiceResource
}
