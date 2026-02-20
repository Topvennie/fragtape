// Package fetch gets automatically new demo files
package fetch

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/config"
	"github.com/topvennie/fragtape/pkg/storage"
	"github.com/topvennie/fragtape/pkg/utils"
	"go.uber.org/zap"
)

type fetchResult struct {
	file     []byte
	source   model.DemoSource
	sourceID string
}

type fetcher interface {
	fetch(context.Context, model.UserDemo) (fetchResult, error)
}

type Manager struct {
	interval time.Duration
	cooldown time.Duration

	fetchers []fetcher

	repo repository.Repository
	demo repository.Demo
	user repository.User
}

func New(repo repository.Repository) *Manager {
	return &Manager{
		interval: config.GetDefaultDurationS("worker.fetcher.interval_s", 60),
		cooldown: config.GetDefaultDurationS("worker.fetcher.cooldown_s", 300),
		fetchers: []fetcher{},
		repo:     repo,
		demo:     *repo.NewDemo(),
		user:     *repo.NewUser(),
	}
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
	users, err := m.user.GetAllRealWithLastDemo(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	lastDemo := now.Add(-1 * m.cooldown)

	users = utils.SliceFilter(users, func(u *model.UserDemo) bool {
		if u.Demo.ID == 0 {
			return true
		}

		return u.Demo.CreatedAt.Before(lastDemo)
	})

	for _, user := range users {
		// Try each fetcher for a new match
		var result fetchResult

		for _, fetcher := range m.fetchers {
			result, err = fetcher.fetch(ctx, *user)
			if err != nil {
				return err
			}
			// Stop when we find a newer demo
			if len(result.file) > 0 && (result.source != user.Demo.Source || result.sourceID != user.Demo.SourceID) {
				break
			}
		}

		// If we have a result save it
		if len(result.file) > 0 && (result.source != user.Demo.Source || result.sourceID != user.Demo.SourceID) {
			demo := model.Demo{
				Source:   result.source,
				SourceID: result.sourceID,
				FileID:   uuid.NewString(),
			}

			if err := m.repo.WithRollback(ctx, func(ctx context.Context) error {
				if err := m.demo.Create(ctx, &demo); err != nil {
					return err
				}

				if err := storage.S.Set(demo.FileID, result.file, 0); err != nil {
					return err
				}

				return nil
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
