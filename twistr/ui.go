package twistr

func GetInput(ui UI, player Aff, message string, inp interface{}) {
	inputStr := ui.Solicit(player, message, nil)
	err := Unmarshal(inputStr, inp)
	for err != nil {
		inputStr = ui.Solicit(player, err.Error()+"\nTry again?", nil)
		err = Unmarshal(inputStr, inp)
	}
}

type UI interface {
	Solicit(player Aff, message string, choices []string) (reply string)
	Message(player Aff, message string)
}
