// Package storage connects with a file / image storage
package storage

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/minio"
)

type MinioCfg struct {
	Bucket    string
	Endpoint  string
	Secure    bool
	AccessKey string
	Secret    string
}

func Minio(cfg MinioCfg) fiber.Storage {
	return minio.New(minio.Config{
		Bucket:   cfg.Bucket,
		Endpoint: cfg.Endpoint,
		Secure:   cfg.Secure,
		Credentials: minio.Credentials{
			AccessKeyID:     cfg.AccessKey,
			SecretAccessKey: cfg.Secret,
		},
	})
}
