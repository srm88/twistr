package twistr

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

// WIP
func PlayCard(s *State, c *CardPlayLog) {
	switch {
	case c.Kind == SPACE:
		next := &SpaceLog{}
		GetInput(s, c.Player, next, "Space roll")
	case c.Kind == OPS && c.Card.Aff == c.Player.Opp():
		// Solicit who goes first
		next := &OpponentOpsLog{}
		GetInput(s, c.Player, next, "Who's next")
	case c.Kind == OPS:
		// Solicit coup/influence/realign/space
		next := &OpsLog{}
		GetInput(s, c.Player, next, "What kinda ops")
	case c.Kind == EVENT:
		panic("Not ready!")
	default:
		panic("WUT R U DOIN")
	}
}
