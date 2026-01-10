package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/topvennie/fragtape/pkg/sqlc"
)

type psql struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

// Interface compliance
var _ DB = (*psql)(nil)

type PostgresCfg struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func NewPSQL(cfg PostgresCfg) (DB, error) {
	pgConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, err
	}

	pgConfig.ConnConfig.Host = cfg.Host
	pgConfig.ConnConfig.Port = uint16(cfg.Port)
	pgConfig.ConnConfig.Database = cfg.Database
	pgConfig.ConnConfig.User = cfg.User
	pgConfig.ConnConfig.Password = cfg.Password

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.TODO()); err != nil {
		return nil, err
	}

	queries := sqlc.New(pool)

	return &psql{pool: pool, queries: queries}, nil
}

func (p *psql) WithRollback(ctx context.Context, fn func(q *sqlc.Queries) error) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		errRollback := tx.Rollback(ctx)
		if !errors.Is(err, pgx.ErrTxClosed) {
			err = errRollback
		}
	}()

	queries := sqlc.New(tx)

	if err := fn(queries); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *psql) Pool() *pgxpool.Pool {
	return p.pool
}

func (p *psql) Queries() *sqlc.Queries {
	return p.queries
}
