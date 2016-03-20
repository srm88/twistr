package twistr

type State struct {
	UI

	VP int

	Defcon int

	MilOps [2]int

	SpaceRace [2]int

	Turn int
	AR   int

	Countries map[CountryId]*Country

	Events map[CardId]Aff

	SREvents map[SpaceId]Aff

	Removed *Deck

	Discard *Deck

	Deck *Deck

	Hands [2]map[CardId]Card

	ChinaCardPlayer Aff
	ChinaCardFaceUp bool
}

func NewState(ui UI) *State {
	return &State{
		UI:              ui,
		VP:              0,
		Defcon:          5,
		Turn:            1,
		AR:              1,
		Countries:       Countries,
		Events:          make(map[CardId]Aff),
		Removed:         NewDeck(),
		Discard:         NewDeck(),
		Deck:            NewDeck(),
		Hands:           [2]map[CardId]Card{make(map[CardId]Card), make(map[CardId]Card)},
		ChinaCardPlayer: SOV,
		ChinaCardFaceUp: true,
		SREvents:        make(map[SpaceId]Aff),
	}
}

func (s *State) Era() Era {
	switch {
	case s.Turn < 4:
		return Early
	case s.Turn < 8:
		return Mid
	default:
		return Late
	}
}

func (s *State) IntoHand(player Aff, cards ...Card) {
	for _, card := range cards {
		s.Hands[player][card.Id] = card
	}
}

func (s *State) HandSize() int {
	if s.Era() == Early {
		return 8
	}
	return 9
}

func (s *State) Effect(which CardId, player ...Aff) bool {
	aff, ok := s.Events[which]
	return ok && (len(player) == 0 || player[0] == aff)
}

// Card management

func (s *State) DrawHand(player Aff, n int) {
	need := n - len(s.Hands[player])
	for _, card := range s.Deck.Draw(need) {
		s.Hands[player][card.Id] = card
	}
}

func (s *State) CardPlayed(player Aff, which CardId, star bool) {
	if which == TheChinaCard {
		s.ChinaCardPlayer = player.Opp()
		s.ChinaCardFaceUp = false
		return
	}
	card := s.Hands[player][which]
	delete(s.Hands[player], which)
	if star {
		s.Removed.Push(card)
	} else {
		s.Discard.Push(card)
	}
}

func (s *State) GainVP(player Aff, n int) {
	switch player {
	case USA:
		s.VP += n
	case SOV:
		s.VP -= n
	}
}
