package main

import (
	"github.com/srm88/twistr/twistr"
)

// Temp:
func main() {
	ui := twistr.MakeTerminalUI()
	state := twistr.NewState(ui)
	// Hacky
	state.Deck.Push(twistr.EarlyWar...)
	state.Deck.Reorder(state.Deck.Shuffle())
	twistr.Deal(state)
	twistr.Action(state)
}
