// Package recorder renders highlights
package recorder

import (
	"context"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/recorder/render"
	"github.com/topvennie/fragtape/pkg/config"
	"go.uber.org/zap"
)

const maxAttempts = 3

type Recorder struct {
	demo      repository.Demo
	highlight repository.Highlight

	interval time.Duration
}

func New(repo repository.Repository) *Recorder {
	render.Init(repo)

	return &Recorder{
		demo:      *repo.NewDemo(),
		highlight: *repo.NewHighlight(),
		interval:  config.GetDefaultDurationS("recorder.interval_s", 60),
	}
}

// Start starts the loop to get new jobs and render them
func (r *Recorder) Start(ctx context.Context) error {
	// Reset stuck demos
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

		err = render.C.Render(ctx, *demo)
		if err != nil {
			// Something failed
			// Reset status
			demo.Error = err.Error()
			demo.Status = model.DemoStatusQueuedRender
			if demo.Attempts > maxAttempts {
				demo.Status = model.DemoStatusFailed
			}
		} else {
			// No errors
			// Go to the next step
			demo.Status = model.DemoStatusRendered
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
