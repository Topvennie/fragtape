package repository

import (
	"context"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/sqlc"
)

type StatsDemo struct {
	repo Repository
}

func (r *Repository) NewStatsDemo() *StatsDemo {
	return &StatsDemo{
		repo: *r,
	}
}

func (s *StatsDemo) Create(ctx context.Context, stat *model.StatsDemo) error {
	id, err := s.repo.queries(ctx).StatsDemoCreate(ctx, sqlc.StatsDemoCreateParams{
		DemoID:   int32(stat.DemoID),
		Map:      stat.Map,
		RoundsCt: int32(stat.RoundsCT),
		RoundsT:  int32(stat.RoundsT),
	})
	if err != nil {
		return fmt.Errorf("create stat demo %+v | %w", *stat, err)
	}

	stat.ID = int(id)

	return nil
}
