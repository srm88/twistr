package twistr

import "fmt"

// All WIP. Maybe obliterate it.

// Realignment
func realignMods(target Country) (mods Influence) {
	switch {
	case target.Inf[USA] > target.Inf[SOV]:
		mods[USA] += 1
	case target.Inf[SOV] > target.Inf[USA]:
		mods[SOV] += 1
	}
	for _, neighbor := range target.AdjCountries {
		control := neighbor.Controlled()
		if control != NEU {
			mods[control] += 1
		}
	}
	return
}

func realign(s *State, target *Country, rollUSA, rollSOV int) {
	mods := realignMods(*target)
	rollUSA += mods[USA]
	rollSOV += mods[SOV]
	initUSA := target.Inf[USA]
	initSOV := target.Inf[SOV]
	switch {
	case rollUSA > rollSOV:
		target.Inf[SOV] -= Min((rollUSA - rollSOV), target.Inf[SOV])

	case rollSOV > rollUSA:
		target.Inf[USA] -= Min((rollSOV - rollUSA), target.Inf[USA])
	}
	if initUSA > target.Inf[USA] {
		MessageBoth(s, fmt.Sprintf("Removed %d influence from USA", initUSA-target.Inf[USA]))
	} else if initSOV > target.Inf[SOV] {
		MessageBoth(s, fmt.Sprintf("Removed %d influence from USSR", initSOV-target.Inf[SOV]))
	} else {
		MessageBoth(s, "No influence removed")
	}
}

func coupBonus(s *State, player Aff, target *Country) (bonus int) {
	if s.Effect(SALTNegotiations) {
		bonus -= 1
	}
	if s.Effect(LatinAmericanDeathSquads, player) {
		bonus += 1
	}
	if s.Effect(LatinAmericanDeathSquads, player.Opp()) {
		bonus -= 1
	}
	return
}

func opsMod(s *State, player Aff, countries []*Country) (mod int) {
	if player == SOV && s.Effect(VietnamRevolts) {
		if AllIn(countries, SoutheastAsia) {
			mod += 1
		}
	}
	return
}

// Coup
func coup(s *State, player Aff, ops int, roll int, target *Country, free bool) bool {
	bonus := coupBonus(s, player, target)
	delta := roll + bonus + ops - (target.Stability * 2)
	if delta <= 0 {
		return false
	}
	oppCurInf := target.Inf[player.Opp()]
	removed := Min(oppCurInf, delta)
	gained := delta - removed
	if target.Battleground {
		// XXX: CubanMissileCrisis, NuclearSubs
		s.DegradeDefcon(1)
	}
	if !free {
		s.MilOps[player] += ops
	}
	target.Inf[player] += gained
	target.Inf[player.Opp()] -= removed
	return true
}

// A country cannot be coup'd if it lacks any of the opponent's influence.
// Some permanent events also impose coup restrictions, e.g. NATO with Europe.
func canCoup(s *State, player Aff, t *Country, free bool) bool {
	switch {
	case t.Inf[player.Opp()] < 1:
		return false
	case natoProtected(s, player, t):
		return false
	case japanProtected(s, player, t):
		return false
	case s.Effect(TheReformer) && player == SOV && t.In(Europe):
		return false
	case defconProtected(s, t) && !free:
		return false
	default:
		return true
	}
}

func canRealign(s *State, player Aff, t *Country, free bool) bool {
	switch {
	case natoProtected(s, player, t):
		return false
	case japanProtected(s, player, t):
		return false
	case t.Inf[player.Opp()] < 1:
		return false
	case defconProtected(s, t) && !free:
		return false
	default:
		return true
	}
}

func defconProtected(s *State, t *Country) bool {
	// asia 3, defcon 5, not protected
	// europe 4, defcon 3, protected
	// middle east 2, defcon 2, protected
	return t.Region.Volatility >= s.Defcon
}

func natoProtected(s *State, player Aff, t *Country) bool {
	return s.Effect(NATO) && player == SOV && t.In(Europe) && t.Controlled() == USA
}

func japanProtected(s *State, player Aff, t *Country) bool {
	return s.Effect(USJapanMutualDefensePact) && t.Id == Japan && player == SOV
}

func influenceCost(player Aff, target Country) int {
	controlled := target.Controlled()
	if controlled != player {
		return 2
	}
	return 1
}
