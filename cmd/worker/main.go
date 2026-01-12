package main

import (
	"context"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/worker/parse"
	"github.com/topvennie/fragtape/pkg/config"
	"github.com/topvennie/fragtape/pkg/db"
	"github.com/topvennie/fragtape/pkg/logger"
	"github.com/topvennie/fragtape/pkg/storage"
	"go.uber.org/zap"
)

func main() {
	if err := config.Init(); err != nil {
		panic(fmt.Errorf("initialize config %w", err))
	}

	zapLogger, err := logger.New()
	if err != nil {
		panic(fmt.Errorf("initialize logger %w", err))
	}
	zap.ReplaceGlobals(zapLogger)

	db, err := db.NewPSQL(db.PostgresCfg{
		Host:     config.GetDefaultString("worker.db.host", "db"),
		Port:     config.GetDefaultInt("worker.db.port", 5432),
		Database: config.GetDefaultString("worker.db.database", "fragtape"),
		User:     config.GetDefaultString("worker.db.user", "postgres"),
		Password: config.GetDefaultString("worker.db.password", "postgres"),
	})
	if err != nil {
		zap.S().Fatalf("Unable to connect to db %v", err)
	}

	storage.Minio(storage.MinioCfg{
		Bucket:    config.GetDefaultString("worker.minio.bucket", "fragtape"),
		Endpoint:  config.GetDefaultString("worker.minio.endpoint", "minio:9000"),
		Secure:    config.GetDefaultBool("worker.minio.secure", false),
		AccessKey: config.GetDefaultString("worker.minio.username", "minio"),
		Secret:    config.GetDefaultString("worker.minio.password", "miniominio"),
	})

	repo := repository.New(db)

	parser, err := parse.New(*repo)
	if err != nil {
		zap.S().Fatalf("Initialize parser %w", err)
	}
	if err := parser.Start(context.Background()); err != nil {
		zap.S().Fatalf("Starting parser failed %v", err)
	}

	zap.S().Info("Worker  is running")

	fmt.Println()
	fmt.Println("┌──────────────────────────────────────────┐")
	fmt.Println("│              Fragtape Worker             │")
	fmt.Println("│                                          │")
	fmt.Printf("│  Interval       %-24s │\n", config.GetDefaultDurationS("worker.interval", 60))
	fmt.Printf("│  Concurrency    %-24d │\n", config.GetDefaultInt("worker.concurrent", 8))
	fmt.Println("└──────────────────────────────────────────┘")
	fmt.Println()

	// Wait indefinitely
	for {
		select {}
	}
}
