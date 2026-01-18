package demo

import (
	"time"

	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
	"github.com/topvennie/fragtape/pkg/utils"
)

func (d *Demo) handleStatRoundStart(p demoinfocs.Parser, _ events.RoundStart) {
	state := p.GameState()
	r := d.getLastRound()

	players := state.Participants().All()

	for _, player := range players {
		player = getPlayer(player) // Safe guard for nil pointers

		if player.SteamID64 == 0 {
			continue
		}
		if player.Team != common.TeamCounterTerrorists && player.Team != common.TeamTerrorists {
			continue
		}

		r.PlayerStats[PlayerID(player.SteamID64)] = &Stat{
			Team:           Team(player.Team),
			MoneyStart:     player.Money(),
			Positions:      []Position{},
			Kills:          []Kill{},
			Assists:        []Assist{},
			DamageReceived: []DamageReceived{},
			DamageDealt:    []DamageDealt{},
			Flashes:        []*Flash{},
			Smokes:         []*Smoke{},
			Hes:            []*He{},
			Decoys:         []*Decoy{},
			Incendiaries:   []*Incendiary{},
			Molotovs:       []*Molotov{},
			Purchases:      []*ItemPurchase{},
			Pickups:        []*ItemPickup{},
			Drops:          []*ItemDrop{},
			Refunds:        []*ItemRefund{},
			Spots:          []*Spot{},
			SpottedBy:      []*SpottedBy{},
			Reloads:        []Reload{},
			Shots:          []Shot{},
			Messages:       []Message{},
			Chickens:       []Chicken{},
		}
	}

	d.playerScoped = map[PlayerID]bool{}
}

func (d *Demo) handleStatRoundEnd(p demoinfocs.Parser, _ events.RoundEnd) {
	state := p.GameState()
	r := d.getLastRound()

	players := state.Participants().All()

	for _, player := range players {
		player = getPlayer(player) // Safe guard for nil pointers

		if stat, ok := r.PlayerStats[PlayerID(player.SteamID64)]; ok {
			stat.MoneySpent = player.MoneySpentThisRound()
		}
	}
}

func (d *Demo) handleStatFrameDone(p demoinfocs.Parser, _ events.FrameDone) {
	state := p.GameState()
	r := d.getLastRound()

	players := state.Participants().All()
	for _, player := range players {
		player = getPlayer(player) // Safe guard

		stat, ok := r.PlayerStats[PlayerID(player.SteamID64)]
		if !ok {
			continue
		}

		// Positions

		lastPos := stat.getLastPosition()
		ticksBetween := Tick(state.IngameTick()) - lastPos.Tick

		// Only add a position once every configured interval
		if ticksBetween >= d.Match.PositionTickInterval {
			stat.Positions = append(stat.Positions, Position{
				Tick:   Tick(state.IngameTick()),
				Vector: Vector(player.Position()),
			})
		}

		// Zooms

		if player.IsScoped() != d.playerScoped[PlayerID(player.SteamID64)] {
			if player.IsScoped() {
				stat.Zooms++
			}
			d.playerScoped[PlayerID(player.SteamID64)] = player.IsScoped()
		}
	}
}

func (d *Demo) handleStatKill(p demoinfocs.Parser, e events.Kill) {
	state := p.GameState()
	r := d.getLastRound()

	// Safe guard
	killer := getPlayer(e.Killer)
	assister := getPlayer(e.Assister)
	victim := getPlayer(e.Victim)

	if statKiller, ok := r.PlayerStats[PlayerID(killer.SteamID64)]; ok {
		statKiller.Kills = append(statKiller.Kills, Kill{
			Tick:           Tick(state.IngameTick()),
			Headshot:       e.IsHeadshot,
			NoScope:        e.NoScope,
			ThroughSmoke:   e.ThroughSmoke,
			Wallbang:       e.IsWallBang(),
			AttackerBlind:  e.AttackerBlind,
			VictimBlind:    e.AssistedFlash,
			Distance:       e.Distance,
			Weapon:         EquipmentType(e.Weapon.Type),
			Victim:         PlayerID(victim.SteamID64),
			KillerPosition: Vector(killer.Position()),
			VictimPosition: Vector(victim.Position()),
		})
	}

	if statAssister, ok := r.PlayerStats[PlayerID(assister.SteamID64)]; ok {
		damageDealt := utils.SliceFilter(statAssister.DamageDealt, func(d DamageDealt) bool { return d.Victim == PlayerID(victim.SteamID64) })
		statAssister.Assists = append(statAssister.Assists, Assist{
			Tick:        Tick(state.IngameTick()),
			Dead:        statAssister.Death != nil,
			Killer:      PlayerID(killer.SteamID64),
			Victim:      PlayerID(victim.SteamID64),
			DamageDealt: utils.SliceMap(damageDealt, func(d DamageDealt) Damage { return d.Damage }),
		})
	}

	if statVictim, ok := r.PlayerStats[PlayerID(victim.SteamID64)]; ok {
		statVictim.Death = &Death{
			Tick:           Tick(state.IngameTick()),
			Weapon:         EquipmentType(e.Weapon.Type),
			Killer:         PlayerID(killer.SteamID64),
			KillerPosition: Vector(killer.Position()),
			VictimPosition: Vector(victim.Position()),
		}
	}
}

func (d *Demo) handleStatPlayerHurt(p demoinfocs.Parser, e events.PlayerHurt) {
	state := p.GameState()
	r := d.getLastRound()

	// Safe guard
	attacker := getPlayer(e.Attacker)
	victim := getPlayer(e.Player)

	var flashDuration time.Duration = 0
	if e.Attacker != nil {
		flashDuration = attacker.FlashDurationTimeRemaining()
	}

	damage := Damage{
		Tick:                    Tick(state.IngameTick()),
		AttackerHealth:          attacker.Health(),
		AttackerArmor:           attacker.Armor(),
		AttackerFlashedDuration: flashDuration,
		VictimHealth:            e.Health,
		VictimArmor:             e.Armor,
		VictimHealthDamage:      e.HealthDamageTaken,
		VictimArmorDamage:       e.ArmorDamageTaken,
		VictimFlashedDuration:   victim.FlashDurationTimeRemaining(),
		Distance:                float32(Vector(e.Player.Position()).Distance(Vector(e.Attacker.Position()))),
		Weapon:                  EquipmentType(e.Weapon.Type),
		HitGroup:                HitGroup(e.HitGroup),
		AttackerPosition:        Vector(attacker.Position()),
		VictimPosition:          Vector(victim.Position()),
	}

	if statAttacker, ok := r.PlayerStats[PlayerID(attacker.SteamID64)]; ok {
		statAttacker.DamageDealt = append(statAttacker.DamageDealt, DamageDealt{
			Victim: PlayerID(victim.SteamID64),
			Damage: damage,
		})
	}

	if statVictim, ok := r.PlayerStats[PlayerID(victim.SteamID64)]; ok {
		statVictim.DamageReceived = append(statVictim.DamageReceived, DamageReceived{
			Attacker: PlayerID(attacker.SteamID64),
			Damage:   damage,
		})
	}
}

var statButtons = []common.ButtonBitMask{
	common.ButtonDuck,
	common.ButtonLookAtWeapon,
	common.ButtonJump,
	common.ButtonScore,
	common.ButtonAttack2,
}

func (d *Demo) handleStatPlayerButtonsStateUpdate(_ demoinfocs.Parser, e events.PlayerButtonsStateUpdate) {
	prev := d.playerButtonMask[PlayerID(e.Player.SteamID64)]
	curr := common.ButtonBitMask(e.ButtonsState)
	d.playerButtonMask[PlayerID(e.Player.SteamID64)] = curr

	pressed := curr &^ prev

	if pressed == 0 {
		return
	}

	r := d.getLastRound()
	stat, ok := r.PlayerStats[PlayerID(e.Player.SteamID64)]
	if !ok {
		return
	}

	for _, b := range statButtons {
		if pressed&b == 0 {
			continue
		}
		switch b {
		case common.ButtonDuck:
			stat.Crouches++
		case common.ButtonLookAtWeapon:
			stat.Inspects++
		case common.ButtonJump:
			stat.Jumps++
		case common.ButtonScore:
			stat.Scores++
		}
	}
}

func (d *Demo) handleStatWeaponReload(p demoinfocs.Parser, e events.WeaponReload) {
	if e.Player.ActiveWeapon().Class() == common.EqClassGrenade || e.Player.ActiveWeapon().Class() == common.EqClassEquipment {
		return
	}

	state := p.GameState()
	r := d.getLastRound()

	player := getPlayer(e.Player)

	if stat, ok := r.PlayerStats[PlayerID(player.SteamID64)]; ok {
		stat.Reloads = append(stat.Reloads, Reload{
			Tick:             Tick(state.IngameTick()),
			Position:         Vector(e.Player.Position()),
			Weapon:           EquipmentType(player.ActiveWeapon().Type),
			BulletsRemaining: player.ActiveWeapon().AmmoInMagazine(),
		})
	}
}

func (d *Demo) handleStatWeaponFire(p demoinfocs.Parser, e events.WeaponFire) {
	if e.Weapon.Class() == common.EqClassGrenade || e.Weapon.Class() == common.EqClassEquipment {
		return
	}

	state := p.GameState()
	r := d.getLastRound()

	player := getPlayer(e.Shooter)

	if stat, ok := r.PlayerStats[PlayerID(player.SteamID64)]; ok {
		stat.Shots = append(stat.Shots, Shot{
			Tick:     Tick(state.IngameTick()),
			Position: Vector(player.Position()),
			Weapon:   EquipmentType(e.Weapon.Type),
		})
	}
}

func (d *Demo) handleStatChatMessage(p demoinfocs.Parser, e events.ChatMessage) {
	state := p.GameState()
	r := d.getLastRound()

	player := getPlayer(e.Sender)

	if stat, ok := r.PlayerStats[PlayerID(player.SteamID64)]; ok {
		stat.Messages = append(stat.Messages, Message{
			Tick: Tick(state.IngameTick()),
			Text: e.Text,
		})
	}
}

func (d *Demo) handleStatOtherDeath(p demoinfocs.Parser, e events.OtherDeath) {
	if e.OtherType != "chicken" {
		return
	}

	state := p.GameState()
	r := d.getLastRound()

	attacker := getPlayer(e.Killer)

	if stat, ok := r.PlayerStats[PlayerID(attacker.SteamID64)]; ok {
		stat.Chickens = append(stat.Chickens, Chicken{
			Damage: Damage{
				Tick:                    Tick(state.IngameTick()),
				AttackerHealth:          attacker.Health(),
				AttackerArmor:           attacker.Armor(),
				AttackerFlashedDuration: attacker.FlashDurationTimeRemaining(),
				VictimHealth:            1,
				VictimArmor:             0,
				VictimHealthDamage:      1,
				VictimArmorDamage:       0,
				VictimFlashedDuration:   0,
				Distance:                float32(Vector(e.OtherPosition).Distance(Vector(attacker.Position()))),
				Weapon:                  EquipmentType(e.Weapon.Type),
				HitGroup:                0,
				AttackerPosition:        Vector(attacker.Position()),
				VictimPosition:          Vector(e.OtherPosition),
			},
		})
	}
}
