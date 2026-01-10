package storage

import (
	"time"

	"github.com/gofiber/storage/postgres/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCfg struct {
	Pool *pgxpool.Pool
}

func NewPostgres(cfg PostgresCfg) {
	S = postgres.New(postgres.Config{
		DB:         cfg.Pool,
		Table:      "fragtape_files",
		Reset:      false,
		GCInterval: 10 * time.Second,
	})
}
