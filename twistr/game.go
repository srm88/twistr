package twistr

// Game-running functions.
// Each function should represent a state in the game.

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
