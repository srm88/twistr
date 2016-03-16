package twistr

// Affiliation
type Aff int

type Era int

type CountryId int

type CardId int

const (
	// Hardcoded instead of iota. We use these as indices into arrays, so the
	// constants' values are no longer arbitrary.
	US  Aff = 0
	Sov Aff = 1
	Neu Aff = 2
)

const (
	Early Era = iota
	Mid
	Late
)

const (
	Mexico CountryId = iota
	Guatemala
	ElSalvador
	Honduras
	CostaRica
	Cuba
	Nicaragua
	Panama
	Haiti
	DominicanRep
	Ecuador
	Peru
	Colombia
	Chile
	Venezuela
	Argentina
	Bolivia
	Paraguay
	Uruguay
	Brazil
	Canada
	UK
	SpainPortugal
	France
	Benelux
	Norway
	Denmark
	WGermany
	EGermany
	Italy
	Austria
	Sweden
	Czechoslovakia
	Yugoslavia
	Poland
	Greece
	Hungary
	Finland
	Romania
	Bulgaria
	Turkey
	Morocco
	WestAfricanStates
	IvoryCoast
	Algeria
	SaharanStates
	Nigeria
	Tunisia
	Cameroon
	Angola
	SouthAfrica
	Zaire
	Botswana
	Zimbabwe
	Sudan
	Ethiopia
	Kenya
	SEAfricanStates
	Somalia
	Libya
	Egypt
	Israel
	Lebanon
	Jordan
	Syria
	Iraq
	SaudiArabia
	GulfStates
	Iran
	Afghanistan
	Pakistan
	India
	Burma
	Thailand
	LaosCambodia
	Vietnam
	Malaysia
	Indonesia
	Australia
	Taiwan
	NKorea
	SKorea
	Philippines
	Japan
)

const (
	AsiaScoring CardId = iota
	EuropeScoring
	MiddleEastScoring
	DuckAndCover
	FiveYearPlan
	TheChinaCard
	SocialistGovernments
	Fidel
	VietnamRevolts
	Blockade
	KoreanWar
	RomanianAbdication
	ArabIsraeliWar
	Comecon
	Nasser
	WarsawPactFormed
	DeGaulleLeadsFrance
	CapturedNaziScientist
	TrumanDoctrine
	OlympicGames
	NATO
	IndependentReds
	MarshallPlan
	IndoPakistaniWar
	Containment
	CIACreated
	USJapanMutualDefensePact
	SuezCrisis
	EastEuropeanUnrest
	Decolonization
	RedScarePurge
	UNIntervention
	DeStalinization
	NuclearTestBan
	FormosanResolution
	Defectors
	BrushWar
	CentralAmericaScoring
	SoutheastAsiaScoring
	ArmsRace
	CubanMissileCrisis
	NuclearSubs
	Quagmire
	SALTNegotiations
	BearTrap
	Summit
	HowILearnedToStopWorrying
	Junta
	KitchenDebates
	MissileEnvy
	WeWillBuryYou
	BrezhnevDoctrine
	PortugueseEmpireCrumbles
	SouthAfricanUnrest
	Allende
	WillyBrandt
	MuslimRevolution
	ABMTreaty
	CulturalRevolution
	FlowerPower
	U2Incident
	OPEC
	LoneGunman
	ColonialRearGuards
	PanamaCanalReturned
	CampDavidAccords
	PuppetGovernments
	GrainSalesToSoviets
	JohnPaulIIElectedPope
	LatinAmericanDeathSquads
	OASFounded
	NixonPlaysTheChinaCard
	SadatExpelsSoviets
	ShuttleDiplomacy
	TheVoiceOfAmerica
	LiberationTheology
	UssuriRiverSkirmish
	AskNotWhatYourCountry
	AllianceForProgress
	AfricaScoring
	OneSmallStep
	SouthAmericaScoring
	IranianHostageCrisis
	TheIronLady
	ReaganBombsLibya
	StarWars
	NorthSeaOil
	TheReformer
	MarineBarracksBombing
	SovietsShootDownKAL007
	Glasnost
	OrtegaElectedInNicaragua
	Terrorism
	IranContraScandal
	Chernobyl
	LatinAmericanDebtCrisis
	TearDownThisWall
	AnEvilEmpire
	AldrichAmesRemix
	PershingIIDeployed
	Wargames
	Solidarity
	IranIraqWar
	TheCambridgeFive
	SpecialRelationship
	NORAD
	Che
	OurManInTehran
	YuriAndSamantha
	AWACSSaleToSaudis
)

var countryIdLookup = []struct {
	Name string
	Id   CountryId
}{
	{"mexico", Mexico},
	{"Guatemala", Guatemala},
	{"elsalvador", ElSalvador},
	{"honduras", Honduras},
	{"costarica", CostaRica},
	{"cuba", Cuba},
	{"nicaragua", Nicaragua},
	{"panama", Panama},
	{"haiti", Haiti},
	{"dominicanrep", DominicanRep},
	{"ecuador", Ecuador},
	{"peru", Peru},
	{"colombia", Colombia},
	{"chile", Chile},
	{"venezuela", Venezuela},
	{"argentina", Argentina},
	{"bolivia", Bolivia},
	{"paraguay", Paraguay},
	{"uruguay", Uruguay},
	{"brazil", Brazil},
	{"canada", Canada},
	{"uk", UK},
	{"spainportugal", SpainPortugal},
	{"france", France},
	{"benelux", Benelux},
	{"norway", Norway},
	{"denmark", Denmark},
	{"wgermany", WGermany},
	{"egermany", EGermany},
	{"italy", Italy},
	{"austria", Austria},
	{"sweden", Sweden},
	{"czechoslovakia", Czechoslovakia},
	{"yugoslavia", Yugoslavia},
	{"poland", Poland},
	{"greece", Greece},
	{"hungary", Hungary},
	{"finland", Finland},
	{"romania", Romania},
	{"bulgaria", Bulgaria},
	{"turkey", Turkey},
	{"morocco", Morocco},
	{"westafricanstates", WestAfricanStates},
	{"ivorycoast", IvoryCoast},
	{"algeria", Algeria},
	{"saharanstates", SaharanStates},
	{"nigeria", Nigeria},
	{"tunisia", Tunisia},
	{"cameroon", Cameroon},
	{"angola", Angola},
	{"southafrica", SouthAfrica},
	{"zaire", Zaire},
	{"botswana", Botswana},
	{"zimbabwe", Zimbabwe},
	{"sudan", Sudan},
	{"ethiopia", Ethiopia},
	{"kenya", Kenya},
	{"seafricanstates", SEAfricanStates},
	{"somalia", Somalia},
	{"libya", Libya},
	{"egypt", Egypt},
	{"israel", Israel},
	{"lebanon", Lebanon},
	{"jordan", Jordan},
	{"syria", Syria},
	{"iraq", Iraq},
	{"saudiarabia", SaudiArabia},
	{"gulfstates", GulfStates},
	{"iran", Iran},
	{"afghanistan", Afghanistan},
	{"pakistan", Pakistan},
	{"india", India},
	{"burma", Burma},
	{"thailand", Thailand},
	{"laoscambodia", LaosCambodia},
	{"vietnam", Vietnam},
	{"malaysia", Malaysia},
	{"indonesia", Indonesia},
	{"australia", Australia},
	{"taiwan", Taiwan},
	{"nkorea", NKorea},
	{"skorea", SKorea},
	{"philippines", Philippines},
	{"Japan", Japan},
}

var cardIdLookup = []struct {
	Name string
	Id   CardId
}{
	{"asiascoring", AsiaScoring},
	{"europescoring", EuropeScoring},
	{"middleeastscoring", MiddleEastScoring},
	{"duckandcover", DuckAndCover},
	{"fiveyearplan", FiveYearPlan},
	{"thechinacard", TheChinaCard},
	{"socialistgovernments", SocialistGovernments},
	{"fidel", Fidel},
	{"vietnamrevolts", VietnamRevolts},
	{"blockade", Blockade},
	{"koreanwar", KoreanWar},
	{"romanianabdication", RomanianAbdication},
	{"arabisraeliwar", ArabIsraeliWar},
	{"comecon", Comecon},
	{"nasser", Nasser},
	{"warsawpactformed", WarsawPactFormed},
	{"degaulleleadsfrance", DeGaulleLeadsFrance},
	{"capturednaziscientist", CapturedNaziScientist},
	{"trumandoctrine", TrumanDoctrine},
	{"olympicgames", OlympicGames},
	{"nato", NATO},
	{"independentreds", IndependentReds},
	{"marshallplan", MarshallPlan},
	{"indopakistaniwar", IndoPakistaniWar},
	{"containment", Containment},
	{"ciacreated", CIACreated},
	{"usjapanmutualdefensepact", USJapanMutualDefensePact},
	{"suezcrisis", SuezCrisis},
	{"easteuropeanunrest", EastEuropeanUnrest},
	{"decolonization", Decolonization},
	{"redscarepurge", RedScarePurge},
	{"unintervention", UNIntervention},
	{"destalinization", DeStalinization},
	{"nucleartestban", NuclearTestBan},
	{"formosanresolution", FormosanResolution},
	{"defectors", Defectors},
	{"brushwar", BrushWar},
	{"centralamericascoring", CentralAmericaScoring},
	{"southeastasiascoring", SoutheastAsiaScoring},
	{"armsrace", ArmsRace},
	{"cubanmissilecrisis", CubanMissileCrisis},
	{"nuclearsubs", NuclearSubs},
	{"quagmire", Quagmire},
	{"saltnegotiations", SALTNegotiations},
	{"beartrap", BearTrap},
	{"summit", Summit},
	{"howilearnedtostopworrying", HowILearnedToStopWorrying},
	{"junta", Junta},
	{"kitchendebates", KitchenDebates},
	{"missileenvy", MissileEnvy},
	{"wewillburyyou", WeWillBuryYou},
	{"brezhnevdoctrine", BrezhnevDoctrine},
	{"portugueseempirecrumbles", PortugueseEmpireCrumbles},
	{"southafricanunrest", SouthAfricanUnrest},
	{"allende", Allende},
	{"willybrandt", WillyBrandt},
	{"muslimrevolution", MuslimRevolution},
	{"abmtreaty", ABMTreaty},
	{"culturalrevolution", CulturalRevolution},
	{"flowerpower", FlowerPower},
	{"u2incident", U2Incident},
	{"opec", OPEC},
	{"lonegunman", LoneGunman},
	{"colonialrearguards", ColonialRearGuards},
	{"panamacanalreturned", PanamaCanalReturned},
	{"campdavidaccords", CampDavidAccords},
	{"puppetgovernments", PuppetGovernments},
	{"grainsalestosoviets", GrainSalesToSoviets},
	{"johnpauliielectedpope", JohnPaulIIElectedPope},
	{"latinamericandeathsquads", LatinAmericanDeathSquads},
	{"oasfounded", OASFounded},
	{"nixonplaysthechinacard", NixonPlaysTheChinaCard},
	{"sadatexpelssoviets", SadatExpelsSoviets},
	{"shuttlediplomacy", ShuttleDiplomacy},
	{"thevoiceofamerica", TheVoiceOfAmerica},
	{"liberationtheology", LiberationTheology},
	{"ussuririverskirmish", UssuriRiverSkirmish},
	{"asknotwhatyourcountry", AskNotWhatYourCountry},
	{"allianceforprogress", AllianceForProgress},
	{"africascoring", AfricaScoring},
	{"onesmallstep", OneSmallStep},
	{"southamericascoring", SouthAmericaScoring},
	{"iranianhostagecrisis", IranianHostageCrisis},
	{"theironlady", TheIronLady},
	{"reaganbombslibya", ReaganBombsLibya},
	{"starwars", StarWars},
	{"northseaoil", NorthSeaOil},
	{"thereformer", TheReformer},
	{"marinebarracksbombing", MarineBarracksBombing},
	{"sovietsshootdownkal007", SovietsShootDownKAL007},
	{"glasnost", Glasnost},
	{"ortegaelectedinnicaragua", OrtegaElectedInNicaragua},
	{"terrorism", Terrorism},
	{"irancontrascandal", IranContraScandal},
	{"chernobyl", Chernobyl},
	{"latinamericandebtcrisis", LatinAmericanDebtCrisis},
	{"teardownthiswall", TearDownThisWall},
	{"anevilempire", AnEvilEmpire},
	{"aldrichamesremix", AldrichAmesRemix},
	{"pershingiideployed", PershingIIDeployed},
	{"wargames", Wargames},
	{"solidarity", Solidarity},
	{"iraniraqwar", IranIraqWar},
	{"thecambridgefive", TheCambridgeFive},
	{"specialrelationship", SpecialRelationship},
	{"norad", NORAD},
	{"che", Che},
	{"ourmanintehran", OurManInTehran},
	{"yuriandsamantha", YuriAndSamantha},
	{"awacssaletosaudis", AWACSSaleToSaudis},
}
