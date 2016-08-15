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
	// XXX: handle running out of cards and shuffling in discards mid-deal
	handSize := s.ActionsPerTurn() + 2
	usDraw := s.Deck.Draw(handSize - len(s.Hands[USA].Cards))
	s.Hands[USA].Push(usDraw...)
	ShowHand(s, USA, USA)
	sovDraw := s.Deck.Draw(handSize - len(s.Hands[SOV].Cards))
	s.Hands[SOV].Push(sovDraw...)
	ShowHand(s, SOV, SOV)
}

func Start(s *State) {
	// Early war cards into the draw deck
	s.Deck.Push(EarlyWar...)
	cards := SelectShuffle(s, s.Deck)
	s.Deck.Reorder(cards)
	s.Commit()
	Deal(s)
	s.Redraw(s.Game)
	// SOV chooses 6 influence in E europe
	cs := SelectInfluenceForce(s, SOV, func() ([]*Country, error) {
		return SelectExactlyNInfluence(s, SOV,
			"6 influence in East Europe", 6,
			InRegion(EastEurope))
	})
	PlaceInfluence(s, SOV, cs)
	s.Commit()
	// US chooses 7 influence in W europe
	csUSA := SelectInfluenceForce(s, USA, func() ([]*Country, error) {
		return SelectExactlyNInfluence(s, USA,
			"7 influence in West Europe", 7,
			InRegion(WestEurope))
	})
	PlaceInfluence(s, USA, csUSA)
	s.Commit()
	// Temporary
	for s.Turn = 1; s.Turn <= 10; s.Turn++ {
		Turn(s)
		EndTurn(s)
	}
}

func ShowHand(s *State, whose, to Aff, showChina ...bool) {
	cs := []Card{}
	for _, c := range s.Hands[whose].Cards {
		cs = append(cs, c)
	}
	if len(showChina) > 0 && showChina[0] && s.ChinaCardPlayer == whose && s.ChinaCardFaceUp {
		cs = append(cs, Cards[TheChinaCard])
	}
	s.Mode = NewCardMode(cs)
	s.Redraw(s.Game)
}

func ShowDiscard(s *State, to Aff) {
	s.Mode = NewCardMode(s.Discard.Cards)
	s.Redraw(s.Game)
}

func ShowCard(s *State, c Card, to Aff) {
	s.Mode = NewCardMode([]Card{c})
	s.Redraw(s.Game)
}

func SelectShuffle(s *State, d *Deck) (cardOrder []Card) {
	// Duplicates what GetOrLog does. It doesn't make sense to reuse GetOrLog
	// because this will never ask for user input.
	if s.ReadInto(&cardOrder) {
		return
	}
	cardOrder = d.Shuffle()
	s.Log(&cardOrder)
	return
}

// Return whether the card is an acceptable choice.
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

func SelectCard(s *State, player Aff, filters ...cardFilter) Card {
	canPlayChina := s.ChinaCardPlayer == player && s.ChinaCardFaceUp
	return selectCardFrom(s, player, s.Hands[player].Cards, canPlayChina, filters...)
}

func hasInHand(s *State, player Aff, filters ...cardFilter) bool {
	for _, c := range s.Hands[player].Cards {
		if passesFilters(c, filters) {
			return true
		}
	}
	return false
}

func SelectDiscarded(s *State, player Aff, filters ...cardFilter) Card {
	return selectCardFrom(s, player, s.Discard.Cards, false, filters...)
}

func selectCardFrom(s *State, player Aff, from []Card, includeChina bool, filters ...cardFilter) (card Card) {
	choices := []string{}
	for _, c := range from {
		if !passesFilters(c, filters) {
			continue
		}
		choices = append(choices, c.Ref())
	}
	if includeChina && passesFilters(Cards[TheChinaCard], filters) {
		choices = append(choices, Cards[TheChinaCard].Ref())
	}
	GetOrLog(s, player, &card, "Choose a card", choices...)
	return
}

func SelectSomeCards(s *State, player Aff, message string, cards []Card) (selected []Card) {
	cardnames := []string{}
	cardSet := make(map[CardId]bool)
	for _, c := range cards {
		cardnames = append(cardnames, c.Ref())
		cardSet[c.Id] = true
	}
	message = fmt.Sprintf("%s: %s", message, strings.Join(cardnames, ", "))
	prefix := ""
retry:
	GetOrLog(s, player, &selected, prefix+message)
	for _, c := range selected {
		if !cardSet[c.Id] {
			prefix = "Invalid choice. "
			goto retry
		}
	}
	return
}

func SelectChoice(s *State, player Aff, message string, choices ...string) (choice string) {
	GetOrLog(s, player, &choice, message, choices...)
	return
}

func GetOrLog(s *State, player Aff, thing interface{}, message string, choices ...string) {
	if s.ReadInto(thing) {
		return
	}
	GetInput(s, player, thing, message, choices...)
	s.Log(thing)
}

func SelectRandomCard(s *State, player Aff) (card Card) {
	if s.ReadInto(&card) {
		return
	}
	n := rng.Intn(len(s.Hands[player].Cards))
	card = s.Hands[player].Cards[n]
	s.Log(&card)
	return
}

func actionsThisTurn(s *State, player Aff) int {
	_, active := s.TurnEvents[NorthSeaOil]
	switch {
	case player == USA && active:
		return 8
	case s.SREvents[ExtraAR] == player:
		return 8
	default:
		return s.ActionsPerTurn()
	}
}

func outOfCards(s *State, player Aff) bool {
	switch {
	case s.ChinaCardPlayer == player && s.ChinaCardFaceUp:
		return false
	case len(s.Hands[player].Cards) > 0:
		return false
	default:
		return true
	}
}

func Turn(s *State) {
	MessageBoth(s, fmt.Sprintf("== Turn %d", s.Turn))
	if s.Turn > 1 {
		Deal(s)
	}
	MessageBoth(s, "= Headline Phase")
	Headline(s)
	usaCap := actionsThisTurn(s, USA)
	sovCap := actionsThisTurn(s, SOV)
	usaDone, sovDone := false, false
	for {
		sovDone = s.AR > sovCap || outOfCards(s, SOV)
		usaDone = s.AR > usaCap || outOfCards(s, USA)
		if sovDone && usaDone {
			return
		}
		if !sovDone {
			MessageBoth(s, fmt.Sprintf("= %s AR %d.", SOV, s.AR))
			s.Phasing = SOV
			Action(s)
		}
		if !usaDone {
			MessageBoth(s, fmt.Sprintf("= %s AR %d.", USA, s.AR))
			s.Phasing = USA
			Action(s)
		}
		s.AR++
	}
}

func awardMilOpsVPs(s *State) {
	// Examples:
	// defcon 5, usa 3, sov 2 => usa scores 1
	// defcon 2, usa 0, sov 2 => sov scores 2
	// defcon 3, usa 5, sov 4 => nobody scores
	// defcon 3, usa 5, sov 0 => usa scores 3
	usaShy := Max(s.Defcon-s.MilOps[USA], 0)
	sovShy := Max(s.Defcon-s.MilOps[SOV], 0)
	switch {
	case usaShy == 0 && sovShy == 0:
		return
	case usaShy > sovShy:
		MessageBoth(s, fmt.Sprintf("%s loses %d VP for not meeting required military operations.", USA, usaShy-sovShy))
		s.GainVP(SOV, usaShy-sovShy)
	case sovShy > usaShy:
		MessageBoth(s, fmt.Sprintf("%s loses %d VP for not meeting required military operations.", SOV, sovShy-usaShy))
		s.GainVP(USA, sovShy-usaShy)
	}
}

// func discardHeldCard performs the space-race ability to discard 1 held card
// if either player has earned this ability.
func discardHeldCard(s *State, player Aff) {
	player, ok := s.SREvents[DiscardHeld]
	if !ok {
		return
	}
	if len(s.Hands[player].Cards) == 0 {
		return
	}
	if SelectChoice(s, player, "Discard one held card?", "yes", "no") != "yes" {
		return
	}
	card := SelectCard(s, player, CardBlacklist(TheChinaCard))
	s.Hands[player].Remove(card)
	s.Discard.Push(card)
}

func EndTurn(s *State) {
	// End turn: milops, defcon, china card, AR reset
	awardMilOpsVPs(s)
	s.MilOps[USA] = 0
	s.MilOps[SOV] = 0
	s.ImproveDefcon(1)
	if !s.ChinaCardFaceUp {
		MessageBoth(s, fmt.Sprintf("The china card is now face up for %s.", s.ChinaCardPlayer))
		s.ChinaCardFaceUp = true
	}
	s.AR = 1
	s.TurnEvents = make(map[CardId]Aff)
	s.ChernobylRegion = Region{}
}

func Action(s *State) {
	card := SelectCard(s, s.Phasing)
	PlayCard(s, s.Phasing, card)
	s.Commit()
}

func PlayCard(s *State, player Aff, card Card) {
	// Safe to remove a card that isn't actually in the hand
	s.Hands[player].Remove(card)
	switch SelectPlay(s, player, card) {
	case SPACE:
		PlaySpace(s, player, card)
	case OPS:
		PlayOps(s, player, card)
	case EVENT:
		PlayEvent(s, player, card)
	}
	if card.Id == TheChinaCard {
		MessageBoth(s, fmt.Sprintf("%s receives the China Card, face down.", player.Opp()))
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
	s.SpaceAttempts[player] += 1
	// China card can be spaced, but Action will take care of moving it to the
	// opponent.
	if card.Id != TheChinaCard {
		s.Discard.Push(card)
	}
}

func SelectRoll(s *State) int {
	// XXX: replay-log
	return Roll()
}

func PlayOps(s *State, player Aff, card Card) {
	MessageBoth(s, fmt.Sprintf("%s plays %s for operations", player, card))
	opp := player.Opp()
	if card.Aff == opp {
		first := SelectFirst(s, player)
		s.Commit()
		if player == first {
			MessageBoth(s, fmt.Sprintf("%s will conduct operations first", player))
			ConductOps(s, player, card)
			s.Redraw(s.Game)
			PlayEvent(s, opp, card)
		} else {
			MessageBoth(s, fmt.Sprintf("%s will implement the event first", opp))
			PlayEvent(s, opp, card)
			s.Redraw(s.Game)
			ConductOps(s, player, card)
		}
	} else {
		ConductOps(s, player, card)
		if card.Id != TheChinaCard {
			s.Discard.Push(card)
		}
	}
}

func ConductOps(s *State, player Aff, card Card, kinds ...OpsKind) {
	switch SelectOps(s, player, card, kinds...) {
	case COUP:
		OpCoup(s, player, card.Ops)
	case REALIGN:
		OpRealign(s, player, card.Ops)
	case INFLUENCE:
		OpInfluence(s, player, card.Ops)
	}
}

func OpRealign(s *State, player Aff, ops int) {
	// XXX; needs opsMod
	for i := 0; i < ops; i++ {
		target := SelectCountry(s, player, "Realign where?")
		for !canRealign(s, player, target, false) {
			target = SelectCountry(s, player, "Oh no you goofed. Realign where?")
		}
		rollUSA := SelectRoll(s)
		rollSOV := SelectRoll(s)
		realign(s, target, rollUSA, rollSOV)
	}
}

func OpCoup(s *State, player Aff, ops int) {
	target := SelectCountry(s, player, "Coup where?")
	for !canCoup(s, player, target, false) {
		target = SelectCountry(s, player, "Oh no you goofed. Coup where?")
	}

	roll := SelectRoll(s)
	ops += opsMod(s, player, []*Country{target})
	coup(s, player, ops, roll, target, false)
}

func DoFreeCoup(s *State, player Aff, card Card, allowedTargets []CountryId) bool {
	targets := []CountryId{}
	for _, t := range allowedTargets {
		if canCoup(s, player, s.Countries[t], true) {
			targets = append(targets, t)
		}
	}
	if len(targets) == 0 {
		// Awkward
		return false
	}
	target := SelectCountry(s, player, "Free coup where?", targets...)
	roll := SelectRoll(s)
	ops := card.Ops + opsMod(s, player, []*Country{target})
	return coup(s, player, ops, roll, target, true)
}

func OpInfluence(s *State, player Aff, ops int) {
	cs := SelectInfluenceForce(s, player, func() ([]*Country, error) {
		return SelectInfluenceOps(s, player, ops)
	})
	PlaceInfluence(s, player, cs)
}

func PlayEvent(s *State, player Aff, card Card) {
	prevented := card.Prevented(s.Game)
	if !prevented {
		// A soviet or US event is *always* played by that player, no matter
		// who causes the event to be played.
		switch card.Aff {
		case USA, SOV:
			MessageBoth(s, fmt.Sprintf("%s implements %s", card.Aff, card))
			card.Impl(s, card.Aff)
		default:
			MessageBoth(s, fmt.Sprintf("%s implements %s", player, card))
			card.Impl(s, player)
		}
	}
	switch {
	case card.Id == MissileEnvy:
		s.Hands[player.Opp()].Push(Cards[MissileEnvy])
		MessageBoth(s, fmt.Sprintf("%s to %s hand", card, player.Opp()))
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
	case card.Prevented(s.Game):
		canEvent = false
	}
	ops := card.Ops + opsMod(s, player, nil)
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

// Caller can pass in an optional whitelist of acceptable kinds.
func SelectOps(s *State, player Aff, card Card, kinds ...OpsKind) (o OpsKind) {
	var message string
	if card.Id == FreeOps {
		message = fmt.Sprintf("Playing a %d ops card", card.Ops)
	} else {
		message = fmt.Sprintf("Playing %s for ops", card.Name)
	}
	var choices []string
	if len(kinds) == 0 {
		choices = []string{COUP.Ref(), REALIGN.Ref(), INFLUENCE.Ref()}
	} else {
		for _, k := range kinds {
			choices = append(choices, k.Ref())
		}
	}
	GetOrLog(s, player, &o, message, choices...)
	return
}

func SelectFirst(s *State, player Aff) (first Aff) {
	GetOrLog(s, player, &first, "Who will play first",
		USA.Ref(), SOV.Ref())
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

func InCountries(countries ...CountryId) countryCheck {
	return func(c *Country) error {
		for _, cid := range countries {
			if cid == c.Id {
				return nil
			}
		}
		return fmt.Errorf("%s not a valid choice", c.Name)
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

func NoInfluence(aff Aff) countryCheck {
	return func(c *Country) error {
		if c.Inf[aff] != 0 {
			return fmt.Errorf("%s has influence in %s", aff, c.Name)
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

func HasInfluence(aff Aff) countryCheck {
	return func(c *Country) error {
		if c.Inf[aff] == 0 {
			return fmt.Errorf("No %s influence in %s", aff, c.Name)
		}
		return nil
	}
}

func CanRemove(aff Aff) countryCheck {
	removed := make(map[CountryId]int)
	return func(c *Country) error {
		removed[c.Id] += 1
		if c.Inf[aff]-removed[c.Id] < 0 {
			return fmt.Errorf("Not enough %s influence in %s", aff, c.Name)
		}
		return nil
	}
}

// SelectNInfluence asks the player to choose a number of countries to receive
// or lose influence, and optional checks to perform on the chosen countries.
func SelectNInfluence(s *State, player Aff, message string, n int, checks ...countryCheck) ([]*Country, error) {
	return selectNInfluence(s, player, message, n, false, checks...)
}

func SelectExactlyNInfluence(s *State, player Aff, message string, n int, checks ...countryCheck) ([]*Country, error) {
	return selectNInfluence(s, player, message, n, true, checks...)
}

func selectNInfluence(s *State, player Aff, message string, n int, exactly bool, checks ...countryCheck) (cs []*Country, err error) {
	cs = SelectInfluence(s, player, message)
	switch {
	case exactly && len(cs) != n:
		err = fmt.Errorf("Select exactly %d", n)
		return
	case !exactly && len(cs) > n:
		err = fmt.Errorf("Too much. Select %d", n)
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

func SelectInfluenceOps(s *State, player Aff, ops int) (cs []*Country, err error) {
	// Compute ops
	ops += opsMod(s, player, cs)
	message := fmt.Sprintf("Place %d influence", ops)
	cs = SelectInfluence(s, player, message)
	for _, c := range cs {
		if !canReach(c, player) {
			return nil, fmt.Errorf("Cannot reach %s", c.Name)
		}
	}
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

func canReach(c *Country, player Aff) bool {
	if c.Inf[player] > 0 {
		return true
	}
	if c.AdjSuper == player {
		return true
	}
	for _, ad := range c.AdjCountries {
		if ad.Inf[player] > 0 {
			return true
		}
	}
	return false
}

// Repeat selectFn until the user's input is acceptible.
// This should be reconsidered once we support log-replay and log-writing.
func SelectInfluenceForce(s *State, player Aff, selectFn func() ([]*Country, error)) []*Country {
	var cs []*Country
	var err error
	cs, err = selectFn()
	for err != nil {
		s.UI.Message(player, err.Error())
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

func SelectRegion(s *State, player Aff, message string) (r Region) {
	choices := make([]string, len(regionIdLookup))
	i := 0
	for name := range regionIdLookup {
		choices[i] = name
		i++
	}
	GetInput(s, player, &r, message, choices...)
	return
}

func PlaceInfluence(s *State, player Aff, cs []*Country) {
	for _, c := range cs {
		c.Inf[player] += 1
	}
}

func RemoveInfluence(s *State, player Aff, cs []*Country) {
	for _, c := range cs {
		c.Inf[player] = Max(0, c.Inf[player]-1)
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
