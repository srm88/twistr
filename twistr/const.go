package twistr

// Affiliation
type Aff uint8

type RegionId uint8

type CountryId uint8

type CardId uint8

const (
    // Hardcoded instead of iota. We use these as indices into arrays, so the
    // constants' values are no longer arbitrary.
    Aff US = 0
    Aff Sov = 1
    Aff Neu = 2
)

const (
    RegionId CentralAmerica = iota
    RegionId SouthAmerica
    RegionId Europe
    RegionId MiddleEast
    RegionId Africa
    RegionId Asia
    RegionId SEAsia
)

const (
    CardId AsiaScoring = iota
    CardId EuropeScoring
    CardId MiddleEastScoring
    CardId DuckAndCover
    CardId FiveYearPlan
    CardId TheChinaCard
    CardId SocialistGovernments
    CardId Fidel
    CardId VietnamRevolts
    CardId Blockade
    CardId KoreanWar
    CardId RomanianAbdication
    CardId ArabIsraeliWar
    CardId Comecon
    CardId Nasser
    CardId WarsawPactFormed
    CardId DeGaulleLeadsFrance
    CardId CapturedNaziScientist
    CardId TrumanDoctrine
    CardId OlympicGames
    CardId NATO
    CardId IndependentReds
    CardId MarshallPlan
    CardId IndoPakistaniWar
    CardId Containment
    CardId CIACreated
    CardId USJapanMutualDefensePact
    CardId SuezCrisis
    CardId EastEuropeanUnrest
    CardId Decolonization
    CardId RedScarePurge
    CardId UNIntervention
    CardId DeStalinization
    CardId NuclearTestBan
    CardId FormosanResolution
    CardId Defectors
    CardId BrushWar
    CardId CentralAmericaScoring
    CardId SoutheastAsiaScoring
    CardId ArmsRace
    CardId CubanMissileCrisis
    CardId NuclearSubs
    CardId Quagmire
    CardId SALTNegotiations
    CardId BearTrap
    CardId Summit
    CardId HowILearnedToStopWorrying
    CardId Junta
    CardId KitchenDebates
    CardId MissileEnvy
    CardId WeWillBuryYou
    CardId BrezhnevDoctrine
    CardId PortugueseEmpireCrumbles
    CardId SouthAfricanUnrest
    CardId Allende
    CardId WillyBrandt
    CardId MuslimRevolution
    CardId ABMTreaty
    CardId CulturalRevolution
    CardId FlowerPower
    CardId U2Incident
    CardId OPEC
    CardId LoneGunman
    CardId ColonialRearGuards
    CardId PanamaCanalReturned
    CardId CampDavidAccords
    CardId PuppetGovernments
    CardId GrainSalesToSoviets
    CardId JohnPaulIIElectedPope
    CardId LatinAmericanDeathSquads
    CardId OASFounded
    CardId NixonPlaysTheChinaCard
    CardId SadatExpelsSoviets
    CardId ShuttleDiplomacy
    CardId TheVoiceOfAmerica
    CardId LiberationTheology
    CardId UssuriRiverSkirmish
    CardId AskNotWhatYourCountry
    CardId AllianceForProgress
    CardId AfricaScoring
    CardId OneSmallStep
    CardId SouthAmericaScoring
    CardId IranianHostageCrisis
    CardId TheIronLady
    CardId ReaganBombsLibya
    CardId StarWars
    CardId NorthSeaOil
    CardId TheReformer
    CardId MarineBarracksBombing
    CardId SovietsShootDownKAL007
    CardId Glasnost
    CardId OrtegaElectedInNicaragua
    CardId Terrorism
    CardId IranContraScandal
    CardId Chernobyl
    CardId LatinAmericanDebtCrisis
    CardId TearDownThisWall
    CardId AnEvilEmpire
    CardId AldrichAmesRemix
    CardId PershingIIDeployed
    CardId Wargames
    CardId Solidarity
    CardId IranIraqWar
    CardId TheCambridgeFive
    CardId SpecialRelationship
    CardId NORAD
    CardId Che
    CardId OurManInTehran
    CardId YuriAndSamantha
    CardId AWACSSaleToSaudis
)

