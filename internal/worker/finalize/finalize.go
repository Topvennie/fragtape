// Package finalize handles demos that have generated highlights
package finalize

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/config"
	"github.com/topvennie/fragtape/pkg/storage"
	"go.uber.org/zap"
)

const maxAttempts = 3

type Finalizer struct {
	demo      repository.Demo
	highlight repository.Highlight

	interval   time.Duration
	concurrent int
}

func New(repo repository.Repository) *Finalizer {
	return &Finalizer{
		demo:       *repo.NewDemo(),
		highlight:  *repo.NewHighlight(),
		interval:   config.GetDefaultDurationS("worker.interval_s.finalizer", 60),
		concurrent: config.GetDefaultInt("worker.concurrent.finalizer", 8),
	}
}

func (f *Finalizer) Start(ctx context.Context) error {
	if err := f.demo.ResetStatusAll(ctx, model.DemoStatusFinalizing, model.DemoStatusQueuedFinalize); err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(f.interval)
		defer ticker.Stop()

		for {
			if err := f.loop(ctx); err != nil {
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
	err error
}

func (f *Finalizer) loop(ctx context.Context) error {
	// Get demos
	// Their attempts counter is increased by the query
	demos, err := f.demo.GetByStatusUpdateAtomic(ctx, model.DemoStatusQueuedFinalize, model.DemoStatusFinalizing, f.concurrent)
	if err != nil {
		return err
	}

	for len(demos) > 0 {
		var wg sync.WaitGroup
		var mu sync.Mutex

		resultMap := make(map[int]loopResult)

		for _, demo := range demos {
			d := demo
			wg.Go(func() {
				err := f.loopOne(ctx, d)

				result := loopResult{
					err: err,
				}

				mu.Lock()
				defer mu.Unlock()

				resultMap[d.ID] = result
			})
		}

		// Wait until everything is finished
		wg.Wait()

		// Update demo status
		for _, demo := range demos {
			result, ok := resultMap[demo.ID]
			if !ok {
				// Shouldn't happen
				// But check just to be safe
				continue
			}

			if result.err != nil {
				// Something failed
				// Reset status
				demo.Error = result.err.Error()
				demo.Status = model.DemoStatusQueuedFinalize
				if demo.Attempts > maxAttempts {
					if err := storage.S.Delete(demo.FileID); err != nil {
						zap.S().Errorf("failed to delete demo file after max attempts reached for demo %+v | %w", *demo, err)
					}

					demo.Status = model.DemoStatusFailed
				}
			} else {
				// No errors
				// Go to the next step
				demo.Status = model.DemoStatusFinished
				demo.Attempts = 0
			}

			if err := f.demo.UpdateStatus(ctx, *demo); err != nil {
				return err
			}
		}

		// Keep finalizing demo's until the database no longer has waiting demo's
		demos, err = f.demo.GetByStatusUpdateAtomic(ctx, model.DemoStatusQueuedFinalize, model.DemoStatusFinalizing, f.concurrent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *Finalizer) loopOne(ctx context.Context, demo *model.Demo) error {
	_, err := f.highlight.GetByDemo(ctx, demo.ID)
	if err != nil {
		return err
	}

	// In the future generate thumbnail, convert to webm
	// Send to discord, ...

	if err := storage.S.Delete(demo.FileID); err != nil {
		return fmt.Errorf("failed to delete demo file for demo %+v | %w", *demo, err)
	}

	return nil
}
