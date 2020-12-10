package resource

import "github.com/linshenqi/collapsar/src/services/base"

func (s *Service) saveResource(resource *base.Resource) error {
	if err := s.db.Save(resource).Error; err != nil {
		return err
	}

	return nil
}

func (s *Service) getResourcesByIDs(ids []string) ([]*base.Resource, error) {
	resources := []*base.Resource{}
	if err := s.db.Where("id in (?)", ids).Order("id asc").Find(&resources).Error; err != nil {
		return nil, err
	}

	return resources, nil
}

func (s *Service) removeResourcesByIDs(ids []string) error {
	if err := s.db.Where("id in (?)", ids).Delete(base.Resource{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *Service) getResourcesByObjectID(objectID string) ([]*base.Resource, error) {
	resources := []*base.Resource{}
	if err := s.db.Where("object_id = ?", objectID).Order("id asc").Find(&resources).Error; err != nil {
		return nil, err
	}

	return resources, nil
}
