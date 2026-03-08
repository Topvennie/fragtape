package parse

import (
	"math/bits"
	"slices"
	"time"

	"github.com/topvennie/fragtape/internal/worker/parse/demo"
	"github.com/topvennie/fragtape/pkg/utils"
)

func weaponPrimary(eq demo.EquipmentType) bool {
	return eq > 100 && eq < 400
}

func weaponSecondary(eq demo.EquipmentType) bool {
	return eq > 0 && eq < 100
}

func weaponSniper(weapon demo.EquipmentType) bool {
	return weapon == demo.EqAWP || weapon == demo.EqSSG08 || weapon == demo.EqG3SG1 || weapon == demo.EqScar20
}

func weaponShotgun(weapon demo.EquipmentType) bool {
	return weapon > 200 && weapon < 205
}

func weaponHeavy(weapon demo.EquipmentType) bool {
	return weapon == demo.EqM249 || weapon == demo.EqNegev
}

func weaponSmg(weapon demo.EquipmentType) bool {
	return weapon > 100 && weapon < 200
}

func weaponZeus(eq demo.EquipmentType) bool {
	return eq == 401
}

func isDirty(kills []kill, count ...int) bool {
	return isX(kills, func(k kill) bool { return k.blind || k.wallbang || k.throughSmoke }, count)
}

func isWeapon(kills []kill, weapon demo.EquipmentType, count ...int) bool {
	return isX(kills, func(k kill) bool { return k.weapon == weapon }, count)
}

func isOneTap(kills []kill, count ...int) bool {
	return isX(kills, func(k kill) bool { return k.oneTap }, count)
}

func isHeadshot(kills []kill, count ...int) bool {
	return isX(kills, func(k kill) bool { return k.headshot }, count)
}

func isNoScope(kills []kill, count ...int) bool {
	return isX(kills, func(k kill) bool { return k.noScope }, count)
}

func isBlind(kills []kill, count ...int) bool {
	return isX(kills, func(k kill) bool { return k.blind }, count)
}

func isWallbang(kills []kill, count ...int) bool {
	return isX(kills, func(k kill) bool { return k.wallbang }, count)
}

func isThroughSmoke(kills []kill, count ...int) bool {
	return isX(kills, func(k kill) bool { return k.throughSmoke }, count)
}

func isPistol(kills []kill, count ...int) bool {
	return isX(kills, func(k kill) bool { return weaponSecondary(k.weapon) }, count)
}

func isSmg(kills []kill, count ...int) bool {
	return isX(kills, func(k kill) bool { return weaponSmg(k.weapon) }, count)
}

func isShotgun(kills []kill, count ...int) bool {
	return isX(kills, func(k kill) bool { return weaponShotgun(k.weapon) }, count)
}

func isHeavy(kills []kill, count ...int) bool {
	return isX(kills, func(k kill) bool { return weaponHeavy(k.weapon) }, count)
}

func isSniper(kills []kill, count ...int) bool {
	return isX(kills, func(k kill) bool { return weaponSniper(k.weapon) }, count)
}

func isEnemyPrimary(kills []kill, count ...int) bool {
	return isX(kills, func(k kill) bool { return weaponPrimary(k.victimWeapon) }, count)
}

func isTeamKill(kills []kill, count ...int) bool {
	return isX(kills, func(k kill) bool { return k.teamKill }, count)
}

// isX is a helper function
// Shouldn't be used outside of this file
func isX[T any](slice []T, cond func(t T) bool, counts []int) bool {
	count := len(slice)
	if len(counts) > 0 {
		count = counts[0]
	}

	return len(utils.SliceFilter(slice, cond)) >= count
}

func isFast(kills []kill, seconds int) bool {
	return kills[len(kills)-1].tickRel-kills[0].tickRel <= time.Duration(seconds)*time.Second
}

func isClutch(r round, p player) bool {
	return r.clutcher == p.id && r.clutchKills == len(p.kills)
}

func isSameWeapon(kills []kill) bool {
	return utils.SliceAll(kills, func(k kill) bool { return k.weapon == kills[0].weapon })
}

func trifectaDirty(kills []kill) bool {
	if len(kills) != 3 {
		return false
	}

	return (kills[0].blind && kills[1].wallbang && kills[2].throughSmoke) ||
		(kills[0].blind && kills[2].wallbang && kills[1].throughSmoke) ||
		(kills[1].blind && kills[0].wallbang && kills[2].throughSmoke) ||
		(kills[1].blind && kills[2].wallbang && kills[0].throughSmoke) ||
		(kills[2].blind && kills[0].wallbang && kills[1].throughSmoke) ||
		(kills[2].blind && kills[1].wallbang && kills[0].throughSmoke)
}

// getIdx returns all possible combinations of indexes with length <= n
func getIdx(n int) [][]int {
	count := (1 << n) - 1

	indexes := make([][]int, 0, count)

	for mask := 1; mask <= count; mask++ {
		vals := make([]int, 0, bits.OnesCount(uint(mask)))
		for i := range n {
			if mask&(1<<i) != 0 {
				vals = append(vals, i)
			}
		}

		indexes = append(indexes, vals)
	}

	slices.SortFunc(indexes, func(a, b []int) int { return len(b) - len(a) })
	return indexes
}
