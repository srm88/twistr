package twistr

import (
	"math/rand"
)

type State struct {
	VP int

	Defcon int

	MilOps [2]int

	SpaceRace [2]int

	Turn int
	AR   int

	Countries map[CountryId]*Country

	Events map[CardId]Aff

	Removed []*Card

	Discard []*Card

	Deck []*Card

	Hands [2]map[CardId]*Card

	ChinaCardPlayer Aff
	ChinaCardFaceUp bool
}

func (s *State) Effect(which CardId, player ...Aff) bool {
	aff, ok := s.Events[which]
	return ok && (len(player) == 0 || player[0] == aff)
}

// Card management
func (s *State) Shuffle() {
	deckLen := len(s.Deck)
	for i := 0; i < 2*deckLen; i++ {
		x := rand.Intn(deckLen)
		s.Deck[i], s.Deck[x] = s.Deck[x], s.Deck[i]
	}
}

func (s *State) ShuffleIn(cards []*Card) {
	s.Deck = append(s.Deck, cards...)
	s.Shuffle()
}

func (s *State) Draw(n int) (draws []*Card) {
	draws, s.Deck = s.Deck[:n], s.Deck[n:]
	return
}

func (s *State) DrawDiscards(n int) (draws []*Card) {
	draws, s.Discard = s.Discard[:n], s.Discard[n:]
	return
}

func (s *State) DrawHand(player Aff, n int) {
	need := n - len(s.Hands[player])
	for _, card := range s.Draw(need) {
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
		s.Removed = append(s.Removed, card)
	} else {
		s.Discard = append(s.Discard, card)
	}
}
