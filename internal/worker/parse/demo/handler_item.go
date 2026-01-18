package demo

import (
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

func (d *Demo) handleItemPickup(p demoinfocs.Parser, e events.ItemPickup) {
	state := p.GameState()
	r := d.getLastRound()

	uniqueID2 := e.Weapon.UniqueID2()

	prevOwner := d.droppedWeapons[uniqueID2]
	if prevOwner == 0 {
		prevOwner = d.weaponOwner[uniqueID2]
	}

	item := Item{
		Tick:      Tick(state.IngameTick()),
		Weapon:    EquipmentType(e.Weapon.Type),
		uniqueID2: uniqueID2,
	}

	if prevOwner != 0 {
		// Item was picked up
		var drop *ItemDrop
		for _, d := range r.PlayerStats[prevOwner].Drops {
			if d.uniqueID2 == uniqueID2 {
				drop = d
				break
			}
		}

		if drop != nil {
			drop.To = PlayerID(e.Player.SteamID64)
		}

		if stat, ok := r.PlayerStats[PlayerID(e.Player.SteamID64)]; ok {
			stat.Pickups = append(stat.Pickups, &ItemPickup{
				Item:     item,
				Position: Vector(e.Player.Position()),
				From:     prevOwner,
			})
		}
	} else {
		// Item was bought
		if stat, ok := r.PlayerStats[PlayerID(e.Player.SteamID64)]; ok {
			stat.Purchases = append(stat.Purchases, &ItemPurchase{
				Item: item,
			})
		}
	}

	// Update cache
	d.weaponOwner[uniqueID2] = PlayerID(e.Player.SteamID64)
}

func (d *Demo) handleItemRefund(p demoinfocs.Parser, e events.ItemRefund) {
	state := p.GameState()
	r := d.getLastRound()

	uniqueID2 := e.Weapon.UniqueID2()

	if stat, ok := r.PlayerStats[PlayerID(e.Player.SteamID64)]; ok {
		stat.Refunds = append(stat.Refunds, &ItemRefund{
			Item: Item{
				Tick:      Tick(state.IngameTick()),
				Weapon:    EquipmentType(e.Weapon.Type),
				uniqueID2: uniqueID2,
			},
		})
	}

	// Update cache
	d.weaponOwner[uniqueID2] = 0
}

func (d *Demo) handleItemFrameDone(p demoinfocs.Parser, _ events.FrameDone) {
	state := p.GameState()
	r := d.getLastRound()

	for _, v := range state.Weapons() {
		if v.Owner != nil {
			continue
		}

		uniqueID2 := v.UniqueID2()
		prevOwner := d.weaponOwner[uniqueID2]

		if prevOwner == 0 {
			// Already dropped
			continue
		}

		players := state.Participants().All()
		playerPos := Vector{}

		for _, p := range players {
			if PlayerID(p.SteamID64) == prevOwner {
				playerPos = Vector(p.Position())
				break
			}
		}

		if stat, ok := r.PlayerStats[prevOwner]; ok {
			stat.Drops = append(stat.Drops, &ItemDrop{
				Item: Item{
					Tick:      Tick(state.IngameTick()),
					Weapon:    EquipmentType(v.Type),
					uniqueID2: uniqueID2,
				},
				Position: playerPos,
			})
		}

		d.droppedWeapons[uniqueID2] = prevOwner
		d.weaponOwner[uniqueID2] = 0
	}
}
