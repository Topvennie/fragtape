// Package finalize handles demos that have generated highlights
package finalize

import (
	"context"
	"sync"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/status"
	"github.com/topvennie/fragtape/pkg/config"
	"github.com/topvennie/fragtape/pkg/storage"
	"go.uber.org/zap"
)

type Manager struct {
	demo      repository.Demo
	highlight repository.Highlight

	interval   time.Duration
	concurrent int

	wg sync.WaitGroup
}

func New(repo repository.Repository) *Manager {
	return &Manager{
		demo:       *repo.NewDemo(),
		highlight:  *repo.NewHighlight(),
		interval:   config.GetDefaultDurationS("worker.finalizer.interval_s", 60),
		concurrent: config.GetDefaultInt("worker.finalizer.concurrent", 8),
		wg:         sync.WaitGroup{},
	}
}

func (m *Manager) Start(ctx context.Context) error {
	if err := status.Demo.Reset(ctx, model.DemoStatusQueuedFinalize); err != nil {
		return err
	}

	for range m.concurrent {
		m.wg.Go(func() {
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
		})
	}

	return nil
}

// loop handles one demo
// It returns a boolean indicating if there are potentially more demos to be handled
func (m *Manager) loop(ctx context.Context) (bool, error) {
	demo, err := status.Demo.Get(ctx, model.DemoStatusQueuedFinalize)
	if err != nil {
		return false, err
	}
	if demo == nil {
		return true, nil
	}

	zap.S().Debug("Finalizing")

	if err = func() error {
		if demo.FileID != "" {
			// Best effort
			_ = storage.S.Delete(demo.FileID)
			demo.FileID = ""
			_ = m.demo.UpdateFile(ctx, *demo)
		}

		return nil
	}(); err != nil {
		return false, status.Demo.Fail(ctx, demo, err)
	}

	return false, status.Demo.Succes(ctx, demo)
}
