package parse

import (
	"context"
	"errors"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/job"
	"github.com/topvennie/fragtape/pkg/storage"
)

// TODO: Take users into account
func (p *Parser) getHighlights(demo model.Demo) (job.Render, error) {
	var zero job.Render

	if demo.FileID == "" {
		return zero, errors.New("demo file deleted")
	}

	_, err := storage.S.Get(demo.FileID)
	if err != nil {
		return zero, fmt.Errorf("get demo file %w", err)
	}

	return zero, nil
}

func (p *Parser) submitHighlights(ctx context.Context, render job.Render) error {
	return p.renderQueue.Enqueue(ctx, render)
}
