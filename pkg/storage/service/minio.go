package minioservice

import (
	"context"
	"fmt"
	"math/rand"
	"mime/multipart"
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
	UploadFile(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error)
}
type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
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

func (s *minioService) UploadFile(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	retryConfig := RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     5 * time.Second,
	}

	var lastErr error
	backoff := retryConfig.InitialBackoff

	for attempt := 0; attempt <= retryConfig.MaxRetries; attempt++ {
		if attempt > 0 {

			jitter := time.Duration(rand.Int63n(int64(backoff / 2)))
			waitTime := backoff + jitter

			select {
			case <-time.After(waitTime):

			case <-ctx.Done():
				return "", fmt.Errorf("upload cancelled: %w", ctx.Err())
			}

			backoff = time.Duration(float64(backoff) * 1.5)
			if backoff > retryConfig.MaxBackoff {
				backoff = retryConfig.MaxBackoff
			}
		}

		if attempt > 0 {
			if _, err := src.Seek(0, 0); err != nil {
				return "", fmt.Errorf("failed to seek file: %w", err)
			}
		}

		info, err := s.minio.PutObject(ctx, s.bucket, objectName, src, file.Size, minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		})
		if err == nil {
			return info.Key, nil
		}

		lastErr = err

		if ctx.Err() != nil {
			return "", fmt.Errorf("upload failed: %w", ctx.Err())
		}

		if attempt < retryConfig.MaxRetries {
			fmt.Printf("Upload attempt %d failed: %v. Retrying...\n", attempt+1, err)
		}
	}

	return "", fmt.Errorf("failed to upload file after %d attempts: %w", retryConfig.MaxRetries+1, lastErr)
}
