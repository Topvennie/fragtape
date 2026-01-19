package parse

import (
	"context"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/worker/parse/demo"
)

func (p *Parser) getHighlights(ctx context.Context, d model.Demo, m demo.Match) ([]*model.Highlight, error) {
	return nil, nil
}
