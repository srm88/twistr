package main

import (
	"github.com/srm88/twistr/twistr"
)

// Temp:
func main() {
	ui := twistr.MakeTerminalUI()
	state := twistr.NewState(ui)
	cpl := &twistr.CardPlayLog{}
	twistr.GetInput(ui, twistr.USA, "player card ops|event", cpl)
	twistr.PlayCard(state, cpl)
}
