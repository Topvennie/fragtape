// Package parse gets unparsed demos, parses them to highlight segments and adds them to the job queue
package parse

import (
	"context"
	"sync"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/config"
	"go.uber.org/zap"
)

const maxAttempts = 3

type Parser struct {
	demo      repository.Demo
	highlight repository.Highlight
	repo      repository.Repository

	interval   time.Duration
	concurrent int
}

func New(repo repository.Repository) *Parser {
	return &Parser{
		demo:       *repo.NewDemo(),
		highlight:  *repo.NewHighlight(),
		repo:       repo,
		interval:   config.GetDefaultDurationS("worker.interval", 60),
		concurrent: config.GetDefaultInt("worker.concurrent", 8),
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
	highlights []model.Highlight
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
			// Get highlights
			d := demo
			wg.Go(func() {
				highlights, err := getHighlights(*d)

				result := loopResult{
					highlights: highlights,
					err:        err,
				}

				mu.Lock()
				defer mu.Unlock()

				resultMap[d.ID] = result
			})
		}

		// Wait until we have all highlights
		wg.Wait()

		// Update demo status and submit job
		for _, demo := range demos {
			result, ok := resultMap[demo.ID]
			if !ok {
				// Shouldn't happen
				// But check just to be safe
				continue
			}

			if result.err != nil {
				demo.Error = result.err.Error()
				demo.Status = model.DemoStatusQueuedParse
				if demo.Attempts > maxAttempts {
					demo.Status = model.DemoStatusFailed
				}

				if err := p.demo.UpdateStatus(ctx, *demo); err != nil {
					zap.S().Error(err)
				}
				continue
			}

			if len(result.highlights) == 0 {
				// No highlights found
				demo.Status = model.DemoStatusCompleted
				if err := p.demo.UpdateStatus(ctx, *demo); err != nil {
					zap.S().Error(err)
				}
				continue
			}

			if err := p.repo.WithRollback(ctx, func(ctx context.Context) error {
				// Prematurely update the status
				// Will be rolled back if anything fails in this function
				demo.Status = model.DemoStatusQueuedRender
				demo.Attempts = 0
				if err := p.demo.UpdateStatus(ctx, *demo); err != nil {
					return err
				}

				for i := range result.highlights {
					if err := p.highlight.Create(ctx, &result.highlights[i]); err != nil {
						return err
					}
				}

				if err := submitHighlights(result.highlights); err != nil {
					return err
				}

				return nil
			}); err != nil {
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
