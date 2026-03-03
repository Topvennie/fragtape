package steam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/config"
)

var S *steam

type steam struct {
	setting repository.SettingUser

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

	if demoResp.Error != "" {
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
		return demo, false, fmt.Errorf("steam service %w", err)
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

	return demo, true, nil
}
