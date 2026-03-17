package fetch

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/steam"
)

type steamFetcher struct {
	setting repository.SettingUser
	user    repository.User

	timeout time.Time
}

// Interface compliance
var _ fetcher = (*steamFetcher)(nil)

func newSteamFetcher(repo repository.Repository) *steamFetcher {
	return &steamFetcher{
		setting: *repo.NewSettingUser(),
		user:    *repo.NewUser(),
		timeout: time.Now(),
	}
}

func (s *steamFetcher) fetch(ctx context.Context, user model.User) (model.Demo, bool, error) {
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

	demoResp, err := steam.S.NextDemo(ctx, steam.NextDemoParams{
		SteamID:                  user.UID,
		SteamAuthenticationToken: user.Setting.SteamAuthenticationToken,
		SteamMatchToken:          user.Setting.SteamMatchToken,
	})
	if err != nil {
		return demo, false, err
	}

	if demoResp.Error != nil {
		// Something went wrong
		// First try to handle known codes
		switch demoResp.Code {
		case http.StatusForbidden:
			fallthrough
		case http.StatusPreconditionFailed:
			// User has an invalid match token / auth token
			user.Setting.SteamAuthenticationToken = ""
			user.Setting.SteamMatchToken = ""
			if err := s.setting.Update(ctx, user.Setting); err != nil {
				return demo, false, err
			}

			return demo, false, nil

		case http.StatusTooManyRequests:
			fallthrough
		case http.StatusServiceUnavailable:
			// Timeout received
			s.timeout = time.Now().Add(10 * time.Second)
			return demo, false, nil

		case http.StatusInternalServerError:
			fallthrough
		case http.StatusGatewayTimeout:
			// Something went wrong on valve's side
			return demo, false, nil

		default:
			return demo, false, fmt.Errorf("steam service %w", demoResp.Error)
		}
	}

	if demoResp.DemoURL == "" || demoResp.Code == 202 {
		// No new demo yet
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
			// Try to find an existing player
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
