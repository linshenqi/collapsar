package oss

import (
	"errors"
	"github.com/linshenqi/collapsar/src/services/base"
	"github.com/linshenqi/collapsar/src/services/qiniu"
	"github.com/linshenqi/sptty"
)

const (
	ServiceName = "oss"
)

type Service struct {
	cfg       Config
	providers map[string]base.IOss
}

func (s *Service) Init(app sptty.Sptty) error {
	if err := app.GetConfig(s.ServiceName(), &s.cfg); err != nil {
		return err
	}

	s.setupProviders()
	s.initProviders()

	return nil
}

func (s *Service) Release() {

}

func (s *Service) Enable() bool {
	return true
}

func (s *Service) ServiceName() string {
	return ServiceName
}

func (s *Service) initProviders() {
	for k, v := range s.cfg.Endpoints {
		provider, err := s.getProvider(v.Provider)
		if err != nil {
			continue
		}

		provider.AddEndpoint(k, v)
	}

	for _, provider := range s.providers {
		provider.Init()
	}
}

func (s *Service) getProvider(providerType string) (base.IOss, error) {
	provider, exist := s.providers[providerType]
	if !exist {
		return nil, errors.New("Provider Not Found ")
	}

	return provider, nil
}

func (s *Service) getEndpoint(endpoint string) (*base.Endpoint, error) {
	ep, exist := s.cfg.Endpoints[endpoint]
	if !exist {
		return nil, errors.New("Endpoint Not Found ")
	}

	return &ep, nil
}

func (s *Service) setupProviders() {
	s.providers = map[string]base.IOss{
		base.Qiniu: &qiniu.Oss{},
	}
}

func (s *Service) Upload(endpoint string, key string, data []byte) error {
	ep, err := s.getEndpoint(endpoint)
	if err != nil {
		return err
	}

	provider, err := s.getProvider(ep.Provider)
	if err != nil {
		return err
	}

	return provider.Upload(endpoint, key, data)
}

func (s *Service) Delete(endpoint string, key string) error {
	ep, err := s.getEndpoint(endpoint)
	if err != nil {
		return err
	}

	provider, err := s.getProvider(ep.Provider)
	if err != nil {
		return err
	}

	return provider.Delete(endpoint, key)
}
