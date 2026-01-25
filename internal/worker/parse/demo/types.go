package demo

import (
	"time"

	"github.com/oklog/ulid/v2"
)

//
// Common
//

type PlayerID uint64

type Tick int

type Vector struct {
	X float32
	Y float32
	Z float32
}

type Damage struct {
	Tick Tick

	AttackerHealth          int           // Can be 0 if it is from world damage
	AttackerArmor           int           // Can be 0 if it is from world damage
	AttackerFlashedDuration time.Duration // Remaining time that the player is flashed

	VictimHealth          int
	VictimArmor           int
	VictimHealthDamage    int           // Excludes over-damage (e.g. if the player has 5 health and is hit for 15 damage then it will be 5)
	VictimArmorDamage     int           // Excludes over-damage (e.g. if the player has 5 health and is hit for 15 damage then it will be 5)
	VictimFlashedDuration time.Duration // Remaining time that the player is flashed

	Distance float32

	Weapon   EquipmentType
	HitGroup HitGroup

	AttackerPosition Vector
	VictimPosition   Vector
}

type DamageReceived struct {
	Damage

	Attacker PlayerID // Can be 0 if the player is taking world damage (e.g. fall damage) or if the demo is partially corrup
}

type DamageDealt struct {
	Damage

	Victim PlayerID // Can be 0 with pov demos or if the demo is partially corrupt
}

//
// Match
//

type Match struct {
	Map      string
	TickRate Tick // Amount of ticks in one second

	Players []*Player
	Rounds  []*Round

	RoundsCT int
	RoundsT  int

	PositionTickInterval Tick    // Amount of ticks between position updates. If set to 0 then no positions are saved
	PositionMinDistance  float32 // Minimum distance in units before a new position is recorded
}

//
// Player
//

type Player struct {
	SteamID       PlayerID
	Name          string
	CrosshairCode string

	Connects    []Tick
	Disconnects []Tick

	Rank RankUpdate
	Won  *bool // nil if it is a draw
}

type RankUpdate struct {
	RankOld  int
	RankNew  int
	WinCount int
}

//
// Round
//

type Round struct {
	Number int

	Start     Tick
	FreezeEnd Tick

	EndAnnouncement Tick // When the round winner is decided, players are still able to move around after
	EndOfficial     Tick // After end, before this players are still able to walk around
	EndReason       RoundEndReason

	Winner Team
	Mvp    PlayerID

	Bomb     *Bomb      // Is nil if it is a hostage map
	Hostages []*Hostage // Has len 0 if it is a bomb map

	PlayerStats map[PlayerID]*Stat
}

//
// Player Round Stats
//

type Stat struct {
	Team Team

	MoneyStart int
	MoneySpent int

	Positions []Position // Poll interval is configurable

	Kills   []Kill
	Death   *Death
	Assists []Assist

	DamageReceived []DamageReceived
	DamageDealt    []DamageDealt

	Flashes      []*Flash
	Smokes       []*Smoke
	Hes          []*He
	Decoys       []*Decoy
	Incendiaries []*Incendiary
	Molotovs     []*Molotov

	Purchases []*ItemPurchase
	Pickups   []*ItemPickup
	Drops     []*ItemDrop
	Refunds   []*ItemRefund

	Spots     []*Spot
	SpottedBy []*SpottedBy

	Crouches int // Amount of times the user crouches
	Scores   int // Amount of times the user looks at the scoreboard
	Inspects int // Amount of times the user presses the inspect button
	Jumps    int // Amount of times the user jumps
	Zooms    int // Amount of times the user scopes

	Reloads  []Reload
	Shots    []Shot
	Messages []Message
	Chickens []Chicken // Chicken killed
}

// Basic stats

type Kill struct {
	Tick Tick

	Headshot      bool
	NoScope       bool
	ThroughSmoke  bool
	Wallbang      bool
	AttackerBlind bool
	VictimBlind   bool
	Distance      float32

	Weapon EquipmentType
	Victim PlayerID // Can be 0 if demo is partially corrupted or the player is 'unconnected'

	KillerPosition Vector // Can be 0 if the player got killed by world damage (e.g. fall damage) or if the demo is partially corrupted or the player is 'unconnected'
	VictimPosition Vector // Can be 0 if demo is partially corrupted or the player is 'unconnected'
}

type Assist struct {
	Tick Tick

	Dead bool // If the assister was dead when it happened

	Killer      PlayerID // Can be 0 if the player got killed by world damage (e.g. fall damage) or if the demo is partially corrupted or the player is 'unconnected'
	Victim      PlayerID // Can be 0 if demo is partially corrupted or the player is 'unconnected'
	DamageDealt []Damage
}

type Death struct {
	Tick   Tick
	Weapon EquipmentType
	Killer PlayerID // Can be 0 if the player got killed by world damage (e.g. fall damage) or if the demo is partially corrupted or the player is 'unconnected'

	KillerPosition Vector // Can be 0 if the player got killed by world damage (e.g. fall damage) or if the demo is partially corrupted or the player is 'unconnected'
	VictimPosition Vector // Can be 0 if demo is partially corrupted or the player is 'unconnected'
}

//
// Nades
//

type Grenade struct {
	Start Tick
	End   Tick

	Trajectory []TrajectoryEntry
	Bounces    int

	uniqueID2 ulid.ULID
}

type TrajectoryEntry struct {
	Tick     Tick
	Position Vector
}

type Flash struct {
	Grenade

	Victims []FlashVictim
}

type FlashVictim struct {
	Tick     Tick
	Position Vector

	Victim   PlayerID
	Duration time.Duration
}

type Smoke struct {
	Grenade

	Bloom Tick
}

type He struct {
	Grenade

	Exploded Tick
	Victims  []DamageDealt
}

type Decoy struct {
	Grenade

	Victims []DamageDealt
}

type Incendiary struct {
	Grenade

	FireStart Tick
	Victims   []DamageDealt

	uniqueID int64 // Needed because for the inferno start event
}

type Molotov struct {
	Grenade

	FireStart Tick
	Victims   []DamageDealt

	uniqueID int64 // Needed because for the inferno start event
}

//
// Bomb
//

type Bomb struct {
	SpawnedWith PlayerID // Can be 0 if it is a hostage map

	Drops []*BombDrop

	Plants  []*BombPlant
	Defuses []*BombDefuse
}

type BombDrop struct {
	DropTick Tick
	Dropper  PlayerID

	Position Vector

	Picker     PlayerID // Can be 0 if no one picked it up
	PickupTick Tick     // Can be 0 if no one picked it up
}

type BombPlant struct {
	Planter PlayerID

	Start Tick
	End   Tick

	Position Vector // Position when the player starts planting
	Site     Bombsite

	Planted bool
}

type BombDefuse struct {
	Defuser PlayerID
	HasKit  bool

	Start Tick
	End   Tick

	Defused bool
}

//
// Hostage
//

type Hostage struct {
	EntityID      int
	SpawnPosition Vector

	Rescue  *HostageRescue // Can be nil if the hostages was not rescued
	Carries []*HostageCarry
}

type HostageRescue struct {
	Tick     Tick
	Position Vector

	Rescuer PlayerID
}

type HostageCarry struct {
	Start Tick
	End   Tick // Can be 0 if the hostage was never dropped (e.g. the round ended with a player carrying the hostage)

	StartPosition Vector
	EndPosition   Vector // Can be 0 if the hostage was never dropped (e.g. the round ended with a player carrying the hostage)

	Carryer PlayerID
}

//
// Item
//

type Item struct {
	Tick   Tick
	Weapon EquipmentType

	uniqueID2 ulid.ULID
}

type ItemPurchase struct {
	Item
}

type ItemDrop struct {
	Item

	Position Vector
	To       PlayerID // Who picked it up, can be 0 if nobody did
}

type ItemPickup struct {
	Item

	Position Vector
	From     PlayerID // Who dropped it
}

type ItemRefund struct {
	Item
}

//
// Spotted
//

type Spot struct {
	Start Tick
	End   Tick

	Spotted PlayerID
}

type SpottedBy struct {
	Start Tick
	End   Tick

	Spotter PlayerID
}

//
// Other
//

type Position struct {
	Tick Tick
	Vector
}

type Reload struct {
	Tick     Tick
	Position Vector

	Weapon           EquipmentType
	BulletsRemaining int
}

type Shot struct {
	Tick     Tick
	Position Vector

	Weapon EquipmentType
}

type Message struct {
	Tick Tick

	Text string
}

type Chicken struct {
	Damage
}
