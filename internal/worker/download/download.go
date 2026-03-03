// Package download downloads demo files
package download

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/status"
	"github.com/topvennie/fragtape/internal/worker/fetch/steam"
	"github.com/topvennie/fragtape/pkg/config"
	"github.com/topvennie/fragtape/pkg/storage"
	"go.uber.org/zap"
)

type Manager struct {
	repo repository.Repository
	demo repository.Demo

	interval   time.Duration
	concurrent int

	wg sync.WaitGroup
}

func New(repo repository.Repository) *Manager {
	return &Manager{
		repo:       repo,
		demo:       *repo.NewDemo(),
		interval:   config.GetDefaultDurationS("worker.downloader.interval_s", 60),
		concurrent: config.GetDefaultInt("worker.downloader.concurrent", 8),
		wg:         sync.WaitGroup{},
	}
}

func (m *Manager) Start(ctx context.Context) error {
	// Reset statusses
	if err := status.Demo.Reset(ctx, model.DemoStatusQueuedDownload); err != nil {
		return err
	}

	// Add the goroutines that will fetch downloads
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
	// Get a new demo
	demo, err := status.Demo.Get(ctx, model.DemoStatusQueuedDownload)
	if err != nil {
		return false, err
	}
	if demo == nil {
		return true, nil
	}

	zap.S().Debug("Downloading")

	// Get the new demo
	if err = func() error {
		if demo.SourceURL == "" {
			return fmt.Errorf("demo without source url  %+v", *demo)
		}

		var file []byte

		switch demo.Source {
		case model.DemoSourceSteam:
			file, err = steam.S.Download(ctx, demo.SourceURL)
		default:
		}

		if err != nil {
			return err
		}
		if len(file) == 0 {
			return fmt.Errorf("no file downloaded for demo %+v", *demo)
		}

		return m.repo.WithRollback(ctx, func(ctx context.Context) error {
			demo.FileID = uuid.NewString()

			if err := m.demo.UpdateFile(ctx, *demo); err != nil {
				return err
			}

			if err := storage.S.Set(demo.FileID, file, 0); err != nil {
				return err
			}

			return nil
		})
	}(); err != nil {
		return false, status.Demo.Fail(ctx, demo, err)
	}

	return false, status.Demo.Succes(ctx, demo)
}
