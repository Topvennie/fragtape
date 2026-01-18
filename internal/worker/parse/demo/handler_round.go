package demo

import (
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

func (d *Demo) handleRoundStart(p demoinfocs.Parser, _ events.RoundStart) {
	state := p.GameState()

	r := &Round{
		Number:      len(d.match.Rounds) + 1,
		Start:       Tick(state.IngameTick()),
		PlayerStats: map[PlayerID]*Stat{},
	}

	d.match.Rounds = append(d.match.Rounds, r)
}

func (d *Demo) handleRoundFreezeEnd(p demoinfocs.Parser, _ events.RoundFreezetimeEnd) {
	state := p.GameState()
	r := d.getLastRound()

	r.FreezeEnd = Tick(state.IngameTick())
}

func (d *Demo) handleRoundEnd(p demoinfocs.Parser, e events.RoundEnd) {
	state := p.GameState()
	r := d.getLastRound()

	// Update round information
	r.EndAnnouncement = Tick(state.IngameTick())
	r.EndReason = RoundEndReason(e.Reason)
	r.Winner = Team(e.Winner)

	// Figure out who got the mvp
	participants := state.Participants().All()
	for _, player := range participants {
		player := getPlayer(player)

		idx := d.playerIndex(PlayerID(player.SteamID64))
		if idx == -1 {
			// Player not found
			continue
		}

		if d.playerMvps[PlayerID(player.SteamID64)] != player.MVPs() {
			// Player just got an MVP
			d.playerMvps[PlayerID(player.SteamID64)]++
			r.Mvp = PlayerID(player.SteamID64)
		}
	}
}

func (d *Demo) handleRoundEndOfficial(p demoinfocs.Parser, _ events.RoundEndOfficial) {
	state := p.GameState()
	r := d.getLastRound()

	r.EndOfficial = Tick(state.IngameTick())
}
