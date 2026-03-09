// Package parse parses demos
package parse

import (
	"context"
	"sync"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/status"
	"github.com/topvennie/fragtape/internal/worker/parse/demo"
	"github.com/topvennie/fragtape/pkg/config"
	"go.uber.org/zap"
)

type Manager struct {
	demo      repository.Demo
	highlight repository.Highlight
	setting   repository.SettingGlobal
	stat      repository.Stat
	statsDemo repository.StatsDemo
	user      repository.User
	repo      repository.Repository

	demoParser demo.Demo

	interval   time.Duration
	concurrent int

	wg sync.WaitGroup
}

func New(repo repository.Repository) *Manager {
	return &Manager{
		demo:      *repo.NewDemo(),
		highlight: *repo.NewHighlight(),
		setting:   *repo.NewSettingGlobal(),
		stat:      *repo.NewStat(),
		statsDemo: *repo.NewStatsDemo(),
		user:      *repo.NewUser(),
		repo:      repo,
		demoParser: *demo.New(
			config.GetDefaultInt("worker.parser.positions_per_second", 4),
			config.GetDefaultInt("worker.parser.positions_min_distance", 10),
		),
		interval:   config.GetDefaultDurationS("worker.parser.interval_s", 60),
		concurrent: config.GetDefaultInt("worker.parser.concurrent", 8),
		wg:         sync.WaitGroup{},
	}
}

// Start starts the loop to fetch and parse new demos
func (m *Manager) Start(ctx context.Context) error {
	// Reset stuck demos
	if err := status.Demo.Reset(ctx, model.DemoStatusQueuedParse); err != nil {
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
	// Get a demo
	demo, err := status.Demo.Get(ctx, model.DemoStatusQueuedParse)
	if err != nil {
		return false, err
	}
	if demo == nil {
		return true, nil
	}

	// Do the logic
	if err = func() error {
		// Parse the match and get the raw match data
		match, err := m.getMatch(ctx, demo)
		if err != nil {
			return err
		}

		// Make sure we have saved all participants
		if err := m.savePlayers(ctx, *match); err != nil {
			return err
		}

		// Save some generic match information
		if err := m.saveStatsDemo(ctx, *demo, *match); err != nil {
			return err
		}

		// Save player statistics
		if err := m.saveStats(ctx, *demo, *match); err != nil {
			return err
		}

		// Save the highight segments
		if err := m.saveHighlights(ctx, *demo, *match); err != nil {
			return err
		}

		return nil
	}(); err != nil {
		return false, status.Demo.Fail(ctx, demo, err)
	}

	return false, status.Demo.Succes(ctx, demo)
}
