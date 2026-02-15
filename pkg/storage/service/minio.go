package minioservice

import (
	"context"
	"fmt"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rakhshon-mirzoev/gps-tracker/config"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"go.uber.org/zap"
)

type minioService struct {
	minio  *minio.Client
	bucket string
}

type Minio interface {
	PresignedUpload(ctx context.Context, objectName string, exp time.Duration) (string, error)
	PresignedDownload(ctx context.Context, objectName string, exp time.Duration) (string, error)
	UploadFile(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error)
	RemoveObject(ctx context.Context, objectName string) error
	MoveObject(ctx context.Context, srcKey, dstKey string) error
	RemoveOldTmpObjects(ctx context.Context, prefix string, olderThan time.Duration) error
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

	ctx := context.Background()

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &minioService{
		minio:  client,
		bucket: cfg.Bucket,
	}, nil
}
func NewMinioService(client *minio.Client, bucket string) Minio {
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
	reqParams := make(url.Values)
	reqParams.Set("host", "minio-gpsttracker.duckdns.org")

	u, err := s.minio.PresignedGetObject(
		ctx,
		s.bucket,
		objectName,
		exp,
		reqParams,
	)
	if err != nil {
		return "", err
	}

	u.Host = "minio-gpsttracker.duckdns.org"
	u.Scheme = "https"

	return u.String(), nil
}

func (s *minioService) MoveObject(ctx context.Context, srcKey, dstKey string) error {
	src := minio.CopySrcOptions{
		Bucket: s.bucket,
		Object: srcKey,
	}

	dst := minio.CopyDestOptions{
		Bucket: s.bucket,
		Object: dstKey,
	}

	_, err := s.minio.CopyObject(ctx, dst, src)
	if err != nil {
		return fmt.Errorf("copy object: %w", err)
	}

	err = s.minio.RemoveObject(ctx, s.bucket, srcKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("remove temp object: %w", err)
	}

	return nil
}

func (s *minioService) RemoveOldTmpObjects(ctx context.Context, prefix string, olderThan time.Duration) error {
	objects := s.minio.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objects {
		if object.Err != nil {
			if object.Err != nil {
				logger.Log.Error("list object error", zap.Error(object.Err))
				continue
			}
		}

		if time.Since(object.LastModified) > olderThan {
			if err := s.minio.RemoveObject(ctx, s.bucket, object.Key, minio.RemoveObjectOptions{}); err != nil {
				logger.Log.Error("failed to remove old object", zap.Error(err))
			}
		}
	}

	return nil
}

func (s *minioService) UploadFile(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read file header: %w", err)
	}

	contentType := http.DetectContentType(buffer[:n])

	if contentType != "image/jpeg" && contentType != "image/png" {
		return "", fmt.Errorf("only JPEG and PNG images are allowed")
	}

	if _, err := src.Seek(0, 0); err != nil {
		return "", fmt.Errorf("failed to seek file: %w", err)
	}

	retryConfig := RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     5 * time.Second,
	}

	var lastErr error
	backoff := retryConfig.InitialBackoff

	for attempt := 0; attempt < retryConfig.MaxRetries; attempt++ {
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

			if _, err := src.Seek(0, 0); err != nil {
				return "", fmt.Errorf("failed to seek file: %w", err)
			}
		}

		info, err := s.minio.PutObject(ctx, s.bucket, objectName, src, file.Size, minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err == nil {
			return info.Key, nil
		}

		lastErr = err

		if ctx.Err() != nil {
			return "", fmt.Errorf("upload failed: %w", ctx.Err())
		}
	}

	return "", fmt.Errorf("failed to upload file after %d attempts: %w", retryConfig.MaxRetries+1, lastErr)
}

func (s *minioService) RemoveObject(ctx context.Context, objectName string) error {
	return s.minio.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
}
