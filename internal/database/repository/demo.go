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

type Demo struct {
	repo Repository
}

func (r *Repository) NewDemo() *Demo {
	return &Demo{
		repo: *r,
	}
}

func (d *Demo) Get(ctx context.Context, demoID int) (*model.Demo, error) {
	demo, err := d.repo.queries(ctx).DemoGet(ctx, int32(demoID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get demo %d | %w", demoID, err)
	}

	return model.DemoModel(demo), nil
}

func (d *Demo) GetByUser(ctx context.Context, userID int) ([]*model.Demo, error) {
	demos, err := d.repo.queries(ctx).DemoGetByUser(ctx, int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get demos by user %d | %w", userID, err)
	}

	return utils.SliceMap(demos, model.DemoModel), nil
}

func (d *Demo) Create(ctx context.Context, demo *model.Demo) error {
	id, err := d.repo.queries(ctx).DemoCreate(ctx, sqlc.DemoCreateParams{
		UserID:     int32(demo.UserID),
		Source:     sqlc.DemoSource(demo.Source),
		SourceID:   toString(demo.SourceID),
		DemoFileID: toString(demo.DemoFileID),
	})
	if err != nil {
		return fmt.Errorf("create demo %+v | %w", *demo, err)
	}

	demo.ID = int(id)

	return nil
}

func (d *Demo) Update(ctx context.Context, demo model.Demo) error {
	if err := d.repo.queries(ctx).DemoUpdateStatus(ctx, sqlc.DemoUpdateStatusParams{
		ID:     int32(demo.ID),
		Status: sqlc.DemoStatus(demo.Status),
	}); err != nil {
		return fmt.Errorf("update demo status %+v | %w", demo, err)
	}

	return nil
}

func (d *Demo) UpdateFile(ctx context.Context, demo model.Demo) error {
	if err := d.repo.queries(ctx).DemoUpdateFile(ctx, sqlc.DemoUpdateFileParams{
		ID:         int32(demo.ID),
		DemoFileID: toString(demo.DemoFileID),
	}); err != nil {
		return fmt.Errorf("update demo file %+v | %w", demo, err)
	}

	return nil
}

func (d *Demo) Delete(ctx context.Context, demoID int) error {
	if err := d.repo.queries(ctx).DemoDelete(ctx, int32(demoID)); err != nil {
		return fmt.Errorf("delete demo %d | %w", demoID, err)
	}

	return nil
}
