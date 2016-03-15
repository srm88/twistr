package twistr

// Affiliation
type Aff int

type Era int

type RegionId int

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
	CentralAmerica RegionId = iota
	SouthAmerica
	Europe
	MiddleEast
	Africa
	Asia
	SEAsia
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
