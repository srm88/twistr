package main

import (
	"fmt"
	"github.com/srm88/twistr/twistr"
)

// Temp:
func main() {
	ui := twistr.MakeTerminalUI()
	input := twistr.HackInput{Ui: ui}
	//state := twistr.NewState(input)
	cpl := &twistr.CardPlayLog{}
	input.GetInput(twistr.US, "player card ops|event", cpl)
	fmt.Println(cpl)
}
