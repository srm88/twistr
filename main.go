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
	il, err := twistr.SelectNInfluenceCheck(state, twistr.Sov, "Place 6 influence in East Europe",
		6, twistr.InRegion(twistr.EastEurope))
	if err != nil {
		fmt.Println("Oh no ya goofed: ", err.Error())
	} else {
		fmt.Println("Influence successful")
		twistr.PlaceInfluence(state, twistr.Sov, il)
	}
}
