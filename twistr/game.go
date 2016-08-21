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
	handSize := s.ActionsPerTurn() + 2
	needCard := func(player Aff) bool {
		return len(s.Hands[player].Cards) == handSize
	}
	drawIfNeeded := func(player Aff) {
		if needCard(player) {
			if len(s.Deck.Cards) == 0 {
				ShuffleInDiscard(s)
			}
			card := s.Deck.Draw(1)[0]
			s.Hands[player].Push(card)
		}
	}
	for needCard(USA) || needCard(SOV) {
		drawIfNeeded(USA)
		drawIfNeeded(SOV)
	}
	ShowHand(s, USA, USA)
	ShowHand(s, SOV, SOV)
}

func ShuffleInDiscard(s *State) {
	s.Transcribe("Deck empty. Shuffling in discard pile ...")
	ShuffleIn(s, s.Discard.Draw(len(s.Discard.Cards)))
}

func Start(s *State) {
	// Early war cards into the draw deck
	ShuffleIn(s, EarlyWar)
	Deal(s)
	s.Redraw(s.Game)
	// SOV chooses 6 influence in E europe
	SelectInfluenceExactly(s, SOV, "6 influence in East Europe",
		PlusInf(SOV, 1), 6, InRegion(EastEurope))
	s.Commit()
	// US chooses 7 influence in W europe
	SelectInfluenceExactly(s, USA, "7 influence in West Europe",
		PlusInf(USA, 1), 7, InRegion(WestEurope))

	s.Commit()
	for s.Turn = 1; s.Turn <= 10; s.Turn++ {
		switch s.Turn {
		case 4:
			s.Transcribe("Shuffling in Mid War.")
			ShuffleIn(s, MidWar)
		case 8:
			s.Transcribe("Shuffling in Late War.")
			ShuffleIn(s, LateWar)
		}
		Turn(s)
		EndTurn(s)
	}
}

func ShuffleIn(s *State, cards []Card) {
	s.Deck.Push(cards...)
	order := SelectShuffle(s, s.Deck)
	s.Deck.Reorder(order)
	s.Commit()
}

func ShowHand(s *State, whose, to Aff, showChina ...bool) {
	cs := []Card{}
	for _, c := range s.Hands[whose].Cards {
		cs = append(cs, c)
	}
	if len(showChina) > 0 && showChina[0] && s.ChinaCardPlayer == whose && s.ChinaCardFaceUp {
		cs = append(cs, Cards[TheChinaCard])
	}
	s.Enter(NewCardMode(cs))
	s.Redraw(s.Game)
}

func ShowDiscard(s *State, to Aff) {
	s.Enter(NewCardMode(s.Discard.Cards))
	s.Redraw(s.Game)
}

func ShowCard(s *State, c Card, to Aff) {
	s.Enter(NewCardMode([]Card{c}))
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

func ExceedsOps(minOps int, s *State, player Aff) cardFilter {
	return func(c Card) bool {
		return ComputeCardOps(s, player, c, nil) > minOps
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
	s.Transcribe(fmt.Sprintf("== Turn %d", s.Turn))
	if s.Turn > 1 {
		Deal(s)
	}
	s.Transcribe("= Headline Phase")
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
			s.Transcribe(fmt.Sprintf("= %s AR %d.", SOV, s.AR))
			s.Phasing = SOV
			Action(s)
		}
		if !usaDone {
			s.Transcribe(fmt.Sprintf("= %s AR %d.", USA, s.AR))
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
		s.Transcribe(fmt.Sprintf("%s loses %d VP for not meeting required military operations.", USA, usaShy-sovShy))
		s.GainVP(SOV, usaShy-sovShy)
	case sovShy > usaShy:
		s.Transcribe(fmt.Sprintf("%s loses %d VP for not meeting required military operations.", SOV, sovShy-usaShy))
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
		s.Transcribe(fmt.Sprintf("The china card is now face up for %s.", s.ChinaCardPlayer))
		s.ChinaCardFaceUp = true
	}
	s.AR = 1
	s.TurnEvents = make(map[CardId]Aff)
	s.ChernobylRegion = Region{}
}

func Action(s *State) {
	switch {
	// BearTrap/Quagmire precede Missile Envy
	case s.Effect(BearTrap, s.Phasing.Opp()):
		TryBearTrap(s)
	case s.Effect(Quagmire, s.Phasing.Opp()):
		TryQuagmire(s)
	// Player forced to play missile envy for ops
	case s.Effect(MissileEnvy, s.Phasing.Opp()):
		card := Cards[MissileEnvy]
		s.Hands[s.Phasing].Remove(card)
		PlayOps(s, s.Phasing, card)
		s.Cancel(MissileEnvy)
	default:
		card := SelectCard(s, s.Phasing)
		PlayCard(s, s.Phasing, card)
		s.Commit()
	}
}

func TryQuagmire(s *State) {
	tryQuagmireBearTrap(s, Quagmire)
}

func TryBearTrap(s *State) {
	tryQuagmireBearTrap(s, BearTrap)
}

func tryQuagmireBearTrap(s *State, event CardId) {
	s.Transcribe(fmt.Sprintf("%s is in %s.", s.Phasing, Cards[event]))
	enoughOps := ExceedsOps(1, s, s.Phasing)
	// Can't discard? Play only scoring cards.
	if !hasInHand(s, s.Phasing, enoughOps) {
		onlyScoring := func(c Card) bool {
			return c.Scoring()
		}
		if !hasInHand(s, s.Phasing, onlyScoring) {
			s.Transcribe(fmt.Sprintf("%s cannot escape the %s and has no scoring cards.", s.Phasing, Cards[event]))
			return
		}
		s.Transcribe(fmt.Sprintf("%s cannot escape the %s and can only play scoring cards.", s.Phasing, Cards[event]))
		card := SelectCard(s, s.Phasing, onlyScoring)
		PlayCard(s, s.Phasing, card)
		return
	}
	// select card and roll
	card := SelectCard(s, s.Phasing, CardBlacklist(TheChinaCard), enoughOps)
	s.Hands[s.Phasing].Remove(card)
	s.Discard.Push(card)
	roll := SelectRoll(s)
	switch roll {
	case 1, 2, 3, 4:
		s.Transcribe(fmt.Sprintf("%s is free of the %s.", s.Phasing, Cards[event]))
		s.Cancel(event)
	default:
		s.Transcribe(fmt.Sprintf("%s is still trapped in the %s.", s.Phasing, Cards[event]))
	}
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
		s.Transcribe(fmt.Sprintf("%s receives the China Card, face down.", player.Opp()))
		s.ChinaCardPlayed()
	}
}

func PlaySpace(s *State, player Aff, card Card) {
	box, _ := nextSRBox(s, player)
	roll := SelectRoll(s)
	s.Transcribe(fmt.Sprintf("%s plays %s for the space race.", player, card))
	if roll <= box.MaxRoll {
		box.Enter(s, player)
		s.Transcribe(fmt.Sprintf("%s rolls %d. Space race attempt success!", player, roll))
	} else {
		s.Transcribe(fmt.Sprintf("%s rolls %d. Space race attempt fails!", player, roll))
	}
	s.SpaceAttempts[player] += 1
	// China card can be spaced, but Action will take care of moving it to the
	// opponent.
	if card.Id != TheChinaCard {
		s.Discard.Push(card)
	}
}

func SelectRoll(s *State) (roll int) {
	if s.ReadInto(&roll) {
		return
	}
	roll = Roll()
	s.Log(roll)
	return
}

func PlayOps(s *State, player Aff, card Card) {
	s.Transcribe(fmt.Sprintf("%s plays %s for operations", player, card))
	opp := player.Opp()
	if card.Aff == opp {
		first := SelectFirst(s, player)
		s.Commit()
		if player == first {
			s.Transcribe(fmt.Sprintf("%s will conduct operations first", player))
			ConductOps(s, player, card)
			s.Redraw(s.Game)
			PlayEvent(s, opp, card)
		} else {
			s.Transcribe(fmt.Sprintf("%s will implement the event first", opp))
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
	conductOps(s, player, card, false, kinds)
}

func ConductOpsFree(s *State, player Aff, card Card, kinds ...OpsKind) {
	conductOps(s, player, card, true, kinds)
}

func conductOps(s *State, player Aff, card Card, free bool, kinds []OpsKind) {
	switch SelectOps(s, player, card, kinds...) {
	case COUP:
		OpCoup(s, player, card, free)
	case REALIGN:
		OpRealign(s, player, card, free)
	case INFLUENCE:
		OpInfluence(s, player, card)
	}
}

func OpRealign(s *State, player Aff, card Card, free bool) {
	// XXX Something here is breaking on replay NOPE It's selectroll
	selectInfluence(s, player, fmt.Sprintf("Realigns with %s (%d)", card.Name, ComputeCardOps(s, player, card, nil)),
		func(c *Country) {
			Realign(s, player, c)
		},
		OpsLimit(s, player, card), false,
		NormalCost,
		CanRealign(s, player, free))
}

func OpCoup(s *State, player Aff, card Card, free bool, checks ...countryCheck) (success bool) {
	var msg string
	if free {
		msg = fmt.Sprintf("Free coup with %s (%d)", card.Name, ComputeCardOps(s, player, card, nil))
	} else {
		msg = fmt.Sprintf("Coup with %s (%d)", card.Name, ComputeCardOps(s, player, card, nil))
	}
	selectInfluence(s, player, msg,
		func(c *Country) {
			success = Coup(s, player, card, c, free)
		},
		LimitN(1), true,
		NormalCost,
		append(checks, CanCoup(s, player, free))...)
	return
}

func OpInfluence(s *State, player Aff, card Card) {
	// XXX chernobyl, etc
	selectInfluence(s, player, fmt.Sprintf("Influence with %s (%d)", card.Name, ComputeCardOps(s, player, card, nil)),
		PlusInf(player, 1),
		OpsLimit(s, player, card), false,
		OpInfluenceCost(player),
		CanReach(s, player))
}

func PlayEvent(s *State, player Aff, card Card) {
	prevented := card.Prevented(s.Game)
	if !prevented {
		// A soviet or US event is *always* played by that player, no matter
		// who causes the event to be played.
		var implementer Aff
		switch card.Aff {
		case USA, SOV:
			implementer = card.Aff
		default:
			implementer = player
		}
		s.Transcribe(fmt.Sprintf("%s implements %s", implementer, card))
		card.Impl(s, implementer)
	} else {
		s.Transcribe(fmt.Sprintf("%s cannot be played as an event", card))
	}
	switch {
	case card.Id == MissileEnvy:
		s.Hands[player.Opp()].Push(Cards[MissileEnvy])
		s.Transcribe(fmt.Sprintf("%s to %s hand", card, player.Opp()))
	case !prevented && card.Star:
		s.Removed.Push(card)
		s.Transcribe(fmt.Sprintf("%s removed", card))
	default:
		s.Discard.Push(card)
		s.Transcribe(fmt.Sprintf("%s to discard", card))
	}
}

func SelectPlay(s *State, player Aff, card Card) (pk PlayKind) {
	if card.Scoring() {
		pk = EVENT
		return
	}
	canEvent, canSpace := true, true
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
	if !CanAdvance(s, player, ComputeCardOps(s, player, card, nil)) {
		canSpace = false
	}
	choices := []string{OPS.Ref()}
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
		message = fmt.Sprintf("Playing a %d ops card (%d)", card.Ops, ComputeCardOps(s, player, card, nil))
	} else {
		message = fmt.Sprintf("Playing %s for ops (%d)", card.Name, ComputeCardOps(s, player, card, nil))
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

func CanReach(s *State, player Aff) countryCheck {
	influenced := make(map[CountryId]bool)
	for cid, c := range s.Countries {
		if c.Inf[player] > 0 {
			influenced[cid] = true
		}
	}
	return func(c *Country) error {
		if influenced[c.Id] {
			return nil
		}
		if c.AdjSuper == player {
			return nil
		}
		for _, ad := range c.AdjCountries {
			if influenced[ad.Id] {
				return nil
			}
		}
		return fmt.Errorf("Cannot reach %s", c.Name)
	}
}

func SelectCountry(s *State, player Aff, message string, countries ...CountryId) (c *Country) {
	// XXX this doesn't permit use of country short codes
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

func Score(s *State, player Aff, region Region) {
	// XXX writeme (#17)
}
