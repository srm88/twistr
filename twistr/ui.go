package twistr

import (
	"fmt"
)

func localInput(ui UI, inp interface{}, message string, choices ...string) {
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
	inputStr := Solicit(ui, message, choices)
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

func Solicit(ui UI, message string, choices []string) (reply string) {
	buf := bytes.NewBufferString(strings.TrimRight(message, "\n"))
	if len(choices) > 0 {
		fmt.Fprintf(buf, " [ %s ]", strings.Join(choices, " "))
	}
	ui.Message(buf.String())
	return ui.Input()
}

type UI interface {
	Input() (string, error)
	Message(message string) error
	Redraw(*State) error
}
