// Package fetch gets automatically new demo files
package fetch

import (
	"context"
	"sync"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/status"
	"github.com/topvennie/fragtape/pkg/config"
	"github.com/topvennie/fragtape/pkg/utils"
	"go.uber.org/zap"
)

type fetcher interface {
	// Fetch gets a new demo file url
	// It returns
	//  model.Demo -> the new demo
	//  bool -> indicating if there is a new demo
	//  error -> error
	fetch(context.Context, model.User) (model.Demo, bool, error)
}

type Manager struct {
	interval time.Duration
	cooldown time.Duration

	fetchers []fetcher

	busyUsers map[int]bool
	mu        sync.Mutex

	repo repository.Repository
	demo repository.Demo
	stat repository.Stat
	user repository.User
}

func New(repo repository.Repository) *Manager {
	return &Manager{
		interval:  config.GetDefaultDurationS("worker.fetcher.interval_s", 300),
		cooldown:  config.GetDefaultDurationS("worker.fetcher.cooldown_s", 600),
		fetchers:  []fetcher{newSteamFetcher(repo), newFaceitFetcher()},
		busyUsers: map[int]bool{},
		repo:      repo,
		demo:      *repo.NewDemo(),
		stat:      *repo.NewStat(),
		user:      *repo.NewUser(),
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

	m.mu.Lock()
	users = utils.SliceFilter(users, func(u *model.User) bool {
		if busy := m.busyUsers[u.ID]; busy {
			return false
		}

		if u.Demo.ID == 0 {
			return true
		}

		return u.Demo.CreatedAt.Before(lastDemo)
	})
	m.mu.Unlock()

	for _, user := range users {
		// Is it a new steam user?
		// Meaning filled in steam credentials but no demos yet
		if user.Demo.ID == 0 && user.Setting.SteamMatchToken != "" && user.Setting.SteamAuthenticationToken != "" {
			// It is a new steam user
			// Mark the user as being processed
			m.mu.Lock()
			m.busyUsers[user.ID] = true
			m.mu.Unlock()

			// Get all the demo's in a seperate thread
			go func() {
				if err := m.loopSteamNew(ctx, user); err != nil {
					zap.S().Error(err)
				}

				// Remove the user from the new users proces queue
				m.mu.Lock()
				delete(m.busyUsers, user.ID)
				m.mu.Unlock()
			}()

			continue
		}

		// Try each fetcher for a new match
		for _, fetcher := range m.fetchers {
			demo, ok, err := fetcher.fetch(ctx, *user)
			if err != nil {
				zap.S().Error(err)
				continue
			}

			// The parser is responsible for making sure it is newer
			if !ok {
				continue
			}

			if err := m.handleNewDemo(ctx, &demo); err != nil {
				zap.S().Error(err)
			}
		}
	}

	return nil
}

// loopSteamNew can be used for a new steam account
// It will keep going over the steam match tokens until it reaches the newest one
func (m *Manager) loopSteamNew(ctx context.Context, user *model.User) error {
	fetcher := newSteamFetcher(m.repo)

	demo, ok, err := fetcher.fetch(ctx, *user)

	for ok {
		if err != nil {
			return err
		}

		if user.Setting.SteamImportOld {
			if err := m.handleNewDemo(ctx, &demo); err != nil {
				return err
			}
		}

		// Wait 5 seconds between each fetch
		time.Sleep(5 * time.Second)
		// Refresh the user to update the steam settings
		user, err = m.user.GetByIDWithSettingLastDemo(ctx, user.ID)
		if err != nil {
			return err
		}

		demo, ok, err = fetcher.fetch(ctx, *user)
	}

	return err
}

func (m *Manager) handleNewDemo(ctx context.Context, demo *model.Demo) error {
	// Does this demo already exist?
	oldDemo, err := m.demo.GetBySourceSourceID(ctx, demo.Source, demo.SourceID)
	if err != nil {
		return err
	}
	if oldDemo != nil {
		// Demo already exists
		// Possible if a different user was also in the match and got handled first
		demo.ID = oldDemo.ID
	} else {
		// New demo
		demo.Status = model.DemoStatusQueuedDownload

		if err := m.demo.Create(ctx, demo); err != nil {
			return err
		}
	}

	// Add all players (if any)
	// Use the no conflict create as we only insert the player id
	// We want to do it atomically and don't overwrite any existing data
	for _, stat := range demo.Stats {
		stat.DemoID = demo.ID
		if err := m.stat.CreateNoConflict(ctx, &stat); err != nil {
			return err
		}
	}

	return nil
}
