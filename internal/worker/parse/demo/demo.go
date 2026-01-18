// Package demo parses a demo
package demo

import (
	"bytes"
	"fmt"

	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	"github.com/oklog/ulid/v2"
)

type Demo struct {
	Match *Match

	samplesPerSecond int

	// Caches for parsing
	weaponOwner      map[ulid.ULID]PlayerID
	droppedWeapons   map[ulid.ULID]PlayerID
	playerSpots      map[PlayerID]map[PlayerID]bool
	playerDeathTick  map[PlayerID]Tick
	playerButtonMask map[PlayerID]common.ButtonBitMask
	playerScoped     map[PlayerID]bool
	playerMvps       map[PlayerID]int
}

func New(samplesPerSecond int) *Demo {
	return &Demo{
		Match: &Match{
			Players: []*Player{},
			Rounds:  []*Round{},
		},
		samplesPerSecond: samplesPerSecond,
		weaponOwner:      map[ulid.ULID]PlayerID{},
		droppedWeapons:   map[ulid.ULID]PlayerID{},
		playerSpots:      map[PlayerID]map[PlayerID]bool{},
		playerDeathTick:  map[PlayerID]Tick{},
		playerButtonMask: map[PlayerID]common.ButtonBitMask{},
		playerScoped:     map[PlayerID]bool{},
		playerMvps:       map[PlayerID]int{},
	}
}

func (d *Demo) Parse(file []byte) error {
	if err := demoinfocs.Parse(bytes.NewReader(file), func(p demoinfocs.Parser) error {
		// Net messages
		p.RegisterNetMessageHandler(wrap(p, d.onNetMessage))

		// Events

		// Announcement
		p.RegisterEventHandler(wrap(p, d.onAnnouncementWinPanelMatch))

		// Round
		p.RegisterEventHandler(wrap(p, d.onRoundStart))
		p.RegisterEventHandler(wrap(p, d.onRoundFreezeEnd))
		p.RegisterEventHandler(wrap(p, d.onRoundEnd))
		p.RegisterEventHandler(wrap(p, d.onRoundEndOfficial))

		// Player
		p.RegisterEventHandler(wrap(p, d.onPlayerConnect))
		p.RegisterEventHandler(wrap(p, d.onPlayerDisconnect))
		p.RegisterEventHandler(wrap(p, d.onPlayerRankUpdate))
		p.RegisterEventHandler(wrap(p, d.onPlayerHurt))
		p.RegisterEventHandler(wrap(p, d.onPlayerSpottersChanged))
		p.RegisterEventHandler(wrap(p, d.onPlayerButtonsStateUpdate))

		// Bomb
		p.RegisterEventHandler(wrap(p, d.onBombPlantBegin))
		p.RegisterEventHandler(wrap(p, d.onBombPlantAborted))
		p.RegisterEventHandler(wrap(p, d.onBombPlanted))
		p.RegisterEventHandler(wrap(p, d.onBombDefuseStart))
		p.RegisterEventHandler(wrap(p, d.onBombDefuseAborted))
		p.RegisterEventHandler(wrap(p, d.onBombDefused))
		p.RegisterEventHandler(wrap(p, d.onBombDrop))
		p.RegisterEventHandler(wrap(p, d.onBombPickup))
		p.RegisterEventHandler(wrap(p, d.onHeExplode))

		// Hostage
		p.RegisterEventHandler(wrap(p, d.onHostageStateChanged))

		// Nade
		p.RegisterEventHandler(wrap(p, d.onGrenadeProjectileThrow))
		p.RegisterEventHandler(wrap(p, d.onGrenadeProjectileBounce))
		p.RegisterEventHandler(wrap(p, d.onGrenadeProjectileDestroy))
		p.RegisterEventHandler(wrap(p, d.onPlayerFlashed))
		p.RegisterEventHandler(wrap(p, d.onSmokeStart))
		p.RegisterEventHandler(wrap(p, d.onInfernoStart))

		// Item
		p.RegisterEventHandler(wrap(p, d.onItemPickup))
		p.RegisterEventHandler(wrap(p, d.onItemRefund))

		// Weapon
		p.RegisterEventHandler(wrap(p, d.onWeaponFire))
		p.RegisterEventHandler(wrap(p, d.onWeaponReload))

		// Other
		p.RegisterEventHandler(wrap(p, d.onFrameDone))
		p.RegisterEventHandler(wrap(p, d.onKill))
		p.RegisterEventHandler(wrap(p, d.onTickRateInfoAvailable))
		p.RegisterEventHandler(wrap(p, d.onChatMessage))
		p.RegisterEventHandler(wrap(p, d.onOtherDeath))

		return nil
	}); err != nil {
		return fmt.Errorf("parse file %w", err)
	}

	return nil
}
