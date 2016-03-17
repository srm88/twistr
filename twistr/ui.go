package twistr

type Input interface {
	GetInput(player Aff, message string, inp interface{})
}

type HackInput struct {
	Ui UI
}

// XXX: won't work. need to get rolls. probably need to have a solicitation
// method per command, which the ReplayInput will always implement by filling
// from the next command in the log.
func (i HackInput) GetInput(player Aff, message string, inp interface{}) {
	inputStr := i.Ui.Solicit(player, message, nil)
	err := Unmarshal(inputStr, inp)
	for err != nil {
		inputStr = i.Ui.Solicit(player, err.Error()+"> Try again?", nil)
		err = Unmarshal(inputStr, inp)
	}
}

type UI interface {
	Solicit(player Aff, message string, choices []string) (reply string)
}

func SelectCountry(player Aff, ui UI, choices ...string) (*Country, error) {
	name := ui.Solicit(player, "Which country?", choices)
	return lookupCountry(name)
}

func SelectCard(player Aff, ui UI, choices ...string) (Card, error) {
	name := ui.Solicit(player, "Which card?", choices)
	return lookupCard(name)
}
