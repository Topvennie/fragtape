package demo

import (
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

func (d *Demo) handleBombRoundStart(p demoinfocs.Parser, _ events.RoundStart) {
	state := p.GameState()
	r := d.getLastRound()

	bomb := state.Bomb()
	if bomb != nil && bomb.Carrier != nil {
		r.Bomb = &Bomb{
			SpawnedWith: PlayerID(bomb.Carrier.SteamID64),
			Drops:       []*BombDrop{},
			Plants:      []*BombPlant{},
			Defuses:     []*BombDefuse{},
		}
	}
}

func (d *Demo) handleBombDefuseStart(p demoinfocs.Parser, e events.BombDefuseStart) {
	state := p.GameState()
	r := d.getLastRound()

	r.Bomb.Defuses = append(r.Bomb.Defuses, &BombDefuse{
		Defuser: PlayerID(e.Player.SteamID64),
		HasKit:  e.HasKit,
		Start:   Tick(state.IngameTick()),
	})
}

func (d *Demo) handleBombDefuseAborted(p demoinfocs.Parser, _ events.BombDefuseAborted) {
	state := p.GameState()
	r := d.getLastRound()

	defuse := r.getLastBombDefuse()
	defuse.End = Tick(state.IngameTick())
}

func (d *Demo) handleBombDefused(p demoinfocs.Parser, _ events.BombDefused) {
	state := p.GameState()
	r := d.getLastRound()

	defuse := r.getLastBombDefuse()
	defuse.End = Tick(state.IngameTick())
	defuse.Defused = true
}

func (d *Demo) handleBombPlantBegin(p demoinfocs.Parser, e events.BombPlantBegin) {
	state := p.GameState()
	r := d.getLastRound()

	player := getPlayer(e.Player)
	if player.SteamID64 == 0 {
		return // Can happen with pov demos
	}

	r.Bomb.Plants = append(r.Bomb.Plants, &BombPlant{
		Planter:  PlayerID(player.SteamID64),
		Start:    Tick(state.IngameTick()),
		Position: Vector(player.Position()),
		Site:     Bombsite(e.Site),
	})
}

func (d *Demo) handleBombPlantAborted(p demoinfocs.Parser, _ events.BombPlantAborted) {
	state := p.GameState()
	r := d.getLastRound()

	plant := r.getLastBombPlant()
	plant.End = Tick(state.IngameTick())
}

func (d *Demo) handleBombPlanted(p demoinfocs.Parser, _ events.BombPlanted) {
	state := p.GameState()
	r := d.getLastRound()

	plant := r.getLastBombPlant()
	plant.End = Tick(state.IngameTick())
	plant.Planted = true
}

func (d *Demo) handleBombDrop(p demoinfocs.Parser, e events.BombDropped) {
	state := p.GameState()
	r := d.getLastRound()

	r.Bomb.Drops = append(r.Bomb.Drops, &BombDrop{
		DropTick: Tick(state.IngameTick()),
		Dropper:  PlayerID(e.Player.SteamID64),
		Position: Vector(e.Player.Position()),
	})
}

func (d *Demo) handleBombPickup(p demoinfocs.Parser, e events.BombPickup) {
	state := p.GameState()
	r := d.getLastRound()
	if r.Number != state.TotalRoundsPlayed()+1 {
		// Sometimes when a player spawns with the bomb this event is fired
		return
	}

	drop := r.getLastBombDrop()
	if drop == nil {
		// Can happen that the drop event hasn't fired at the start of the game
		return
	}

	drop.Picker = PlayerID(e.Player.SteamID64)
	drop.PickupTick = Tick(state.IngameTick())
}
