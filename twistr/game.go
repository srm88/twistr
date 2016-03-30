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
	cs := SelectInfluenceForce(s, SOV, func() ([]*Country, error) {
		return SelectNInfluenceCheck(s, SOV,
			"6 influence in East Europe", 6,
			InRegion(EastEurope))
	})
	PlaceInfluence(s, SOV, cs)
	ShowHand(s, USA, USA)
	// US chooses 7 influence in W europe
	csUSA := SelectInfluenceForce(s, USA, func() ([]*Country, error) {
		return SelectNInfluenceCheck(s, USA,
			"7 influence in West Europe", 7,
			InRegion(WestEurope))
	})
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

// Return whether the card is an acceptible choice.
type cardFilter func(Card) bool

func ExceedsOps(minOps int) cardFilter {
	return func(c Card) bool {
		return c.Ops > minOps
	}
}

func passesFilters(c Card, filters []cardFilter) bool {
	for _, filter := range filters {
		if !filter(c) {
			return false
		}
	}
	return true
}

func SelectCard(s *State, player Aff, filters ...cardFilter) (c Card) {
	canPlayChina := s.ChinaCardPlayer == player && s.ChinaCardFaceUp
	choices := []string{}
	for _, c := range s.Hands[player].Cards {
		if !passesFilters(c, filters) {
			continue
		}
		choices = append(choices, c.Ref())
	}

	if canPlayChina && passesFilters(Cards[TheChinaCard], filters) {
		choices = append(choices, Cards[TheChinaCard].Ref())
	}
	GetOrLog(s, player, &c, "Choose a card", choices...)
	return
}

func SelectDiscarded(s *State, player Aff, filters ...cardFilter) (c Card) {
	choices := []string{}
	for _, c := range s.Discard.Cards {
		if !passesFilters(c, filters) {
			continue
		}
		choices = append(choices, c.Ref())
	}
	GetOrLog(s, player, &c, "Choose a discarded card", choices...)
	return
}

func SelectChoice(s *State, player Aff, message string, choices ...string) (choice string) {
	GetOrLog(s, player, &choice, message, choices...)
	return
}

func GetOrLog(s *State, player Aff, thing interface{}, message string, choices ...string) {
	if s.Aof.ReadInto(thing) {
		return
	}
	GetInput(s, player, thing, message, choices...)
	s.Aof.Log(thing)
}

func SelectRandomCard(s *State, player Aff) Card {
	n := rng.Intn(len(s.Hands[player].Cards))
	return s.Hands[player].Cards[n]
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
	roll := SelectRoll(s)
	MessageBoth(s, fmt.Sprintf("%s plays %s for the space race.", player, card))
	if roll <= box.MaxRoll {
		box.Enter(s, player)
		MessageBoth(s, fmt.Sprintf("%s rolls %d. Space race attempt success!", player, roll))
	} else {
		MessageBoth(s, fmt.Sprintf("%s rolls %d. Space race attempt fails!", player, roll))
	}
	s.Discard.Push(card)
}

func SelectRoll(s *State) int {
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

func ConductOps(s *State, player Aff, card Card, exclude ...OpsKind) {
	switch SelectOps(s, player, card, exclude...) {
	case COUP:
		MessageBoth(s, "coup not implemented")
	case REALIGN:
		MessageBoth(s, "realign not implemented")
	case INFLUENCE:
		MessageBoth(s, "influence not implemented")
	}
}

func DoFreeCoup(s *State, player Aff, card Card, allowedTargets []CountryId) bool {
	targets := []CountryId{}
	for _, t := range allowedTargets {
		if canCoup(s, player, s.Countries[t]) {
			targets = append(targets, t)
		}
	}
	if len(targets) == 0 {
		// Awkward
		return false
	}
	target := SelectCountry(s, player, "Free coup where?", targets...)
	roll := SelectRoll(s)
	ops := card.Ops + opsMod(s, player, card, []*Country{target})
	return coup(s, player, ops, roll, target, true)
}

func PlayEvent(s *State, player Aff, card Card) {
	MessageBoth(s, fmt.Sprintf("%s implements %s", player, card))
	prevented := card.Prevented(s)
	if !prevented {
		card.Impl(s, player)
	}
	switch {
	case !prevented && card.Star:
		s.Removed.Push(card)
		MessageBoth(s, fmt.Sprintf("%s removed", card))
	default:
		s.Discard.Push(card)
		MessageBoth(s, fmt.Sprintf("%s to discard", card))
	}
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
		choices = append(choices, OPS.Ref())
	}
	if canEvent {
		choices = append(choices, EVENT.Ref())
	}
	if canSpace {
		choices = append(choices, SPACE.Ref())
	}
	GetOrLog(s, player, &pk, fmt.Sprintf("Playing %s", card.Name), choices...)
	return
}

// We only support excluding 0 or 1 opskind. Excluding more than one means
// there is only one remaining kind, so there isn't a choice.
func SelectOps(s *State, player Aff, card Card, exclude ...OpsKind) (o OpsKind) {
	var message string
	if card.Id == FreeOps {
		message = fmt.Sprintf("Playing a %d ops card", card.Ops)
	} else {
		message = fmt.Sprintf("Playing %s for ops", card.Name)
	}
	var choices []string
	switch {
	case len(exclude) == 0:
		choices = []string{COUP.Ref(), REALIGN.Ref(), INFLUENCE.Ref()}
	case exclude[0] == COUP:
		choices = []string{REALIGN.Ref(), INFLUENCE.Ref()}
	case exclude[0] == REALIGN:
		choices = []string{COUP.Ref(), INFLUENCE.Ref()}
	case exclude[0] == INFLUENCE:
		choices = []string{COUP.Ref(), REALIGN.Ref()}
	}
	GetOrLog(s, player, &o, message, choices...)
	return
}

func SelectFirst(s *State, player Aff) (first Aff) {
	GetOrLog(s, player, &first, "Who will play first",
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

func ControlledBy(aff Aff) countryCheck {
	return func(c *Country) error {
		if c.Controlled() != aff {
			return fmt.Errorf("%s is not %s-controlled", c, aff)
		}
		return nil
	}
}

func NotControlledBy(aff Aff) countryCheck {
	return func(c *Country) error {
		if c.Controlled() == aff {
			return fmt.Errorf("%s is %s-controlled", c, aff)
		}
		return nil
	}
}

func MaxPerCountry(n int) countryCheck {
	counts := make(map[CountryId]int)
	return func(c *Country) error {
		counts[c.Id] += 1
		if counts[c.Id] > n {
			return fmt.Errorf("Too much in %s", c.Name)
		}
		return nil
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

// Repeat selectFn until the user's input is acceptible.
// This should be reconsidered once we support log-replay and log-writing.
func SelectInfluenceForce(s *State, player Aff, selectFn func() ([]*Country, error)) []*Country {
	var cs []*Country
	var err error
	cs, err = selectFn()
	for err != nil {
		s.Message(player, err.Error())
		cs, err = selectFn()
	}
	return cs
}

func SelectInfluence(s *State, player Aff, message string) (cs []*Country) {
	GetOrLog(s, player, &cs, message)
	return
}

func SelectCountry(s *State, player Aff, message string, countries ...CountryId) (c *Country) {
	choices := make([]string, len(countries))
	for i, cn := range countries {
		choices[i] = s.Countries[cn].Ref()
	}
	GetOrLog(s, player, &c, message, choices...)
	return
}

func PlaceInfluence(s *State, player Aff, cs []*Country) {
	for _, c := range cs {
		c.Inf[player] += 1
	}
}

func RemoveInfluence(s *State, player Aff, cs []*Country) {
	for _, c := range cs {
		c.Inf[player] -= 1
	}
}

// PseudoCard returns a card struct that can be used for events with text like
// "then the player may conduct ops as if they played an N ops card"
func PseudoCard(ops int) Card {
	return Card{
		Id:   FreeOps,
		Name: "-",
		Aff:  NEU,
		Ops:  ops,
	}
}

func score(s *State, player Aff, region Region) {
	// XXX writeme (#17)
}
