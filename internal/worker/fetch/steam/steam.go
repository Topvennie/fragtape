package steam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/config"
)

var S *steam

type steam struct {
	setting repository.SettingUser
	user    repository.User

	timeout time.Time

	steamServiceURL string
	webAPIKey       string
}

func Init(repo repository.Repository) error {
	if S != nil {
		return nil
	}

	steamService := config.GetDefaultString("worker.fetcher.steam.service_url", "")
	webAPIKey := config.GetDefaultString("server.auth.steam.api_key", "")

	if steamService == "" {
		return errors.New("no steam service url set")
	}
	if webAPIKey == "" {
		return errors.New("no web api key set")
	}

	S = &steam{
		setting:         *repo.NewSettingUser(),
		user:            *repo.NewUser(),
		timeout:         time.Now(),
		steamServiceURL: steamService,
		webAPIKey:       webAPIKey,
	}

	return nil
}

func (s *steam) Fetch(ctx context.Context, user model.User) (model.Demo, bool, error) {
	demo := model.Demo{
		Source: model.DemoSourceSteam,
	}

	if time.Now().Before(s.timeout) {
		// We're still waiting a bit
		return demo, false, nil
	}

	if user.Setting.SteamAuthenticationToken == "" || user.Setting.SteamMatchToken == "" {
		// User doesn't have steam configured
		return demo, false, nil
	}

	// Get next demo from the steam service
	demoResp, err := s.NextDemo(ctx, user)
	if err != nil {
		return demo, false, err
	}

	if demoResp.Error != nil {
		// Something went wrong
		// First try to handle known codes
		switch demoResp.Code {
		case http.StatusForbidden:
		case http.StatusPreconditionFailed:
			// User has an invalid match token / auth token
			user.Setting.SteamAuthenticationToken = ""
			user.Setting.SteamMatchToken = ""
			if err := s.setting.Update(ctx, user.Setting); err != nil {
				return demo, false, err
			}

			return demo, false, nil

		case http.StatusTooManyRequests:
		case http.StatusServiceUnavailable:
			// Timeout received
			s.timeout = time.Now().Add(10 * time.Second)
			return demo, false, nil

		case http.StatusInternalServerError:
		case http.StatusGatewayTimeout:
			// Something went wrong on valve's side
			return demo, false, nil
		}

		// We have no idea what went wrong
		return demo, false, fmt.Errorf("steam service %w", demoResp.Error)
	}

	if demoResp.DemoURL == "" || demoResp.Code == 202 {
		// No new demo yet
		return demo, false, nil
	}

	// Sanity check
	if demoResp.NextCode == user.Demo.SourceID {
		// User already has this demo
		return demo, false, nil
	}

	// New demo!
	// Update next demo match code
	user.Setting.SteamMatchToken = demoResp.NextCode
	if err := s.setting.Update(ctx, user.Setting); err != nil {
		return demo, false, err
	}

	demo.SourceID = demoResp.NextCode
	demo.SourceURL = demoResp.DemoURL
	demo.PlayedAt = demoResp.MatchTime

	if len(demoResp.Players) > 0 {
		demo.Stats = make([]model.Stat, 0, len(demoResp.Players))

		users, err := s.user.GetByUIDs(ctx, demoResp.Players)
		if err != nil {
			return demo, false, err
		}

		for _, player := range demoResp.Players {
			if idx := slices.IndexFunc(users, func(u *model.User) bool { return u.UID == player }); idx != -1 {
				demo.Stats = append(demo.Stats, model.Stat{
					UserID: users[idx].ID,
				})
				continue
			}

			user := &model.User{
				UID: player,
			}
			if err := s.user.Create(ctx, user); err != nil {
				return demo, false, err
			}

			demo.Stats = append(demo.Stats, model.Stat{
				UserID: user.ID,
			})
		}
	}

	return demo, true, nil
}
