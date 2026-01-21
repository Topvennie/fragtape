package repository

import (
	"context"
	"database/sql"
	"errors"
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

func (s *StatsDemo) GetByDemo(ctx context.Context, id int) (*model.StatsDemo, error) {
	stat, err := s.repo.queries(ctx).StatsDemoGetByDemo(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get stats demo with demo id %d | %w", id, err)
	}

	return model.StatsDemoModel(stat), nil
}

func (s *StatsDemo) Create(ctx context.Context, stat *model.StatsDemo) error {
	id, err := s.repo.queries(ctx).StatsDemoCreate(ctx, sqlc.StatsDemoCreateParams{
		DemoID:   int32(stat.DemoID),
		Map:      toString(stat.Map),
		RoundsCt: toInt(stat.RoundsCT),
		RoundsT:  toInt(stat.RoundsT),
	})
	if err != nil {
		return fmt.Errorf("create stat demo %+v | %w", *stat, err)
	}

	stat.ID = int(id)

	return nil
}

func (s *StatsDemo) Update(ctx context.Context, stat model.StatsDemo) error {
	if err := s.repo.queries(ctx).StatsDemoUpdate(ctx, sqlc.StatsDemoUpdateParams{
		ID:       int32(stat.ID),
		Map:      toString(stat.Map),
		RoundsCt: toInt(stat.RoundsCT),
		RoundsT:  toInt(stat.RoundsT),
	}); err != nil {
		return fmt.Errorf("update stat demo %+v | %w", stat, err)
	}

	return nil
}
