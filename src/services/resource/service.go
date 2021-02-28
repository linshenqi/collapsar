package resource

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/linshenqi/collapsar/src/services/base"
	"github.com/linshenqi/sptty"
)

type Service struct {
	sptty.BaseService
	base.IResourceService

	db  *gorm.DB
	oss base.IOssService
	cfg Config
}

func (s *Service) ServiceName() string {
	return base.ServiceResource
}

func (s *Service) Init(app sptty.ISptty) error {
	app.AddModel(&base.Resource{})

	if err := app.GetConfig(base.ServiceResource, &s.cfg); err != nil {
		return err
	}

	s.db = app.Model().(*sptty.ModelService).DB()
	if s.db == nil {
		return fmt.Errorf("Model Service Is Required")
	}

	s.oss = app.GetService(base.ServiceOss).(base.IOssService)
	if s.oss == nil {
		return fmt.Errorf("Oss Service Is Required")
	}

	return nil
}

func (s *Service) CreateResources(resources []*base.Resource) error {
	for k := range resources {
		if err := s.createResource(resources[k]); err != nil {
			sptty.Log(sptty.ErrorLevel, fmt.Sprintf("createResource Error: %s", err.Error()), s.ServiceName())
		}
	}

	return nil
}

func (s *Service) createResource(resource *base.Resource) error {
	resource.ID = sptty.GenerateUID()

	if err := s.saveResource(resource); err != nil {
		return err
	}

	if err := s.oss.Upload(s.cfg.OssEndpoint, resource.Name, resource.Data); err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveResourcesByIDs(ids []string) error {
	resources, err := s.getResourcesByIDs(ids)
	if err != nil {
		return err
	}

	for k := range resources {
		if err := s.removeResource(resources[k]); err != nil {
			sptty.Log(sptty.ErrorLevel, fmt.Sprintf("removeResource Error: %s", err.Error()), s.ServiceName())
		}
	}

	if err := s.removeResourcesByIDs(ids); err != nil {
		return err
	}

	return nil
}

func (s *Service) removeResource(resource *base.Resource) error {
	if err := s.oss.Delete(s.cfg.OssEndpoint, resource.Name); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetResourcesByIDs(ids []string) ([]*base.Resource, error) {
	resources, err := s.getResourcesByIDs(ids)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (s *Service) GetResourcesByObjectID(objectID string) ([]*base.Resource, error) {
	resources, err := s.getResourcesByObjectID(objectID)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

// decrepit
func (s *Service) SetOssEndpoint(endpoint string) {
}

func (s *Service) GetResourceUrl() string {
	return s.cfg.ResourceUrl
}
