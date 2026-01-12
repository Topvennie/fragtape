package repository

import (
	"context"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/sqlc"
)

type Highlight struct {
	repo Repository
}

func (r *Repository) NewHighlight() *Highlight {
	return &Highlight{
		repo: *r,
	}
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
