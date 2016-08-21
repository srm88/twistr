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
	case CentralAmerica.Name:
		switch level {
		case Presence:
			return 1
		case Domination:
			return 3
		case Control:
			return 5
		}
	case SouthAmerica.Name:
		switch level {
		case Presence:
			return 2
		case Domination:
			return 5
		case Control:
			return 6
		}
	case Europe.Name:
		switch level {
		case Presence:
			return 3
		case Domination:
			return 7
		case Control:
			return WIN
		}
	case MiddleEast.Name:
		switch level {
		case Presence:
			return 3
		case Domination:
			return 5
		case Control:
			return 7
		}
	case Africa.Name:
		switch level {
		case Presence:
			return 1
		case Domination:
			return 4
		case Control:
			return 6
		}
	case Asia.Name:
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
	// Track separately the country that was nullified by shuttle diplomacy
	// for messaging.
	ShuttleDiplomacyNullified *Country
}

func ScoreRegion(g *Game, r Region) ScoreResult {
	result := ScoreResult{
		Levels:                    [2]ScoreLevel{Nothing, Nothing},
		Battlegrounds:             [2][]*Country{[]*Country{}, []*Country{}},
		AdjSuper:                  [2][]*Country{[]*Country{}, []*Country{}},
		ShuttleDiplomacyNullified: nil,
	}
	counts := [2]int{0, 0}
	allBattlegrounds := 0
	isBattleground := func(c *Country) bool {
		return (c.Battleground ||
			(c.Id == Taiwan && g.Effect(FormosanResolution) && c.Controlled() == USA))
	}
	for _, cid := range r.Countries {
		c := g.Countries[cid]
		if isBattleground(c) {
			allBattlegrounds += 1
		}
		aff := c.Controlled()
		if aff == NEU {
			continue
		}
		counts[aff] += 1
		if isBattleground(c) {
			if aff == SOV && g.Effect(ShuttleDiplomacy) && (&r == &Asia || &r == &MiddleEast) && result.ShuttleDiplomacyNullified == nil {
				result.ShuttleDiplomacyNullified = c
				counts[aff] -= 1
			} else {
				result.Battlegrounds[aff] = append(result.Battlegrounds[aff], c)
			}
		}
		if c.AdjSuper == aff.Opp() {
			result.AdjSuper[aff] = append(result.AdjSuper[aff], c)
		}
	}
	score := func(aff Aff) ScoreLevel {
		opp := aff.Opp()
		switch {
		case counts[aff] == 0:
			return Nothing
		case counts[aff] <= counts[opp]:
			return Presence
		case len(result.Battlegrounds[aff]) == allBattlegrounds:
			return Control
		case len(result.Battlegrounds[aff]) > len(result.Battlegrounds[opp]) && counts[aff] > len(result.Battlegrounds[aff]):
			return Domination
		default:
			// Case where you have more countries but fewer battlegrounds
			return Presence
		}
	}
	result.Levels[USA] = score(USA)
	result.Levels[SOV] = score(SOV)
	return result
}
