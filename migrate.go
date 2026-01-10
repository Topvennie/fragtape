package main

import (
	"embed"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/topvennie/fragtape/pkg/config"
	"github.com/topvennie/fragtape/pkg/db"
)

//go:embed db/migrations/*.sql
var embedMigrations embed.FS

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	// setup database
	db, err := db.NewPSQL(db.PostgresCfg{
		Host:     config.GetDefaultString("server.db.host", "db"),
		Port:     config.GetDefaultInt("server.db.port", 5432),
		Database: config.GetDefaultString("server.db.database", "fragtape"),
		User:     config.GetDefaultString("server.db.user", "postgres"),
		Password: config.GetDefaultString("server.db.password", "postgres"),
	})
	if err != nil {
		panic(err)
	}
	conn := stdlib.OpenDBFromPool(db.Pool())

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(conn, "db/migrations"); err != nil {
		panic(err)
	}

	// run app
}
