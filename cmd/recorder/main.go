package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/recorder"
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

	loggerFile := config.GetDefaultString("recorder.logger.file", "recorder")
	zapLogger, err := logger.New(logger.Config{
		Console: true,
		File:    loggerFile,
	})
	if err != nil {
		panic(fmt.Errorf("initialize logger %w", err))
	}
	zap.ReplaceGlobals(zapLogger)

	// Make sure that panics are logged to a file
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("panic",
				zap.Any("recover", r),
				zap.ByteString("stack", debug.Stack()),
			)

			_ = zap.L().Sync()

			os.Exit(1)
		}
	}()

	db, err := db.NewPSQL(db.PostgresCfg{
		Host:     config.GetDefaultString("recorder.db.host", "db"),
		Port:     config.GetDefaultInt("recorder.db.port", 5432),
		Database: config.GetDefaultString("recorder.db.database", "fragtape"),
		User:     config.GetDefaultString("recorder.db.user", "postgres"),
		Password: config.GetDefaultString("recorder.db.password", "postgres"),
	})
	if err != nil {
		zap.S().Fatalf("Unable to connect to db %v", err)
	}

	storage.Minio(storage.MinioCfg{
		Bucket:    config.GetDefaultString("recorder.minio.bucket", "fragtape"),
		Endpoint:  config.GetDefaultString("recorder.minio.endpoint", "minio:9000"),
		Secure:    config.GetDefaultBool("recorder.minio.secure", false),
		AccessKey: config.GetDefaultString("recorder.minio.username", "minio"),
		Secret:    config.GetDefaultString("recorder.minio.password", "miniominio"),
	})

	repo := repository.New(db)

	recorder, err := recorder.New(*repo)
	if err != nil {
		zap.S().Fatalf("Creating recorder failed %v", err)
	}
	if err := recorder.Start(context.Background()); err != nil {
		zap.S().Fatalf("Starting recorder failed %v", err)
	}

	zap.S().Info("Recorder is running")

	zap.S().Info()
	zap.S().Info("┌─────────────────────────────────────────┐")
	zap.S().Info("│            Fragtape Recorder            │")
	zap.S().Info("│                                         │")
	zap.S().Infof("│  Interval       %-23s │", config.GetDefaultDurationS("recorder.interval_s", 60))
	zap.S().Infof("│  Dummy          %-23t │", config.GetDefaultBool("recorder.dummy_data", false))
	zap.S().Info("└─────────────────────────────────────────┘")
	zap.S().Info()

	// Wait indefinitely
	for {
		select {}
	}
}
