package demo

import (
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

func (d *Demo) handlePlayerConnect(p demoinfocs.Parser, e events.PlayerConnect) {
	player := getPlayer(e.Player)

	idx := d.playerIndex(PlayerID(player.SteamID64))
	if idx == -1 {
		d.Match.Players = append(d.Match.Players, &Player{
			SteamID:       PlayerID(player.SteamID64),
			Name:          player.Name,
			Connects:      []Tick{},
			Disconnects:   []Tick{},
			CrosshairCode: player.CrosshairCode(),
		})
		idx = len(d.Match.Players) - 1
	}

	d.Match.Players[idx].Connects = append(d.Match.Players[idx].Connects, Tick(p.GameState().IngameTick()))
}

func (d *Demo) handlePlayerDisconnect(p demoinfocs.Parser, e events.PlayerDisconnected) {
	player := getPlayer(e.Player)

	if idx := d.playerIndex(PlayerID(player.SteamID64)); idx != -1 {
		d.Match.Players[idx].Disconnects = append(d.Match.Players[idx].Disconnects, Tick(p.GameState().IngameTick()))
	}
}

func (d *Demo) handlePlayerRankUpdate(_ demoinfocs.Parser, e events.RankUpdate) {
	player := getPlayer(e.Player)

	if idx := d.playerIndex(PlayerID(player.SteamID64)); idx != -1 {
		d.Match.Players[idx].Rank = RankUpdate{
			RankOld:  e.RankOld,
			RankNew:  e.RankNew,
			WinCount: e.WinCount,
		}
	}
}
