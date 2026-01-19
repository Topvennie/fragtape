package repository

import (
	"context"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/sqlc"
)

type Stat struct {
	repo Repository
}

func (r *Repository) NewStat() *Stat {
	return &Stat{
		repo: *r,
	}
}

func (s *Stat) Create(ctx context.Context, stat *model.Stat) error {
	id, err := s.repo.queries(ctx).StatCreate(ctx, sqlc.StatCreateParams{
		DemoID:  int32(stat.DemoID),
		UserID:  int32(stat.UserID),
		Kills:   int32(stat.Kills),
		Assists: int32(stat.Assists),
		Deaths:  int32(stat.Deaths),
	})
	if err != nil {
		return fmt.Errorf("create stat %+v | %w", *stat, err)
	}

	stat.ID = int(id)

	return nil
}
