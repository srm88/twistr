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
	SouthEastAsia Region = Region{
		Name: "SouthEastAsia",
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
	Countries = make(map[CountryId]*Country)
	for _, c := range countryTable {
		Countries[c.Id] = &Country{
			Id:           c.Id,
			Name:         c.Name,
			Inf:          Influence{c.USInf, c.SovInf},
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
	USInf        int
	SovInf       int
	Stability    int
	Battleground bool
	AdjSuper     Aff
	RegionId     RegionId
}{
	{Mexico, "Mexico", 0, 0, 2, true, US, CAM},
	{Guatemala, "Guatemala", 0, 0, 1, false, Neu, CAM},
	{ElSalvador, "ElSalvador", 0, 0, 1, false, Neu, CAM},
	{Honduras, "Honduras", 0, 0, 2, false, Neu, CAM},
	{CostaRica, "CostaRica", 0, 0, 3, false, Neu, CAM},
	{Cuba, "Cuba", 0, 0, 3, true, US, CAM},
	{Nicaragua, "Nicaragua", 0, 0, 1, false, Neu, CAM},
	{Panama, "Panama", 1, 0, 2, true, Neu, CAM},
	{Haiti, "Haiti", 0, 0, 1, false, Neu, CAM},
	{DominicanRep, "DominicanRep", 0, 0, 1, false, Neu, CAM},
	{Ecuador, "Ecuador", 0, 0, 2, false, Neu, SAM},
	{Peru, "Peru", 0, 0, 2, false, Neu, SAM},
	{Colombia, "Colombia", 0, 0, 1, false, Neu, SAM},
	{Chile, "Chile", 0, 0, 3, true, Neu, SAM},
	{Venezuela, "Venezuela", 0, 0, 2, true, Neu, SAM},
	{Argentina, "Argentina", 0, 0, 2, true, Neu, SAM},
	{Bolivia, "Bolivia", 0, 0, 2, false, Neu, SAM},
	{Paraguay, "Paraguay", 0, 0, 2, false, Neu, SAM},
	{Uruguay, "Uruguay", 0, 0, 2, false, Neu, SAM},
	{Brazil, "Brazil", 0, 0, 2, true, Neu, SAM},
	{Canada, "Canada", 0, 0, 4, false, US, EUR},
	{UK, "UK", 5, 0, 5, false, Neu, EUR},
	{SpainPortugal, "SpainPortugal", 0, 0, 2, false, Neu, EUR},
	{France, "France", 0, 0, 3, true, Neu, EUR},
	{Benelux, "Benelux", 0, 0, 3, false, Neu, EUR},
	{Norway, "Norway", 0, 0, 4, false, Neu, EUR},
	{Denmark, "Denmark", 0, 0, 3, false, Neu, EUR},
	{WGermany, "WGermany", 0, 0, 4, true, Neu, EUR},
	{EGermany, "EGermany", 3, 0, 3, true, Neu, EUR},
	{Italy, "Italy", 0, 0, 2, true, Neu, EUR},
	{Austria, "Austria", 0, 0, 4, false, Neu, EUR},
	{Sweden, "Sweden", 0, 0, 4, false, Neu, EUR},
	{Czechoslovakia, "Czechoslovakia", 0, 0, 3, false, Neu, EUR},
	{Yugoslavia, "Yugoslavia", 0, 0, 3, false, Neu, EUR},
	{Poland, "Poland", 0, 0, 3, true, Sov, EUR},
	{Greece, "Greece", 0, 0, 2, false, Neu, EUR},
	{Hungary, "Hungary", 0, 0, 3, false, Neu, EUR},
	{Finland, "Finland", 1, 0, 4, false, Sov, EUR},
	{Romania, "Romania", 0, 0, 3, false, Sov, EUR},
	{Bulgaria, "Bulgaria", 0, 0, 3, false, Neu, EUR},
	{Turkey, "Turkey", 0, 0, 2, false, Neu, EUR},
	{Morocco, "Morocco", 0, 0, 3, false, Neu, AFR},
	{WestAfricanStates, "WestAfricanStates", 0, 0, 2, false, Neu, AFR},
	{IvoryCoast, "IvoryCoast", 0, 0, 2, false, Neu, AFR},
	{Algeria, "Algeria", 0, 0, 2, true, Neu, AFR},
	{SaharanStates, "SaharanStates", 0, 0, 1, false, Neu, AFR},
	{Nigeria, "Nigeria", 0, 0, 1, true, Neu, AFR},
	{Tunisia, "Tunisia", 0, 0, 2, false, Neu, AFR},
	{Cameroon, "Cameroon", 0, 0, 1, false, Neu, AFR},
	{Angola, "Angola", 0, 0, 1, true, Neu, AFR},
	{SouthAfrica, "SouthAfrica", 1, 0, 3, true, Neu, AFR},
	{Zaire, "Zaire", 0, 0, 1, true, Neu, AFR},
	{Botswana, "Botswana", 0, 0, 2, false, Neu, AFR},
	{Zimbabwe, "Zimbabwe", 0, 0, 1, false, Neu, AFR},
	{Sudan, "Sudan", 0, 0, 1, false, Neu, AFR},
	{Ethiopia, "Ethiopia", 0, 0, 1, false, Neu, AFR},
	{Kenya, "Kenya", 0, 0, 2, false, Neu, AFR},
	{SEAfricanStates, "SEAfricanStates", 0, 0, 1, false, Neu, AFR},
	{Somalia, "Somalia", 0, 0, 2, false, Neu, AFR},
	{Libya, "Libya", 0, 0, 2, true, Neu, MDE},
	{Egypt, "Egypt", 0, 0, 2, true, Neu, MDE},
	{Israel, "Israel", 1, 0, 4, true, Neu, MDE},
	{Lebanon, "Lebanon", 0, 0, 1, false, Neu, MDE},
	{Jordan, "Jordan", 0, 0, 2, false, Neu, MDE},
	{Syria, "Syria", 1, 0, 2, false, Neu, MDE},
	{Iraq, "Iraq", 1, 0, 3, true, Neu, MDE},
	{SaudiArabia, "SaudiArabia", 0, 0, 3, true, Neu, MDE},
	{GulfStates, "GulfStates", 0, 0, 3, false, Neu, MDE},
	{Iran, "Iran", 1, 0, 2, true, Neu, MDE},
	{Afghanistan, "Afghanistan", 0, 0, 2, false, Sov, ASI},
	{Pakistan, "Pakistan", 0, 0, 2, true, Neu, ASI},
	{India, "India", 0, 0, 3, true, Neu, ASI},
	{Burma, "Burma", 0, 0, 2, false, Neu, ASI},
	{Thailand, "Thailand", 0, 0, 2, true, Neu, ASI},
	{LaosCambodia, "LaosCambodia", 0, 0, 1, false, Neu, ASI},
	{Vietnam, "Vietnam", 0, 0, 1, false, Neu, ASI},
	{Malaysia, "Malaysia", 0, 0, 2, false, Neu, ASI},
	{Indonesia, "Indonesia", 0, 0, 1, false, Neu, ASI},
	{Australia, "Australia", 4, 0, 4, false, Neu, ASI},
	{Taiwan, "Taiwan", 0, 0, 3, false, Neu, ASI},
	{NKorea, "NKorea", 3, 0, 3, true, Sov, ASI},
	{SKorea, "SKorea", 1, 0, 3, true, Neu, ASI},
	{Philippines, "Philippines", 1, 0, 2, false, Neu, ASI},
	{Japan, "Japan", 1, 0, 4, true, US, ASI},
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
