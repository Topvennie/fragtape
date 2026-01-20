package parse

import (
	"context"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/worker/parse/demo"
)

func (p *Parser) getHighlights(ctx context.Context, d model.Demo, m demo.Match) ([]*model.Highlight, error) {
	// TODO: Delete existing highlights

	// Very simplified right now
	// It gets all 4k's
	highlights := []*model.Highlight{}

	for _, r := range m.Rounds {
		for player, stat := range r.PlayerStats {
			user, err := p.user.GetByUID(ctx, int(player))
			if err != nil {
				return nil, err
			}
			if user == nil {
				// Shouldn't be possible
				continue
			}

			if len(stat.Kills) >= 4 {
				duration := 0

				segments := make([]model.HighlightSegment, 0, len(stat.Kills))
				for _, k := range stat.Kills {
					start := int(k.Tick) - 265
					end := int(k.Tick) + 265

					segments = append(segments, model.HighlightSegment{
						StartTick: start,
						EndTick:   end,
					})

					duration += end - start
				}

				highlights = append(highlights, &model.Highlight{
					DemoID:   d.ID,
					UserID:   user.ID,
					Title:    "4k",
					Round:    r.Number,
					Duration: time.Duration(duration/int(m.TickRate)) * time.Second,
					Segments: segments,
				})
			}
		}
	}

	return highlights, nil
}
