package demo

type Team byte

const (
	TeamUnassigned Team = iota
	TeamSpectators
	TeamTerrorists
	TeamCounterTerrorists
)

type RoundEndReason byte

const (
	RoundEndReasonStillInProgress RoundEndReason = iota
	RoundEndReasonTargetBombed
	RoundEndReasonVIPEscaped
	RoundEndReasonVIPKilled
	RoundEndReasonTerroristsEscaped
	RoundEndReasonCTStoppedEscape
	RoundEndReasonTerroristsStopped
	RoundEndReasonBombDefused
	RoundEndReasonCTWin
	RoundEndReasonTerroristsWin
	RoundEndReasonDraw
	RoundEndReasonHostagesRescued
	RoundEndReasonTargetSaved
	RoundEndReasonHostagesNotRescued
	RoundEndReasonTerroristsNotEscaped
	RoundEndReasonVIPNotEscaped
	RoundEndReasonGameStart
	RoundEndReasonTerroristsSurrender
	RoundEndReasonCTSurrender
	RoundEndReasonTerroristsPlanted
	RoundEndReasonCTsReachedHostage
)

type RoundMVPReason byte

const (
	MVPReasonMostEliminations RoundMVPReason = iota + 1
	MVPReasonBombDefused
	MVPReasonBombPlanted
)

type Bombsite byte

const (
	BomsiteUnknown Bombsite = 0
	BombsiteA      Bombsite = 'A'
	BombsiteB      Bombsite = 'B'
)

type EquipmentType int

const (
	EqUnknown EquipmentType = 0

	EqP2000        EquipmentType = 1
	EqGlock        EquipmentType = 2
	EqP250         EquipmentType = 3
	EqDeagle       EquipmentType = 4
	EqFiveSeven    EquipmentType = 5
	EqDualBerettas EquipmentType = 6
	EqTec9         EquipmentType = 7
	EqCZ           EquipmentType = 8
	EqUSP          EquipmentType = 9
	EqRevolver     EquipmentType = 10

	EqMP7   EquipmentType = 101
	EqMP9   EquipmentType = 102
	EqBizon EquipmentType = 103
	EqMac10 EquipmentType = 104
	EqUMP   EquipmentType = 105
	EqP90   EquipmentType = 106
	EqMP5   EquipmentType = 107

	EqSawedOff EquipmentType = 201
	EqNova     EquipmentType = 202
	EqMag7     EquipmentType = 203 // You should consider using EqSwag7 instead
	EqSwag7    EquipmentType = 203
	EqXM1014   EquipmentType = 204
	EqM249     EquipmentType = 205
	EqNegev    EquipmentType = 206

	EqGalil  EquipmentType = 301
	EqFamas  EquipmentType = 302
	EqAK47   EquipmentType = 303
	EqM4A4   EquipmentType = 304
	EqM4A1   EquipmentType = 305
	EqScout  EquipmentType = 306
	EqSSG08  EquipmentType = 306
	EqSG556  EquipmentType = 307
	EqSG553  EquipmentType = 307
	EqAUG    EquipmentType = 308
	EqAWP    EquipmentType = 309
	EqScar20 EquipmentType = 310
	EqG3SG1  EquipmentType = 311

	EqZeus                     EquipmentType = 401
	EqKevlar                   EquipmentType = 402
	EqHelmet                   EquipmentType = 403
	EqBomb                     EquipmentType = 404
	EqKnife                    EquipmentType = 405
	EqDefuseKit                EquipmentType = 406
	EqWorld                    EquipmentType = 407
	EqZoneRepulsor             EquipmentType = 408
	EqShield                   EquipmentType = 409
	EqHeavyAssaultSuit         EquipmentType = 410
	EqNightVision              EquipmentType = 411
	EqHealthShot               EquipmentType = 412
	EqTacticalAwarenessGrenade EquipmentType = 413
	EqFists                    EquipmentType = 414
	EqBreachCharge             EquipmentType = 415
	EqTablet                   EquipmentType = 416
	EqAxe                      EquipmentType = 417
	EqHammer                   EquipmentType = 418
	EqWrench                   EquipmentType = 419
	EqSnowball                 EquipmentType = 420
	EqBumpMine                 EquipmentType = 421

	EqDecoy      EquipmentType = 501
	EqMolotov    EquipmentType = 502
	EqIncendiary EquipmentType = 503
	EqFlash      EquipmentType = 504
	EqSmoke      EquipmentType = 505
	EqHE         EquipmentType = 506
)

type HitGroup byte

const (
	HitGroupGeneric  HitGroup = 0
	HitGroupHead     HitGroup = 1
	HitGroupChest    HitGroup = 2
	HitGroupStomach  HitGroup = 3
	HitGroupLeftArm  HitGroup = 4
	HitGroupRightArm HitGroup = 5
	HitGroupLeftLeg  HitGroup = 6
	HitGroupRightLeg HitGroup = 7
	HitGroupNeck     HitGroup = 8
	HitGroupGear     HitGroup = 10
)
