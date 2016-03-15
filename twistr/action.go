package twistr

// All WIP. Maybe obliterate it.

// Coup
func coup(s *State, player Aff, bonus, ops int, target CountryId) bool {
	roll := Roll()
	country := s.Countries[target]
	delta := roll + bonus + ops - (country.Stability * 2)
	if delta <= 0 {
		return false
	}
	oppCurInf := country.Inf[Opp(player)]
	removed := Min(oppCurInf, delta)
	gained := delta - removed
	if country.Battleground {
		s.Defcon -= 1
	}
	s.MilOps[player] += ops
	country.Inf[player] += gained
	country.Inf[Opp(player)] -= removed
	return true
}

func influenceCost(player Aff, target Country) int {
	controlled := target.Controlled()
	if controlled != player {
		return 2
	}
	return 1
}
