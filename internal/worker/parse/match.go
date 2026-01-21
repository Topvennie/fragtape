package parse

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/worker/parse/demo"
	"github.com/topvennie/fragtape/pkg/storage"
	"github.com/topvennie/fragtape/pkg/utils"
)

func (p *Parser) getMatch(ctx context.Context, d *model.Demo) (*demo.Match, error) {
	if d.FileID == "" {
		return nil, errors.New("demo file deleted")
	}

	file, err := storage.S.Get(d.FileID)
	if err != nil {
		return nil, fmt.Errorf("get demo file %w", err)
	}

	var match *demo.Match

	if d.DataID == "" {
		match, err = p.demoParser.Parse(file)
		if err != nil {
			return nil, fmt.Errorf("parse demo file %w", err)
		}
		compressed, err := match.Compress()
		if err != nil {
			return nil, fmt.Errorf("compress parsed demo %w", err)
		}

		if err := p.repo.WithRollback(ctx, func(ctx context.Context) error {
			d.DataID = uuid.NewString()
			if err := p.demo.UpdateData(ctx, *d); err != nil {
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

func (p *Parser) savePlayers(ctx context.Context, match demo.Match) error {
	for _, player := range match.Players {
		user, err := p.user.GetByUID(ctx, int(player.SteamID))
		if err != nil {
			return err
		}

		if user == nil {
			user = &model.User{
				UID:         int(player.SteamID),
				DisplayName: player.Name,
				Crosshair:   player.CrosshairCode,
			}
			if err := p.user.Create(ctx, user); err != nil {
				return err
			}
		} else if user.DisplayName != player.Name || user.Crosshair != player.CrosshairCode {
			user.DisplayName = player.Name
			user.Crosshair = player.CrosshairCode

			if err := p.user.Update(ctx, *user); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Parser) getStatsDemo(ctx context.Context, d model.Demo, m demo.Match) (*model.StatsDemo, error) {
	statDB, err := p.statsDemo.GetByDemo(ctx, d.ID)
	if err != nil {
		return nil, err
	}

	stat := &model.StatsDemo{
		DemoID:   d.ID,
		Map:      m.Map,
		RoundsCT: m.RoundsCT,
		RoundsT:  m.RoundsT,
	}

	// Add id if it already exists
	if statDB != nil {
		stat.ID = statDB.ID
	}

	return stat, nil
}

func (p *Parser) getStats(ctx context.Context, d model.Demo, m demo.Match) ([]*model.Stat, error) {
	if len(m.Rounds) == 0 {
		return nil, nil
	}

	statsDB, err := p.stat.GetByDemo(ctx, d.ID)
	if err != nil {
		return nil, err
	}

	stats := make(map[demo.PlayerID]*model.Stat)

	for _, player := range m.Players {
		// Only add players that are in the ct or t team in the first round
		demoStat, ok := m.Rounds[0].PlayerStats[player.SteamID]
		if !ok {
			continue
		}
		if demoStat.Team != demo.TeamCounterTerrorists && demoStat.Team != demo.TeamTerrorists {
			continue
		}

		user, err := p.user.GetByUID(ctx, int(player.SteamID))
		if err != nil {
			return nil, err
		}

		if user == nil {
			return nil, errors.New("user not found")
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

		// Add id if it already exists
		if statDB, ok := utils.SliceFind(statsDB, func(s *model.Stat) bool { return s.UserID == user.ID }); ok {
			stat.ID = (*statDB).ID
		}

		stats[player.SteamID] = stat
	}

	for _, r := range m.Rounds {
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

	return utils.MapValues(stats), nil
}
