package demo

import (
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/msg"
)

func wrap[E any](p demoinfocs.Parser, fn func(demoinfocs.Parser, E)) func(E) {
	return func(e E) {
		fn(p, e)
	}
}

// Net messages

func (d *Demo) onNetMessage(_ demoinfocs.Parser, msg *msg.CSVCMsg_ServerInfo) {
	d.match.Map = msg.GetMapName()
}

// Events

//
// Announcement
//

func (d *Demo) onAnnouncementWinPanelMatch(p demoinfocs.Parser, e events.AnnouncementWinPanelMatch) {
	d.handleMatchAnnouncementWinPanel(p, e)
}

//
// Round
//

func (d *Demo) onRoundStart(p demoinfocs.Parser, e events.RoundStart) {
	d.handleMatchRoundStarted(p, e)

	if !d.started {
		return
	}

	d.handleRoundStart(p, e)

	d.handleBombRoundStart(p, e)
	d.handleHostageRoundStart(p, e)
	d.handleStatRoundStart(p, e)
	d.handleSpotRoundStart(p, e)
}

func (d *Demo) onRoundFreezeEnd(p demoinfocs.Parser, e events.RoundFreezetimeEnd) {
	if !d.started {
		return
	}

	d.handleRoundFreezeEnd(p, e)
}

func (d *Demo) onRoundEnd(p demoinfocs.Parser, e events.RoundEnd) {
	if !d.started {
		return
	}

	d.handleStatRoundEnd(p, e)
	d.handleRoundEnd(p, e)
}

func (d *Demo) onRoundEndOfficial(p demoinfocs.Parser, e events.RoundEndOfficial) {
	if !d.started {
		return
	}

	d.handleRoundEndOfficial(p, e)
}

//
// Player
//

func (d *Demo) onPlayerConnect(p demoinfocs.Parser, e events.PlayerConnect) {
	if getPlayer(e.Player).SteamID64 == 0 {
		return
	}

	d.handlePlayerConnect(p, e)
}

func (d *Demo) onPlayerDisconnect(p demoinfocs.Parser, e events.PlayerDisconnected) {
	if getPlayer(e.Player).SteamID64 == 0 {
		return
	}

	d.handlePlayerDisconnect(p, e)
}

func (d *Demo) onPlayerRankUpdate(p demoinfocs.Parser, e events.RankUpdate) {
	if getPlayer(e.Player).SteamID64 == 0 {
		return
	}

	d.handlePlayerRankUpdate(p, e)
}

func (d *Demo) onPlayerHurt(p demoinfocs.Parser, e events.PlayerHurt) {
	if !d.started {
		return
	}

	d.handleStatPlayerHurt(p, e)
	d.handleNadePlayerHurt(p, e)
}

func (d *Demo) onPlayerSpottersChanged(p demoinfocs.Parser, e events.PlayerSpottersChanged) {
	if !d.started {
		return
	}

	d.handleSpotPlayerSpottersChanged(p, e)
}

func (d *Demo) onPlayerButtonsStateUpdate(p demoinfocs.Parser, e events.PlayerButtonsStateUpdate) {
	if !d.started {
		return
	}

	d.handleStatPlayerButtonsStateUpdate(p, e)
}

//
// Bomb
//

func (d *Demo) onBombDefuseStart(p demoinfocs.Parser, e events.BombDefuseStart) {
	if !d.started {
		return
	}

	d.handleBombDefuseStart(p, e)
}

func (d *Demo) onBombDefuseAborted(p demoinfocs.Parser, e events.BombDefuseAborted) {
	if !d.started {
		return
	}

	d.handleBombDefuseAborted(p, e)
}

func (d *Demo) onBombDefused(p demoinfocs.Parser, e events.BombDefused) {
	if !d.started {
		return
	}

	d.handleBombDefused(p, e)
}

func (d *Demo) onBombPlantBegin(p demoinfocs.Parser, e events.BombPlantBegin) {
	if !d.started {
		return
	}

	d.handleBombPlantBegin(p, e)
}

func (d *Demo) onBombPlantAborted(p demoinfocs.Parser, e events.BombPlantAborted) {
	if !d.started {
		return
	}

	d.handleBombPlantAborted(p, e)
}

func (d *Demo) onBombPlanted(p demoinfocs.Parser, e events.BombPlanted) {
	if !d.started {
		return
	}

	d.handleBombPlanted(p, e)
}

func (d *Demo) onBombDrop(p demoinfocs.Parser, e events.BombDropped) {
	if !d.started {
		return
	}

	d.handleBombDrop(p, e)
}

func (d *Demo) onBombPickup(p demoinfocs.Parser, e events.BombPickup) {
	if !d.started {
		return
	}

	d.handleBombPickup(p, e)
}

//
// Hostage
//

func (d *Demo) onHostageStateChanged(p demoinfocs.Parser, e events.HostageStateChanged) {
	if !d.started {
		return
	}

	d.handleHostageStateChanged(p, e)
}

//
// Grenade
//

func (d *Demo) onGrenadeProjectileThrow(p demoinfocs.Parser, e events.GrenadeProjectileThrow) {
	if !d.started {
		return
	}

	d.handleNadeThrow(p, e)
}

func (d *Demo) onGrenadeProjectileBounce(p demoinfocs.Parser, e events.GrenadeProjectileBounce) {
	if !d.started {
		return
	}

	d.handleNadeBounce(p, e)
}

func (d *Demo) onGrenadeProjectileDestroy(p demoinfocs.Parser, e events.GrenadeProjectileDestroy) {
	if !d.started {
		return
	}

	d.handleNadeDestroy(p, e)
}

func (d *Demo) onPlayerFlashed(p demoinfocs.Parser, e events.PlayerFlashed) {
	if !d.started {
		return
	}

	d.handleNadePlayerFlashed(p, e)
}

func (d *Demo) onSmokeStart(p demoinfocs.Parser, e events.SmokeStart) {
	if !d.started {
		return
	}

	d.handleNadeSmokeStart(p, e)
}

func (d *Demo) onInfernoStart(p demoinfocs.Parser, e events.InfernoStart) {
	if !d.started {
		return
	}

	d.handleNadeInfernoStart(p, e)
}

func (d *Demo) onHeExplode(p demoinfocs.Parser, e events.HeExplode) {
	if !d.started {
		return
	}

	d.handleNadeHeExplode(p, e)
}

//
// Item
//

func (d *Demo) onItemPickup(p demoinfocs.Parser, e events.ItemPickup) {
	if !d.started {
		return
	}

	d.handleItemPickup(p, e)
}

func (d *Demo) onItemRefund(p demoinfocs.Parser, e events.ItemRefund) {
	if !d.started {
		return
	}

	d.handleItemRefund(p, e)
}

//
// Weapon
//

func (d *Demo) onWeaponFire(p demoinfocs.Parser, e events.WeaponFire) {
	if !d.started {
		return
	}

	d.handleStatWeaponFire(p, e)
}

func (d *Demo) onWeaponReload(p demoinfocs.Parser, e events.WeaponReload) {
	if !d.started {
		return
	}

	d.handleStatWeaponReload(p, e)
}

//
// Other
//

func (d *Demo) onFrameDone(p demoinfocs.Parser, e events.FrameDone) {
	if !d.started {
		return
	}

	d.handleStatFrameDone(p, e)
	d.handleItemFrameDone(p, e)
}

func (d *Demo) onKill(p demoinfocs.Parser, e events.Kill) {
	if !d.started {
		return
	}

	d.handleStatKill(p, e)
	d.handleSpotKill(p, e)
}

func (d *Demo) onTickRateInfoAvailable(p demoinfocs.Parser, e events.TickRateInfoAvailable) {
	d.handleMatchTickRateInfoAvailable(p, e)
}

func (d *Demo) onChatMessage(p demoinfocs.Parser, e events.ChatMessage) {
	if !d.started {
		return
	}

	d.handleStatChatMessage(p, e)
}

func (d *Demo) onOtherDeath(p demoinfocs.Parser, e events.OtherDeath) {
	if !d.started {
		return
	}

	d.handleStatOtherDeath(p, e)
}
