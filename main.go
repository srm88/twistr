package main

import (
	"fmt"
	"twistr"
)

// Temp:
func main() {
	ui := twistr.MakeTerminalUI()
	for {
		switch ui.Solicit(">", []string{"card", "country", "quit"}) {
		case "card":
			fmt.Println(twistr.LookupCard(ui.Solicit("Which?", []string{"wargames", "socialistgovernments", "etc"})))
		case "country":
			fmt.Println(twistr.SelectCountry(ui, "skorea", "uk", "etc"))
		case "quit":
			return
		default:
			fmt.Println("What?")
		}
	}
}
