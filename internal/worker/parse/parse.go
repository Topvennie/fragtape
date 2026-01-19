// Package parse parses demos
package parse

// TODO: Cleanup data when the entire pipeline fails after x amount of attempts
// Right now it leaves some data behind

// TODO: Attempts go up one even when going to the next stage

import (
	"context"
	"sync"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/worker/parse/demo"
	"github.com/topvennie/fragtape/pkg/config"
	"go.uber.org/zap"
)

const maxAttempts = 3

type Parser struct {
	demo      repository.Demo
	highlight repository.Highlight
	stat      repository.Stat
	statsDemo repository.StatsDemo
	user      repository.User
	repo      repository.Repository

	demoParser demo.Demo

	interval   time.Duration
	concurrent int
}

func New(repo repository.Repository) *Parser {
	return &Parser{
		demo:      *repo.NewDemo(),
		highlight: *repo.NewHighlight(),
		stat:      *repo.NewStat(),
		statsDemo: *repo.NewStatsDemo(),
		user:      *repo.NewUser(),
		repo:      repo,
		demoParser: *demo.New(
			config.GetDefaultInt("worker.parser.positions_per_second", 4),
			config.GetDefaultInt("worker.parser.positions_min_distance", 10),
		),
		interval:   config.GetDefaultDurationS("worker.parser.interval_s", 60),
		concurrent: config.GetDefaultInt("worker.parser.concurrent", 8),
	}
}

// Start starts the loop to fetch and parse new demos
func (p *Parser) Start(ctx context.Context) error {
	// Reset stuck demos
	if err := p.demo.ResetStatusAll(ctx, model.DemoStatusParsing, model.DemoStatusQueuedParse); err != nil {
		return err
	}

	// Start the loop
	go func() {
		ticker := time.NewTicker(p.interval)
		defer ticker.Stop()

		for {
			if err := p.loop(ctx); err != nil {
				zap.S().Error(err)
			}

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()

	return nil
}

type loopResult struct {
	statsDemo  *model.StatsDemo
	stats      []*model.Stat
	highlights []*model.Highlight
	err        error
}

func (p *Parser) loop(ctx context.Context) error {
	// Get demos
	// Their attempts counter is increased by the query
	demos, err := p.demo.GetByStatusUpdateAtomic(ctx, model.DemoStatusQueuedParse, model.DemoStatusParsing, p.concurrent)
	if err != nil {
		return err
	}

	for len(demos) > 0 {
		var wg sync.WaitGroup
		var mu sync.Mutex

		resultMap := make(map[int]loopResult)

		for _, demo := range demos {
			// Multithread parsing the demo file and resulting data
			d := demo
			wg.Go(func() {
				save := func(result loopResult) {
					mu.Lock()
					defer mu.Unlock()

					resultMap[d.ID] = result
				}

				match, err := p.getMatch(ctx, d)
				if err != nil {
					save(loopResult{err: err})
					return
				}

				if err := p.savePlayers(ctx, *match); err != nil {
					save(loopResult{err: err})
					return
				}

				statsDemo := p.getStatsDemo(*d, *match)

				stats, err := p.getStats(ctx, *d, *match)
				if err != nil {
					save(loopResult{err: err})
					return
				}

				highlights, err := p.getHighlights(ctx, *d, *match)
				if err != nil {
					save(loopResult{err: err})
					return
				}

				save(loopResult{
					statsDemo:  statsDemo,
					stats:      stats,
					highlights: highlights,
				})
			})
		}

		// Wait until we have all highlights
		wg.Wait()

		// Update demo status
		for _, demo := range demos {
			result, ok := resultMap[demo.ID]
			if !ok {
				// Shouldn't happen
				// But check just to be safe
				continue
			}

			if err := p.repo.WithRollback(ctx, func(ctx context.Context) error {
				// DB transactions so unlikely something will fail
				// But nice safety just in case
				// Create stat and highlight db entries
				if result.err == nil {
					if err := p.statsDemo.Create(ctx, result.statsDemo); err != nil {
						return err
					}

					for i := range result.stats {
						if err := p.stat.Create(ctx, result.stats[i]); err != nil {
							return err
						}
					}

					for i := range result.highlights {
						if err := p.highlight.Create(ctx, result.highlights[i]); err != nil {
							return err
						}
						for j := range result.highlights[i].Segments {
							result.highlights[i].Segments[j].HighlightID = result.highlights[i].ID
							if err := p.highlight.CreateSegment(ctx, &result.highlights[i].Segments[j]); err != nil {
								return err
							}
						}
					}
				}

				if result.err != nil {
					// Something failed
					// Reset status
					demo.Error = result.err.Error()
					demo.Status = model.DemoStatusQueuedParse
					if demo.Attempts > maxAttempts {
						demo.Status = model.DemoStatusFailed
					}
				} else {
					// No errors
					// Go to the next step
					demo.Status = model.DemoStatusQueuedRender
					demo.Attempts = 0
				}

				if err := p.demo.UpdateStatus(ctx, *demo); err != nil {
					return err
				}

				return nil
			}); err != nil {
				// TODO: Fail if this errors???
				zap.S().Error(err)
			}
		}

		// Keep parsing demo's until the database no longer has unparsed demo's
		demos, err = p.demo.GetByStatusUpdateAtomic(ctx, model.DemoStatusQueuedParse, model.DemoStatusParsing, p.concurrent)
		if err != nil {
			return err
		}
	}

	return nil
}
