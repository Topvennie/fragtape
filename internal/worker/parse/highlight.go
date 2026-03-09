package parse

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/worker/parse/demo"
	"github.com/topvennie/fragtape/pkg/utils"
)

func (m *Manager) saveHighlights(ctx context.Context, d model.Demo, match demo.Match) error {
	if err := m.highlight.DeleteByDemo(ctx, d.ID); err != nil {
		return err
	}

	// Get all users
	users, err := m.user.GetByUIDs(ctx, utils.SliceMap(match.Players, func(p *demo.Player) int64 { return int64(p.SteamID) }))
	if err != nil {
		return err
	}
	// Only keep the real users
	users = utils.SliceFilter(users, func(u *model.User) bool { return u.IsReal() })

	settings, err := m.setting.Get(ctx)
	if err != nil {
		return err
	}

	// Get converted round data
	rounds := convertMatch(match)

	highlights := make([]*model.Highlight, 0)

	// Go over round
	for _, round := range rounds {
		for _, user := range users {
			player, ok := round.players[demo.PlayerID(user.UID)]
			if !ok {
				// User didnt play this round
				continue
			}

			// Determine if the user gets a highlight
			// If title != "" then the user get's a highlight of this round
			// All kills of that round are recorded, regardless of why the highlight was triggered
			title := ""

			indexes := getIdx(len(player.kills))
			for _, idx := range indexes {
				if title != "" {
					// we found clip
					break
				}

				switch len(idx) {
				case 1:
					title = killOne(round, player, idx)
				case 2:
					title = killTwo(round, player, idx)
				case 3:
					title = killThree(round, player, idx)
				case 4:
					title = killFour(round, player, idx)
				case 5:
					title = killFive(round, player, idx)
				}
			}

			if title == "" {
				// No clip found
				// Maybe a special case?
				title = special(player, *settings)

				if title == "" {
					// Nope still no clip
					continue
				}
			}

			// Construct highlight
			highlight := constructHighlight(*user, d, match, round, player, title)
			highlights = append(highlights, highlight)
		}
	}

	for _, highlight := range highlights {
		if err := m.highlight.Create(ctx, highlight); err != nil {
			return err
		}

		for _, segment := range highlight.Segments {
			segment.HighlightID = highlight.ID
			if err := m.highlight.CreateSegment(ctx, &segment); err != nil {
				return err
			}
		}
	}

	return nil
}

func killOne(_ round, _ player, _ []int) string {
	// A 1k never leads to a highlight
	return ""
}

func killTwo(round round, player player, killIdx []int) string {
	kills := utils.SliceMap(killIdx, func(idx int) kill { return player.kills[idx] })

	// Not all 2k's lead to a highlight

	// 1 v 2 clutch
	if isClutch(round, player) {
		return utils.SliceRandom([]string{
			"1v2 clutch",
			"Ice Cold 1v2",
			"Clutch Or Kick",
			"Ice In The Veins",
			"Clutch Minister",
			"Last Man Standing",
			"Calm Under Pressure",
			"Saved The Round",
			"Clutch Gene",
		})
	}

	// Collateral
	if kills[0].tick == kills[1].tick {
		return utils.SliceRandom([]string{
			"Collateral",
			"Two For One",
			"One Bullet, Two Bodies",
			"Lineup",
			"You Meant To Do That?",
		})
	}

	// All one taps and no sniper
	if isOneTap(kills) && !isSniper(kills, 1) {
		return utils.SliceRandom([]string{
			"2K One Taps",
			"Tap Tap",
			"Scream Moment",
			"Aim Diff",
			"NiKo Mode",
			"Straight To Demo Review",
		})
	}

	// No scopes
	if isNoScope(kills) {
		return utils.SliceRandom([]string{
			"RNG On Your Side",
			"No Scope Double",
			"Broky Immitation",
			"What Is A Scope",
		})
	}

	// Fast 2k (<= 3 seconds)
	if isFast(kills, 3) {
		// <= 1 second
		if isFast(kills, 1) {
			return utils.SliceRandom([]string{
				"Lightning Fast 2k",
				"Instant Double",
				"Snap Double",
				"Blink And It's Two",
				"Two In A HeartBeat",
			})
		}

		return utils.SliceRandom([]string{
			"Fast 2K",
			"Quick Double",
			"Back To Back",
			"Two Quick Picks",
		})
	}

	// Heavy weapon
	if isHeavy(kills) {
		return utils.SliceRandom([]string{
			"Adding Insult To Injury",
			"Hold Left Click",
			"Mobile Turret",
			"Noise Complaint",
		})
	}

	// Fully blind
	if isBlind(kills) {
		return utils.SliceRandom([]string{
			"Spray And Pray",
			"No Vision Needed",
			"Guessing Correctly",
			"Flash? Still Shooting",
		})
	}

	// Wallbang
	if isWallbang(kills) {
		return utils.SliceRandom([]string{
			"No Hiding",
			"Paper Walls",
			"X-Ray Vision",
			"Through The Wall",
		})
	}

	// Player doesn't have a sniper and enemies had an awp both times
	if !isSniper(kills, 1) && utils.SliceAll(kills, func(k kill) bool { return k.victimWeapon == demo.EqAWP }) {
		return utils.SliceRandom([]string{
			"Anti Awp",
			"AWP Confiscation",
			"Sniper Shutdown",
		})
	}

	// Just an ordinary 2k so no clip
	return ""
}

func killThree(round round, p player, killIdx []int) string {
	kills := utils.SliceMap(killIdx, func(idx int) kill { return p.kills[idx] })

	// Not all 3k's lead to a highlight

	// 1 v 3 clutch
	if isClutch(round, p) {
		return utils.SliceRandom([]string{
			"1v3 Clutch",
			"Ice Cold 1v3",
			"Clutch Or Kick",
			"Ice In The Veins",
			"Clutch Minister",
			"Last Man Standing",
			"Calm Under Pressure",
			"Saved The Round",
			"Clutch Gene",
		})
	}

	// All one taps and no sniper
	if isOneTap(kills) && !isSniper(kills, 1) {
		return utils.SliceRandom([]string{
			"3K One Taps",
			"Tap Tap Tap",
			"Scream Moment",
			"Aim Diff",
			"NiKo Mode",
			"Straight To Demo Review",
		})
	}

	// Fast 3k (<= 8 seconds)
	if isFast(kills, 8) {
		// <= 4 seconds
		if isFast(kills, 4) {
			return utils.SliceRandom([]string{
				"Lightning Fast 3k",
				"Instant Triple",
				"Snap Triple",
				"Blink And It's Three",
				"Three In A HeartBeat",
			})
		}

		return utils.SliceRandom([]string{
			"Fast 3K",
			"Quick Triple",
			"Back To Back To Back",
			"Three Quick Picks",
		})
	}

	// Dirty kills
	if isDirty(kills) {
		// One blind, one wallbang and one through a smoke
		if trifectaDirty(kills) {
			return utils.SliceRandom([]string{
				"Trifecta dirty",
				"Full Cheat Sheet",
				"All The Sins",
				"Blind. Wall. Smoke.",
				"VAC Speed-Dial",
				"Report Incoming",
			})
		}

		// Blind
		if isBlind(kills) {
			return utils.SliceRandom([]string{
				"Spray And Pray",
				"No Vision Needed",
				"Guessing Correctly",
				"Flash? Still Shooting",
			})
		}

		// Through smokes
		if isThroughSmoke(kills) {
			return utils.SliceRandom([]string{
				"Peekaboo",
				"Smoke Spam Triple",
				"Through The Grey",
				"X-Ray Vision",
				"Smoked 'em",
			})
		}

		return utils.SliceRandom([]string{
			"Filthy 3K",
			"No Honor",
			"Dirty Work",
			"Zero Shame",
			"Unclean Triple",
		})
	}

	// Pistol 3k
	if isPistol(kills) {
		// Deagle
		if isWeapon(kills, demo.EqDeagle) {
			return utils.SliceRandom([]string{
				"Deag Clan Application",
				"Hand Cannon Triple",
				"Deagle Demon",
			})
		}

		// Dualies
		if isWeapon(kills, demo.EqDualBerettas) {
			return utils.SliceRandom([]string{
				"Wild West 3K",
				"Gunslinger",
			})
		}

		// Revolver
		if isWeapon(kills, demo.EqRevolver) {
			return utils.SliceRandom([]string{
				"Cowboy",
				"Sheriff In Town",
				"Spin And Win",
				"Old West Triple",
				"Revolver Roulette",
			})
		}

		// Usp and one taps
		if isWeapon(kills, demo.EqUSP) && isOneTap(kills) {
			return utils.SliceRandom([]string{
				"007 3K",
				"Silent Triple",
				"Agent Mode",
				"Tap Tap Tap",
			})
		}

		// Default
		return utils.SliceRandom([]string{
			"Pistol 3K",
			"Sidearm Triple",
			"Pocket Rocket",
		})
	}

	// Entry kill
	allKills := utils.SliceFlatten(utils.SliceMap(utils.MapValues(round.players), func(p player) []kill { return p.kills }))
	if utils.SliceAll(allKills, func(k kill) bool { return k.tickRel >= kills[0].tickRel }) {
		return utils.SliceRandom([]string{
			"Entry 3K",
			"First Blood To Triple",
			"Entry Fragger",
			"Round Opener",
			"Space Creator",
			"Donk?!",
		})
	}

	// First kill with pistol and enemy always had a primary
	if weaponSecondary(kills[0].weapon) && isEnemyPrimary(kills) {
		return utils.SliceRandom([]string{
			"Eco 3K",
			"Budget Triple",
			"Eco To Hero",
			"David vs Goliath",
			"Round Theft",
		})
	}

	// All kills with shotgun and enemy always had a primary
	if isShotgun(kills) && isEnemyPrimary(kills) {
		return utils.SliceRandom([]string{
			"Don't Come Near Me",
			"Close Quarters",
			"Boomstick Business",
		})
	}

	// All kills with smg and enemy always had a primary
	if isSmg(kills) && isEnemyPrimary(kills) {
		return utils.SliceRandom([]string{
			"Run And Gun",
			"SMG Farming",
		})
	}

	// All headshots
	if isHeadshot(kills) {
		return utils.SliceRandom([]string{
			"Barber 3K",
			"Headhunter",
			"Aim Lab Graduate",
			"Precision Triple",
		})
	}

	// Ordinary 3k so no clip
	return ""
}

func killFour(round round, player player, killIdx []int) string {
	kills := utils.SliceMap(killIdx, func(idx int) kill { return player.kills[idx] })

	// 4 kills always lead to a highlight
	// We just need to figure out a name

	// 1 v 4 clutch
	if isClutch(round, player) {
		return utils.SliceRandom([]string{
			"1v4 clutch",
			"Xyp9x Mode",
		})
	}

	// All one taps and no sniper
	if isOneTap(kills) && !isSniper(kills, 1) {
		return utils.SliceRandom([]string{
			"4K One Taps",
			"Quad Dink",
			"Tap Tap Tap Tap",
			"Straight To Reddit",
		})
	}

	// Fast 4k (<= 15 seconds)
	if isFast(kills, 15) {
		return utils.SliceRandom([]string{
			"Lightning Fast 4K",
			"4k Blink And It's Over",
			"Speedrun 4K",
			"Rapid Quad",
			"Instant Impact",
		})
	}

	// >= 3 dirty kills
	if isDirty(kills, 3) {
		return utils.SliceRandom([]string{
			"VAC Check",
			"Unreal 4K",
			"Certified Suspicious",
			"Report Incoming",
			"Overwatch Highlight",
		})
	}

	// >= 3 noscopes
	if isNoScope(kills, 3) {
		return utils.SliceRandom([]string{
			"Sniper Casino 4K",
			"Scope? Optional.",
			"RNG Masterclass",
		})
	}

	// Heavy
	if isHeavy(kills) {
		return utils.SliceRandom([]string{
			"Heavy 4K",
			"Bullet Hose",
			"Mobile Turret",
			"Surpressing Fire",
			"Hold Left Click",
			"Noise Complaint",
		})
	}

	// All headshot
	if isHeadshot(kills) {
		return utils.SliceRandom([]string{
			"Surgical 4K",
			"Barber Shop",
			"Head Collector",
		})
	}

	// Same weapon
	if isSameWeapon(kills) {
		// USP
		if isWeapon(kills, demo.EqUSP) {
			return utils.SliceRandom([]string{
				"Silent quad",
				"Pistol Diff",
			})
		}

		// Deagle
		if isWeapon(kills, demo.EqDeagle) {
			return utils.SliceRandom([]string{
				"Hand Cannon Quad",
				"Deagle Demon",
				"One Deag, Four Down",
			})
		}

		// Dualies
		if isWeapon(kills, demo.EqDualBerettas) {
			return utils.SliceRandom([]string{
				"Wild West 3K",
				"Gunslinger",
				"Two Guns, Four Down",
			})
		}

		// Revolver
		if isWeapon(kills, demo.EqRevolver) {
			return utils.SliceRandom([]string{
				"Yeeh Haw",
				"Cowboy Quad",
				"Spin And Win",
				"Old West Quad",
			})
		}

		// Scout
		if isWeapon(kills, demo.EqScout) {
			return utils.SliceRandom([]string{
				"Reported",
				"Scout Quad",
				"What is an AWP",
			})
		}

		// AWP
		if isWeapon(kills, demo.EqAWP) {
			return utils.SliceRandom([]string{
				"Kennys 4K",
				"AWP Quad",
				"Big Green Quad",
				"S1mple?!",
			})
		}

		// Default
		return utils.SliceRandom([]string{
			fmt.Sprintf("%s 4K", kills[0].weapon),
			fmt.Sprintf("%s Quad", kills[0].weapon),
			fmt.Sprintf("%s Masterclass", kills[0].weapon),
			fmt.Sprintf("4K With %s", kills[0].weapon),
			fmt.Sprintf("%s Only", kills[0].weapon),
		})
	}

	// No special name
	return utils.SliceRandom([]string{
		"4K",
		"Quad Kill",
	})
}

func killFive(round round, player player, killIdx []int) string {
	kills := utils.SliceMap(killIdx, func(idx int) kill { return player.kills[idx] })

	// 5 kills always lead to a highlight
	// We just need to figure a name

	// 1 v 5 clutch
	if isClutch(round, player) {
		return utils.SliceRandom([]string{
			"1v5 ACE CLUTCH",
		})
	}

	// All one taps and no sniper
	if isOneTap(kills) && !isSniper(kills, 1) {
		return utils.SliceRandom([]string{
			"SCREAM ACE",
			"ONE TAP ACE",
		})
	}

	// Fast ace (<= 20 seconds)
	if isFast(kills, 20) {
		return utils.SliceRandom([]string{
			"LIGHTNING FAST ACE",
			"SPEEDRUN ACE",
		})
	}

	// >= 3 dirty kills
	if isDirty(kills, 3) {
		return utils.SliceRandom([]string{
			"ACE VAC CHECK",
		})
	}

	// >= 3 noscopes
	if isNoScope(kills, 3) {
		return utils.SliceRandom([]string{
			"ACE SNIPER CASINO",
		})
	}

	// Heavy
	if isHeavy(kills) {
		return utils.SliceRandom([]string{
			"HEAVY 5K",
		})
	}

	// All headshot
	if isHeadshot(kills) {
		return utils.SliceRandom([]string{
			"ACE HEADSHOT",
		})
	}

	// Same weapon
	if isSameWeapon(kills) {
		// USP
		if isWeapon(kills, demo.EqUSP) {
			return utils.SliceRandom([]string{
				"JAMES 'ACE' BOND",
			})
		}

		// Deagle
		if isWeapon(kills, demo.EqDeagle) {
			return utils.SliceRandom([]string{
				"HAND CANNON ACE",
			})
		}

		// Dualies
		if isWeapon(kills, demo.EqDualBerettas) {
			return utils.SliceRandom([]string{
				"GUNSLINGER ACE",
			})
		}

		// Revolver
		if isWeapon(kills, demo.EqRevolver) {
			return utils.SliceRandom([]string{
				"COWBOY ACE",
			})
		}

		// Default
		return utils.SliceRandom([]string{
			fmt.Sprintf("%s ACE", kills[0].weapon),
			fmt.Sprintf("ACE WITH %s", kills[0].weapon),
		})
	}

	// No special name
	return "ACE"
}

func special(player player, setting model.SettingGlobal) string {
	kills := player.kills

	// Knife
	if isWeapon(kills, demo.EqKnife, 1) {
		return utils.SliceRandom([]string{
			"Knife Kill",
			"Stabber",
			"Get Leetify",
			"Cut Throat",
		})
	}

	// Zeus
	if isWeapon(kills, demo.EqZeus, 1) {
		return utils.SliceRandom([]string{
			"Zeus Kill",
			"Taser'd",
			"Shock Therapy",
			"ThunderStruck",
		})
	}

	// Team kill
	if isTeamKill(kills, 1) {
		return utils.SliceRandom([]string{
			"Teamkill",
			"Friendly Fire",
			"Misclick",
			"Oops",
			"Imposter",
			"Communication Issue",
		})
	}

	// Kill with smoke, flashbang or decoy
	if isWeapon(kills, demo.EqDecoy, 1) || isWeapon(kills, demo.EqFlash, 1) || isWeapon(kills, demo.EqSmoke, 1) {
		return utils.SliceRandom([]string{
			"Utility Kill",
			"Tactical Masterpiece",
			"Support Player",
		})
	}

	// In chat
	if setting.ChatCommand {
		trigger := strings.ToLower(setting.ChatTrigger)

		for _, msg := range player.messages {
			if strings.Contains(strings.ToLower(msg), trigger) {
				return utils.SliceRandom([]string{
					"Clip it",
				})
			}
		}
	}

	// Nothing found
	return ""
}

func constructHighlight(user model.User, demo model.Demo, match demo.Match, round round, player player, title string) *model.Highlight {
	// Get all segments
	segments := make([]model.HighlightSegment, 0, len(player.kills))

	for _, k := range player.kills {
		// Use a 5 second before and after the kill buffer
		start := int(k.tick - 5*match.TickRate)
		end := int(k.tick + 5*match.TickRate)

		// Get the end of the previous segment
		prevEnd := 0
		if len(segments) > 0 {
			prevEnd = segments[len(segments)-1].EndTick
		}

		// If the previous segment overlaps with this one
		// or if there's less than 2 seconds between them
		// merge them
		if start-int(2*match.TickRate) <= prevEnd {
			segments[len(segments)-1].EndTick = max(end, prevEnd)
			continue
		}

		// It's a new segment so add it
		segments = append(segments, model.HighlightSegment{
			StartTick: start,
			EndTick:   end,
		})
	}

	// We have all the segments
	// Construct the highlight
	var duration time.Duration = 0
	for _, seg := range segments {
		duration += time.Second * time.Duration((seg.EndTick-seg.StartTick)/int(match.TickRate))
	}

	return &model.Highlight{
		DemoID:   demo.ID,
		UserID:   user.ID,
		Title:    title,
		Round:    round.number,
		Duration: duration,
		Segments: segments,
	}
}
