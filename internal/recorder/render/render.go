// Package render converts a highlight to an actual video
package render

import (
	"context"

	"github.com/topvennie/fragtape/internal/database/model"
)

func (c *client) Render(ctx context.Context, demo model.Demo) error {
	highlights, err := c.highlight.GetByDemo(ctx, demo.ID)
	if err != nil {
		return err
	}

	if len(highlights) == 0 {
		// No highlights
		return nil
	}

	// Actually render it now

	return nil
}
