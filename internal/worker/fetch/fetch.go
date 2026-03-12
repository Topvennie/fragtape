// Package fetch gets automatically new demo files
package fetch

import (
	"context"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/status"
	"github.com/topvennie/fragtape/internal/worker/fetch/faceit"
	"github.com/topvennie/fragtape/internal/worker/fetch/steam"
	"github.com/topvennie/fragtape/pkg/config"
	"github.com/topvennie/fragtape/pkg/utils"
	"go.uber.org/zap"
)

type Fetcher interface {
	// Fetch gets a new demo file url
	// It returns
	//  model.Demo -> the new demo
	//  bool -> indicating if there is a new demo
	//  error -> error
	Fetch(context.Context, model.User) (model.Demo, bool, error)
}

type Manager struct {
	interval time.Duration
	cooldown time.Duration

	fetchers []Fetcher

	repo repository.Repository
	demo repository.Demo
	stat repository.Stat
	user repository.User
}

func New(repo repository.Repository) *Manager {
	return &Manager{
		interval: config.GetDefaultDurationS("worker.fetcher.interval_s", 300),
		cooldown: config.GetDefaultDurationS("worker.fetcher.cooldown_s", 600),
		fetchers: []Fetcher{steam.S, faceit.F},
		repo:     repo,
		demo:     *repo.NewDemo(),
		stat:     *repo.NewStat(),
		user:     *repo.NewUser(),
	}
}

func (m *Manager) Start(ctx context.Context) error {
	// Reset statusses
	if err := status.Demo.Reset(ctx, model.DemoStatusQueuedParse); err != nil {
		return err
	}

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
		// Try each fetcher for a new match
		for _, fetcher := range m.fetchers {
			demo, ok, err := fetcher.Fetch(ctx, *user)
			if err != nil {
				zap.S().Error(err)
				continue
			}

			// The parser is responsible for making sure it is newer
			if !ok {
				continue
			}

			// Does this demo already exist?
			oldDemo, err := m.demo.GetBySourceSourceID(ctx, demo.Source, demo.SourceID)
			if err != nil {
				zap.S().Error(err)
				continue
			}
			if oldDemo != nil {
				// Demo already exists
				// Possible if a different user was also in the match and got handled first
				demo.ID = oldDemo.ID
			} else {
				// New demo
				demo.Status = model.DemoStatusQueuedDownload

				if err := m.demo.Create(ctx, &demo); err != nil {
					zap.S().Error(err)
					continue
				}
			}

			// Add all players (if any)
			// Use the no conflict create as we only insert the player id
			// We want to do it atomically and don't overwrite any existing data
			for _, stat := range demo.Stats {
				stat.DemoID = demo.ID
				if err := m.stat.CreateNoConflict(ctx, &stat); err != nil {
					zap.S().Error(err)
				}
			}
		}
	}

	return nil
}
