package mino

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rakhshon-mirzoev/gps-tracker/config"
)

func NewMinio(cfg config.MinioConfig) (*minio.Client, error) {
	return minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		),
		Secure: cfg.UseSSL,
	})
}
