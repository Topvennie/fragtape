package faceit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/topvennie/fragtape/internal/database/model"
	"go.uber.org/zap"
)

func (f *faceit) GetUserID(ctx context.Context, user model.User) (string, error) {
	// First try for cs2
	zap.S().Debug("trying cs2")
	resp, err := f.getUser(ctx, user, "cs2")
	if err != nil {
		return "", err
	}
	if resp.PlayerID != "" {
		return resp.PlayerID, nil
	}

	// Not found for cs2, let's try csgo
	zap.S().Debug("trying csgo")
	resp, err = f.getUser(ctx, user, "csgo")
	if err != nil {
		return "", err
	}

	return resp.PlayerID, nil
}

type userResponse struct {
	PlayerID string `json:"player_id"`
}

func (f *faceit) getUser(ctx context.Context, user model.User, gameID string) (userResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL+"players", http.NoBody)
	if err != nil {
		return userResponse{}, fmt.Errorf("new request %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+f.webAPIKey)

	q := req.URL.Query()
	q.Add("game", gameID)
	q.Add("game_player_id", strconv.FormatInt(user.UID, 10))
	req.URL.RawQuery = q.Encode()

	zap.S().Debugf("%+v", req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return userResponse{}, fmt.Errorf("do request %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	zap.S().Debugf("%+v", *resp)

	if resp.StatusCode != http.StatusOK {
		// Some kind of error
		switch resp.StatusCode {
		case http.StatusBadRequest:
			// Bad request
			return userResponse{}, nil
		case http.StatusForbidden:
			// Bad api key
			return userResponse{}, nil
		case http.StatusNotFound:
			// User not found
			return userResponse{}, nil
		case http.StatusTooManyRequests:
			// Unlucky
			return userResponse{}, nil
		case http.StatusServiceUnavailable:
			// Unlucky
			return userResponse{}, nil
		default:
			return userResponse{}, fmt.Errorf("unexpected status code %s", resp.Status)
		}
	}

	var userResp userResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return userResponse{}, fmt.Errorf("decode response to body %w", err)
	}

	return userResp, nil
}
