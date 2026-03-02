// Package fetch gets automatically new demo files
package fetch

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/worker/fetch/steam"
	"github.com/topvennie/fragtape/pkg/config"
	"github.com/topvennie/fragtape/pkg/storage"
	"github.com/topvennie/fragtape/pkg/utils"
	"go.uber.org/zap"
)

type Result struct {
	File     []byte
	Source   model.DemoSource
	SourceID string
}

type Fetcher interface {
	// Fetch gets a new demo file
	// It should only return a file if it is newer
	Fetch(context.Context, model.User) ([]byte, model.DemoSource, string, error)
}

type Manager struct {
	interval time.Duration
	cooldown time.Duration

	fetchers []Fetcher

	repo repository.Repository
	demo repository.Demo
	user repository.User
}

func New(repo repository.Repository) (*Manager, error) {
	if err := steam.Init(repo); err != nil {
		return nil, fmt.Errorf("init steam %w", err)
	}

	return &Manager{
		interval: config.GetDefaultDurationS("worker.fetcher.interval_s", 60),
		cooldown: config.GetDefaultDurationS("worker.fetcher.cooldown_s", 300),
		fetchers: []Fetcher{steam.S},
		repo:     repo,
		demo:     *repo.NewDemo(),
		user:     *repo.NewUser(),
	}, nil
}

func (m *Manager) Start(ctx context.Context) error {
	// Start the loop
	go func() {
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		for {
			if err := m.loop(ctx); err != nil {
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

func (m *Manager) loop(ctx context.Context) error {
	users, err := m.user.GetAllRealWithSettingLastDemo(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	lastDemo := now.Add(-1 * m.cooldown)

	users = utils.SliceFilter(users, func(u *model.User) bool {
		if u.Demo.ID == 0 {
			return true
		}

		return u.Demo.CreatedAt.Before(lastDemo)
	})

	for _, user := range users {
		zap.S().Debugf("Fetching for %s", user.DisplayName)
		// Try each fetcher for a new match
		var file []byte
		var source model.DemoSource
		var sourceID string

		for _, fetcher := range m.fetchers {
			file, source, sourceID, err = fetcher.Fetch(ctx, *user)
			if err != nil {
				zap.S().Error(err)
				continue
			}

			// Stop when we find a demo
			// The parser is responsible for making sure it is newer
			if len(file) > 0 {
				break
			}
		}

		// If we have a result save it
		if len(file) > 0 {
			zap.S().Debug("Demo received, saving")
			demo := model.Demo{
				Source:   source,
				SourceID: sourceID,
				FileID:   uuid.NewString(),
			}

			if err := m.repo.WithRollback(ctx, func(ctx context.Context) error {
				if err := m.demo.Create(ctx, &demo); err != nil {
					return err
				}

				if err := storage.S.Set(demo.FileID, file, 0); err != nil {
					return err
				}

				zap.S().Debug("All good")

				return nil
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
