package demo

import (
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

func (d *Demo) handleMatchRoundStarted(p demoinfocs.Parser, _ events.RoundStart) {
	if d.started {
		// Already started
		return
	}

	if !p.GameState().IsMatchStarted() {
		return
	}

	d.started = true
}

func (d *Demo) handleMatchTickRateInfoAvailable(_ demoinfocs.Parser, e events.TickRateInfoAvailable) {
	d.match.TickRate = int(e.TickRate)
	if d.samplesPerSecond > 0 {
		d.match.PositionTickInterval = Tick(d.match.TickRate / d.samplesPerSecond)
	}
}

func (d *Demo) handleMatchAnnouncementWinPanel(p demoinfocs.Parser, _ events.AnnouncementWinPanelMatch) {
	state := p.GameState()

	participants := state.Participants().All()

	ctScore := state.TeamCounterTerrorists().Score()
	tScore := state.TeamTerrorists().Score()

	winner := common.TeamUnassigned
	if ctScore > tScore {
		winner = common.TeamCounterTerrorists
	} else if ctScore < tScore {
		winner = common.TeamTerrorists
	}

	if winner == common.TeamUnassigned {
		// Tie
		return
	}

	trueTmp := true
	falseTmp := false

	for _, player := range participants {
		player := getPlayer(player)

		if idx := d.playerIndex(PlayerID(player.SteamID64)); idx != -1 {
			if player.Team == winner {
				d.match.Players[idx].Won = &trueTmp
			} else {
				d.match.Players[idx].Won = &falseTmp
			}
		}
	}
}
