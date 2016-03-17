package twistr

import (
	"bytes"
)

var (
	countries map[CountryId]*Country
	SEAsia    Region = Region{
		Countries: []CountryId{
			Burma, LaosCambodia, Vietnam, Thailand, Malaysia, Indonesia, Philippines,
		},
		Volatility: 2,
	}
)

// Temp:
func ByName(name string) *Country {
	for _, c := range countries {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func CountryNames(cs []*Country) string {
	var b bytes.Buffer
	for _, c := range cs {
		b.WriteString(c.Name)
		b.WriteString(" ")
	}
	return b.String()
}

func init() {
	countries = make(map[CountryId]*Country)
	for _, c := range countryTable {
		countries[c.Id] = &Country{
			Id:           c.Id,
			Name:         c.Name,
			Inf:          Influence{c.USInf, c.SovInf},
			Stability:    c.Stability,
			Battleground: c.Battleground,
			AdjSuper:     c.AdjSuper,
		}
	}
	for _, link := range countryLinks {
		foo := countries[link[0]]
		bar := countries[link[1]]
		foo.AdjCountries = append(foo.AdjCountries, bar)
		bar.AdjCountries = append(bar.AdjCountries, foo)
	}
}

var countryTable = []struct {
	Id           CountryId
	Name         string
	USInf        int
	SovInf       int
	Stability    int
	Battleground bool
	AdjSuper     Aff
}{
	{Mexico, "Mexico", 0, 0, 2, true, US},
	{Guatemala, "Guatemala", 0, 0, 1, false, Neu},
	{ElSalvador, "ElSalvador", 0, 0, 1, false, Neu},
	{Honduras, "Honduras", 0, 0, 2, false, Neu},
	{CostaRica, "CostaRica", 0, 0, 3, false, Neu},
	{Cuba, "Cuba", 0, 0, 3, true, US},
	{Nicaragua, "Nicaragua", 0, 0, 1, false, Neu},
	{Panama, "Panama", 1, 0, 2, true, Neu},
	{Haiti, "Haiti", 0, 0, 1, false, Neu},
	{DominicanRep, "DominicanRep", 0, 0, 1, false, Neu},
	{Ecuador, "Ecuador", 0, 0, 2, false, Neu},
	{Peru, "Peru", 0, 0, 2, false, Neu},
	{Colombia, "Colombia", 0, 0, 1, false, Neu},
	{Chile, "Chile", 0, 0, 3, true, Neu},
	{Venezuela, "Venezuela", 0, 0, 2, true, Neu},
	{Argentina, "Argentina", 0, 0, 2, true, Neu},
	{Bolivia, "Bolivia", 0, 0, 2, false, Neu},
	{Paraguay, "Paraguay", 0, 0, 2, false, Neu},
	{Uruguay, "Uruguay", 0, 0, 2, false, Neu},
	{Brazil, "Brazil", 0, 0, 2, true, Neu},
	{Canada, "Canada", 0, 0, 4, false, US},
	{UK, "UK", 5, 0, 5, false, Neu},
	{SpainPortugal, "SpainPortugal", 0, 0, 2, false, Neu},
	{France, "France", 0, 0, 3, true, Neu},
	{Benelux, "Benelux", 0, 0, 3, false, Neu},
	{Norway, "Norway", 0, 0, 4, false, Neu},
	{Denmark, "Denmark", 0, 0, 3, false, Neu},
	{WGermany, "WGermany", 0, 0, 4, true, Neu},
	{EGermany, "EGermany", 3, 0, 3, true, Neu},
	{Italy, "Italy", 0, 0, 2, true, Neu},
	{Austria, "Austria", 0, 0, 4, false, Neu},
	{Sweden, "Sweden", 0, 0, 4, false, Neu},
	{Czechoslovakia, "Czechoslovakia", 0, 0, 3, false, Neu},
	{Yugoslavia, "Yugoslavia", 0, 0, 3, false, Neu},
	{Poland, "Poland", 0, 0, 3, true, Sov},
	{Greece, "Greece", 0, 0, 2, false, Neu},
	{Hungary, "Hungary", 0, 0, 3, false, Neu},
	{Finland, "Finland", 1, 0, 4, false, Sov},
	{Romania, "Romania", 0, 0, 3, false, Sov},
	{Bulgaria, "Bulgaria", 0, 0, 3, false, Neu},
	{Turkey, "Turkey", 0, 0, 2, false, Neu},
	{Morocco, "Morocco", 0, 0, 3, false, Neu},
	{WestAfricanStates, "WestAfricanStates", 0, 0, 2, false, Neu},
	{IvoryCoast, "IvoryCoast", 0, 0, 2, false, Neu},
	{Algeria, "Algeria", 0, 0, 2, true, Neu},
	{SaharanStates, "SaharanStates", 0, 0, 1, false, Neu},
	{Nigeria, "Nigeria", 0, 0, 1, true, Neu},
	{Tunisia, "Tunisia", 0, 0, 2, false, Neu},
	{Cameroon, "Cameroon", 0, 0, 1, false, Neu},
	{Angola, "Angola", 0, 0, 1, true, Neu},
	{SouthAfrica, "SouthAfrica", 1, 0, 3, true, Neu},
	{Zaire, "Zaire", 0, 0, 1, true, Neu},
	{Botswana, "Botswana", 0, 0, 2, false, Neu},
	{Zimbabwe, "Zimbabwe", 0, 0, 1, false, Neu},
	{Sudan, "Sudan", 0, 0, 1, false, Neu},
	{Ethiopia, "Ethiopia", 0, 0, 1, false, Neu},
	{Kenya, "Kenya", 0, 0, 2, false, Neu},
	{SEAfricanStates, "SEAfricanStates", 0, 0, 1, false, Neu},
	{Somalia, "Somalia", 0, 0, 2, false, Neu},
	{Libya, "Libya", 0, 0, 2, true, Neu},
	{Egypt, "Egypt", 0, 0, 2, true, Neu},
	{Israel, "Israel", 1, 0, 4, true, Neu},
	{Lebanon, "Lebanon", 0, 0, 1, false, Neu},
	{Jordan, "Jordan", 0, 0, 2, false, Neu},
	{Syria, "Syria", 1, 0, 2, false, Neu},
	{Iraq, "Iraq", 1, 0, 3, true, Neu},
	{SaudiArabia, "SaudiArabia", 0, 0, 3, true, Neu},
	{GulfStates, "GulfStates", 0, 0, 3, false, Neu},
	{Iran, "Iran", 1, 0, 2, true, Neu},
	{Afghanistan, "Afghanistan", 0, 0, 2, false, Sov},
	{Pakistan, "Pakistan", 0, 0, 2, true, Neu},
	{India, "India", 0, 0, 3, true, Neu},
	{Burma, "Burma", 0, 0, 2, false, Neu},
	{Thailand, "Thailand", 0, 0, 2, true, Neu},
	{LaosCambodia, "LaosCambodia", 0, 0, 1, false, Neu},
	{Vietnam, "Vietnam", 0, 0, 1, false, Neu},
	{Malaysia, "Malaysia", 0, 0, 2, false, Neu},
	{Indonesia, "Indonesia", 0, 0, 1, false, Neu},
	{Australia, "Australia", 4, 0, 4, false, Neu},
	{Taiwan, "Taiwan", 0, 0, 3, false, Neu},
	{NKorea, "NKorea", 3, 0, 3, true, Sov},
	{SKorea, "SKorea", 1, 0, 3, true, Neu},
	{Philippines, "Philippines", 1, 0, 2, false, Neu},
	{Japan, "Japan", 1, 0, 4, true, US},
}

var countryLinks = [][2]CountryId{
	// Central America
	{Mexico, Guatemala},
	{Guatemala, ElSalvador},
	{Guatemala, Honduras},
	{ElSalvador, Honduras},
	{Honduras, CostaRica},
	{Honduras, Nicaragua},
	{CostaRica, Nicaragua},
	{Cuba, Nicaragua},
	{Cuba, Haiti},
	{Haiti, DominicanRep},
	{CostaRica, Panama},
	// South America
	{Panama, Colombia},
	{Colombia, Ecuador},
	{Ecuador, Peru},
	{Peru, Chile},
	{Peru, Bolivia},
	{Chile, Argentina},
	{Bolivia, Paraguay},
	{Paraguay, Argentina},
	{Paraguay, Uruguay},
	{Uruguay, Brazil},
	{Brazil, Venezuela},
	{Venezuela, Colombia},
	// Europe
	{Canada, UK},
	{UK, France},
	{UK, Norway},
	{SpainPortugal, France},
	{SpainPortugal, Italy},
	{France, Italy},
	{France, WGermany},
	{Benelux, WGermany},
	{Norway, Sweden},
	{Sweden, Denmark},
	{Sweden, Finland},
	{Denmark, WGermany},
	{WGermany, Austria},
	{Austria, Italy},
	{Austria, EGermany},
	{Austria, Hungary},
	{EGermany, WGermany},
	{EGermany, Poland},
	{EGermany, Czechoslovakia},
	{Poland, Czechoslovakia},
	{Czechoslovakia, Hungary},
	{Hungary, Yugoslavia},
	{Hungary, Romania},
	{Romania, Turkey},
	{Romania, Yugoslavia},
	{Yugoslavia, Italy},
	{Yugoslavia, Greece},
	{Greece, Italy},
	{Greece, Bulgaria},
	{Greece, Turkey},
	// Middle east
	{Turkey, Syria},
	{Syria, Lebanon},
	{Syria, Israel},
	{Lebanon, Jordan},
	{Lebanon, Israel},
	{Israel, Egypt},
	{Israel, Jordan},
	{Egypt, Libya},
	{Jordan, Iraq},
	{Jordan, SaudiArabia},
	{Iraq, SaudiArabia},
	{Iraq, GulfStates},
	{Iraq, Iran},
	{GulfStates, SaudiArabia},
	// Africa
	{Egypt, Sudan},
	{Libya, Tunisia},
	{Algeria, France},
	{Morocco, SpainPortugal},
	{Algeria, Tunisia},
	{Morocco, Algeria},
	{Morocco, WestAfricanStates},
	{WestAfricanStates, IvoryCoast},
	{IvoryCoast, Nigeria},
	{Nigeria, SaharanStates},
	{SaharanStates, Algeria},
	{Nigeria, Cameroon},
	{Cameroon, Zaire},
	{Zaire, Angola},
	{Zaire, Zimbabwe},
	{Angola, Botswana},
	{Angola, SouthAfrica},
	{Botswana, SouthAfrica},
	{Botswana, Zimbabwe},
	{Zimbabwe, SEAfricanStates},
	{SEAfricanStates, Kenya},
	{Kenya, Somalia},
	{Somalia, Ethiopia},
	{Ethiopia, Sudan},
	// Asia
	{Iran, Afghanistan},
	{Iran, Pakistan},
	{Pakistan, Afghanistan},
	{Pakistan, India},
	{India, Burma},
	{Burma, LaosCambodia},
	{LaosCambodia, Thailand},
	{LaosCambodia, Vietnam},
	{Vietnam, Thailand},
	{Thailand, Malaysia},
	{Malaysia, Australia},
	{Malaysia, Indonesia},
	{Indonesia, Philippines},
	{Philippines, Japan},
	{Japan, Taiwan},
	{Japan, SKorea},
	{Taiwan, SKorea},
	{SKorea, NKorea},
}
