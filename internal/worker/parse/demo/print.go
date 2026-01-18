package demo

import (
	"fmt"
	"maps"
	"slices"
	"time"
)

func (m Match) formatDuration(t Tick) time.Duration {
	return time.Second * time.Duration(t/Tick(m.TickRate))
}

func (m Match) formatRelative(zero Tick, t Tick) time.Duration {
	return m.formatDuration(t) - m.formatDuration(zero)
}

func fmtDur(d time.Duration) string {
	return d.Round(10 * time.Millisecond).String()
}

func add[T Sortable](acc []Sortable, items []T) []Sortable {
	for _, i := range items {
		acc = append(acc, i)
	}
	return acc
}

func (m Match) Print(rounds ...int) {
	playerMap := map[PlayerID]Player{}
	for _, p := range m.Players {
		playerMap[p.SteamID] = *p
	}

	fmt.Printf("\n----- Match -----\n\n")
	fmt.Printf("Map: %s | TickRate: %d\n", m.Map, m.TickRate)
	if len(m.Rounds) > 0 {
		fmt.Printf("Bomb: %t | Hostages: %t\n", m.Rounds[0].Bomb != nil, len(m.Rounds[0].Hostages) > 0)
	}

	fmt.Printf("\n----- Players -----\n\n")
	m.printPlayers(playerMap)

	fmt.Printf("\n----- Rounds -----\n\n")
	m.printRounds(playerMap, rounds)
}

func (m Match) printPlayers(playerMap map[PlayerID]Player) {
	for _, p := range m.Players {
		wonStr := "tied"
		if p.Won != nil {
			wonStr = fmt.Sprintf("%t", *p.Won)
		}

		fmt.Printf("ID %d (%s)\n", p.SteamID, playerMap[p.SteamID].Name)
		fmt.Printf("\tRank: %d -> %d (Wins: %d)\n", p.Rank.RankOld, p.Rank.RankNew, p.Rank.WinCount)
		fmt.Printf("\tWon: %s | Crosshair: %s\n", wonStr, p.CrosshairCode)
		fmt.Print("\tConnects: ")
		for _, c := range p.Connects {
			fmt.Printf("%s ", m.formatDuration(c))
		}
		fmt.Print("\n\tDisconnects: ")
		for _, d := range p.Disconnects {
			fmt.Printf("%s ", m.formatDuration(d))
		}
		fmt.Println()
	}
}

func (m Match) printRounds(playerMap map[PlayerID]Player, rounds []int) {
	for _, r := range m.Rounds {
		if len(rounds) > 0 && !slices.Contains(rounds, r.Number) {
			continue
		}
		fmt.Printf("Round %d\n", r.Number)
		fmt.Printf("\tStart %s | Freeze End %s\n", m.formatDuration(r.Start), m.formatRelative(r.Start, r.FreezeEnd))
		fmt.Printf("\tEnd Announcement %s | End Official %s | End Reason %s\n", m.formatRelative(r.Start, r.EndAnnouncement), m.formatRelative(r.Start, r.EndOfficial), r.EndReason)
		fmt.Printf("\tWinner %s\n", r.Winner)
		fmt.Printf("\tMVP %s\n", playerMap[r.Mvp].Name)

		m.printBomb(playerMap, r)
		m.printHostages(playerMap, r)
		m.printStats(playerMap, r)
		m.printChat(playerMap, r)
		fmt.Println()
	}
}

func (m Match) printBomb(playerMap map[PlayerID]Player, r *Round) {
	if r.Bomb == nil {
		return
	}

	events := make([]Sortable, 0,
		len(r.Bomb.Drops)+len(r.Bomb.Plants)+len(r.Bomb.Defuses),
	)

	events = add(events, r.Bomb.Drops)
	events = add(events, r.Bomb.Plants)
	events = add(events, r.Bomb.Defuses)

	Sort(events)

	fmt.Printf("\tBomb (Spawned with %s)\n", playerMap[r.Bomb.SpawnedWith].Name)

	for _, e := range events {
		t := e.GetTick()
		ts := m.formatRelative(r.Start, t)

		switch v := e.(type) {
		case *BombDrop:
			fmt.Printf("\t\t[%s] Dropped by %s at %s\n", ts, playerMap[v.Dropper].Name, v.Position)
			if v.Picker != 0 {
				fmt.Printf("\t\t[%s] Picked up by %s\n", ts, playerMap[v.Picker].Name)
			}

		case *BombPlant:
			fmt.Printf("\t\t[%s] Plant started by %s site=%s end=%s planted=%t pos=%s\n",
				ts, playerMap[v.Planter].Name, v.Site, m.formatRelative(t, v.End), v.Planted, v.Position)

		case *BombDefuse:
			fmt.Printf("\t\t[%s] Defuse started by %s kit=%t end=%s defused=%t\n",
				ts, playerMap[v.Defuser].Name, v.HasKit, m.formatRelative(t, v.End), v.Defused)
		}
	}
}

func (m Match) printHostages(playerMap map[PlayerID]Player, r *Round) {
	if len(r.Hostages) == 0 {
		return
	}

	fmt.Println("\tHostages")

	for _, h := range r.Hostages {
		fmt.Printf("\t\tEntity %d (Spawn %s)\n", h.EntityID, h.SpawnPosition)

		if len(h.Carries) > 0 {
			events := make([]Sortable, 0, len(h.Carries))

			events = add(events, h.Carries)

			Sort(events)

			for _, e := range events {
				switch v := e.(type) {
				case *HostageCarry:
					fmt.Printf("\t\t\t[%s] Picked up by %s at %s\n", m.formatRelative(r.Start, v.Start), playerMap[v.Carryer].Name, v.StartPosition)
					if v.End != 0 {
						fmt.Printf("\t\t\t[%s] Dropped at %s\n", m.formatRelative(r.Start, v.End), v.EndPosition)
					}
				}
			}

			if h.Rescue != nil {
				fmt.Printf("\t\t\t[%s] Rescued by %s at %s\n", m.formatRelative(r.Start, h.Rescue.Tick), playerMap[h.Rescue.Rescuer].Name, h.Rescue.Position)
			}
		}
	}
}

func (m Match) printStats(playerMap map[PlayerID]Player, r *Round) {
	fmt.Println("\tPlayer stats")

	sortedPlayers := slices.Sorted(maps.Keys(r.PlayerStats))

	for _, k := range sortedPlayers {
		v := r.PlayerStats[k]

		fmt.Printf("\t\tPlayer %s\n", playerMap[k].Name)

		fmt.Printf("\t\t\tMoney start %d | Spent %d\n", v.MoneyStart, v.MoneySpent)
		fmt.Printf("\t\t\tK/A/D: %d/%d/%t\n", len(v.Kills), len(v.Assists), v.Death != nil)
		fmt.Printf("\t\t\tMechanics: Jumps %d, Crouches %d, Zooms %d, Inspects %d\n", v.Jumps, v.Crouches, v.Zooms, v.Inspects)

		m.printPlayerItems(r, v, playerMap)
		m.printPlayerActions(r, v, playerMap)

	}
}

func (m Match) printPlayerItems(r *Round, v *Stat, playerMap map[PlayerID]Player) {
	fmt.Println("\t\t\tItems")

	events := make([]Sortable, 0)

	events = add(events, v.Purchases)
	events = add(events, v.Refunds)
	events = add(events, v.Pickups)
	events = add(events, v.Drops)

	Sort(events)

	for _, e := range events {
		t := e.GetTick()
		ts := m.formatRelative(r.Start, t)

		switch v := e.(type) {
		case *ItemPurchase:
			fmt.Printf("\t\t\t\t[%s] Purchased %s\n", ts, v.Weapon)

		case *ItemRefund:
			fmt.Printf("\t\t\t\t[%s] Refunded %s\n", ts, v.Weapon)

		case *ItemDrop:
			to := ""
			if v.To != 0 {
				to = " to " + playerMap[v.To].Name
			}

			fmt.Printf("\t\t\t\t[%s] Dropped %s%s\n", ts, v.Weapon, to)

		case *ItemPickup:
			from := ""
			if v.From != 0 {
				from = " from " + playerMap[v.From].Name
			}

			fmt.Printf("\t\t\t\t[%s] Picked up %s%s\n", ts, v.Weapon, from)
		}
	}
}

func (m Match) printPlayerActions(r *Round, v *Stat, playerMap map[PlayerID]Player) {
	fmt.Println("\t\t\tActions")

	events := make([]Sortable, 0)

	events = add(events, v.Kills)
	events = add(events, v.Assists)
	events = add(events, v.Flashes)
	events = add(events, v.Smokes)
	events = add(events, v.Hes)
	events = add(events, v.Decoys)
	events = add(events, v.Incendiaries)
	events = add(events, v.Molotovs)
	events = add(events, v.Spots)
	events = add(events, v.SpottedBy)
	events = add(events, v.Reloads)
	events = add(events, v.Shots)
	events = add(events, v.Chickens)
	if v.Death != nil {
		events = append(events, *v.Death)
	}

	Sort(events)

	var shotCount int
	var shotWeapon EquipmentType

	flushShots := func() {
		if shotCount > 0 {
			plural := "time"
			if shotCount > 1 {
				plural = "times"
			}
			fmt.Printf("%d %s\n", shotCount, plural)
			shotCount = 0
		}
	}

	for _, e := range events {
		t := e.GetTick()
		ts := m.formatRelative(r.Start, t)

		if s, ok := e.(Shot); ok {
			if shotCount > 0 && s.Weapon != shotWeapon {
				flushShots() // Weapon changed
			}
			shotWeapon = s.Weapon
			shotCount++

			if shotCount == 1 {
				fmt.Printf("\t\t\t\t[%s] Shoots %s ", m.formatRelative(r.Start, s.Tick), s.Weapon)
			}
			continue
		}

		flushShots()

		switch v := e.(type) {
		case Kill:
			flags := ""
			if v.Headshot {
				flags += "[HS] "
			}
			if v.NoScope {
				flags += "[NoScope] "
			}
			if v.ThroughSmoke {
				flags += "[Smoke] "
			}
			if v.Wallbang {
				flags += "[Wallbang] "
			}
			if v.AttackerBlind {
				flags += "[Blind] "
			}

			fmt.Printf("\t\t\t\t[%s] Kills %s (%.2fm) with %s %s\n", ts, playerMap[v.Victim].Name, v.Distance, v.Weapon, flags)

		case Assist:
			damage := 0
			for _, d := range v.DamageDealt {
				damage += d.VictimHealthDamage
			}

			fmt.Printf("\t\t\t\t[%s] Assist on %s with %d damage (killed by %s)\n", ts, playerMap[v.Victim].Name, damage, playerMap[v.Killer].Name)

		case Death:
			fmt.Printf("\t\t\t\t[%s] Died to %s (%s)\n", ts, playerMap[v.Killer].Name, v.Weapon)

		case *Flash:
			fmt.Printf("\t\t\t\t[%s] Flash bounces=%d pops=%s trajectory=%d", ts, v.Bounces, m.formatRelative(v.Start, v.End), len(v.Trajectory))
			if len(v.Victims) > 0 {
				fmt.Printf(" hits=")
				for _, v := range v.Victims {
					fmt.Printf("%s (%s) ", playerMap[v.Victim].Name, fmtDur(v.Duration))
				}
			}
			fmt.Println()

		case *Smoke:
			fmt.Printf("\t\t\t\t[%s] Smoke bounces=%d blooms=%s pops=%s trajectory=%d\n", ts, v.Bounces, m.formatRelative(v.Start, v.Bloom), m.formatRelative(v.Start, v.End), len(v.Trajectory))

		case *He:
			fmt.Printf("\t\t\t\t[%s] HE bounces=%d pops=%s trajectory=%d", m.formatRelative(r.Start, v.Start), v.Bounces, m.formatRelative(v.Start, v.End), len(v.Trajectory))
			if len(v.Victims) > 0 {
				fmt.Printf(" hits=")
				for _, v := range v.Victims {
					fmt.Printf("%s (%d) ", playerMap[v.Victim].Name, v.VictimHealthDamage)
				}
			}
			fmt.Println()

		case *Decoy:
			fmt.Printf("\t\t\t\t[%s] Decoy bounces=%d pops=%s trajectory=%d", ts, v.Bounces, m.formatRelative(v.Start, v.End), len(v.Trajectory))
			if len(v.Victims) > 0 {
				fmt.Printf(" hits=")
				for _, v := range v.Victims {
					fmt.Printf("%s (%d) ", playerMap[v.Victim].Name, v.VictimHealthDamage)
				}
			}
			fmt.Println()

		case *Incendiary:
			fmt.Printf("\t\t\t\t[%s] Incendiary bounces=%d pops=%s trajectory=%d", ts, v.Bounces, m.formatRelative(v.Start, v.End), len(v.Trajectory))
			if len(v.Victims) > 0 {
				damageMap := make(map[PlayerID]int)
				for _, v := range v.Victims {
					damageMap[v.Victim] += v.VictimHealthDamage
				}
				fmt.Printf(" hits=")
				for k, v := range damageMap {
					fmt.Printf("%s (%d) ", playerMap[k].Name, v)
				}
			}
			fmt.Println()

		case *Molotov:
			fmt.Printf("\t\t\t\t[%s] Molotov bounces=%d pops=%s trajectory=%d", ts, v.Bounces, m.formatRelative(v.Start, v.End), len(v.Trajectory))
			if len(v.Victims) > 0 {
				damageMap := make(map[PlayerID]int)
				for _, v := range v.Victims {
					damageMap[v.Victim] += v.VictimHealthDamage
				}
				fmt.Printf(" hits=")
				for k, v := range damageMap {
					fmt.Printf("%s (%d) ", playerMap[k].Name, v)
				}
			}
			fmt.Println()

		case *Spot:
			fmt.Printf("\t\t\t\t[%s] Spot %s (%s)\n", ts, playerMap[v.Spotted].Name, m.formatRelative(v.Start, v.End))

		case *SpottedBy:
			fmt.Printf("\t\t\t\t[%s] Spotted by %s (%s)\n", ts, playerMap[v.Spotter].Name, m.formatRelative(v.Start, v.End))

		case Reload:
			fmt.Printf("\t\t\t\t[%s] reload %s with %d bullets remaining\n", ts, v.Weapon, v.BulletsRemaining)

		case Chicken:
			fmt.Printf("\t\t\t\t[%s] Kills a chicken with %s\n", ts, v.Weapon)
		}
	}

	flushShots()
}

type message struct {
	tick   Tick
	text   string
	author Player
}

func (m message) GetTick() Tick {
	return m.tick
}

func (m Match) printChat(playerMap map[PlayerID]Player, r *Round) {
	events := make([]Sortable, 0)

	for k, v := range r.PlayerStats {
		for _, m := range v.Messages {
			events = append(events, message{
				tick:   m.Tick,
				text:   m.Text,
				author: playerMap[k],
			})
		}
	}

	Sort(events)

	if len(events) == 0 {
		return
	}

	fmt.Println("\tChat")

	for _, e := range events {
		v := e.(message)

		fmt.Printf("\t\t[%s] %s: %s\n", m.formatRelative(r.Start, v.tick), v.author.Name, v.text)
	}
}

// String representations of some types

func (t Team) String() string {
	switch t {
	case TeamUnassigned:
		return "Unassigned"
	case TeamSpectators:
		return "Spectators"
	case TeamTerrorists:
		return "T"
	case TeamCounterTerrorists:
		return "CT"
	default:
		return "Unknown"
	}
}

func (r RoundEndReason) String() string {
	switch r {
	case RoundEndReasonStillInProgress:
		return "Still in progress"
	case RoundEndReasonTargetBombed:
		return "Target bombed"
	case RoundEndReasonVIPEscaped:
		return "VIP escaped"
	case RoundEndReasonVIPKilled:
		return "VIP killed"
	case RoundEndReasonTerroristsEscaped:
		return "T's escaped"
	case RoundEndReasonCTStoppedEscape:
		return "CT's stopped escaped"
	case RoundEndReasonTerroristsStopped:
		return "T's stopped"
	case RoundEndReasonBombDefused:
		return "Bomb defused"
	case RoundEndReasonCTWin:
		return "CT win"
	case RoundEndReasonTerroristsWin:
		return "T win"
	case RoundEndReasonDraw:
		return "Draw"
	case RoundEndReasonHostagesRescued:
		return "Hostages rescued"
	case RoundEndReasonTargetSaved:
		return "Target saved"
	case RoundEndReasonHostagesNotRescued:
		return "Hostages not rescued"
	case RoundEndReasonTerroristsNotEscaped:
		return "T's not escaped"
	case RoundEndReasonVIPNotEscaped:
		return "VIP not escaped"
	case RoundEndReasonGameStart:
		return "Game start"
	case RoundEndReasonTerroristsSurrender:
		return "T's surrender"
	case RoundEndReasonCTSurrender:
		return "CT's surrender"
	case RoundEndReasonTerroristsPlanted:
		return "T's planted"
	case RoundEndReasonCTsReachedHostage:
		return "CT's reached hostage"
	default:
		return fmt.Sprintf("Reason %d", r)
	}
}

func (b Bombsite) String() string {
	switch b {
	case BombsiteA:
		return "A"
	case BombsiteB:
		return "B"
	default:
		return "Unknown"
	}
}

func (e EquipmentType) String() string {
	switch e {
	case EqUnknown:
		return "Unknown"
	case EqP2000:
		return "P2000"
	case EqGlock:
		return "Glock-18"
	case EqP250:
		return "P250"
	case EqDeagle:
		return "Desert Eagle"
	case EqFiveSeven:
		return "Five-SeveN"
	case EqDualBerettas:
		return "Dual Berettas"
	case EqTec9:
		return "Tec-9"
	case EqCZ:
		return "CZ75-Auto"
	case EqUSP:
		return "USP-S"
	case EqRevolver:
		return "R8 Revolver"
	case EqMP7:
		return "MP7"
	case EqMP9:
		return "MP9"
	case EqBizon:
		return "PP-Bizon"
	case EqMac10:
		return "MAC-10"
	case EqUMP:
		return "UMP-45"
	case EqP90:
		return "P90"
	case EqMP5:
		return "MP5-SD"
	case EqSawedOff:
		return "Sawed-Off"
	case EqNova:
		return "Nova"
	case EqMag7:
		return "MAG-7"
	case EqXM1014:
		return "XM1014"
	case EqM249:
		return "M249"
	case EqNegev:
		return "Negev"
	case EqGalil:
		return "Galil AR"
	case EqFamas:
		return "FAMAS"
	case EqAK47:
		return "AK-47"
	case EqM4A4:
		return "M4A4"
	case EqM4A1:
		return "M4A1-S"
	case EqSSG08:
		return "SSG 08"
	case EqSG553:
		return "SG 553"
	case EqAUG:
		return "AUG"
	case EqAWP:
		return "AWP"
	case EqScar20:
		return "SCAR-20"
	case EqG3SG1:
		return "G3SG1"
	case EqZeus:
		return "Zeus x27"
	case EqKevlar:
		return "Kevlar Vest"
	case EqHelmet:
		return "Kevlar + Helmet"
	case EqBomb:
		return "C4 Explosive"
	case EqKnife:
		return "Knife"
	case EqDefuseKit:
		return "Defuse Kit"
	case EqWorld:
		return "World"
	case EqZoneRepulsor:
		return "Zone Repulsor"
	case EqShield:
		return "Ballistic Shield"
	case EqHeavyAssaultSuit:
		return "Heavy Assault Suit"
	case EqNightVision:
		return "Night Vision Goggles"
	case EqHealthShot:
		return "Medi-Shot"
	case EqTacticalAwarenessGrenade:
		return "Tactical Awareness Grenade"
	case EqFists:
		return "Fists"
	case EqBreachCharge:
		return "Breach Charge"
	case EqTablet:
		return "Tablet"
	case EqAxe:
		return "Axe"
	case EqHammer:
		return "Hammer"
	case EqWrench:
		return "Wrench"
	case EqSnowball:
		return "Snowball"
	case EqBumpMine:
		return "Bump Mine"
	case EqDecoy:
		return "Decoy Grenade"
	case EqMolotov:
		return "Molotov"
	case EqIncendiary:
		return "Incendiary Grenade"
	case EqFlash:
		return "Flashbang"
	case EqSmoke:
		return "Smoke Grenade"
	case EqHE:
		return "HE Grenade"
	default:
		return "unknown"
	}
}

func (v Vector) String() string {
	return fmt.Sprintf("X %.2f Y %.2f Z %.2f", v.X, v.Y, v.Z)
}
