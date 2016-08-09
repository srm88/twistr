package twistr

import (
	"fmt"
)

func GetInput(ui UI, inp interface{}, message string, choices ...string) {
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

func RemoteInput(ui UI, inp interface{}) {
	// XXX who knows
	var err error
	var inputStr string
	inputStr, err = ui.Input()
	if err != nil {
		panic(fmt.Sprintf("Failed getting remote input %s\n", err.Error()))
	}
	err = Unmarshal(inputStr, inp)
	if err != nil {
		panic(fmt.Sprintf("Wonky remote input %s\n", err.Error()))
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
	Commit(*State) error
}
