package qiniu

import (
	"bytes"
	"context"
	"errors"

	"github.com/linshenqi/collapsar/src/services/base"
	"github.com/linshenqi/sptty"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

var Zones = map[string]storage.Zone{
	"huadong":  storage.ZoneHuadong,
	"huabei":   storage.ZoneHuabei,
	"huanan":   storage.ZoneHuanan,
	"beimei":   storage.ZoneBeimei,
	"xinjiapo": storage.ZoneXinjiapo,
}

type Client struct {
	Uploader *storage.FormUploader
	Manager  *storage.BucketManager
}

type Oss struct {
	base.BaseOss
	clients map[string]Client
}

func (s *Oss) getZone(zone string) (*storage.Zone, error) {
	z, exist := Zones[zone]
	if !exist {
		return nil, errors.New("Zone Not Found ")
	}

	return &z, nil
}

func (s *Oss) Init() {
	s.clients = map[string]Client{}
	for name, endpoint := range s.Endpoints {
		zone, err := s.getZone(endpoint.Zone)
		if err != nil {
			sptty.Log(sptty.ErrorLevel, err.Error(), base.Qiniu)
			continue
		}

		uploader := storage.NewFormUploader(&storage.Config{
			Zone:          zone,
			UseHTTPS:      false,
			UseCdnDomains: false,
		})

		mac := qbox.NewMac(endpoint.AppKey, endpoint.AppSecret)
		manager := storage.NewBucketManager(mac, &storage.Config{
			UseHTTPS: false,
		})

		s.clients[name] = Client{
			Uploader: uploader,
			Manager:  manager,
		}
	}
}

func (s *Oss) Upload(endpoint string, key string, data []byte) error {
	ep, err := s.GetEndpoint(endpoint)
	if err != nil {
		return err
	}

	client, exist := s.clients[endpoint]
	if !exist {
		return errors.New("Client Not Found ")
	}

	putPolicy := storage.PutPolicy{
		Scope: ep.Bucket,
	}

	mac := qbox.NewMac(ep.AppKey, ep.AppSecret)
	upToken := putPolicy.UploadToken(mac)

	l := int64(len(data))
	ret := storage.PutRet{}
	return client.Uploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(data), l, &storage.PutExtra{})
}

func (s *Oss) Delete(endpoint string, key string) error {
	ep, err := s.GetEndpoint(endpoint)
	if err != nil {
		return err
	}

	client, exist := s.clients[endpoint]
	if !exist {
		return errors.New("Client Not Found ")
	}

	return client.Manager.Delete(ep.Bucket, key)
}
