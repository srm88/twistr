package main

import (
	"fmt"
	"github.com/srm88/twistr/twistr"
)

// Temp:
func main() {
	ui := twistr.MakeTerminalUI()
	input := twistr.HackInput{Ui: ui}
	state := twistr.NewState(input)
	twistr.Start(state)
	fmt.Println("Nice.")
}
