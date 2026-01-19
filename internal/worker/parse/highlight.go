package parse

import (
	"context"

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
				start := stat.Kills[0].Tick - 20
				end := stat.Kills[len(stat.Kills)-1].Tick + 20

				highlights = append(highlights, &model.Highlight{
					DemoID: d.ID,
					UserID: user.ID,
					Title:  "4k",
					Segments: []model.HighlightSegment{
						{
							StartTick: int(start),
							EndTick:   int(end),
						},
					},
				})
			}
		}
	}

	return highlights, nil
}
