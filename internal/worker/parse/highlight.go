package parse

import (
	"context"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/worker/parse/demo"
)

func (m *Manager) saveHighlights(ctx context.Context, d model.Demo, match demo.Match) error {
	if err := m.highlight.DeleteByDemo(ctx, d.ID); err != nil {
		return err
	}

	// Very simplified right now
	// It gets all 4k's
	highlights := []*model.Highlight{}

	for _, r := range match.Rounds {
		for player, stat := range r.PlayerStats {
			user, err := m.user.GetByUID(ctx, int(player))
			if err != nil {
				return err
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
					Duration: time.Duration(duration/int(match.TickRate)) * time.Second,
					Segments: segments,
				})
			}
		}
	}

	for _, highlight := range highlights {
		if err := m.highlight.Create(ctx, highlight); err != nil {
			return err
		}

		for _, segment := range highlight.Segments {
			segment.HighlightID = highlight.ID
			if err := m.highlight.CreateSegment(ctx, &segment); err != nil {
				return err
			}
		}
	}

	return nil
}
