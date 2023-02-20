package service

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Service struct {
	cfg      client.ConfigProvider
	bucket   string
	endpoint string
}

func NewS3Service(region, endpoint, secretId, secretKey, bucket string) (*S3Service, error) {
	conf, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Endpoint:    &endpoint,
		Credentials: credentials.NewStaticCredentials(secretId, secretKey, ""),
	})
	if err != nil {
		return nil, err
	}

	return &S3Service{
		cfg:      conf,
		bucket:   bucket,
		endpoint: endpoint,
	}, nil
}

// SaveFile 保存文件至S3
func (s *S3Service) SaveFile(filePath string, file []byte) (string, error) {
	service := s3.New(s.cfg)
	_, err := service.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
		Body:   bytes.NewReader(file),
	})

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s/%s", s.bucket, s.endpoint, filePath), nil
}

// GetFile 获取文件
func (s *S3Service) GetFile(filePath string) (io.Reader, error) {
	service := s3.New(s.cfg)

	object, err := service.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return nil, err
	}

	return object.Body, nil
}

func (s *S3Service) DeleteFile(filePath string) error {
	service := s3.New(s.cfg)

	res, err := service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return err
	}

	if res.DeleteMarker != nil {
		return fmt.Errorf("delete file failed")
	}
	return nil
}
