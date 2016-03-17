package main

import (
	"fmt"
	"github.com/srm88/twistr/twistr"
)

func getThing(ui twistr.UI) (t interface{}, err error) {
	switch ui.Solicit(twistr.US, ">", []string{"card", "country"}) {
	case "card":
		t, err = twistr.SelectCard(twistr.US, ui, "wargames", "socialistgovernments", "etc")
	case "country":
		t, err = twistr.SelectCountry(twistr.US, ui, "skorea", "uk", "etc")
	}
	return
}

// Temp:
func main() {
	ui := twistr.MakeTerminalUI()
	for {
		switch ui.Solicit(twistr.US, ">", []string{"lookup", "command", "quit"}) {
		case "lookup":
			t, err := getThing(ui)
			if err != nil {
				fmt.Println("Bad: ", err.Error())
				continue
			}
			fmt.Println(t)
		case "command":
			switch ui.Solicit(twistr.US, ">", []string{"coup", "realign", "influence"}) {
			case "influence":
				var c twistr.InfluenceCommand
				cmd := ui.Solicit(twistr.US, ">", []string{"[ country ... ]"})
				err := twistr.Unmarshal(cmd, &c)
				if err != nil {
					fmt.Println("Error: ", err.Error())
				}
				fmt.Println("Parsed: ", c)
			case "coup":
				var c twistr.CoupCommand
				cmd := ui.Solicit(twistr.US, ">", []string{"[country]", "[roll]"})
				err := twistr.Unmarshal(cmd, &c)
				if err != nil {
					fmt.Println("Error: ", err.Error())
				}
				fmt.Println("Parsed: ", c)
			case "realign":
				var c twistr.RealignCommand
				cmd := ui.Solicit(twistr.US, ">", []string{"[country]", "[rollUS]", "[rollSov]"})
				err := twistr.Unmarshal(cmd, &c)
				if err != nil {
					fmt.Println("Error: ", err.Error())
				}
				fmt.Println("Parsed: ", c)
			}
		case "quit":
			return
		default:
			fmt.Println("What?")
		}
	}
}
