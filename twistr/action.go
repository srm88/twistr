package twistr

// All WIP. Maybe obliterate it.

// Coup
func coup(s *State, player Aff, bonus, ops int, target CountryId) bool {
	roll := Roll()
	country := s.Countries[target]
	delta := roll + bonus + ops - (country.Stab * 2)
	if delta <= 0 {
		return false
	}
	old := country.Inf[player]
	removed := Min(old, delta)
	gained := delta - removed
	if country.Battleground {
		s.Defcon -= 1
	}
	s.MilOps[player] += ops
	country.Inf[player] += gained
	country.Inf[Other(player)] -= removed
	return true
}

func influenceCost(player Aff, target Country) int {
	controlled := target.Controlled()
	if controlled != player {
		return 2
	}
	return 1
}
