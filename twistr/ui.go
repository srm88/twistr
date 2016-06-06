package twistr

import (
	"fmt"
)

func GetInput(ui UI, player Aff, inp interface{}, message string, choices ...string) {
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
	inputStr := ui.Solicit(player, message, choices)
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
