package twistr

import (
	"fmt"
	"strings"
)

// Game-running functions.
// Each function should represent a state in the game.

// Deck / Hand states.
// Special cards:
// StarWars: search discard
// SALTNegotiations: search discard
// AskNotWhatYourCountry: discard up to hand, draw replacements
// OurManInTehran: draw top 5, return or discard, reshuffle

// WIP
func PlayCard(s *State, c *CardPlayLog) {
	switch {
	case c.Kind == SPACE:
		next := &SpaceLog{}
		s.Input.GetInput(c.Player, "Space roll", next)
	case c.Kind == OPS && c.Card.Aff == c.Player.Opp():
		// Solicit who goes first
		next := &OpponentOpsLog{}
		s.Input.GetInput(c.Player, "Who's next", next)
	case c.Kind == OPS:
		// Solicit coup/influence/realign/space
		next := &OpsLog{}
		s.Input.GetInput(c.Player, "What kinda ops", next)
	case c.Kind == EVENT:
		panic("Not ready!")
	default:
		panic("WUT R U DOIN")
	}
}

type countryCheck func(*Country) error

// InRegion returns a countryCheck that will reject any country that is not in
// at least one of the given regions.
func InRegion(regions ...Region) countryCheck {
	return func(c *Country) error {
		for _, r := range regions {
			if c.In(r) {
				return nil
			}
		}
		rNames := make([]string, len(regions))
		for i, r := range regions {
			rNames[i] = r.Name
		}
		return fmt.Errorf("%s not in %s", c.Name, strings.Join(rNames, " or "))
	}
}

// SelectNInfluenceCheck asks the player to choose a number of countries to
// receive influence, and optional checks to perform on the chosen countries.
func SelectNInfluenceCheck(s *State, player Aff, message string, n int, checks ...countryCheck) (il *InfluenceLog, err error) {
	il = SelectInfluence(s, player, message)
	if len(il.Countries) != n {
		err = fmt.Errorf("Place %d influence", n)
		return
	}
	for _, placement := range il.Countries {
		for _, check := range checks {
			if err = check(placement); err != nil {
				return
			}
		}
	}
	return
}

func SelectInfluenceOps(s *State, player Aff, card Card) (il *InfluenceLog, err error) {
	message := "Place influence"
	il = SelectInfluence(s, player, message)
	// Compute ops
	ops := card.Ops + opsMod(s, player, card, il.Countries)
	// Compute cost. Copy each country so that we can update its influence
	// as we go. E.g. two ops are spent breaking control, then the next
	// influence place costs one op.
	cost := 0
	workingCountries := make(map[CountryId]Country)
	for _, c := range il.Countries {
		workingCountries[c.Id] = *c
	}
	for _, c := range il.Countries {
		cost += influenceCost(player, workingCountries[c.Id])
		tmp := workingCountries[c.Id]
		tmp.Inf[player] += 1
		workingCountries[c.Id] = tmp
	}
	switch {
	case cost > ops:
		err = fmt.Errorf("Overspent ops by %d", (cost - ops))
	case cost < ops:
		err = fmt.Errorf("Underspent ops by %d", (ops - cost))
	}
	return
}

func SelectInfluence(s *State, player Aff, message string) *InfluenceLog {
	il := &InfluenceLog{}
	s.Input.GetInput(player, message, il)
	return il
}

func PlaceInfluence(s *State, player Aff, il *InfluenceLog) {
	for _, c := range il.Countries {
		c.Inf[player] += 1
	}
}
