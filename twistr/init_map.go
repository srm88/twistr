package twistr

var (
	Countries      map[CountryId]*Country
	CentralAmerica Region = Region{
		Name: "CentralAmerica",
		Countries: []CountryId{
			Mexico, Guatemala, ElSalvador, Honduras, CostaRica, Cuba, Nicaragua, Panama, Haiti, DominicanRep,
		},
	}
	SouthAmerica Region = Region{
		Name: "SouthAmerica",
		Countries: []CountryId{
			Ecuador, Peru, Colombia, Chile, Venezuela, Argentina, Bolivia, Paraguay, Uruguay, Brazil,
		},
	}
	Europe Region = Region{
		Name: "Europe",
		Countries: []CountryId{
			Canada, UK, SpainPortugal, France, Benelux, Norway, Denmark, WGermany, EGermany, Italy, Austria, Sweden, Czechoslovakia, Yugoslavia, Poland, Greece, Hungary, Finland, Romania, Bulgaria, Turkey,
		},
		Volatility: 4,
	}
	MiddleEast Region = Region{
		Name: "MiddleEast",
		Countries: []CountryId{
			Libya, Egypt, Israel, Lebanon, Jordan, Syria, Iraq, SaudiArabia, GulfStates, Iran,
		},
		Volatility: 2,
	}
	Africa Region = Region{
		Name: "Africa",
		Countries: []CountryId{
			Morocco, WestAfricanStates, IvoryCoast, Algeria, SaharanStates, Nigeria, Tunisia, Cameroon, Angola, SouthAfrica, Zaire, Botswana, Zimbabwe, Sudan, Ethiopia, Kenya, SEAfricanStates, Somalia,
		},
	}
	Asia Region = Region{
		Name: "Asia",
		Countries: []CountryId{
			Afghanistan, Pakistan, India, Burma, Thailand, LaosCambodia, Vietnam, Malaysia, Indonesia, Australia, Taiwan, NKorea, SKorea, Philippines, Japan,
		},
		Volatility: 3,
	}
	// Sub regions
	WestEurope Region = Region{
		Name: "WestEurope",
		Countries: []CountryId{
			Canada, UK, SpainPortugal, France, Benelux, Norway, Denmark, WGermany, Italy, Greece, Austria, Sweden, Finland, Turkey,
		},
	}
	EastEurope Region = Region{
		Name: "EastEurope",
		Countries: []CountryId{
			EGermany, Austria, Yugoslavia, Czechoslovakia, Poland, Hungary, Finland, Bulgaria, Romania,
		},
	}
	SoutheastAsia Region = Region{
		Name: "SoutheastAsia",
		Countries: []CountryId{
			Burma, LaosCambodia, Vietnam, Thailand, Malaysia, Indonesia, Philippines,
		},
	}
	Regions map[RegionId]Region = map[RegionId]Region{
		CAM: CentralAmerica,
		SAM: SouthAmerica,
		EUR: Europe,
		MDE: MiddleEast,
		AFR: Africa,
		ASI: Asia,
	}
)

func init() {
	resetCountries()
}

func resetCountries() {
	Countries = make(map[CountryId]*Country)
	for _, c := range countryTable {
		Countries[c.Id] = &Country{
			Id:           c.Id,
			Name:         c.Name,
			Inf:          Influence{c.USAInf, c.SOVInf},
			Stability:    c.Stability,
			Battleground: c.Battleground,
			AdjSuper:     c.AdjSuper,
			Region:       Regions[c.RegionId],
		}
	}
	for _, link := range countryLinks {
		foo := Countries[link[0]]
		bar := Countries[link[1]]
		foo.AdjCountries = append(foo.AdjCountries, bar)
		bar.AdjCountries = append(bar.AdjCountries, foo)
	}
}

var countryTable = []struct {
	Id           CountryId
	Name         string
	USAInf       int
	SOVInf       int
	Stability    int
	Battleground bool
	AdjSuper     Aff
	RegionId     RegionId
}{
	{Mexico, "Mexico", 0, 0, 2, true, USA, CAM},
	{Guatemala, "Guatemala", 0, 0, 1, false, NEU, CAM},
	{ElSalvador, "ElSalvador", 0, 0, 1, false, NEU, CAM},
	{Honduras, "Honduras", 0, 0, 2, false, NEU, CAM},
	{CostaRica, "CostaRica", 0, 0, 3, false, NEU, CAM},
	{Cuba, "Cuba", 0, 0, 3, true, USA, CAM},
	{Nicaragua, "Nicaragua", 0, 0, 1, false, NEU, CAM},
	{Panama, "Panama", 1, 0, 2, true, NEU, CAM},
	{Haiti, "Haiti", 0, 0, 1, false, NEU, CAM},
	{DominicanRep, "DominicanRep", 0, 0, 1, false, NEU, CAM},
	{Ecuador, "Ecuador", 0, 0, 2, false, NEU, SAM},
	{Peru, "Peru", 0, 0, 2, false, NEU, SAM},
	{Colombia, "Colombia", 0, 0, 1, false, NEU, SAM},
	{Chile, "Chile", 0, 0, 3, true, NEU, SAM},
	{Venezuela, "Venezuela", 0, 0, 2, true, NEU, SAM},
	{Argentina, "Argentina", 0, 0, 2, true, NEU, SAM},
	{Bolivia, "Bolivia", 0, 0, 2, false, NEU, SAM},
	{Paraguay, "Paraguay", 0, 0, 2, false, NEU, SAM},
	{Uruguay, "Uruguay", 0, 0, 2, false, NEU, SAM},
	{Brazil, "Brazil", 0, 0, 2, true, NEU, SAM},
	{Canada, "Canada", 2, 0, 4, false, USA, EUR},
	{UK, "UK", 5, 0, 5, false, NEU, EUR},
	{SpainPortugal, "SpainPortugal", 0, 0, 2, false, NEU, EUR},
	{France, "France", 0, 0, 3, true, NEU, EUR},
	{Benelux, "Benelux", 0, 0, 3, false, NEU, EUR},
	{Norway, "Norway", 0, 0, 4, false, NEU, EUR},
	{Denmark, "Denmark", 0, 0, 3, false, NEU, EUR},
	{WGermany, "WGermany", 0, 0, 4, true, NEU, EUR},
	{EGermany, "EGermany", 0, 3, 3, true, NEU, EUR},
	{Italy, "Italy", 0, 0, 2, true, NEU, EUR},
	{Austria, "Austria", 0, 0, 4, false, NEU, EUR},
	{Sweden, "Sweden", 0, 0, 4, false, NEU, EUR},
	{Czechoslovakia, "Czechoslovakia", 0, 0, 3, false, NEU, EUR},
	{Yugoslavia, "Yugoslavia", 0, 0, 3, false, NEU, EUR},
	{Poland, "Poland", 0, 0, 3, true, SOV, EUR},
	{Greece, "Greece", 0, 0, 2, false, NEU, EUR},
	{Hungary, "Hungary", 0, 0, 3, false, NEU, EUR},
	{Finland, "Finland", 0, 1, 4, false, SOV, EUR},
	{Romania, "Romania", 0, 0, 3, false, SOV, EUR},
	{Bulgaria, "Bulgaria", 0, 0, 3, false, NEU, EUR},
	{Turkey, "Turkey", 0, 0, 2, false, NEU, EUR},
	{Morocco, "Morocco", 0, 0, 3, false, NEU, AFR},
	{WestAfricanStates, "WestAfricanStates", 0, 0, 2, false, NEU, AFR},
	{IvoryCoast, "IvoryCoast", 0, 0, 2, false, NEU, AFR},
	{Algeria, "Algeria", 0, 0, 2, true, NEU, AFR},
	{SaharanStates, "SaharanStates", 0, 0, 1, false, NEU, AFR},
	{Nigeria, "Nigeria", 0, 0, 1, true, NEU, AFR},
	{Tunisia, "Tunisia", 0, 0, 2, false, NEU, AFR},
	{Cameroon, "Cameroon", 0, 0, 1, false, NEU, AFR},
	{Angola, "Angola", 0, 0, 1, true, NEU, AFR},
	{SouthAfrica, "SouthAfrica", 1, 0, 3, true, NEU, AFR},
	{Zaire, "Zaire", 0, 0, 1, true, NEU, AFR},
	{Botswana, "Botswana", 0, 0, 2, false, NEU, AFR},
	{Zimbabwe, "Zimbabwe", 0, 0, 1, false, NEU, AFR},
	{Sudan, "Sudan", 0, 0, 1, false, NEU, AFR},
	{Ethiopia, "Ethiopia", 0, 0, 1, false, NEU, AFR},
	{Kenya, "Kenya", 0, 0, 2, false, NEU, AFR},
	{SEAfricanStates, "SEAfricanStates", 0, 0, 1, false, NEU, AFR},
	{Somalia, "Somalia", 0, 0, 2, false, NEU, AFR},
	{Libya, "Libya", 0, 0, 2, true, NEU, MDE},
	{Egypt, "Egypt", 0, 0, 2, true, NEU, MDE},
	{Israel, "Israel", 1, 0, 4, true, NEU, MDE},
	{Lebanon, "Lebanon", 0, 0, 1, false, NEU, MDE},
	{Jordan, "Jordan", 0, 0, 2, false, NEU, MDE},
	{Syria, "Syria", 0, 1, 2, false, NEU, MDE},
	{Iraq, "Iraq", 0, 1, 3, true, NEU, MDE},
	{SaudiArabia, "SaudiArabia", 0, 0, 3, true, NEU, MDE},
	{GulfStates, "GulfStates", 0, 0, 3, false, NEU, MDE},
	{Iran, "Iran", 1, 0, 2, true, NEU, MDE},
	{Afghanistan, "Afghanistan", 0, 0, 2, false, SOV, ASI},
	{Pakistan, "Pakistan", 0, 0, 2, true, NEU, ASI},
	{India, "India", 0, 0, 3, true, NEU, ASI},
	{Burma, "Burma", 0, 0, 2, false, NEU, ASI},
	{Thailand, "Thailand", 0, 0, 2, true, NEU, ASI},
	{LaosCambodia, "LaosCambodia", 0, 0, 1, false, NEU, ASI},
	{Vietnam, "Vietnam", 0, 0, 1, false, NEU, ASI},
	{Malaysia, "Malaysia", 0, 0, 2, false, NEU, ASI},
	{Indonesia, "Indonesia", 0, 0, 1, false, NEU, ASI},
	{Australia, "Australia", 4, 0, 4, false, NEU, ASI},
	{Taiwan, "Taiwan", 0, 0, 3, false, NEU, ASI},
	{NKorea, "NKorea", 0, 3, 3, true, SOV, ASI},
	{SKorea, "SKorea", 1, 0, 3, true, NEU, ASI},
	{Philippines, "Philippines", 1, 0, 2, false, NEU, ASI},
	{Japan, "Japan", 1, 0, 4, true, USA, ASI},
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
