package twistr

type ScoreLevel int

const (
	Nothing ScoreLevel = iota
	Presence
	Domination
	Control
)

const WIN int = -1

func VPAward(level ScoreLevel, r Region) int {
	if level == Nothing {
		return 0
	}
	switch r.Name {
	case "CentralAmerica":
		switch level {
		case Presence:
			return 1
		case Domination:
			return 3
		case Control:
			return 5
		}
	case "SouthAmerica":
		switch level {
		case Presence:
			return 2
		case Domination:
			return 5
		case Control:
			return 6
		}
	case "Europe":
		switch level {
		case Presence:
			return 3
		case Domination:
			return 7
		case Control:
			return WIN
		}
	case "MiddleEast":
		switch level {
		case Presence:
			return 3
		case Domination:
			return 5
		case Control:
			return 7
		}
	case "Africa":
		switch level {
		case Presence:
			return 1
		case Domination:
			return 4
		case Control:
			return 6
		}
	case "Asia":
		switch level {
		case Presence:
			return 3
		case Domination:
			return 7
		case Control:
			return 9
		}
	}
	return 0
}

type ScoreResult struct {
	Levels        [2]ScoreLevel
	Battlegrounds [2][]*Country
	AdjSuper      [2][]*Country
}

func ScoreRegion(s *State, r Region) ScoreResult {
	result := ScoreResult{
		Levels:        [2]ScoreLevel{Nothing, Nothing},
		Battlegrounds: [2][]*Country{[]*Country{}, []*Country{}},
		AdjSuper:      [2][]*Country{[]*Country{}, []*Country{}},
	}
	counts := [2]int{0, 0}
	allBattlegrounds := 0
	for _, cid := range r.Countries {
		c := s.Countries[cid]
		if c.Battleground {
			allBattlegrounds += 1
		}
		aff := c.Controlled()
		if aff == NEU {
			continue
		}
		counts[aff] += 1
		if c.Battleground {
			result.Battlegrounds[aff] = append(result.Battlegrounds[aff], c)
		}
		if c.AdjSuper == aff.Opp() {
			result.AdjSuper[aff] = append(result.AdjSuper[aff], c)
		}
	}
	// Determine scoring levels. Goofy for-loop.
	for aff := USA; aff < NEU; aff++ {
		opp := aff.Opp()
		switch {
		case counts[aff] == 0:
			result.Levels[aff] = Nothing
		case counts[aff] <= counts[opp]:
			result.Levels[aff] = Presence
		case len(result.Battlegrounds[aff]) == allBattlegrounds:
			result.Levels[aff] = Control
		case len(result.Battlegrounds[aff]) > len(result.Battlegrounds[opp]) && counts[aff] > len(result.Battlegrounds[aff]):
			result.Levels[aff] = Domination
		default:
			// Case where you have more countries but fewer battlegrounds
			result.Levels[aff] = Presence
		}
	}
	return result
}
