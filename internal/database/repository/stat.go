package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/sqlc"
	"github.com/topvennie/fragtape/pkg/utils"
)

type Stat struct {
	repo Repository
}

func (r *Repository) NewStat() *Stat {
	return &Stat{
		repo: *r,
	}
}

func (s *Stat) GetByDemos(ctx context.Context, demoIDs []int) ([]*model.Stat, error) {
	stats, err := s.repo.queries(ctx).StatGetByDemos(ctx, utils.SliceMap(demoIDs, func(id int) int32 { return int32(id) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get stats with demo id %+v | %w", demoIDs, err)
	}

	return utils.SliceMap(stats, model.StatModel), nil
}

func (s *Stat) Create(ctx context.Context, stat *model.Stat) error {
	id, err := s.repo.queries(ctx).StatCreate(ctx, sqlc.StatCreateParams{
		DemoID:    int32(stat.DemoID),
		UserID:    int32(stat.UserID),
		Result:    sqlc.Result(stat.Result),
		StartTeam: sqlc.Team(stat.StartTeam),
		Kills:     int32(stat.Kills),
		Assists:   int32(stat.Assists),
		Deaths:    int32(stat.Deaths),
	})
	if err != nil {
		return fmt.Errorf("create stat %+v | %w", *stat, err)
	}

	stat.ID = int(id)

	return nil
}
