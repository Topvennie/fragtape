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

func (d *Demo) GetByUserFiltered(ctx context.Context, filter model.DemoFilter) (*model.DemoFilterResult, error) {
	params := sqlc.DemoGetByUserFilteredParams{
		UserID: int32(filter.UserID),
		Limit:  int32(filter.Limit),
		Offset: int32(filter.Offset),
	}

	demosDB, err := d.repo.queries(ctx).DemoGetByUserFiltered(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get demos by user filtered %+v | %w", filter, err)
	}

	demos := utils.SliceMap(demosDB, func(d sqlc.DemoGetByUserFilteredRow) model.Demo { return *model.DemoModel(d.Demo) })
	if len(demos) == 0 {
		return nil, nil
	}

	return &model.DemoFilterResult{
		Demos: demos,
		Total: int(demosDB[0].TotalCount),
	}, nil
}

func (d *Demo) GetBySourceSourceID(ctx context.Context, source model.DemoSource, sourceID string) (*model.Demo, error) {
	demo, err := d.repo.queries(ctx).DemoGetBySourceSourceID(ctx, sqlc.DemoGetBySourceSourceIDParams{
		Source:   sqlc.DemoSource(source),
		SourceID: sourceID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get demo by source %s and source id %s | %w", source, sourceID, err)
	}

	return model.DemoModel(demo), nil
}

func (d *Demo) GetByStatus(ctx context.Context, status model.DemoStatus) ([]*model.Demo, error) {
	demos, err := d.repo.queries(ctx).DemoGetByStatus(ctx, sqlc.DemoStatus(status))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get demos by status %s | %w", status, err)
	}

	return utils.SliceMap(demos, model.DemoModel), nil
}

func (d *Demo) GetByStatusUpdateAtomic(ctx context.Context, oldStatus, newStatus model.DemoStatus, amount int) ([]*model.Demo, error) {
	demos, err := d.repo.queries(ctx).DemoGetByStatusUpdateAtomic(ctx, sqlc.DemoGetByStatusUpdateAtomicParams{
		OldStatus: sqlc.DemoStatus(oldStatus),
		NewStatus: sqlc.DemoStatus(newStatus),
		Amount:    int32(amount),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get demos by status and update atomically %s -> %s | %d | %w", oldStatus, newStatus, amount, err)
	}

	return utils.SliceMap(demos, model.DemoModel), nil
}

func (d *Demo) Create(ctx context.Context, demo *model.Demo) error {
	id, err := d.repo.queries(ctx).DemoCreate(ctx, sqlc.DemoCreateParams{
		Source:    sqlc.DemoSource(demo.Source),
		SourceID:  demo.SourceID,
		SourceUrl: toString(demo.SourceURL),
		Status:    sqlc.DemoStatus(demo.Status),
		FileID:    toString(demo.FileID),
		PlayedAt:  toTime(demo.PlayedAt),
	})
	if err != nil {
		return fmt.Errorf("create demo %+v | %w", *demo, err)
	}

	demo.ID = int(id)

	return nil
}

func (d *Demo) UpdateStatus(ctx context.Context, demo model.Demo) error {
	if err := d.repo.queries(ctx).DemoUpdateStatus(ctx, sqlc.DemoUpdateStatusParams{
		ID:       int32(demo.ID),
		Status:   sqlc.DemoStatus(demo.Status),
		Error:    toString(demo.Error),
		Attempts: int32(demo.Attempts),
	}); err != nil {
		return fmt.Errorf("update demo status %+v | %w", demo, err)
	}

	return nil
}

func (d *Demo) UpdateFile(ctx context.Context, demo model.Demo) error {
	if err := d.repo.queries(ctx).DemoUpdateFile(ctx, sqlc.DemoUpdateFileParams{
		ID:     int32(demo.ID),
		FileID: toString(demo.FileID),
	}); err != nil {
		return fmt.Errorf("update demo file %+v | %w", demo, err)
	}

	return nil
}

func (d *Demo) UpdateData(ctx context.Context, demo model.Demo) error {
	if err := d.repo.queries(ctx).DemoUpdateData(ctx, sqlc.DemoUpdateDataParams{
		ID:     int32(demo.ID),
		DataID: toString(demo.DataID),
	}); err != nil {
		return fmt.Errorf("update demo data %+v | %w", demo, err)
	}

	return nil
}

func (d *Demo) ResetStatusAll(ctx context.Context, oldStatus, newStatus model.DemoStatus) error {
	if err := d.repo.queries(ctx).DemoResetStatusAll(ctx, sqlc.DemoResetStatusAllParams{
		OldStatus: sqlc.DemoStatus(oldStatus),
		NewStatus: sqlc.DemoStatus(newStatus),
	}); err != nil {
		return fmt.Errorf("reset demo status from %s to %s | %w", oldStatus, newStatus, err)
	}

	return nil
}
