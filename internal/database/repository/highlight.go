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

type Highlight struct {
	repo Repository
}

func (r *Repository) NewHighlight() *Highlight {
	return &Highlight{
		repo: *r,
	}
}

func (h *Highlight) GetByDemo(ctx context.Context, demoID int) ([]*model.Highlight, error) {
	highlights, err := h.repo.queries(ctx).HighlightGetByDemo(ctx, int32(demoID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get highlights by demo %d | %w", demoID, err)
	}

	return utils.SliceMap(highlights, model.HighlightModel), nil
}

func (h *Highlight) GetByDemos(ctx context.Context, demoIDs []int) ([]*model.Highlight, error) {
	highlights, err := h.repo.queries(ctx).HighlightGetByDemos(ctx, utils.SliceMap(demoIDs, func(id int) int32 { return int32(id) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get highlights by demos %+v | %w", demoIDs, err)
	}

	return utils.SliceMap(highlights, model.HighlightModel), nil
}

func (h *Highlight) Create(ctx context.Context, highlight *model.Highlight) error {
	id, err := h.repo.queries(ctx).HighlightCreate(ctx, sqlc.HighlightCreateParams{
		UserID: int32(highlight.UserID),
		DemoID: int32(highlight.DemoID),
		Title:  highlight.Title,
	})
	if err != nil {
		return fmt.Errorf("create highlight %+v | %w", *highlight, err)
	}

	highlight.ID = int(id)

	return nil
}

func (h *Highlight) CreateSegment(ctx context.Context, segment *model.HighlightSegment) error {
	id, err := h.repo.queries(ctx).HighlightSegmentCreate(ctx, sqlc.HighlightSegmentCreateParams{
		HighlightID: int32(segment.HighlightID),
		StartTick:   int32(segment.StartTick),
		EndTick:     int32(segment.EndTick),
	})
	if err != nil {
		return fmt.Errorf("create highlight segment %+v | %w", *segment, err)
	}

	segment.ID = int(id)

	return nil
}

func (h *Highlight) Update(ctx context.Context, highlight model.Highlight) error {
	if err := h.repo.queries(ctx).HighlightUpdate(ctx, sqlc.HighlightUpdateParams{
		ID:        int32(highlight.ID),
		DemoID:    toInt(highlight.DemoID),
		FileID:    toString(highlight.FileID),
		FileWebID: toString(highlight.FileWebID),
		Title:     toString(highlight.Title),
	}); err != nil {
		return fmt.Errorf("update highlight %+v | %w", highlight, err)
	}

	return nil
}

func (h *Highlight) DeleteFile(ctx context.Context, highlightID int) error {
	if err := h.repo.queries(ctx).HighlightDeleteFile(ctx, int32(highlightID)); err != nil {
		return fmt.Errorf("delete highlight file %d | %w", highlightID, err)
	}

	return nil
}
