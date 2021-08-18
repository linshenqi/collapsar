package s3

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/linshenqi/collapsar/src/base"
	"github.com/linshenqi/sptty"
)

type Oss struct {
	base.BaseOss

	clients map[string]*s3.S3
}

func (s *Oss) Init() {
	s.clients = map[string]*s3.S3{}
	for name, endpoint := range s.Endpoints {
		session, err := session.NewSession(&aws.Config{
			Region:           aws.String(endpoint.Zone),
			Credentials:      credentials.NewStaticCredentials(endpoint.AppKey, endpoint.AppSecret, ""),
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(false),
		})

		if err != nil {
			sptty.Log(sptty.ErrorLevel, fmt.Sprintf("s3.Init.NewSession Failed: %s", err.Error()), "S3")
			continue
		}

		s.clients[name] = s3.New(session)
	}
}

func (s *Oss) Upload(endpoint string, key string, data []byte) error {
	ep, err := s.GetEndpoint(endpoint)
	if err != nil {
		return err
	}

	client, exist := s.clients[endpoint]
	if !exist {
		return fmt.Errorf("Client Not Found ")
	}

	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(ep.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Oss) Delete(endpoint string, key string) error {
	ep, err := s.GetEndpoint(endpoint)
	if err != nil {
		return err
	}

	client, exist := s.clients[endpoint]
	if !exist {
		return fmt.Errorf("Client Not Found ")
	}

	_, err = client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &ep.Bucket,
		Key:    aws.String(key),
	})

	if err != nil {
		return err
	}

	return nil
}