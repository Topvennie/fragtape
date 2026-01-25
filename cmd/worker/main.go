package main

import (
	"context"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/worker/finalize"
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

	loggerFile := config.GetDefaultString("worker.logger.file", "")
	zapLogger, err := logger.New(logger.Config{
		Console: true,
		File:    loggerFile,
	})
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

	parser := parse.New(*repo)
	if err := parser.Start(context.Background()); err != nil {
		zap.S().Fatalf("Starting parser failed %v", err)
	}

	finalizer := finalize.New(*repo)
	if err := finalizer.Start(context.Background()); err != nil {
		zap.S().Fatalf("Starting finalizer failed %v", err)
	}

	zap.S().Info("Worker is running")

	zap.S().Info()
	zap.S().Info("┌──────────────────────────────────────────┐")
	zap.S().Info("│              Fragtape Worker             │")
	zap.S().Info("│                                          │")
	zap.S().Info("│  Parser                                  │")
	zap.S().Infof("│    Interval       %-22s │\n", config.GetDefaultDurationS("worker.parser.interval_s", 60))
	zap.S().Infof("│    Concurrency    %-22d │\n", config.GetDefaultInt("worker.parser.concurrent", 8))
	zap.S().Info("│  Finalizer                               │")
	zap.S().Infof("│    Interval       %-22s │\n", config.GetDefaultDurationS("worker.finalizer.interval_s", 60))
	zap.S().Infof("│    Concurrency    %-22d │\n", config.GetDefaultInt("worker.finalizer.concurrent", 8))
	zap.S().Info("└──────────────────────────────────────────┘")
	zap.S().Info()

	// Wait indefinitely
	for {
		select {}
	}
}
