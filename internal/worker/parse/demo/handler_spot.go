package demo

import (
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

func (d *Demo) handleSpotPlayerSpottersChanged(p demoinfocs.Parser, _ events.PlayerSpottersChanged) {
	state := p.GameState()
	r := d.getLastRound()
	tick := Tick(state.IngameTick())

	ts := state.TeamTerrorists().Members()
	cts := state.TeamCounterTerrorists().Members()

	for _, t := range ts {
		for _, ct := range cts {
			d.processSpot(tick, r, t, ct)
			d.processSpot(tick, r, ct, t)
		}
	}
}

func (d *Demo) processSpot(tick Tick, r *Round, spotter, spotted *common.Player) {
	if spotter == nil {
		return
	}

	spotterID := PlayerID(spotter.SteamID64)
	spottedID := PlayerID(spotted.SteamID64)

	if spotterID == 0 {
		return
	}

	isSpotted := spotter.HasSpotted(spotted)

	if isSpotted != d.getPlayerSpots(spotterID)[spottedID] {
		spotterStat, spotterOk := r.PlayerStats[spotterID]
		spottedStat, spottedOk := r.PlayerStats[spottedID]
		if !spotterOk || !spottedOk {
			return
		}

		if isSpotted {
			var spotsAfterDeath []*Spot
			if !spotted.IsAlive() && d.playerDeathTick[spottedID] > 0 {
				// Allow for a max of 1 spot after the spotted died
				// The player's death might not have been recorded so check if it's bigger than 0
				spotsAfterDeath = spotterStat.getSpottedAfter(d.playerDeathTick[spottedID], spottedID)
			}

			// If the spotted is dead immediatly close the spot
			end := Tick(0)
			if !spotted.IsAlive() {
				end = tick + 1
				isSpotted = false
			}

			// Only open if we haven't already recorded a post dead spot
			if len(spotsAfterDeath) == 0 {
				spotterStat.openSpot(tick, end, spottedID)
				spottedStat.openSpottedBy(tick, end, spotterID)
			}
		} else {
			spotterStat.closeSpot(tick, spottedID)
			spottedStat.closeSpottedBy(tick, spotterID)
		}

		d.playerSpots[spotterID][spottedID] = isSpotted
	}
}

func (d *Demo) handleSpotKill(p demoinfocs.Parser, e events.Kill) {
	state := p.GameState()
	r := d.getLastRound()

	tick := Tick(state.IngameTick())
	victim := getPlayer(e.Victim)

	// Record death time
	d.playerDeathTick[PlayerID(victim.SteamID64)] = tick

	// Close any remaining spots / spotted by
	if stat, ok := r.PlayerStats[PlayerID(victim.SteamID64)]; ok {
		for _, s := range stat.Spots {
			if s.End == 0 {
				stat.closeSpot(tick, s.Spotted)
				if spotStat, ok := r.PlayerStats[s.Spotted]; ok {
					spotStat.closeSpottedBy(tick, PlayerID(victim.SteamID64))
				}

				// Leave the cache to true
				// This eliminates a few post dead spots
			}
		}

		for _, s := range stat.SpottedBy {
			if s.End == 0 {
				stat.closeSpottedBy(tick, s.Spotter)
				if spotStat, ok := r.PlayerStats[s.Spotter]; ok {
					spotStat.closeSpot(tick, PlayerID(victim.SteamID64))
				}

				// Leave the cache to true
				// This eliminates a few post dead spots
			}
		}
	}
}

func (d *Demo) handleSpotRoundStart(p demoinfocs.Parser, _ events.RoundStart) {
	d.playerSpots = map[PlayerID]map[PlayerID]bool{}
	for _, player := range p.GameState().Participants().All() {
		player := getPlayer(player)
		if player.SteamID64 == 0 {
			continue
		}
		if player.Team != common.TeamCounterTerrorists && player.Team != common.TeamTerrorists {
			continue
		}

		d.playerSpots[PlayerID(player.SteamID64)] = map[PlayerID]bool{}
	}

	d.playerDeathTick = map[PlayerID]Tick{}
}
