package main

import (
	"context"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/repository"
	renderer "github.com/topvennie/fragtape/internal/renderer/render"
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
		Host:     config.GetDefaultString("renderer.db.host", "db"),
		Port:     config.GetDefaultInt("renderer.db.port", 5432),
		Database: config.GetDefaultString("renderer.db.database", "fragtape"),
		User:     config.GetDefaultString("renderer.db.user", "postgres"),
		Password: config.GetDefaultString("renderer.db.password", "postgres"),
	})
	if err != nil {
		zap.S().Fatalf("Unable to connect to db %v", err)
	}

	storage.Minio(storage.MinioCfg{
		Bucket:    config.GetDefaultString("renderer.minio.bucket", "fragtape"),
		Endpoint:  config.GetDefaultString("renderer.minio.endpoint", "minio:9000"),
		Secure:    config.GetDefaultBool("renderer.minio.secure", false),
		AccessKey: config.GetDefaultString("renderer.minio.username", "minio"),
		Secret:    config.GetDefaultString("renderer.minio.password", "miniominio"),
	})

	repo := repository.New(db)

	renderer := renderer.New(*repo)
	if err := renderer.Start(context.Background()); err != nil {
		zap.S().Fatalf("Starting renderer failed %v", err)
	}

	zap.S().Info("Renderer is running")

	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────┐")
	fmt.Println("│            Fragtape Renderer            │")
	fmt.Println("│                                         │")
	fmt.Printf("│  Interval       %-23s │\n", config.GetDefaultDurationS("renderer.interval_s", 60))
	fmt.Println("└──────────────────────────────────────────┘")
	fmt.Println()

	// Wait indefinitely
	for {
		select {}
	}
}
