package twistr

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
