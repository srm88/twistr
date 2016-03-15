package twistr

// Affiliation
type Aff int

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
	CentralAmerica RegionId = iota
	SouthAmerica
	Europe
	MiddleEast
	Africa
	Asia
	SEAsia
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
