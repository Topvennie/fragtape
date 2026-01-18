package demo

import (
	"math"
	"slices"

	"github.com/golang/geo/r3"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	"github.com/oklog/ulid/v2"
	"github.com/topvennie/fragtape/pkg/utils"
)

// Generic

func closest[T any](center Vector, items []T, fn func(t T) Vector) T {
	var minDist float32 = -1
	var minItem T

	for _, item := range items {
		dist := center.Distance(fn(item))
		if minDist == -1 || dist < minDist {
			minDist = dist
			minItem = item
		}
	}

	return minItem
}

// Demo

func (d *Demo) playerIndex(id PlayerID) int {
	return slices.IndexFunc(d.match.Players, func(p *Player) bool { return p.SteamID == id })
}

func (d *Demo) getLastRound() *Round {
	return utils.SliceLast(d.match.Rounds)
}

func (d *Demo) getPlayerSpots(id PlayerID) map[PlayerID]bool {
	if spots, ok := d.playerSpots[id]; ok {
		return spots
	}

	d.playerSpots[id] = map[PlayerID]bool{}
	return d.playerSpots[id]
}

// Round

func (r *Round) getLastBombDefuse() *BombDefuse {
	return utils.SliceLast(r.Bomb.Defuses)
}

func (r *Round) getLastBombPlant() *BombPlant {
	return utils.SliceLast(r.Bomb.Plants)
}

func (r *Round) getLastBombDrop() *BombDrop {
	return utils.SliceLast(r.Bomb.Drops)
}

func (r *Round) getHostage(id int) *Hostage {
	if idx := slices.IndexFunc(r.Hostages, func(h *Hostage) bool { return h.EntityID == id }); idx != -1 {
		return r.Hostages[idx]
	}

	return nil
}

// Vector

func toVector(v r3.Vector) Vector {
	return Vector{
		X: float32(v.X),
		Y: float32(v.Y),
		Z: float32(v.Z),
	}
}

func (v Vector) Distance(v2 Vector) float32 {
	return float32(math.Sqrt(float64((v.X-v2.X)*(v.X-v2.X) + (v.Y-v2.Y)*(v.Y-v2.Y) + (v.Z-v2.Z)*(v.Z-v2.Z))))
}

func (v Vector) IsZero() bool {
	return v.X == 0 && v.Y == 0 && v.Z == 0
}

// Stat

func (s Stat) getLastPosition() Position {
	return utils.SliceLast(s.Positions)
}

func (s Stat) getNade(uniqueID2 ulid.ULID) any {
	for _, f := range s.Flashes {
		if f.uniqueID2 == uniqueID2 {
			return f
		}
	}
	for _, s := range s.Smokes {
		if s.uniqueID2 == uniqueID2 {
			return s
		}
	}
	for _, h := range s.Hes {
		if h.uniqueID2 == uniqueID2 {
			return h
		}
	}
	for _, d := range s.Decoys {
		if d.uniqueID2 == uniqueID2 {
			return d
		}
	}
	for _, i := range s.Incendiaries {
		if i.uniqueID2 == uniqueID2 {
			return i
		}
	}
	for _, m := range s.Molotovs {
		if m.uniqueID2 == uniqueID2 {
			return m
		}
	}

	return nil
}

func (s Stat) getSpottedAfter(start Tick, spotted PlayerID) []*Spot {
	spots := make([]*Spot, 0)

	for i := len(s.Spots) - 1; i >= 0; i-- {
		if s.Spots[i].Start < start {
			break
		}

		if s.Spots[i].Spotted == spotted {
			spots = append(spots, s.Spots[i])
		}
	}

	return spots
}

func (s *Stat) openSpot(start, end Tick, spotted PlayerID) {
	for i := len(s.Spots) - 1; i >= 0; i-- {
		spot := s.Spots[i]
		if spot.Spotted == spotted {
			if spot.End == 0 {
				return
			}
			break
		}
	}

	s.Spots = append(s.Spots, &Spot{Start: start, End: end, Spotted: spotted})
}

func (s *Stat) closeSpot(end Tick, spotted PlayerID) {
	for i := len(s.Spots) - 1; i >= 0; i-- {
		spot := s.Spots[i]
		if spot.Spotted == spotted && spot.End == 0 {
			spot.End = end
			return
		}
	}
}

func (s *Stat) openSpottedBy(start, end Tick, spotter PlayerID) {
	for i := len(s.SpottedBy) - 1; i >= 0; i-- {
		spot := s.SpottedBy[i]
		if spot.Spotter == spotter {
			if spot.End == 0 {
				return
			}
			break
		}
	}

	s.SpottedBy = append(s.SpottedBy, &SpottedBy{Start: start, End: end, Spotter: spotter})
}

func (s *Stat) closeSpottedBy(end Tick, spotter PlayerID) {
	for i := len(s.SpottedBy) - 1; i >= 0; i-- {
		spot := s.SpottedBy[i]
		if spot.Spotter == spotter && spot.End == 0 {
			spot.End = end
			return
		}
	}
}

// TrajectoryEntry

func toTrajectoryEntrySlice(entries []common.TrajectoryEntry) []TrajectoryEntry {
	return utils.SliceMap(entries, func(e common.TrajectoryEntry) TrajectoryEntry {
		return TrajectoryEntry{
			Tick:     Tick(e.Tick),
			Position: toVector(e.Position),
		}
	})
}

// Other

// getPlayer populates the player with zero values if it is nil
func getPlayer(player *common.Player) *common.Player {
	if player == nil {
		return &common.Player{}
	}

	return player
}
