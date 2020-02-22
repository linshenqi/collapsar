package oss

import "github.com/linshenqi/collapsar/src/services/base"

type Config struct {
	Endpoints map[string]base.Endpoint `yaml:"endpoints"`
}

func (s *Config) ConfigName() string {
	return ServiceName
}

func (s *Config) Validate() error {
	return nil
}

func (s *Config) Default() interface{} {
	return &Config{}
}
