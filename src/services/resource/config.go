package resource

import (
	"github.com/linshenqi/collapsar/src/services/base"
	"github.com/linshenqi/sptty"
)

type Config struct {
	sptty.BaseConfig
	Url         string `yaml:"url"`
	OssEndpoint string `yaml:"oss_endpoint"`
}

func (s *Config) ConfigName() string {
	return base.ServiceResource
}
