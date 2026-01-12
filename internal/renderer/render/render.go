// Package renderer converts a highligh job to an actual video
package renderer

import (
	"context"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/config"
	"go.uber.org/zap"
)

const maxAttempts = 3

type Render struct {
	demo      repository.Demo
	highlight repository.Highlight

	interval time.Duration
	dummy    bool
}

func New(repo repository.Repository) *Render {
	return &Render{
		demo:      *repo.NewDemo(),
		highlight: *repo.NewHighlight(),
		interval:  config.GetDefaultDurationS("renderer.interval_s", 60),
		dummy:     config.GetDefaultBool("renderer.dummy_data", false),
	}
}

// Start starts the loop to get new jobs and render them
func (r *Render) Start(ctx context.Context) error {
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

func (r *Render) loop(ctx context.Context) error {
	// Get demos
	// Their attemps counter is increased by the query
	demos, err := r.demo.GetByStatusUpdateAtomic(ctx, model.DemoStatusQueuedRender, model.DemoStatusRendering, 1)
	if err != nil {
		return err
	}

	for len(demos) > 0 {
		demo := demos[0]

		err := r.render(ctx, *demo)
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

// render will render the clips, save them to minio and save the file id
// to the right highlight in the render struct
func (r *Render) render(ctx context.Context, demo model.Demo) error {
	highlights, err := r.highlight.GetByDemo(ctx, demo.ID)
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
