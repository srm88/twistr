package twistr

type UI interface {
	Solicit(message string, choices []string) (reply string)
}

func SelectCountry(ui UI, choices ...string) Country {
	name := ui.Solicit("Which country?", choices)
	c := LookupCountry(name)
	return c
}

func LookupCountry(name string) Country {
	cid := lookupCountryId(name)
	return *countries[cid]
}

func lookupCountryId(name string) CountryId {
	for _, mapping := range countryIdLookup {
		if mapping.Name == name {
			return mapping.Id
		}
	}
	return -1
}
func LookupCard(name string) Card {
	cid := lookupCardId(name)
	for _, c := range EarlyWar {
		if c.Id == cid {
			return c
		}
	}
	for _, c := range MidWar {
		if c.Id == cid {
			return c
		}
	}
	for _, c := range LateWar {
		if c.Id == cid {
			return c
		}
	}
	panic("Oh god")
}

func lookupCardId(name string) CardId {
	for _, mapping := range cardIdLookup {
		if mapping.Name == name {
			return mapping.Id
		}
	}
	return -1
}
