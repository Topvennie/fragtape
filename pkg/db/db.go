// Package db connects with the databank
package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/topvennie/flashback/pkg/sqlc"
)

type DB interface {
	WithRollback(ctx context.Context, fn func(q *sqlc.Queries) error) error
	Pool() *pgxpool.Pool
	Queries() *sqlc.Queries
}
