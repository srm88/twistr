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
func Deal(s *State) {
	hs := s.HandSize()
	usDraw := s.Deck.Draw(hs - len(s.Hands[USA].Cards))
	s.Hands[USA].Push(usDraw...)
	sovDraw := s.Deck.Draw(hs - len(s.Hands[SOV].Cards))
	s.Hands[SOV].Push(sovDraw...)
}

func Start(s *State) {
	// Early war cards into the draw deck
	s.Deck.Push(EarlyWar...)
	cards := SelectShuffle(s.Deck)
	s.Deck.Reorder(cards)
	// Deal out players' hands
	Deal(s)
	// SOV chooses 6 influence in E europe
	ShowHand(s, SOV, SOV)
	cs, err := SelectNInfluenceCheck(s, SOV, "6 influence in East Europe", 6,
		InRegion(EastEurope))
	for err != nil {
		cs, err = SelectNInfluenceCheck(s, SOV, err.Error(), 6,
			InRegion(EastEurope))
	}
	PlaceInfluence(s, SOV, cs)
	ShowHand(s, USA, USA)
	// US chooses 7 influence in W europe
	csUSA, err := SelectNInfluenceCheck(s, USA, "7 influence in West Europe", 7,
		InRegion(WestEurope))
	for err != nil {
		csUSA, err = SelectNInfluenceCheck(s, USA, err.Error(), 7,
			InRegion(WestEurope))
	}
	PlaceInfluence(s, USA, csUSA)
	// Temporary
	Turn(s)
}

func ShowHand(s *State, whose, to Aff) {
	s.Message(to, fmt.Sprintf("%s hand: %s\n", whose, strings.Join(s.Hands[whose].Names(), ", ")))
}

func SelectShuffle(d *Deck) []Card {
	// XXX: replay-log
	return d.Shuffle()
}

func SelectCard(s *State, player Aff) (c Card) {
	canPlayChina := s.ChinaCardPlayer == player && s.ChinaCardFaceUp
	choices := make([]string, len(s.Hands[player].Cards))
	for i, c := range s.Hands[player].Cards {
		choices[i] = c.Name
	}
	if canPlayChina {
		choices = append(choices, Cards[TheChinaCard].Name)
	}
	GetInput(s, player, &c, "Choose a card", choices...)
	return
}

func Turn(s *State) {
	// Stub: awaiting implementation in issue#13
	MessageBoth(s, fmt.Sprintf("TURN %d", s.Turn))
	s.Phasing = SOV
	Action(s)
	s.Phasing = USA
	Action(s)
	s.Turn++
}

func Action(s *State) {
	p := s.Phasing
	card := SelectCard(s, p)
	// Safe to remove a card that isn't actually in the hand
	s.Hands[p].Remove(card)
	switch SelectPlay(s, p, card) {
	case SPACE:
		PlaySpace(s, p, card)
	case OPS:
		PlayOps(s, p, card)
	case EVENT:
		PlayEvent(s, p, card)
	}
	if card.Id == TheChinaCard {
		s.ChinaCardPlayed()
	}
}

func PlaySpace(s *State, player Aff, card Card) {
	box, _ := nextSRBox(s, player)
	roll := SelectSpaceRoll(s)
	MessageBoth(s, fmt.Sprintf("%s plays %s for the space race.", player, card))
	if roll <= box.MaxRoll {
		box.Enter(s, player)
		MessageBoth(s, fmt.Sprintf("%s rolls %d. Space race attempt success!", player, roll))
	} else {
		MessageBoth(s, fmt.Sprintf("%s rolls %d. Space race attempt fails!", player, roll))
	}
	s.Discard.Push(card)
}

func SelectSpaceRoll(s *State) int {
	// XXX: replay-log
	return Roll()
}

func PlayOps(s *State, player Aff, card Card) {
	MessageBoth(s, fmt.Sprintf("%s plays %s for operations", player, card))
	opp := player.Opp()
	if card.Aff == opp {
		if player == SelectFirst(s, player) {
			MessageBoth(s, fmt.Sprintf("%s will conduct operations first", player))
			ConductOps(s, player, card)
			PlayEvent(s, opp, card)
		} else {
			MessageBoth(s, fmt.Sprintf("%s will implement the event first", opp))
			PlayEvent(s, opp, card)
			ConductOps(s, player, card)
		}
	} else {
		ConductOps(s, player, card)
		if card.Id != TheChinaCard {
			s.Discard.Push(card)
		}
	}
}

func ConductOps(s *State, player Aff, card Card) {
	switch SelectOps(s, player, card) {
	case COUP:
		panic("Not implemented")
	case REALIGN:
		panic("Not implemented")
	case INFLUENCE:
		panic("Not implemented")
	}
}

func PlayEvent(s *State, player Aff, card Card) {
	MessageBoth(s, fmt.Sprintf("%s implements %s", player, card))
	switch {
	case !card.Prevented(s) && card.Star:
		s.Removed.Push(card)
	default:
		s.Discard.Push(card)
	}
	panic("Not implemented")
}

func SelectPlay(s *State, player Aff, card Card) (pk PlayKind) {
	canEvent, canSpace := true, true
	// Scoring cards cannot be played for ops
	canOps := card.Ops > 0
	switch {
	case card.Id == TheChinaCard:
		canEvent = false
	case card.Aff == player.Opp():
		// It isn't clear from the rules that playing your opponent's card as
		// an event is forbidden, but it is always a strictly worse move than
		// playing it for ops, and the rules don't prevent you from flipping
		// the table either ...
		canEvent = false
	case card.Prevented(s):
		canEvent = false
	}
	ops := card.Ops + opsMod(s, player, card, nil)
	if !CanAdvance(s, player, ops) {
		canSpace = false
	}
	choices := []string{}
	if canOps {
		choices = append(choices, OPS.String())
	}
	if canEvent {
		choices = append(choices, EVENT.String())
	}
	if canSpace {
		choices = append(choices, SPACE.String())
	}
	GetInput(s, player, &pk, fmt.Sprintf("Playing %s", card.Name), choices...)
	return
}

func SelectOps(s *State, player Aff, card Card) (o OpsKind) {
	GetInput(s, player, &o, fmt.Sprintf("Playing %s for ops", card.Name),
		COUP.String(), REALIGN.String(), INFLUENCE.String())
	return
}

func SelectFirst(s *State, player Aff) (first Aff) {
	GetInput(s, player, &first, "Who will play first",
		USA.String(), SOV.String())
	return
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
func SelectNInfluenceCheck(s *State, player Aff, message string, n int, checks ...countryCheck) (cs []*Country, err error) {
	cs = SelectInfluence(s, player, message)
	if len(cs) != n {
		err = fmt.Errorf("Select %d influence", n)
		return
	}
	for _, placement := range cs {
		for _, check := range checks {
			if err = check(placement); err != nil {
				return
			}
		}
	}
	return
}

func SelectInfluenceOps(s *State, player Aff, card Card) (cs []*Country, err error) {
	message := "Place influence"
	cs = SelectInfluence(s, player, message)
	// Compute ops
	ops := card.Ops + opsMod(s, player, card, cs)
	// Compute cost. Copy each country so that we can update its influence
	// as we go. E.g. two ops are spent breaking control, then the next
	// influence place costs one op.
	cost := 0
	workingCountries := make(map[CountryId]Country)
	for _, c := range cs {
		workingCountries[c.Id] = *c
	}
	for _, c := range cs {
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

func SelectInfluence(s *State, player Aff, message string) (cs []*Country) {
	GetInput(s, player, &cs, message)
	return
}

func PlaceInfluence(s *State, player Aff, cs []*Country) {
	for _, c := range cs {
		c.Inf[player] += 1
	}
}
