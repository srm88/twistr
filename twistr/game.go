package twistr

import (
	"fmt"
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
	if roll <= box.MaxRoll {
		box.Enter(s, player)
		s.Message(player, "Space race attempt succeeded.")
	} else {
		s.Message(player, "Space race attempt failed.")
	}
	s.Discard.Push(card)
}

func SelectSpaceRoll(s *State) *SpaceLog {
	// XXX: replay-log
	return &SpaceLog{Roll: Roll()}
}

func PlayOps(s *State, player Aff, card Card) {
	if card.Aff == player.Opp() {
		if player == SelectFirst(s, player).First {
			ConductOps(s, player, card)
			PlayEvent(s, player.Opp(), card)
		} else {
			PlayEvent(s, player.Opp(), card)
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
	// XXX: not positive this belongs here
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
	GetInput(s, player, fl, "Who goes first", USA.String(), SOV.String())
	return fl
}
