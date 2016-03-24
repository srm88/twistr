package twistr

import (
	"fmt"
)

func Debug(s *State) {
	for {
		switch s.Solicit(USA, "S3KR1T", []string{"influence", "dump", "score", "back"}) {
		case "back":
			return
		case "influence":
			DebugInfluence(s)
		case "dump":
			DebugDump(s)
		case "score":
			DebugScore(s)
		}
	}
}

func DebugInfluence(s *State) {
	var player Aff
	switch s.Solicit(USA, "Whose", []string{"usa", "sov"}) {
	case "usa":
		player = USA
	case "sov":
		player = SOV
	default:
		return
	}
	what := s.Solicit(USA, "What", []string{"add", "remove"})
	if what != "add" && what != "remove" {
		return
	}
	il := SelectInfluence(s, player, what)
	switch what {
	case "add":
		PlaceInfluence(s, player, il)
	case "remove":
		RemoveInfluence(s, player, il)
	}
}

func DebugDump(s *State) {
	switch s.Solicit(USA, "What", []string{"country"}) {
	case "country":
		c, err := lookupCountry(s.Solicit(USA, "Which", nil))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("%#v\n", c)
	default:
		return
	}
}

func DebugScore(s *State) {
	var region Region
	switch s.Solicit(USA, "Region?", []string{"cam", "sam", "eur", "mde", "afr", "asi", "sea"}) {
	case "cam":
		region = CentralAmerica
	case "sam":
		region = SouthAmerica
	case "eur":
		region = Europe
	case "mde":
		region = MiddleEast
	case "afr":
		region = Africa
	case "asi":
		region = Asia
	case "sea":
		region = SouthEastAsia
	default:
		return
	}
	result := ScoreRegion(s, region)
	fmt.Printf("%#v\n", result)
}
