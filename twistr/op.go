package twistr

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
	switch {
	case rollUSA > rollSOV:
		target.Inf[SOV] -= Min((rollUSA - rollSOV), target.Inf[SOV])
	case rollSOV > rollUSA:
		target.Inf[USA] -= Min((rollSOV - rollUSA), target.Inf[USA])
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

func opsMod(s *State, player Aff, card Card, countries []*Country) (mod int) {
	if player == SOV && s.Effect(VietnamRevolts) {
		if AllIn(countries, SoutheastAsia) {
			mod += 1
		}
	}
	return
}

// Coup
func coup(s *State, player Aff, card Card, roll int, target *Country) bool {
	bonus := coupBonus(s, player, target)
	ops := card.Ops + opsMod(s, player, card, []*Country{target})
	delta := roll + bonus + ops - (target.Stability * 2)
	if delta <= 0 {
		return false
	}
	oppCurInf := target.Inf[player.Opp()]
	removed := Min(oppCurInf, delta)
	gained := delta - removed
	if target.Battleground {
		// XXX: CubanMissileCrisis, NuclearSubs
		s.Defcon -= 1
	}
	s.MilOps[player] += ops
	target.Inf[player] += gained
	target.Inf[player.Opp()] -= removed
	return true
}

func influenceCost(player Aff, target Country) int {
	controlled := target.Controlled()
	if controlled != player {
		return 2
	}
	return 1
}
