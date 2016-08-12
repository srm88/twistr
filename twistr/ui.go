package twistr

import (
	"fmt"
)

func GetInput(g *Game, player Aff, inp interface{}, message string, choices ...string) {
	var err error
	validChoice := func(in string) bool {
		if len(choices) == 0 {
			return true
		}
		for _, choice := range choices {
			if choice == in {
				return true
			}
		}
		return false
	}
retry:
	inputStr := g.Solicit(player, message, choices)
	switch inputStr {
	case "canrewind":
		message = fmt.Sprintf("%v\n", g.CanRewind())
		goto retry
	case "barf":
		g.History.Dump()
		goto retry
	case "rewind":
		if !g.CanRewind() {
			message = "Cannot rewind."
			goto retry
		}
		g.Rewind()
		panic("Nope")
	}
	if len(choices) > 0 && !validChoice(inputStr) {
		err = fmt.Errorf("'%s' is not a valid choice", inputStr)
	} else {
		err = Unmarshal(inputStr, inp)
	}
	if err != nil {
		message = err.Error() + ". Try again?"
		goto retry
	}
}

type UI interface {
	Solicit(player Aff, message string, choices []string) (reply string)
	Message(player Aff, message string)
	Redraw(*State)
	Close() error
}
