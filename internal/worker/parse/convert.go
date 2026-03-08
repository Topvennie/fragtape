package parse

import (
	"slices"
	"time"

	"github.com/topvennie/fragtape/internal/worker/parse/demo"
	"github.com/topvennie/fragtape/pkg/utils"
)

type round struct {
	number      int
	clutcher    demo.PlayerID // Is 0 if there is no clutcher
	clutchKills int
	players     map[demo.PlayerID]player
}

type player struct {
	id       demo.PlayerID
	kills    []kill
	messages []string
}

type kill struct {
	tick         demo.Tick     // Absolute tick when it happened
	tickRel      time.Duration // Duration since the start of the round when the kill happened
	oneTap       bool
	noScope      bool
	wallbang     bool
	blind        bool
	throughSmoke bool
	headshot     bool
	teamKill     bool
	weapon       demo.EquipmentType
	victimWeapon demo.EquipmentType
}

// Convert a match to rounds
func convertMatch(match demo.Match) []round {
	rounds := make([]round, 0, len(match.Rounds))

	for idx := range match.Rounds {
		rounds = append(rounds, convertRound(match, idx))
	}

	return rounds
}

// Convert a single round
func convertRound(match demo.Match, roundIdx int) round {
	demoRound := match.Rounds[roundIdx]

	round := round{
		number:  demoRound.Number,
		players: make(map[demo.PlayerID]player),
	}

	players := getRoundPlayers(*demoRound, match.Players)

	clutcher, clutchKills := getClutcher(*demoRound, players)
	round.clutcher = clutcher
	round.clutchKills = clutchKills

	for _, demoPlayer := range players {
		stat := demoRound.PlayerStats[demoPlayer.SteamID]

		// Get all kills
		kills := make([]kill, 0, len(stat.Kills))
		for killIdx, k := range stat.Kills {
			teamKill := false
			if victimStat, ok := demoRound.PlayerStats[k.Victim]; ok {
				teamKill = victimStat.Team == stat.Team
			}

			victimWeapon := demo.EqUnknown
			if k.Victim != 0 {
				victimWeapon = getWeaponOnDeath(match.Rounds, roundIdx, k.Victim)
			}

			kills = append(kills, kill{
				tick:         k.Tick,
				tickRel:      match.TickDurationRel(demoRound.FreezeEnd, k.Tick),
				oneTap:       killOneTap(match, *stat, killIdx),
				noScope:      k.NoScope,
				wallbang:     k.Wallbang,
				blind:        k.AttackerBlind,
				throughSmoke: k.ThroughSmoke,
				headshot:     k.Headshot,
				teamKill:     teamKill,
				weapon:       k.Weapon,
				victimWeapon: victimWeapon,
			})
		}

		// Get all messages
		messages := make([]string, 0, len(stat.Messages))
		for _, m := range stat.Messages {
			messages = append(messages, m.Text)
		}

		round.players[demoPlayer.SteamID] = player{
			id:       demoPlayer.SteamID,
			kills:    kills,
			messages: messages,
		}
	}

	return round
}

// isOneTap returns if the kill was a one tap
func killOneTap(match demo.Match, stat demo.Stat, killIdx int) bool {
	kill := stat.Kills[killIdx]

	// Get a range of 2 seconds between the kill
	startTick := kill.Tick - 2*match.TickRate
	endTick := kill.Tick + 2*match.TickRate

	// If the player got another kill in that range change the range
	if killIdx > 0 {
		startTick = max(startTick, stat.Kills[killIdx-1].Tick)
	}
	if killIdx < len(stat.Kills)-1 {
		endTick = min(endTick, stat.Kills[killIdx+1].Tick)
	}

	// Check if there was only one bullet shot in that range
	shots := utils.SliceFilter(stat.Shots, func(s demo.Shot) bool { return s.Tick >= startTick && s.Tick <= endTick })

	return len(shots) == 1
}

// getWeaponOnDeath returns the best weapon the player had before he died
// It assumes the player died that round
// It favours in the following order
// - Primary weapons
// - Secondary
// - Zeus
// - Knife
func getWeaponOnDeath(rounds []*demo.Round, roundIdx int, playerID demo.PlayerID) demo.EquipmentType {
	var primary demo.EquipmentType = 0
	var secondary demo.EquipmentType = 0
	hasZeus := false

	currentIdx := 0
	if roundIdx >= 11 {
		currentIdx = 11
	}

	for {
		// New round
		stat, ok := rounds[currentIdx].PlayerStats[playerID]
		if !ok {
			// No stats for this round
			// So the user disconnected
			// Reset their equipment
			primary = 0
			secondary = 0
			hasZeus = false

			continue
		}

		// Give the user a default pistol if he doesn't have one
		if secondary == 0 {
			switch stat.Team {
			case demo.TeamCounterTerrorists:
				secondary = demo.EqUSP // Skip p2000
			case demo.TeamTerrorists:
				secondary = demo.EqGlock
			}
		}

		died := stat.Death != nil

		// Go over each purchase, refund drop and pickup in chronological order
		events := make([]demo.Sortable, 0, len(stat.Purchases)+len(stat.Drops)+len(stat.Pickups)+len(stat.Refunds))
		for _, p := range stat.Purchases {
			events = append(events, p)
		}
		for _, d := range stat.Drops {
			events = append(events, d)
		}
		for _, p := range stat.Pickups {
			events = append(events, p)
		}
		for _, r := range stat.Refunds {
			events = append(events, r)
		}
		demo.Sort(events)

		for _, e := range events {
			switch v := e.(type) {
			case *demo.ItemPurchase:
			case *demo.ItemPickup:
				if weaponPrimary(v.Weapon) {
					primary = v.Weapon
				}
				if weaponSecondary(v.Weapon) {
					secondary = v.Weapon
				}
				if weaponZeus(v.Weapon) {
					hasZeus = true
				}
			case *demo.ItemRefund:
			case *demo.ItemDrop:
				// Only remove weapons if it happens before the player died (in that round)
				if died && stat.Death.Tick > v.Tick {
					if weaponPrimary(v.Weapon) {
						primary = 0
					}
					if weaponSecondary(v.Weapon) {
						secondary = 0
					}
					if weaponZeus(v.Weapon) {
						hasZeus = false
					}
				}
			}
		}

		if currentIdx == roundIdx {
			// We're done
			break
		}

		if died {
			// Reset
			primary = 0
			secondary = 0
			hasZeus = false
		}

		currentIdx++
	}

	switch {
	case primary != 0:
		return primary

	case secondary != 0:
		return secondary

	case hasZeus:
		return demo.EqZeus

	default:
		return demo.EqKnife
	}
}

// getRoundPlayers returns all players that were connected for that round
func getRoundPlayers(round demo.Round, players []*demo.Player) []demo.Player {
	connected := make([]demo.Player, 0, 10)
	for _, p := range players {
		stat := round.PlayerStats[p.SteamID]
		if stat == nil || stat.Team == demo.TeamSpectators {
			continue
		}

		// Get the last connect and disconnect before the round start
		// For users that don't disconnect this is equal to 0
		var lastConnect demo.Tick = 0
		var lastDisconnect demo.Tick = 0
		for _, c := range p.Connects {
			if c <= round.FreezeEnd && c > lastConnect {
				lastConnect = c
			}
		}
		for _, d := range p.Disconnects {
			if d <= round.FreezeEnd && d > lastDisconnect {
				lastDisconnect = d
			}
		}

		if lastConnect > lastDisconnect {
			connected = append(connected, *p)
		}
	}

	return connected
}

// getClutcher returns the id of the clutcher and the amount of kills he got
// If no one clutched than the id == 0
func getClutcher(round demo.Round, players []demo.Player) (demo.PlayerID, int) {
	// Get all stats
	statMap := make(map[demo.PlayerID]demo.Stat, 0)
	for _, p := range players {
		if stat, ok := round.PlayerStats[p.SteamID]; ok {
			statMap[stat.SteamID] = *stat
		}
	}

	// Get the amount of living players for each team
	cts := make([]demo.PlayerID, 0, 5)
	ts := make([]demo.PlayerID, 0, 5)
	for k, v := range statMap {
		if v.Team == demo.TeamCounterTerrorists {
			cts = append(cts, k)
		} else {
			ts = append(ts, k)
		}
	}

	// Get all kills
	kills := make([]demo.Kill, 0, 10)
	for _, stat := range statMap {
		kills = append(kills, stat.Kills...)
	}

	// Sort the kills in order that they happened
	slices.SortFunc(kills, func(a, b demo.Kill) int { return int(a.Tick - b.Tick) })

	// Go over each kill and check when we have a clutcher
	var ctClutcher demo.PlayerID = 0
	ctClutchKills := 0
	var tClutcher demo.PlayerID = 0
	tClutchKills := 0

	for _, k := range kills {
		stat, ok := statMap[k.Victim]
		if !ok {
			continue
		}

		if stat.Team == demo.TeamCounterTerrorists {
			utils.SliceRemove(cts, k.Victim)

			// Only one ct left
			if len(cts) == 1 {
				ctClutcher = cts[0]
			}

			// T is in a clutch
			if len(ts) == 1 {
				tClutchKills++
			}
		} else {
			utils.SliceRemove(ts, k.Victim)

			// Only one t left
			if len(ts) == 1 {
				tClutcher = ts[0]
			}

			// CT is in a clutch
			if len(cts) == 1 {
				ctClutchKills++
			}
		}
	}

	if round.Winner == demo.TeamCounterTerrorists {
		return ctClutcher, ctClutchKills
	}

	return tClutcher, tClutchKills
}
