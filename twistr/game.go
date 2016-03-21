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
	dsl := SelectShuffle(s.Deck)
	s.Deck.Reorder(dsl.Cards)
	// Deal out players' hands
	Deal(s)
	// China card handled in NewState
	// SOV chooses 6 influence in E europe
	ShowHand(s, SOV, SOV)
	il, err := SelectNInfluenceCheck(s, SOV, "6 influence in East Europe", 6,
		InRegion(EastEurope))
	for err != nil {
		il, err = SelectNInfluenceCheck(s, SOV, err.Error(), 6,
			InRegion(EastEurope))
	}
	PlaceInfluence(s, SOV, il)
	ShowHand(s, USA, USA)
	// US chooses 7 influence in W europe
	ilUSA, err := SelectNInfluenceCheck(s, USA, "7 influence in West Europe", 7,
		InRegion(WestEurope))
	for err != nil {
		ilUSA, err = SelectNInfluenceCheck(s, USA, err.Error(), 7,
			InRegion(WestEurope))
	}
	PlaceInfluence(s, USA, ilUSA)
}

func ShowHand(s *State, whose, to Aff) {
	s.Message(to, fmt.Sprintf("%s hand: %s\n", whose, strings.Join(s.Hands[whose].Names(), ", ")))
}

func SelectShuffle(d *Deck) *DeckShuffleLog {
	// XXX: replay-log
	return &DeckShuffleLog{d.Shuffle()}
}

func Turn(s *State) {
	MessageBoth(s, fmt.Sprintf("TURN %d", s.Turn))
	s.Phasing = SOV
	Action()
	s.Phasing = US
	Action()
	s.Turn++
}

func Action(s *State) {
	p := s.Phasing
	card := SelectCard(s, p).Card
	// Safe to remove a card that isn't actually in the hand
	s.Hands[p].Remove(card)
	switch SelectPlay(s, p, card).Kind {
	case SPACE:
		PlaySpace(s, p, card)
	case OPS:
		PlayOps(s, p, card)
	case EVENT:
		PlayEvent(s, p, card)
	}
}

func PlaySpace(s *State, player Aff, card Card) {
	box, _ := nextSRBox(s, player)
	roll := SelectSpaceRoll(s).Roll
	MessageBoth(s, fmt.Sprintf("%s plays %s for the space race.", player, card))
	if roll <= box.MaxRoll {
		box.Enter(s, player)
		MessageBoth(s, fmt.Sprintf("%s rolls %d. Space race attempt success!", player, roll))
	} else {
		MessageBoth(s, fmt.Sprintf("%s rolls %d. Space race attempt fails!", player, roll))
	}
	s.Discard.Push(card)
}

func SelectSpaceRoll(s *State) *SpaceLog {
	// XXX: replay-log
	return &SpaceLog{Roll: Roll()}
}

func PlayOps(s *State, player Aff, card Card) {
	MessageBoth(s, fmt.Sprintf("%s plays %s for operations", player, card))
	opp := player.Opp()
	if card.Aff == opp {
		if player == SelectFirst(s, player).First {
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
		s.Discard.Push(card)
	}
}

func ConductOps(s *State, player Aff, card Card) {
	switch SelectOps(s, player, card).Kind {
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
	if card.Star {
		s.Removed.Push(card)
	} else {
		s.Discard.Push(card)
	}
	panic("Not implemented")
}

func SelectCard(s *State, player Aff) *CardLog {
	canPlayChina := s.ChinaCardPlayer == player && s.ChinaCardFaceUp
	choices := make([]string, len(s.Hands[player].Cards))
	for i, c := range s.Hands[player].Cards {
		choices[i] = c.Name
	}
	if canPlayChina {
		choices = append(choices, Cards[TheChinaCard].Name)
	}
	cl := &CardLog{}
	GetInput(s, player, cl, "Choose a card", choices...)
	return cl
}

func SelectPlay(s *State, player Aff, card Card) *PlayLog {
	// XXX: not all cards can be played as either ops or event, e.g. china card
	pl := &PlayLog{}
	// XXX: calculate ops at this point, SPACE should or should not be a choice
	// SPACE should also be omitted if the player is at the end of the track.
	GetInput(s, player, pl, fmt.Sprintf("Playing %s", card.Name),
		SPACE.String(), OPS.String(), EVENT.String())
	return pl
}

func SelectOps(s *State, player Aff, card Card) *OpsLog {
	ol := &OpsLog{}
	GetInput(s, player, ol, fmt.Sprintf("Playing %s for ops", card.Name),
		COUP.String(), REALIGN.String(), INFLUENCE.String())
	return ol
}

func SelectFirst(s *State, player Aff) *FirstLog {
	fl := &FirstLog{}
	GetInput(s, player, fl, "Who will play first",
		USA.String(), SOV.String())
	return fl
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
		err = fmt.Errorf("Select %d influence", n)
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
	GetInput(s, player, il, message)
	return il
}

func PlaceInfluence(s *State, player Aff, il *InfluenceLog) {
	for _, c := range il.Countries {
		c.Inf[player] += 1
	}
}
