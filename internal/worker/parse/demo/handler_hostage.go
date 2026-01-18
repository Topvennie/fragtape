package demo

import (
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
	"github.com/topvennie/fragtape/pkg/utils"
)

func (d *Demo) handleHostageRoundStart(p demoinfocs.Parser, _ events.RoundStart) {
	state := p.GameState()
	r := d.getLastRound()

	hostages := state.Hostages()
	if len(hostages) > 0 {
		r.Hostages = []*Hostage{}
		for _, h := range hostages {
			r.Hostages = append(r.Hostages, &Hostage{
				EntityID:      h.Entity.ID(),
				SpawnPosition: Vector(h.Entity.Position()),
				Carries:       []*HostageCarry{},
			})
		}
	}
}

func (d *Demo) handleHostageStateChanged(p demoinfocs.Parser, e events.HostageStateChanged) {
	state := p.GameState()
	r := d.getLastRound()

	h := r.getHostage(e.Hostage.Entity.ID())
	if h == nil {
		return
	}

	// Find who is interacting with the hostage
	// Is probably the CT the closest to the hostage
	ct := closest(Vector(e.Hostage.Position()), state.TeamCounterTerrorists().Members(), func(p *common.Player) Vector { return Vector(p.Position()) })

	switch e.NewState {
	case common.HostageStateBeingCarried:
		// Player picked hostage up
		h.Carries = append(h.Carries, &HostageCarry{
			Start:         Tick(state.IngameTick()),
			StartPosition: Vector(e.Hostage.Position()),
			Carryer:       PlayerID(ct.SteamID64),
		})
	case common.HostageStateGettingDropped:
		// Player dropped hostage
		carry := utils.SliceLast(h.Carries)
		carry.End = Tick(state.IngameTick())
		carry.EndPosition = Vector(e.Hostage.Position())
	case common.HostageStateRescued:
		// Player rescued hostage
		h.Rescue = &HostageRescue{
			Tick:     Tick(state.IngameTick()),
			Position: Vector(e.Hostage.Position()),
			Rescuer:  PlayerID(ct.SteamID64),
		}
	}
}
