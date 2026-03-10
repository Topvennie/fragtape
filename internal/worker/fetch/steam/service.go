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
	"github.com/topvennie/fragtape/pkg/utils"
)

type nextDemoReq struct {
	WebAPIKey  string `json:"webApiKey"`
	SteamID    int    `json:"steamId"`
	AuthToken  string `json:"authToken"`
	MatchToken string `json:"matchToken"`
}

type nextDemoResp struct {
	NextCode  string  `json:"nextCode"`
	DemoURL   string  `json:"demoUrl"`
	MatchTime int     `json:"matchTime"`
	Players   []int32 `json:"players"`
	Code      int     `json:"code"`
	Error     string  `json:"error"`
}

type NextDemo struct {
	NextCode  string
	DemoURL   string
	MatchTime time.Time
	Players   []int64
	Code      int
	Error     error
}

// NextDemo communicates with the steam service to get the next match
// If error != nil then something unexpected happened on the golang side
// If the resp.Error != nil then something unexpected happened on the typescript side
func (s *steam) NextDemo(ctx context.Context, user model.User) (NextDemo, error) {
	body := nextDemoReq{
		WebAPIKey:  s.webAPIKey,
		SteamID:    int(int32(user.UID)), // Steam expects 32 bit version
		AuthToken:  user.Setting.SteamAuthenticationToken,
		MatchToken: user.Setting.SteamMatchToken,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return NextDemo{}, fmt.Errorf("encode body %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		time.Sleep(10 * time.Second)
		cancel()
	}()
	req, err := http.NewRequestWithContext(ctx, "POST", s.steamServiceURL+"/steam/next-demo", &buf)
	if err != nil {
		return NextDemo{}, fmt.Errorf("create request %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return NextDemo{}, errors.New("steam service timed out")
		}
		return NextDemo{}, fmt.Errorf("do request %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return NextDemo{}, fmt.Errorf("unknown response code %s", resp.Status)
	}

	var demoRespService nextDemoResp
	if err := json.NewDecoder(resp.Body).Decode(&demoRespService); err != nil {
		return NextDemo{}, fmt.Errorf("decode response %w", err)
	}

	demoResp := NextDemo{
		NextCode: demoRespService.NextCode,
		DemoURL:  demoRespService.DemoURL,
		Players:  utils.SliceMap(demoRespService.Players, func(p int32) int64 { return steamID32To64(p) }),
		Code:     demoRespService.Code,
	}

	if demoRespService.MatchTime != 0 {
		demoResp.MatchTime = time.Unix(int64(demoRespService.MatchTime), 0)
	}
	if demoRespService.Error != "" {
		demoResp.Error = errors.New(demoRespService.Error)
	}

	return demoResp, nil
}

const steamID64Base int64 = 76561197960265728

func steamID32To64(id int32) int64 {
	return steamID64Base + int64(id)
}
