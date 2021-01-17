package base

import (
	"bytes"
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/linshenqi/sptty"
)

const (
	ServiceResource = "resource"
)

type IResourceService interface {
	CreateResources(resources []*Resource) error
	RemoveResourcesByIDs(ids []string) error
	GetResourcesByIDs(ids []string) ([]*Resource, error)
	GetResourcesByObjectID(objectID string) ([]*Resource, error)
	SetOssEndpoint(endpoint string)
	GetResourceUrl() string
}

type Resource struct {
	ID string `gorm:"primary_key,size:32" json:"id"`

	ObjectID string `gorm:"size:32" json:"-"`
	Name     string `json:"name"`
	Mime     string `json:"mime"`
	Data     []byte `gorm:"-" json:"-"`
}

func (s *Resource) Serialize(resUrl string) *Resource {

	s.Name = fmt.Sprintf("%s/%s", resUrl, s.Name)
	return s
}

func SerializeResources(resources []*Resource, resUrl string) []*Resource {
	for k := range resources {
		resources[k] = resources[k].Serialize(resUrl)
	}

	return resources
}

func GetResourcesFromRequestForm(ctx iris.Context) ([]*Resource, error) {
	form := ctx.Request().MultipartForm
	resources := []*Resource{}
	var err error = fmt.Errorf("No File")
	for _, v := range form.File {
		for _, file := range v {
			src, _ := file.Open()
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(src)
			if err != nil {
				break
			}

			mime := file.Header["Content-Type"][0]
			resources = append(resources, &Resource{
				Name: sptty.RandomFilename(file.Filename),
				Mime: mime,
				Data: buf.Bytes(),
			})
		}
	}

	if resources == nil {
		return nil, err
	}

	return resources, nil
}
