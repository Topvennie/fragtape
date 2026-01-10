package storage

import (
	"github.com/gofiber/storage/minio"
)

type MinioCfg struct {
	Bucket    string
	Endpoint  string
	Secure    bool
	AccessKey string
	Secret    string
}

func Minio(cfg MinioCfg) {
	S = minio.New(minio.Config{
		Bucket:   cfg.Bucket,
		Endpoint: cfg.Endpoint,
		Secure:   cfg.Secure,
		Credentials: minio.Credentials{
			AccessKeyID:     cfg.AccessKey,
			SecretAccessKey: cfg.Secret,
		},
	})
}
