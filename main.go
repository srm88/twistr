package main

import (
	"fmt"
	"github.com/srm88/twistr/twistr"
)

// Temp:
func main() {
	ui := twistr.MakeTerminalUI()
	state := twistr.NewState(ui)
	il, err := twistr.SelectNInfluenceCheck(state, twistr.SOV, "Place 6 influence in East Europe",
		6, twistr.InRegion(twistr.EastEurope))
	if err != nil {
		fmt.Println("Oh no ya goofed: ", err.Error())
	} else {
		fmt.Println("Influence successful")
		twistr.PlaceInfluence(state, twistr.SOV, il)
	}
}
