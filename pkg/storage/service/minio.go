package minioservice

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rakhshon-mirzoev/gps-tracker/config"
)

type minioService struct {
	minio  *minio.Client
	bucket string
}

type Minio interface {
	PresignedUpload(ctx context.Context, objectName string, exp time.Duration) (string, error)
	PresignedDownload(ctx context.Context, objectName string, exp time.Duration) (string, error)
}

func NewMinio(cfg config.MinioConfig) (Minio, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		),
		Secure: cfg.UseSSL,
	})

	if err != nil {
		return nil, err
	}

	minio := InitMinioService(client, cfg.Bucket)

	return minio, nil
}

func InitMinioService(client *minio.Client, bucket string) Minio {
	return &minioService{
		minio:  client,
		bucket: bucket,
	}
}

func (s *minioService) PresignedUpload(ctx context.Context, objectName string, exp time.Duration) (string, error) {
	url, err := s.minio.PresignedPutObject(
		ctx,
		s.bucket,
		objectName,
		exp,
	)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s *minioService) PresignedDownload(ctx context.Context, objectName string, exp time.Duration) (string, error) {
	url, err := s.minio.PresignedGetObject(
		ctx,
		s.bucket,
		objectName,
		exp,
		nil,
	)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
