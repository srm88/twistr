package twistr

// Game-running functions.
// Each function should represent a state in the game.

// WIP
func PlayCard(s *State, c *CardPlayLog) {
	switch {
	case c.Kind == SPACE:
		next := &SpaceLog{}
		GetInput(s, c.Player, "Space roll", next)
	case c.Kind == OPS && c.Card.Aff == c.Player.Opp():
		// Solicit who goes first
		next := &OpponentOpsLog{}
		GetInput(s, c.Player, "Who's next", next)
	case c.Kind == OPS:
		// Solicit coup/influence/realign/space
		next := &OpsLog{}
		GetInput(s, c.Player, "What kinda ops", next)
	case c.Kind == EVENT:
		panic("Not ready!")
	default:
		panic("WUT R U DOIN")
	}
}
