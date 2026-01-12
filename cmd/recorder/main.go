package main

import (
	"context"
	"fmt"

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

	zapLogger, err := logger.New()
	if err != nil {
		panic(fmt.Errorf("initialize logger %w", err))
	}
	zap.ReplaceGlobals(zapLogger)

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

	recorder := recorder.New(*repo)
	if err := recorder.Start(context.Background()); err != nil {
		zap.S().Fatalf("Starting recorder failed %v", err)
	}

	zap.S().Info("Recorder is running")

	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────┐")
	fmt.Println("│            Fragtape Recorder            │")
	fmt.Println("│                                         │")
	fmt.Printf("│  Interval       %-23s │\n", config.GetDefaultDurationS("recorder.interval_s", 60))
	fmt.Println("└──────────────────────────────────────────┘")
	fmt.Println()

	// Wait indefinitely
	for {
		select {}
	}
}
