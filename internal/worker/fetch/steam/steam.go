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
	"go.uber.org/zap"
)

const steamSource = model.DemoSourceSteam

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

func (s *steam) Fetch(ctx context.Context, user model.User) ([]byte, model.DemoSource, string, error) {
	zap.S().Debug("Fetching steam")
	if time.Now().Before(s.timeout) {
		// We're still waiting a bit
		zap.S().Debug("Too early")
		return nil, steamSource, "", nil
	}

	if user.Setting.SteamAuthenticationToken == "" || user.Setting.SteamMatchToken == "" {
		zap.S().Debug("No steam configured")
		return nil, steamSource, "", nil
	}

	demoResp, err := s.NextDemo(ctx, user)
	if err != nil {
		zap.S().Debug("Couldnt fetch")
		return nil, steamSource, "", err
	}

	if demoResp.Error != "" {
		zap.S().Debug("Error with demo")
		// Something went wrong
		// First try to handle known codes
		switch demoResp.Code {
		case http.StatusForbidden:
		case http.StatusPreconditionFailed:
			// User has an invalid match token / auth token
			user.Setting.SteamAuthenticationToken = ""
			user.Setting.SteamMatchToken = ""
			if err := s.setting.Update(ctx, user.Setting); err != nil {
				return nil, steamSource, "", err
			}

			return nil, steamSource, "", nil

		case http.StatusTooManyRequests:
		case http.StatusServiceUnavailable:
			// Timeout received
			s.timeout = time.Now().Add(10 * time.Second)
			return nil, steamSource, "", nil

		case http.StatusInternalServerError:
		case http.StatusGatewayTimeout:
			// Something went wrong on valve's side
			return nil, steamSource, "", nil
		}

		// We have no idea what went wrong
		return nil, steamSource, "", fmt.Errorf("steam service %w", err)
	}

	if demoResp.DemoURL == "" || demoResp.Code == 202 {
		// No new demo yet
		return nil, steamSource, "", nil
	}

	// New demo!
	demo, err := s.downloadDemo(ctx, demoResp.DemoURL)
	if err != nil {
		zap.S().Debug("Download failed")
		return nil, steamSource, "", err
	}

	zap.S().Debug("Setting next code")

	// Update next demo match code
	user.Setting.SteamMatchToken = demoResp.NextCode
	if err := s.setting.Update(ctx, user.Setting); err != nil {
		return nil, steamSource, "", nil
	}

	return demo, steamSource, demoResp.NextCode, nil
}
