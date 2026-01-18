package service

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

type minioService struct {
	minio  *minio.Client
	bucket string
}

type Minio interface {
	Upload(ctx context.Context, name string, r io.Reader, size int64, contentType string) error
	PresignedUpload(ctx context.Context, objectName string, exp time.Duration) (string, error)
	PresignedDownload(ctx context.Context, objectName string, exp time.Duration) (string, error)
}

func InitMinioService(client *minio.Client, bucket string) Minio {
	return &minioService{
		minio:  client,
		bucket: bucket,
	}
}

func (s *minioService) Upload(ctx context.Context, name string, r io.Reader, size int64, contentType string) error {
	_, err := s.minio.PutObject(
		ctx,
		s.bucket,
		name,
		r,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	return err
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
