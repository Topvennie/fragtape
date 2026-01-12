package main

import (
	"fmt"

	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/server"
	"github.com/topvennie/fragtape/internal/server/service"
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
		Host:     config.GetDefaultString("server.db.host", "db"),
		Port:     config.GetDefaultInt("server.db.port", 5432),
		Database: config.GetDefaultString("server.db.database", "fragtape"),
		User:     config.GetDefaultString("server.db.user", "postgres"),
		Password: config.GetDefaultString("server.db.password", "postgres"),
	})
	if err != nil {
		zap.S().Fatalf("Unable to connect to db %v", err)
	}

	storage.Minio(storage.MinioCfg{
		Bucket:    config.GetDefaultString("server.minio.bucket", "fragtape"),
		Endpoint:  config.GetDefaultString("server.minio.endpoint", "minio:9000"),
		Secure:    config.GetDefaultBool("server.minio.secure", false),
		AccessKey: config.GetDefaultString("server.minio.username", "minio"),
		Secret:    config.GetDefaultString("server.minio.password", "miniominio"),
	})

	repo := repository.New(db)
	service := service.New(*repo)

	api := server.New(*service, db.Pool())

	zap.S().Infof("Server is running on %s", api.Addr)
	if err := api.Listen(api.Addr); err != nil {
		zap.S().Fatalf("Failure while running the server %v", err)
	}
}
