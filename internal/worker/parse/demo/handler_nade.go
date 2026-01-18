package demo

import (
	"slices"

	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

func (d *Demo) handleNadeThrow(p demoinfocs.Parser, e events.GrenadeProjectileThrow) {
	state := p.GameState()
	r := d.getLastRound()

	thrower := getPlayer(e.Projectile.Thrower)

	// We get the entire trajectory when the nades gets destroyed
	grenade := Grenade{
		Start:     Tick(state.IngameTick()),
		uniqueID2: e.Projectile.WeaponInstance.UniqueID2(),
	}

	stat := r.PlayerStats[PlayerID(thrower.SteamID64)]

	switch EquipmentType(e.Projectile.WeaponInstance.Type) {
	case EqDecoy:
		stat.Decoys = append(stat.Decoys, &Decoy{
			Grenade: grenade,
			Victims: []DamageDealt{},
		})
	case EqMolotov:
		stat.Molotovs = append(stat.Molotovs, &Molotov{
			Grenade:  grenade,
			Victims:  []DamageDealt{},
			uniqueID: e.Projectile.UniqueID(),
		})
	case EqIncendiary:
		stat.Incendiaries = append(stat.Incendiaries, &Incendiary{
			Grenade:  grenade,
			Victims:  []DamageDealt{},
			uniqueID: e.Projectile.UniqueID(),
		})
	case EqFlash:
		stat.Flashes = append(stat.Flashes, &Flash{
			Grenade: grenade,
			Victims: []FlashVictim{},
		})
	case EqSmoke:
		stat.Smokes = append(stat.Smokes, &Smoke{
			Grenade: grenade,
			Bloom:   Tick(state.IngameTick()),
		})
	case EqHE:
		stat.Hes = append(stat.Hes, &He{
			Grenade: grenade,
			Victims: []DamageDealt{},
		})
	}
}

func (d *Demo) handleNadeBounce(_ demoinfocs.Parser, e events.GrenadeProjectileBounce) {
	r := d.getLastRound()

	thrower := getPlayer(e.Projectile.Thrower)

	if stat, ok := r.PlayerStats[PlayerID(thrower.SteamID64)]; ok {
		grenade := stat.getNade(e.Projectile.WeaponInstance.UniqueID2())
		if grenade == nil {
			return
		}

		switch v := grenade.(type) {
		case *Flash:
			v.Bounces++
		case *Smoke:
			v.Bounces++
		case *He:
			v.Bounces++
		case *Decoy:
			v.Bounces++
		case *Incendiary:
			v.Bounces++
		case *Molotov:
			v.Bounces++
		}
	}
}

func (d *Demo) handleNadeDestroy(p demoinfocs.Parser, e events.GrenadeProjectileDestroy) {
	state := p.GameState()
	r := d.getLastRound()

	thrower := getPlayer(e.Projectile.Thrower)

	if stat, ok := r.PlayerStats[PlayerID(thrower.SteamID64)]; ok {
		grenade := stat.getNade(e.Projectile.WeaponInstance.UniqueID2())
		if grenade == nil {
			return
		}

		switch v := grenade.(type) {
		case *Flash:
			v.End = Tick(state.IngameTick())
			v.Trajectory = toTrajectoryEntrySlice(e.Projectile.Trajectory)
		case *Smoke:
			v.End = Tick(state.IngameTick())
			v.Trajectory = toTrajectoryEntrySlice(e.Projectile.Trajectory)
		case *He:
			v.End = Tick(state.IngameTick())
			v.Trajectory = toTrajectoryEntrySlice(e.Projectile.Trajectory)
		case *Decoy:
			v.End = Tick(state.IngameTick())
			v.Trajectory = toTrajectoryEntrySlice(e.Projectile.Trajectory)
		case *Incendiary:
			v.End = Tick(state.IngameTick())
			v.Trajectory = toTrajectoryEntrySlice(e.Projectile.Trajectory)
		case *Molotov:
			v.End = Tick(state.IngameTick())
			v.Trajectory = toTrajectoryEntrySlice(e.Projectile.Trajectory)
		}
	}
}

func (d *Demo) handleNadePlayerFlashed(p demoinfocs.Parser, e events.PlayerFlashed) {
	state := p.GameState()
	r := d.getLastRound()

	player := getPlayer(e.Player)
	thrower := getPlayer(e.Attacker)

	if stat, ok := r.PlayerStats[PlayerID(thrower.SteamID64)]; ok {
		for _, f := range stat.Flashes {
			if f.uniqueID2 == e.Projectile.WeaponInstance.UniqueID2() {
				f.Victims = append(f.Victims, FlashVictim{
					Tick:     Tick(state.IngameTick()),
					Position: Vector(player.Position()),
					Victim:   PlayerID(player.SteamID64),
					Duration: player.FlashDurationTimeRemaining(),
				})
				return
			}
		}
	}
}

func (d *Demo) handleNadeSmokeStart(p demoinfocs.Parser, e events.SmokeStart) {
	state := p.GameState()
	r := d.getLastRound()

	thrower := getPlayer(e.Thrower)

	if stat, ok := r.PlayerStats[PlayerID(thrower.SteamID64)]; ok {
		for _, s := range stat.Smokes {
			if s.uniqueID2 == e.Grenade.UniqueID2() {
				s.Bloom = Tick(state.IngameTick())
				return
			}
		}
	}
}

func (d *Demo) handleNadePlayerHurt(p demoinfocs.Parser, e events.PlayerHurt) {
	weaponType := EquipmentType(e.Weapon.Type)

	nades := []EquipmentType{EqHE, EqDecoy, EqIncendiary, EqMolotov}
	if idx := slices.Index(nades, weaponType); idx == -1 {
		return
	}

	state := p.GameState()
	r := d.getLastRound()

	attacker := getPlayer(e.Attacker)
	victim := getPlayer(e.Player)

	damage := Damage{
		Tick:                    Tick(state.IngameTick()),
		AttackerHealth:          attacker.Health(),
		AttackerArmor:           attacker.Armor(),
		AttackerFlashedDuration: attacker.FlashDurationTimeRemaining(),
		VictimHealth:            e.Health,
		VictimArmor:             e.Armor,
		VictimHealthDamage:      e.HealthDamageTaken,
		VictimArmorDamage:       e.ArmorDamageTaken,
		VictimFlashedDuration:   victim.FlashDurationTimeRemaining(),
		Distance:                float32(Vector(victim.Position()).Distance(Vector(attacker.Position()))),
		Weapon:                  EquipmentType(e.Weapon.Type),
		HitGroup:                HitGroup(e.HitGroup),
		AttackerPosition:        Vector(attacker.Position()),
		VictimPosition:          Vector(victim.Position()),
	}

	uniqueID2 := e.Weapon.UniqueID2()

	if stat, ok := r.PlayerStats[PlayerID(attacker.SteamID64)]; ok {
		switch weaponType {
		case EqHE:
			for _, h := range stat.Hes {
				if h.uniqueID2 == uniqueID2 {
					h.Victims = append(h.Victims, DamageDealt{
						Victim: PlayerID(victim.SteamID64),
						Damage: damage,
					})
				}
			}

		case EqDecoy:
			for _, d := range stat.Decoys {
				if d.uniqueID2 == uniqueID2 {
					d.Victims = append(d.Victims, DamageDealt{
						Victim: PlayerID(victim.SteamID64),
						Damage: damage,
					})
				}
			}

		case EqIncendiary:
			for _, i := range stat.Incendiaries {
				if i.uniqueID2 == uniqueID2 {
					i.Victims = append(i.Victims, DamageDealt{
						Victim: PlayerID(victim.SteamID64),
						Damage: damage,
					})
				}
			}

		case EqMolotov:
			for _, m := range stat.Molotovs {
				if m.uniqueID2 == uniqueID2 {
					m.Victims = append(m.Victims, DamageDealt{
						Victim: PlayerID(victim.SteamID64),
						Damage: damage,
					})
				}
			}
		}
	}
}

func (d *Demo) handleNadeInfernoStart(p demoinfocs.Parser, e events.InfernoStart) {
	state := p.GameState()
	r := d.getLastRound()

	thrower := getPlayer(e.Inferno.Thrower())

	if stat, ok := r.PlayerStats[PlayerID(thrower.SteamID64)]; ok {
		for _, i := range stat.Incendiaries {
			if i.uniqueID == e.Inferno.UniqueID() {
				i.FireStart = Tick(state.IngameTick())
				return
			}
		}
		for _, m := range stat.Molotovs {
			if m.uniqueID == e.Inferno.UniqueID() {
				m.FireStart = Tick(state.IngameTick())
				return
			}
		}
	}
}

func (d *Demo) handleNadeHeExplode(p demoinfocs.Parser, e events.HeExplode) {
	state := p.GameState()
	r := d.getLastRound()

	thrower := getPlayer(e.Thrower)

	if stat, ok := r.PlayerStats[PlayerID(thrower.SteamID64)]; ok {
		for _, h := range stat.Hes {
			if h.uniqueID2 == e.Grenade.UniqueID2() {
				h.Exploded = Tick(state.IngameTick())
			}
		}
	}
}
