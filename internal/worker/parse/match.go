package parse

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/worker/parse/demo"
	"github.com/topvennie/fragtape/pkg/storage"
)

func (m *Manager) getMatch(ctx context.Context, d *model.Demo) (*demo.Match, error) {
	if d.FileID == "" {
		return nil, errors.New("demo file deleted")
	}

	file, err := storage.S.Get(d.FileID)
	if err != nil {
		return nil, fmt.Errorf("get demo file %w", err)
	}

	var match *demo.Match

	// Take into account that we might already have the data
	// if the pipeline failed later on
	if d.DataID == "" {
		match, err = m.demoParser.Parse(file)
		if err != nil {
			return nil, fmt.Errorf("parse demo file %w", err)
		}
		compressed, err := match.Compress()
		if err != nil {
			return nil, fmt.Errorf("compress parsed demo %w", err)
		}

		if err := m.repo.WithRollback(ctx, func(ctx context.Context) error {
			d.DataID = uuid.NewString()
			if err := m.demo.UpdateData(ctx, *d); err != nil {
				return err
			}
			if err := storage.S.Set(d.DataID, compressed, 0); err != nil {
				return fmt.Errorf("save compressed match %w", err)
			}

			return nil
		}); err != nil {
			return nil, err
		}
	} else {
		compressed, err := storage.S.Get(d.DataID)
		if err != nil {
			return nil, fmt.Errorf("get demo data from storage %w", err)
		}
		match, err = demo.Decompress(compressed)
		if err != nil {
			return nil, fmt.Errorf("decompress match data %w", err)
		}
	}

	return match, nil
}

func (m *Manager) savePlayers(ctx context.Context, match demo.Match) error {
	for _, player := range match.Players {
		user, err := m.user.GetByUID(ctx, int64(player.SteamID))
		if err != nil {
			return err
		}

		if user == nil {
			user = &model.User{
				UID:         int64(player.SteamID),
				DisplayName: player.Name,
				Crosshair:   player.CrosshairCode,
			}
			if err := m.user.Create(ctx, user); err != nil {
				return err
			}
		} else if user.DisplayName != player.Name || user.Crosshair != player.CrosshairCode {
			user.DisplayName = player.Name
			user.Crosshair = player.CrosshairCode

			if err := m.user.Update(ctx, *user); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Manager) saveStatsDemo(ctx context.Context, d model.Demo, match demo.Match) error {
	statDB, err := m.statsDemo.GetByDemo(ctx, d.ID)
	if err != nil {
		return err
	}

	stat := &model.StatsDemo{
		DemoID:   d.ID,
		Map:      match.Map,
		RoundsCT: match.RoundsCT,
		RoundsT:  match.RoundsT,
	}

	// Add id if it already exists
	if statDB != nil {
		stat.ID = statDB.ID

		return m.statsDemo.Update(ctx, *stat)
	}

	return m.statsDemo.Create(ctx, stat)
}

func (m *Manager) saveStats(ctx context.Context, d model.Demo, match demo.Match) error {
	if len(match.Rounds) == 0 {
		return nil
	}

	stats := make(map[demo.PlayerID]*model.Stat)

	for _, player := range match.Players {
		// Only add players that are in the ct or t team in the first round
		demoStat, ok := match.Rounds[0].PlayerStats[player.SteamID]
		if !ok {
			continue
		}
		if demoStat.Team != demo.TeamCounterTerrorists && demoStat.Team != demo.TeamTerrorists {
			continue
		}

		user, err := m.user.GetByUID(ctx, int64(player.SteamID))
		if err != nil {
			return err
		}
		if user == nil {
			return errors.New("user not found")
		}

		result := model.ResultTie
		if player.Won != nil {
			if *player.Won {
				result = model.ResultWin
			} else {
				result = model.ResultLoss
			}
		}

		stat := &model.Stat{
			DemoID: d.ID,
			UserID: user.ID,
			Result: result,
		}

		stats[player.SteamID] = stat
	}

	for _, r := range match.Rounds {
		for player, s := range stats {
			if stat, ok := r.PlayerStats[player]; ok {
				if r.Number == 1 {
					s.StartTeam = model.TeamCT
					if stat.Team == demo.TeamTerrorists {
						s.StartTeam = model.TeamT
					}
				}

				s.Kills += len(stat.Kills)
				s.Assists += len(stat.Assists)
				if stat.Death != nil {
					s.Deaths++
				}
			}
		}
	}

	for _, stat := range stats {
		if err := m.stat.CreateUpdateAtomic(ctx, stat); err != nil {
			return err
		}
	}

	return nil
}
