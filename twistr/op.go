package twistr

// All WIP. Maybe obliterate it.

// Realignment
func realignMods(target Country) (mods Influence) {
	switch {
	case target.Inf[US] > target.Inf[Sov]:
		mods[US] += 1
	case target.Inf[Sov] > target.Inf[US]:
		mods[Sov] += 1
	}
	for _, neighbor := range target.AdjCountries {
		control := neighbor.Controlled()
		if control != Neu {
			mods[control] += 1
		}
	}
	return
}

func realign(s *State, target *Country, rollUS, rollSov int) {
	mods := realignMods(*target)
	rollUS += mods[US]
	rollSov += mods[Sov]
	switch {
	case rollUS > rollSov:
		target.Inf[Sov] -= Min((rollUS - rollSov), target.Inf[Sov])
	case rollSov > rollUS:
		target.Inf[US] -= Min((rollSov - rollUS), target.Inf[US])
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
	if player == Sov && s.Effect(VietnamRevolts) {
		if AllIn(countries, SouthEastAsia) {
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
