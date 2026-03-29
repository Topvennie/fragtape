// Package recorder renders highlights
package recorder

import (
	"context"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/recorder/capture"
	"github.com/topvennie/fragtape/internal/status"
	"github.com/topvennie/fragtape/pkg/config"
	"github.com/topvennie/fragtape/pkg/storage"
	"go.uber.org/zap"
)

type Manager struct {
	capturer capture.Capturer

	demo      repository.Demo
	highlight repository.Highlight

	interval time.Duration
}

func New(repo repository.Repository) (*Manager, error) {
	capturer, err := capture.New(repo)
	if err != nil {
		return nil, err
	}

	return &Manager{
		capturer:  *capturer,
		demo:      *repo.NewDemo(),
		highlight: *repo.NewHighlight(),
		interval:  config.GetDefaultDurationS("recorder.interval_s", 60),
	}, nil
}

// Start starts the loop to get new jobs and render them
func (m *Manager) Start(ctx context.Context) error {
	if err := status.Demo.Reset(ctx, model.DemoStatusQueuedRender); err != nil {
		return err
	}

	// Start the loop
	go func() {
		for {
			empty, err := m.loop(ctx)
			if err != nil {
				zap.S().Error(err)
			}

			if empty {
				select {
				case <-ctx.Done():
					return
				case <-time.After(m.interval):
				}
			}

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	return nil
}

// loop handles one demo
// It returns true if there are no longer any demos that need to be handled
func (m *Manager) loop(ctx context.Context) (bool, error) {
	// Get demos
	// Their attemps counter is increased by the query
	demo, err := status.Demo.Get(ctx, model.DemoStatusQueuedRender)
	if err != nil {
		return false, err
	}
	if demo == nil {
		return true, nil
	}

	// Do the logic
	if err = func() error {
		if captureErr := m.capturer.Capture(ctx, *demo); captureErr != nil {
			// Best effort to clean up anything made
			highlights, err := m.highlight.GetByDemo(ctx, demo.ID)
			if err != nil {
				return captureErr
			}

			for _, h := range highlights {
				if h.FileID != "" {
					_ = storage.S.Delete(h.FileID)
					_ = m.highlight.DeleteFile(ctx, h.ID)
				}
			}

			return captureErr
		}

		return nil
	}(); err != nil {
		return false, status.Demo.Fail(ctx, demo, err)
	}

	return false, status.Demo.Succes(ctx, demo)
}
