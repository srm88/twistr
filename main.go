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
	cpi := &twistr.CardPlayInput{}
	input.GetInput(twistr.US, "player card ops|event", cpi)
	fmt.Println(cpi)
}
