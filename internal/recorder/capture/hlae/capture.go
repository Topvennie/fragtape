package hlae

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"text/template"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/storage"
	"github.com/topvennie/fragtape/pkg/utils"
)

func (h *Hlae) Capture(ctx context.Context, demo model.Demo, highlights []model.Highlight) error {
	data, err := storage.S.Get(demo.FileID)
	if err != nil {
		return fmt.Errorf("failed to get demo file from storage %w", err)
	}

	// Save demo file somewhere accessible to CS2
	if err := os.WriteFile(h.cs2Demo(demo), data, 0o644); err != nil {
		return fmt.Errorf("failed to write demo file %w", err)
	}

	// Create script
	if err := h.buildScript(ctx, h.hlaeScript(demo), demo, highlights); err != nil {
		return fmt.Errorf("failed to build script %w", err)
	}

	// Create cfg
	if err := h.buildCfg(h.cs2Cfg(demo), demo); err != nil {
		return fmt.Errorf("failed to build cfg %w", err)
	}

	// Cleanup
	// TODO: Uncomment
	// defer func() {
	// 	_ = os.Remove(h.cs2Demo(demo))
	// 	_ = os.Remove(h.hlaeScript(demo))
	// 	_ = os.Remove(h.cs2Cfg(demo))
	// }()

	if err := h.launch(ctx, demo); err != nil {
		return fmt.Errorf("launch cs2 %w", err)
	}

	return h.convert(ctx, demo, highlights)
}

//go:embed script.js.tmpl
var scriptTemplate string

type scriptData struct {
	RecordingDir string
	DemoID       int
	Players      []player
	PlayersJSON  string
}

type player struct {
	ID         int         `json:"id"`
	SteamID    int64       `json:"steamId"`
	Highlights []highlight `json:"highlights"`
}

type highlight struct {
	ID       int       `json:"id"`
	Segments []segment `json:"segments"`
}

type segment struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func (h *Hlae) buildScript(ctx context.Context, scriptPath string, demo model.Demo, highlights []model.Highlight) error {
	users, err := h.user.GetByIDs(ctx, utils.SliceUnique(utils.SliceMap(highlights, func(h model.Highlight) int { return h.UserID })))
	if err != nil {
		return err
	}

	// Only keep the highlights with segments
	// In theorie it shouldn't be possible to have a highlight without segments
	highlights = utils.SliceFilter(highlights, func(h model.Highlight) bool { return len(h.Segments) > 0 })

	// Sort all highlights and segments
	for _, h := range highlights {
		slices.SortFunc(h.Segments, func(a, b model.HighlightSegment) int { return a.StartTick - b.StartTick })
	}
	slices.SortFunc(highlights, func(a, b model.Highlight) int { return a.Segments[0].StartTick - b.Segments[0].StartTick })

	tmpl, err := template.New("script").Parse(scriptTemplate)
	if err != nil {
		return fmt.Errorf("parse script template %w", err)
	}

	playerMap := make(map[int]player)

	for _, h := range highlights {
		userIdx := slices.IndexFunc(users, func(u *model.User) bool { return u.ID == h.UserID })
		if userIdx == -1 {
			continue
		}

		pl, ok := playerMap[h.UserID]

		if !ok {
			pl = player{
				ID:         h.UserID,
				SteamID:    users[userIdx].UID,
				Highlights: []highlight{},
			}
		}

		pl.Highlights = append(pl.Highlights, highlight{
			ID: h.ID,
			Segments: utils.SliceMap(h.Segments, func(s model.HighlightSegment) segment {
				return segment{
					Start: s.StartTick,
					End:   s.EndTick,
				}
			}),
		})

		playerMap[h.UserID] = pl
	}

	recordingPath := h.cs2Video()
	recordingPath = filepath.ToSlash(recordingPath)

	data := scriptData{
		RecordingDir: recordingPath,
		DemoID:       demo.ID,
		Players:      utils.MapValues(playerMap),
	}

	playersJSON, err := json.Marshal(data.Players)
	if err != nil {
		return fmt.Errorf("marshal players data %+v | %w", data, err)
	}
	data.PlayersJSON = string(playersJSON)

	var out bytes.Buffer
	if err := tmpl.Execute(&out, data); err != nil {
		return fmt.Errorf("execute script template %w", err)
	}

	if err := os.WriteFile(scriptPath, out.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write script file %w", err)
	}

	return nil
}

//go:embed config.cfg.tmpl
var configTemplate string

type configData struct {
	DemoID int
}

func (h *Hlae) buildCfg(cfgPath string, demo model.Demo) error {
	tmpl, err := template.New("cfg").Parse(configTemplate)
	if err != nil {
		return fmt.Errorf("parse config template %w", err)
	}

	data := configData{
		DemoID: demo.ID,
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, data); err != nil {
		return fmt.Errorf("failed to write config file %w", err)
	}

	if err := os.WriteFile(cfgPath, out.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write cfg file %w", err)
	}

	return nil
}
