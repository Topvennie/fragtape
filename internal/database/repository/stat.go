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

func (s *Stat) GetByDemo(ctx context.Context, id int) ([]*model.Stat, error) {
	stats, err := s.repo.queries(ctx).StatGetByDemo(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get stats with demo id %d | %w", id, err)
	}

	return utils.SliceMap(stats, model.StatModel), nil
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
	result := sqlc.NullResult{Valid: false}
	if stat.Result != model.ResultUnknown {
		result.Result = sqlc.Result(stat.Result)
		result.Valid = true
	}
	startTeam := sqlc.NullTeam{Valid: false}
	if stat.StartTeam != model.TeamUnknown {
		startTeam.Team = sqlc.Team(stat.StartTeam)
		startTeam.Valid = true
	}

	id, err := s.repo.queries(ctx).StatCreate(ctx, sqlc.StatCreateParams{
		DemoID:    int32(stat.DemoID),
		UserID:    int32(stat.UserID),
		Result:    result,
		StartTeam: startTeam,
		Kills:     toInt(stat.Kills),
		Assists:   toInt(stat.Assists),
		Deaths:    toInt(stat.Deaths),
	})
	if err != nil {
		return fmt.Errorf("create stat %+v | %w", *stat, err)
	}

	stat.ID = int(id)

	return nil
}

func (s *Stat) Update(ctx context.Context, stat model.Stat) error {
	result := sqlc.NullResult{Valid: false}
	if stat.Result != model.ResultUnknown {
		result.Result = sqlc.Result(stat.Result)
		result.Valid = true
	}
	startTeam := sqlc.NullTeam{Valid: false}
	if stat.StartTeam != model.TeamUnknown {
		startTeam.Team = sqlc.Team(stat.StartTeam)
		startTeam.Valid = true
	}

	if err := s.repo.queries(ctx).StatUpdate(ctx, sqlc.StatUpdateParams{
		ID:        int32(stat.ID),
		Result:    result,
		StartTeam: startTeam,
		Kills:     toInt(stat.Kills),
		Assists:   toInt(stat.Assists),
		Deaths:    toInt(stat.Deaths),
	}); err != nil {
		return fmt.Errorf("update stat %+v | %w", stat, err)
	}

	return nil
}
