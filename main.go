package main

import (
	"github.com/srm88/twistr/twistr"
)

// Temp:
func main() {
	ui := twistr.MakeTerminalUI()
	state := twistr.NewState(ui)
	twistr.Start(state)
}
