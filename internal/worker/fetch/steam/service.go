package steam

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
)

type nextDemoReq struct {
	WebAPIKey  string `json:"webApiKey"`
	SteamID    int    `json:"steamId"`
	AuthToken  string `json:"authToken"`
	MatchToken string `json:"matchToken"`
}

type nextDemoResp struct {
	NextCode string `json:"nextCode"`
	DemoURL  string `json:"demoUrl"`
	Code     int    `json:"code"`
	Error    string `json:"error"`
}

func (s *steam) NextDemo(ctx context.Context, user model.User) (nextDemoResp, error) {
	body := nextDemoReq{
		WebAPIKey:  s.webAPIKey,
		SteamID:    user.UID,
		AuthToken:  user.Setting.SteamAuthenticationToken,
		MatchToken: user.Setting.SteamMatchToken,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nextDemoResp{}, fmt.Errorf("encode body %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		time.Sleep(10 * time.Second)
		cancel()
	}()
	req, err := http.NewRequestWithContext(ctx, "POST", s.steamServiceURL+"/steam/next-demo", &buf)
	if err != nil {
		return nextDemoResp{}, fmt.Errorf("create request %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nextDemoResp{}, errors.New("steam service timed out")
		}
		return nextDemoResp{}, fmt.Errorf("do request %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return nextDemoResp{}, fmt.Errorf("unknown response code %s", resp.Status)
	}

	var demoResp nextDemoResp
	if err := json.NewDecoder(resp.Body).Decode(&demoResp); err != nil {
		return nextDemoResp{}, fmt.Errorf("decode response %w", err)
	}

	return demoResp, nil
}
