// Package recorder renders highlights
package recorder

import (
	"context"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/recorder/capture"
	"github.com/topvennie/fragtape/pkg/config"
	"github.com/topvennie/fragtape/pkg/storage"
	"go.uber.org/zap"
)

const maxAttempts = 3

type Recorder struct {
	capturer capture.Capturer

	demo      repository.Demo
	highlight repository.Highlight

	interval time.Duration
}

func New(repo repository.Repository) (*Recorder, error) {
	capturer, err := capture.New(repo)
	if err != nil {
		return nil, err
	}

	return &Recorder{
		capturer:  *capturer,
		demo:      *repo.NewDemo(),
		highlight: *repo.NewHighlight(),
		interval:  config.GetDefaultDurationS("recorder.interval_s", 60),
	}, nil
}

// Start starts the loop to get new jobs and render them
func (r *Recorder) Start(ctx context.Context) error {
	if err := r.demo.ResetStatusAll(ctx, model.DemoStatusRendering, model.DemoStatusQueuedRender); err != nil {
		return err
	}

	// Start the loop
	go func() {
		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()

		for {
			if err := r.loop(ctx); err != nil {
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

func (r *Recorder) loop(ctx context.Context) error {
	// Get demos
	// Their attemps counter is increased by the query
	demos, err := r.demo.GetByStatusUpdateAtomic(ctx, model.DemoStatusQueuedRender, model.DemoStatusRendering, 1)
	if err != nil {
		return err
	}

	for len(demos) > 0 {
		demo := demos[0]

		err = r.capturer.Capture(ctx, *demo)
		if err != nil {
			// Something failed
			// Clean up
			if highlights, err := r.highlight.GetByDemo(ctx, demo.ID); err == nil {
				for _, h := range highlights {
					if h.FileID != "" {
						_ = r.highlight.DeleteFile(ctx, h.ID)
						_ = storage.S.Delete(h.FileID)
					}
				}
			}

			// Reset status
			demo.Error = err.Error()
			demo.Status = model.DemoStatusQueuedRender
			if demo.Attempts > maxAttempts {
				if err := storage.S.Delete(demo.FileID); err != nil {
					zap.S().Errorf("failed to delete demo file after max attempts reached for demo %+v | %w", *demo, err)
				}

				demo.Status = model.DemoStatusFailed
			}
		} else {
			// No errors
			// Go to the next step
			demo.Status = model.DemoStatusQueuedFinalize
			demo.Attempts = 0
		}

		if err := r.demo.UpdateStatus(ctx, *demo); err != nil {
			return err
		}

		demos, err = r.demo.GetByStatusUpdateAtomic(ctx, model.DemoStatusQueuedRender, model.DemoStatusRendering, 1)
		if err != nil {
			return err
		}
	}

	return nil
}
