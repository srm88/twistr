package twistr

var (
	Cards    map[CardId]Card
	EarlyWar []Card
	MidWar   []Card
	LateWar  []Card
)

func init() {
	Cards = make(map[CardId]Card)
	// XXX: always includes optional cards atm
	for _, c := range cardTable {
		card := Card{
			Id:   c.Id,
			Aff:  c.Aff,
			Ops:  c.Ops,
			Name: c.Name,
			Text: c.Text,
			Star: c.Name[len(c.Name)-1] == '*',
			Era:  c.Era,
			Impl: c.Impl,
		}
		Cards[c.Id] = card
		// China card is never in a deck
		if c.Id == TheChinaCard {
			continue
		}
		switch c.Era {
		case Early:
			EarlyWar = append(EarlyWar, card)
		case Mid:
			MidWar = append(MidWar, card)
		case Late:
			LateWar = append(LateWar, card)
		}
	}
}

var cardTable = []struct {
	Id       CardId
	Era      Era
	Optional bool
	Aff      Aff
	Num      int
	Name     string
	Ops      int
	Text     string
	Impl     func(*State, Aff)
}{
	{AsiaScoring, Early, false, NEU, 1, "Asia Scoring", 0, "Presence: 3; Domination: 7; Control: 9; +1 VP per controlled Battleground country in Region; +1 VP per country controlled that is adjacent to enemy superpower; MAY NOT BE HELD!", PlayAsiaScoring},
	{EuropeScoring, Early, false, NEU, 2, "Europe Scoring", 0, "Presence: 3; Domination: 7; Control: Automatic Victory; +1 VP per controlled Battleground country in Region; +1 VP per country controlled that is adjacent to enemy superpower; MAY NOT BE HELD!", PlayEuropeScoring},
	{MiddleEastScoring, Early, false, NEU, 3, "Middle East Scoring", 0, "Presence: 3; Domination: 5; Control: 7; +1 VP per controlled Battleground country in Region; MAY NOT BE HELD!", PlayMiddleEastScoring},
	{DuckAndCover, Early, false, USA, 4, "Duck and Cover", 3, "Degrade the DEFCON level by 1. The US receives VP equal to 5 minus the current DEFCON level.", PlayDuckAndCover},
	{FiveYearPlan, Early, false, USA, 5, "Five Year Plan", 3, "The USSR must randomly discard a card. If the card has a US associated Event, the Event occurs immediately. If the card has a USSR associated Event or an Event applicable to both players, then the card must be discarded without triggering the Event.", PlayFiveYearPlan},
	{TheChinaCard, Early, false, NEU, 6, "The China Card", 4, "This card begins the game with the USSR. When played, the player receives +1 Operations to the Operations value of this card if it uses all its Operations in Asia. It is passed to the opponent once played. A player receives 1 VP for holding this card at the end of Turn 10.", nil},
	{SocialistGovernments, Early, false, SOV, 7, "Socialist Governments", 3, "Remove a total of 3 US Influence from any countries in Western Europe (removing no more than 2 Influence per country). This Event cannot be used after the \"#83 - The Iron Lady\" Event has been played.", PlaySocialistGovernments},
	{Fidel, Early, false, SOV, 8, "Fidel*", 2, "Remove all US Influence from Cuba. USSR adds sufficient Influence in Cuba for Control.", PlayFidel},
	{VietnamRevolts, Early, false, SOV, 9, "Vietnam Revolts*", 2, "Add 2 USSR Influence to Vietnam. For the remainder of the turn, the USSR receives +1 Operations to the Operations value of a card that uses all its Operations in Southeast Asia.", PlayVietnamRevolts},
	{Blockade, Early, false, SOV, 10, "Blockade*", 1, "Unless the US immediately discards a card with an Operations value of 3 or more, remove all US Influence from West Germany.", PlayBlockade},
	{KoreanWar, Early, false, SOV, 11, "Korean War*", 2, "North Korea invades South Korea. Roll a die and subtract (-1) from the die roll for every US controlled country adjacent to South Korea. On a modified die roll of 4-6, the USSR receives 2 VP and replaces all US Influence in South Korea with USSR Influence. The USSR adds 2 to its Military Operations Track.", PlayKoreanWar},
	{RomanianAbdication, Early, false, SOV, 12, "Romanian Abdication*", 1, "Remove all US Influence from Romania. The USSR adds sufficient Influence to Romania for Control.", PlayRomanianAbdication},
	{ArabIsraeliWar, Early, false, SOV, 13, "Arab-Israeli War", 2, "Pan-Arab Coalition invades Israel. Roll a die and subtract (-1) from the die roll for Israel, if it is US controlled, and for every US controlled country adjacent to Israel. On a modified die roll of 4-6, the USSR receives 2 VP and replaces all US Influence in Israel with USSR Influence. The USSR adds 2 to its Military Operations Track. This Event cannot be used after the \"#65 - Camp David Accords\" Event has been played.", PlayArabIsraeliWar},
	{Comecon, Early, false, SOV, 14, "Comecon*", 3, "Add 1 USSR Influence to each of 4 non-US controlled countries of Eastern Europe.", PlayComecon},
	{Nasser, Early, false, SOV, 15, "Nasser*", 1, "Add 2 USSR Influence to Egypt. The US removes half, rounded up, of its Influence from Egypt.", PlayNasser},
	{WarsawPactFormed, Early, false, SOV, 16, "Warsaw Pact Formed*", 3, "Remove all US influence from 4 countries in Eastern Europe or add 5 USSR Influence to any countries in Eastern Europe (adding no more than 2 Influence per country). This Event allows the \"#21 - NATO\" card to be played as an Event.", PlayWarsawPactFormed},
	{DeGaulleLeadsFrance, Early, false, SOV, 17, "De Gaulle Leads France*", 3, "Remove 2 US Influence from France and add 1 USSR Influence to France. This Event cancels the effect(s) of the \"#21 - NATO\" Event for France only.", PlayDeGaulleLeadsFrance},
	{CapturedNaziScientist, Early, false, NEU, 18, "Captured Nazi Scientist*", 1, "Move the Space Race Marker ahead by 1 space.", PlayCapturedNaziScientist},
	{TrumanDoctrine, Early, false, USA, 19, "Truman Doctrine*", 1, "Remove all USSR Influence from a single uncontrolled country in Europe.", PlayTrumanDoctrine},
	{OlympicGames, Early, false, NEU, 20, "Olympic Games", 2, "This player sponsors the Olympics. The opponent must either participate or boycott. If the opponent participates, each player rolls a die and the sponsor adds 2 to their roll. The player with the highest modified die roll receives 2 VP (reroll ties). If the opponent boycotts, degrade the DEFCON level by 1 and the sponsor may conduct Operations as if they played a 4 Ops card.", PlayOlympicGames},
	{NATO, Early, false, USA, 21, "NATO*", 4, "The USSR cannot make Coup Attempts or Realignment rolls against any US controlled countries in Europe. US controlled countries in Europe cannot be attacked by play of the \"#36 - Brush War\" Event. This card requires prior play of either the \"#16 - Warsaw Pact Formed\" or \"#23 - Marshall Plan\" Event(s) in order to be played as an Event.", PlayNATO},
	{IndependentReds, Early, false, USA, 22, "Independent Reds*", 2, "Add US Influence to either Yugoslavia, Romania, Bulgaria, Hungary, or Czechoslovakia so that it equals the USSR Influence in that country.", PlayIndependentReds},
	{MarshallPlan, Early, false, USA, 23, "Marshall Plan*", 4, "Add 1 US Influence to each of any 7 non-USSR controlled countries in Western Europe. This Event allows the \"#21 - NATO\" card to be played as an Event.", PlayMarshallPlan},
	{IndoPakistaniWar, Early, false, NEU, 24, "Indo-Pakistani War", 2, "India invades Pakistan or vice versa (player's choice). Roll a die and subtract (-1) from the die roll for every enemy controlled country adjacent to the target of the invasion (India or Pakistan). On a modified die roll of 4-6, the player receives 2 VP and replaces all the opponent's Influence in the target country with their Influence. The player adds 2 to its Military Operations Track.", PlayIndoPakistaniWar},
	{Containment, Early, false, USA, 25, "Containment*", 3, "All Operations cards played by the US, for the remainder of this turn, receive +1 to their Operations value (to a maximum of 4 Operations per card).", PlayContainment},
	{CIACreated, Early, false, USA, 26, "CIA Created*", 1, "The USSR reveals their hand of cards for this turn. The US may use the Operations value of this card to conduct Operations.", PlayCIACreated},
	{USJapanMutualDefensePact, Early, false, USA, 27, "US/Japan Mutual Defense Pact*", 4, "The US adds sufficient Influence to Japan for Control. The USSR cannot make Coup Attempts or Realignment rolls against Japan.", PlayUSJapanMutualDefensePact},
	{SuezCrisis, Early, false, SOV, 28, "Suez Crisis*", 3, "Remove a total of 4 US Influence from France, the United Kingdom and Israel (removing no more than 2 Influence per country).", PlaySuezCrisis},
	{EastEuropeanUnrest, Early, false, USA, 29, "East European Unrest", 3, "Early or Mid War: Remove 1 USSR Influence from 3 countries in Eastern Europe. Late War: Remove 2 USSR Influence from 3 countries in Eastern Europe.", PlayEastEuropeanUnrest},
	{Decolonization, Early, false, SOV, 30, "Decolonization", 2, "Add 1 USSR Influence to each of any 4 countries in Africa and/or Southeast Asia.", PlayDecolonization},
	{RedScarePurge, Early, false, NEU, 31, "Red Scare/Purge", 4, "All Operations cards played by the opponent, for the remainder of this turn, receive -1 to their Operations value (to a minimum value of 1 Operations point).", PlayRedScarePurge},
	{UNIntervention, Early, false, NEU, 32, "UN Intervention", 1, "Play this card simultaneously with a card containing an opponent's associated Event. The opponent's associated Event is canceled but you may use the Operations value of the opponent's card to conduct Operations. This Event cannot be played during the Headline Phase.", PlayUNIntervention},
	{DeStalinization, Early, false, SOV, 33, "De-Stalinization*", 3, "The USSR may reallocate up to a total of 4 Influence from one or more countries to any non-US controlled countries (adding no more than 2 Influence per country).", PlayDeStalinization},
	{NuclearTestBan, Early, false, NEU, 34, "Nuclear Test Ban", 4, "The player receives VP equal to the current DEFCON level minus 2 then improves the DEFCON level by 2.", PlayNuclearTestBan},
	{FormosanResolution, Early, false, USA, 35, "Formosan Resolution*", 2, "If this card's Event is in effect, Taiwan will be treated as a Battleground country, for scoring purposes only, if Taiwan is US controlled when the Asia Scoring Card is played. This Event is cancelled after the US has played the \"#6 - The China Card\" card.", PlayFormosanResolution},
	{Defectors, Early, false, USA, 103, "Defectors", 2, "The US may play this card during the Headline Phase in order to cancel the USSR Headline Event (including a scoring card). The canceled card is placed into the discard pile. If this card is played by the USSR during its action round, the US gains 1 VP.", PlayDefectors},
	{BrushWar, Mid, false, NEU, 36, "Brush War", 3, "The player attacks any country with a stability number of 1 or 2. Roll a die and subtract (-1) from the die roll for every adjacent enemy controlled country. On a modified die roll of 3-6, the player receives 1 VP and replaces all the opponent's Influence in the target country with their Influence. The player adds 3 to its Military Operations Track.", PlayBrushWar},
	{CentralAmericaScoring, Mid, false, NEU, 37, "Central America Scoring", 0, "Presence: 1; Domination: 3; Control: 5; +1 VP per controlled Battleground country in Region; +1 VP per country controlled that is adjacent to enemy superpower; MAY NOT BE HELD!", PlayCentralAmericaScoring},
	{SoutheastAsiaScoring, Mid, false, NEU, 38, "Southeast Asia Scoring*", 0, "1 VP each for Control of Burma, Cambodia/Laos, Vietnam, Malaysia, Indonesia and the Philippines. 2 VP for Control of Thailand; MAY NOT BE HELD!", PlaySoutheastAsiaScoring},
	{ArmsRace, Mid, false, NEU, 39, "Arms Race", 3, "Compare each player's value on the Military Operations Track. If the phasing player has a higher value than their opponent on the Military Operations Track, that player receives 1 VP. If the phasing player has a higher value than their opponent, and has met the \"required\" amount, on the Military Operations Track, that player receives 3 VP instead.", PlayArmsRace},
	{CubanMissileCrisis, Mid, false, NEU, 40, "Cuban Missile Crisis*", 3, "Set the DEFCON level to 2. Any Coup Attempts by your opponent, for the remainder of this turn, will result in Global Thermonuclear War. Your opponent will lose the game. This card's Event may be canceled, at any time, if the USSR removes 2 Influence from Cuba or the US removes 2 Influence from West Germany or Turkey.", PlayCubanMissileCrisis},
	{NuclearSubs, Mid, false, USA, 41, "Nuclear Subs*", 2, "US Operations used for Coup Attempts in Battleground countries, for the remainder of this turn, do not degrade the DEFCON level. This card's Event does not apply to any Event that would affect the DEFCON level (ex. the \"#40 - Cuban Missile Crisis\" Event).", PlayNuclearSubs},
	{Quagmire, Mid, false, SOV, 42, "Quagmire*", 3, "On the US's next action round, it must discard an Operations card with a value of 2 or more and roll 1-4 on a die to cancel this Event. Repeat this Event for each US action round until the US successfully rolls 1-4 on a die. If the US is unable to discard an Operations card, it must play all of its scoring cards and then skip each action round for the rest of the turn. This Event cancels the effect(s) of the \"#106 - NORAD\" Event (if applicable).", PlayQuagmire},
	{SALTNegotiations, Mid, false, NEU, 43, "SALT Negotiations*", 3, "Improve the DEFCON level by 2. For the remainder of the turn, both players receive -1 to all Coup Attempt rolls. The player of this card's Event may look through the discard pile, pick any 1 non-scoring card, reveal it to their opponent and then place the drawn card into their hand.", PlaySALTNegotiations},
	{BearTrap, Mid, false, USA, 44, "Bear Trap*", 3, "On the USSR's next action round, it must discard an Operations card with a value of 2 or more and roll 1-4 on a die to cancel this Event. Repeat this Event for each USSR action round until the USSR successfully rolls 1-4 on a die. If the USSR is unable to discard an Operations card, it must play all of its scoring cards and then skip each action round for the rest of the turn.", PlayBearTrap},
	{Summit, Mid, false, NEU, 45, "Summit", 1, "Both players roll a die. Each player receives +1 to the die roll for each Region (Europe, Asia, etc.) they Dominate or Control. The player with the highest modified die roll receives 2 VP and may degrade or improve the DEFCON level by 1 (do not reroll ties).", PlaySummit},
	{HowILearnedToStopWorrying, Mid, false, NEU, 46, "How I Learned to Stop Worrying*", 2, "Set the DEFCON level to any level desired (1-5). The player adds 5 to its Military Operations Track.", PlayHowILearnedToStopWorrying},
	{Junta, Mid, false, NEU, 47, "Junta", 2, "Add 2 Influence to a single country in Central or South America. The player may make free Coup Attempts or Realignment rolls in either Central or South America using the Operations value of this card.", PlayJunta},
	{KitchenDebates, Mid, false, USA, 48, "Kitchen Debates*", 1, "If the US controls more Battleground countries than the USSR, the US player uses this Event to poke their opponent in the chest and receive 2 VP!", PlayKitchenDebates},
	{MissileEnvy, Mid, false, NEU, 49, "Missile Envy", 2, "Exchange this card for your opponent's highest value Operations card. If 2 or more cards are tied, opponent chooses. If the exchanged card contains an Event applicable to yourself or both players, it occurs immediately. If it contains an opponent's Event, use the Operations value (no Event). The opponent must use this card for Operations during their next action round.", PlayMissileEnvy},
	{WeWillBuryYou, Mid, false, SOV, 50, "\"We Will Bury You\"*", 4, "Degrade the DEFCON level by 1. Unless the #32 UN Intervention card is played as an Event on the US's next action round, the USSR receives 3 VP.", PlayWeWillBuryYou},
	{BrezhnevDoctrine, Mid, false, SOV, 51, "Brezhnev Doctrine*", 3, "All Operations cards played by the USSR, for the remainder of this turn, receive +1 to their Operations value (to a maximum of 4 Operations per card).", PlayBrezhnevDoctrine},
	{PortugueseEmpireCrumbles, Mid, false, SOV, 52, "Portuguese Empire Crumbles*", 2, "Add 2 USSR Influence to Angola and the SE African States.", PlayPortugueseEmpireCrumbles},
	{SouthAfricanUnrest, Mid, false, SOV, 53, "South African Unrest", 2, "The USSR either adds 2 Influence to South Africa or adds 1 Influence to South Africa and 2 Influence to a single country adjacent to South Africa.", PlaySouthAfricanUnrest},
	{Allende, Mid, false, SOV, 54, "Allende*", 1, "Add 2 USSR Influence to Chile.", PlayAllende},
	{WillyBrandt, Mid, false, SOV, 55, "Willy Brandt*", 2, "The USSR receives 1 VP and adds 1 Influence to West Germany. This Event cancels the effect(s) of the \"#21 - NATO\" Event for West Germany only. This Event is prevented / canceled by the \"#96 - Tear Down this Wall\" Event.", PlayWillyBrandt},
	{MuslimRevolution, Mid, false, SOV, 56, "Muslim Revolution", 4, "Remove all US Influence from 2 of the following countries: Sudan, Iran, Iraq, Egypt, Libya, Saudi Arabia, Syria, Jordan. This Event cannot be used after the \"#110 - AWACS Sale to Saudis\" Event has been played.", PlayMuslimRevolution},
	{ABMTreaty, Mid, false, NEU, 57, "ABM Treaty", 4, "Improve the DEFCON level by 1 and then conduct Operations using the Operations value of this card.", PlayABMTreaty},
	{CulturalRevolution, Mid, false, SOV, 58, "Cultural Revolution*", 3, "If the US has the \"#6 - The China Card\" card, the US must give the card to the USSR (face up and available to be played). If the USSR already has \"#6 - The China Card\" card, the USSR receives 1 VP.", PlayCulturalRevolution},
	{FlowerPower, Mid, false, SOV, 59, "Flower Power*", 4, "The USSR receives 2 VP for every US played \"War\" card (Arab-Israeli War, Korean War, Brush War, Indo-Pakistani War, Iran-Iraq War), used for Operations or an Event, after this card is played. This Event is prevented / canceled by the \"#97 - 'An Evil Empire'\" Event.", PlayFlowerPower},
	{U2Incident, Mid, false, SOV, 60, "U2 Incident*", 3, "The USSR receives 1 VP. If the \"#32 - UN Intervention\" Event is played later this turn, either by the US or the USSR, the USSR receives an additional 1 VP.", PlayU2Incident},
	{OPEC, Mid, false, SOV, 61, "OPEC", 3, "The USSR receives 1 VP for Control of each of the following countries: Egypt, Iran, Libya, Saudi Arabia, Iraq, Gulf States, Venezuela. This Event cannot be used after the \"#86 - North Sea Oil\" Event has been played.", PlayOPEC},
	{LoneGunman, Mid, false, SOV, 62, "\"Lone Gunman\"*", 1, "The US reveals their hand of cards. The USSR may use the Operations value of this card to conduct Operations.", PlayLoneGunman},
	{ColonialRearGuards, Mid, false, USA, 63, "Colonial Rear Guards", 2, "Add 1 US Influence to each of any 4 countries in Africa and/or Southeast Asia.", PlayColonialRearGuards},
	{PanamaCanalReturned, Mid, false, USA, 64, "Panama Canal Returned*", 1, "Add 1 US Influence to Panama, Costa Rica and Venezuela.", PlayPanamaCanalReturned},
	{CampDavidAccords, Mid, false, USA, 65, "Camp David Accords*", 2, "The US receives 1 VP and adds 1 Influence to Israel, Jordan and Egypt. This Event prevents the \"#13 - Arab-Israeli War\" card from being played as an Event.", PlayCampDavidAccords},
	{PuppetGovernments, Mid, false, USA, 66, "Puppet Governments*", 2, "The US may add 1 Influence to 3 countries that do not contain Influence from either the US or USSR.", PlayPuppetGovernments},
	{GrainSalesToSoviets, Mid, false, USA, 67, "Grain Sales to Soviets", 2, "The US randomly selects 1 card from the USSR's hand (if available). The US must either play the card or return it to the USSR. If the card is returned, or the USSR has no cards, the US may use the Operations value of this card to conduct Operations.", PlayGrainSalesToSoviets},
	{JohnPaulIIElectedPope, Mid, false, USA, 68, "John Paul II Elected Pope*", 2, "Remove 2 USSR Influence from Poland and add 1 US Influence to Poland. This Event allows the \"#101 - Solidarity\" card to be played as an Event.", PlayJohnPaulIIElectedPope},
	{LatinAmericanDeathSquads, Mid, false, NEU, 69, "Latin American Death Squads", 2, "All of the phasing player's Coup Attempts in Central and South America, for the remainder of this turn, receive +1 to their die roll. All of the opponent's Coup Attempts in Central and South America, for the remainder of this turn, receive -1 to their die roll.", PlayLatinAmericanDeathSquads},
	{OASFounded, Mid, false, USA, 70, "OAS Founded*", 1, "Add a total of 2 US Influence to any countries in Central or South America.", PlayOASFounded},
	{NixonPlaysTheChinaCard, Mid, false, USA, 71, "Nixon Plays the China Card*", 2, "If the USSR has the \"#6 - The China Card\" card, the USSR must give the card to the US (face down and unavailable for immediate play). If the US already has the \"#6 - The China Card\" card, the US receives 2 VP.", PlayNixonPlaysTheChinaCard},
	{SadatExpelsSoviets, Mid, false, USA, 72, "Sadat Expels Soviets*", 1, "Remove all USSR Influence from Egypt and add 1 US Influence to Egypt.", PlaySadatExpelsSoviets},
	{ShuttleDiplomacy, Mid, false, USA, 73, "Shuttle Diplomacy", 3, "If this card's Event is in effect, subtract (-1) a Battleground country from the USSR total and then discard this card during the next scoring of the Middle East or Asia (which ever comes first).", PlayShuttleDiplomacy},
	{TheVoiceOfAmerica, Mid, false, USA, 74, "The Voice of America", 2, "Remove 4 USSR Influence from any countries NOT in Europe (removing no more than 2 Influence per country).", PlayTheVoiceOfAmerica},
	{LiberationTheology, Mid, false, SOV, 75, "Liberation Theology", 2, "Add a total of 3 USSR Influence to any countries in Central America (adding no more than 2 Influence per country).", PlayLiberationTheology},
	{UssuriRiverSkirmish, Mid, false, USA, 76, "Ussuri River Skirmish*", 3, "If the USSR has the \"#6 - The China Card\" card, the USSR must give the card to the US (face up and available for play). If the US already has the \"#6 - The China Card\" card, add a total of 4 US Influence to any countries in Asia (adding no more than 2 Influence per country).", PlayUssuriRiverSkirmish},
	{AskNotWhatYourCountry, Mid, false, USA, 77, "\"Ask Not What Your Country...\"*", 3, "The US may discard up to their entire hand of cards (including scoring cards) to the discard pile and draw replacements from the draw pile. The number of cards to be discarded must be decided before drawing any replacement cards from the draw pile.", PlayAskNotWhatYourCountry},
	{AllianceForProgress, Mid, false, USA, 78, "Alliance for Progress*", 3, "The US receives 1 VP for each US controlled Battleground country in Central and South America.", PlayAllianceForProgress},
	{AfricaScoring, Mid, false, NEU, 79, "Africa Scoring", 0, "Presence: 1; Domination: 4; Control: 6; +1 VP per controlled Battleground country in Region; MAY NOT BE HELD!", PlayAfricaScoring},
	{OneSmallStep, Mid, false, NEU, 80, "\"One Small Step...\"", 2, "If you are behind on the Space Race Track, the player uses this Event to move their marker 2 spaces forward on the Space Race Track. The player receives VP only from the final space moved into.", PlayOneSmallStep},
	{SouthAmericaScoring, Mid, false, NEU, 81, "South America Scoring", 0, "Presence: 2; Domination: 5; Control: 6; +1 VP per controlled Battleground country in Region; MAY NOT BE HELD!", PlaySouthAmericaScoring},
	{IranianHostageCrisis, Late, false, SOV, 82, "Iranian Hostage Crisis*", 3, "Remove all US Influence and add 2 USSR Influence to Iran. This card's Event requires the US to discard 2 cards, instead of 1 card, if the \"#92 - Terrorism\" Event is played.", PlayIranianHostageCrisis},
	{TheIronLady, Late, false, USA, 83, "The Iron Lady*", 3, "Add 1 USSR Influence to Argentina and remove all USSR Influence from the United Kingdom. The US receives 1 VP. This Event prevents the \"#7 - Socialist Governments\" card from being played as an Event.", PlayTheIronLady},
	{ReaganBombsLibya, Late, false, USA, 84, "Reagan Bombs Libya*", 2, "The US receives 1 VP for every 2 USSR Influence in Libya.", PlayReaganBombsLibya},
	{StarWars, Late, false, USA, 85, "Star Wars*", 2, "If the US is ahead on the Space Race Track, the US player uses this Event to look through the discard pile, pick any 1 non-scoring card and play it immediately as an Event.", PlayStarWars},
	{NorthSeaOil, Late, false, USA, 86, "North Sea Oil*", 3, "The US may play 8 cards (in 8 action rounds) for this turn only. This Event prevents the \"#61 - OPEC\" card from being played as an Event.", PlayNorthSeaOil},
	{TheReformer, Late, false, SOV, 87, "The Reformer*", 3, "Add 4 USSR Influence to Europe (adding no more than 2 Influence per country). If the USSR is ahead of the US in VP, 6 Influence may be added to Europe instead. The USSR may no longer make Coup Attempts in Europe.", PlayTheReformer},
	{MarineBarracksBombing, Late, false, SOV, 88, "Marine Barracks Bombing*", 2, "Remove all US Influence in Lebanon and remove a total of 2 US Influence from any countries in the Middle East.", PlayMarineBarracksBombing},
	{SovietsShootDownKAL007, Late, false, USA, 89, "Soviets Shoot Down KAL-007*", 4, "Degrade the DEFCON level by 1 and the US receives 2 VP. The US may place influence or make Realignment rolls, using this card, if South Korea is US controlled.", PlaySovietsShootDownKAL007},
	{Glasnost, Late, false, SOV, 90, "Glasnost*", 4, "Improve the DEFCON level by 1 and the USSR receives 2 VP. The USSR may make Realignment rolls or add Influence, using this card, if the \"#87 - The Reformer\" Event has already been played.", PlayGlasnost},
	{OrtegaElectedInNicaragua, Late, false, SOV, 91, "Ortega Elected in Nicaragua*", 2, "Remove all US Influence from Nicaragua. The USSR may make a free Coup Attempt, using this card's Operations value, in a country adjacent to Nicaragua.", PlayOrtegaElectedInNicaragua},
	{Terrorism, Late, false, NEU, 92, "Terrorism", 2, "The player's opponent must randomly discard 1 card from their hand. If the \"#82 - Iranian Hostage Crisis\" Event has already been played, a US player (if applicable) must randomly discard 2 cards from their hand.", PlayTerrorism},
	{IranContraScandal, Late, false, SOV, 93, "Iran-Contra Scandal*", 2, "All US Realignment rolls, for the remainder of this turn, receive -1 to their die roll.", PlayIranContraScandal},
	{Chernobyl, Late, false, USA, 94, "Chernobyl*", 3, "The US must designate a single Region (Europe, Asia, etc.) that, for the remainder of the turn, the USSR cannot add Influence to using Operations points.", PlayChernobyl},
	{LatinAmericanDebtCrisis, Late, false, SOV, 95, "Latin American Debt Crisis", 2, "The US must immediately discard a card with an Operations value of 3 or more or the USSR may double the amount of USSR Influence in 2 countries in South America.", PlayLatinAmericanDebtCrisis},
	{TearDownThisWall, Late, false, USA, 96, "Tear Down this Wall*", 3, "Add 3 US Influence to East Germany. The US may make free Coup Attempts or Realignment rolls in Europe using the Operations value of this card. This Event prevents / cancels the effect(s) of the \"#55 - Willy Brandt\" Event.", PlayTearDownThisWall},
	{AnEvilEmpire, Late, false, USA, 97, "\"An Evil Empire\"*", 3, "The US receives 1 VP. This Event prevents / cancels the effect(s) of the \"#59 - Flower Power\" Event.", PlayAnEvilEmpire},
	{AldrichAmesRemix, Late, false, SOV, 98, "Aldrich Ames Remix*", 3, "The US reveals their hand of cards, face-up, for the remainder of the turn and the USSR discards a card from the US hand.", PlayAldrichAmesRemix},
	{PershingIIDeployed, Late, false, SOV, 99, "Pershing II Deployed*", 3, "The USSR receives 1 VP. Remove 1 US Influence from any 3 countries in Western Europe.", PlayPershingIIDeployed},
	{Wargames, Late, false, NEU, 100, "Wargames*", 4, "If the DEFCON level is 2, the player may immediately end the game after giving their opponent 6 VP. How about a nice game of chess?", PlayWargames},
	{Solidarity, Late, false, USA, 101, "Solidarity*", 2, "Add 3 US Influence to Poland. This card requires prior play of the \"#68 - John Paul II Elected Pope\" Event in order to be played as an Event.", PlaySolidarity},
	{IranIraqWar, Late, false, NEU, 102, "Iran-Iraq War*", 2, "Iran invades Iraq or vice versa (player's choice). Roll a die and subtract (-1) from the die roll for every enemy controlled country adjacent to the target of the invasion (Iran or Iraq). On a modified die roll of 4-6, the player receives 2 VP and replaces all the opponent's Influence in the target country with their Influence. The player adds 2 to its Military Operations Track.", PlayIranIraqWar},
	{TheCambridgeFive, Early, true, SOV, 104, "The Cambridge Five", 2, "The US reveals all scoring cards in their hand of cards. The USSR player may add 1 USSR Influence to a single Region named on one of the revealed scoring cards. This card can not be played as an Event during the Late War.", PlayTheCambridgeFive},
	{SpecialRelationship, Early, true, USA, 105, "Special Relationship", 2, "Add 1 US Influence to a single country adjacent to the U.K. if the U.K. is US-controlled but NATO is not in effect. Add 2 US Influence to a single country in Western Europe, and the US gains 2 VP, if the U.K. is US-controlled and NATO is in effect.", PlaySpecialRelationship},
	{NORAD, Early, true, USA, 106, "NORAD", 3, "Add 1 US Influence to a single country containing US Influence, at the end of each Action Round, if Canada is US-controlled and the DEFCON level moved to 2 during that Action Round. This Event is canceled by the \"#42 - Quagmire\" Event.", PlayNORAD},
	{Che, Mid, true, SOV, 107, "Che", 3, "The USSR may perform a Coup Attempt, using this card's Operations value, against a non-Battleground country in Central America, South America or Africa. The USSR may perform a second Coup Attempt, against a different non-Battleground country in Central America, South America or Africa, if the first Coup Attempt removed any US Influence from the target country.", PlayChe},
	{OurManInTehran, Mid, true, USA, 108, "Our Man in Tehran*", 2, "If the US controls at least one Middle East country, the US player uses this Event to draw the top 5 cards from the draw pile. The US may discard any or all of the drawn cards, after revealing the discarded card(s) to the USSR player, without triggering the Event(s). Any remaining drawn cards are returned to the draw pile and the draw pile is reshuffled.", PlayOurManInTehran},
	{YuriAndSamantha, Late, true, SOV, 109, "Yuri and Samantha*", 2, "The USSR receives 1 VP for each US Coup Attempt performed during the remainder of the Turn.", PlayYuriAndSamantha},
	{AWACSSaleToSaudis, Late, true, USA, 110, "AWACS Sale to Saudis*", 3, "Add 2 US Influence to Saudi Arabia. This Event prevents the \"#56 - Muslim Revolution\" card from being played as an Event.", PlayAWACSSaleToSaudis},
}
