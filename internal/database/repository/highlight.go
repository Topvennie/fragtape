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

func (h *Highlight) Create(ctx context.Context, highlight *model.Highlight) error {
	id, err := h.repo.queries(ctx).HighlightCreate(ctx, sqlc.HighlightCreateParams{
		DemoID: int32(highlight.DemoID),
		FileID: toString(highlight.FileID),
		Title:  highlight.Title,
	})
	if err != nil {
		return fmt.Errorf("create highlight %+v | %w", *highlight, err)
	}

	highlight.ID = int(id)

	return nil
}

func (h *Highlight) UpdateFile(ctx context.Context, highlight model.Highlight) error {
	if err := h.repo.queries(ctx).HighlightUpdateFile(ctx, sqlc.HighlightUpdateFileParams{
		ID:     int32(highlight.ID),
		FileID: toString(highlight.FileID),
	}); err != nil {
		return fmt.Errorf("update highlight %+v | %w", highlight, err)
	}

	return nil
}
