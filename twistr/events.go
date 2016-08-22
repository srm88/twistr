package twistr

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*
 * Early War
 * ---------
 */

func PlayAsiaScoring(s *State, player Aff) {
	/* Presence: 3; Domination: 7; Control: 9; +1 VP per controlled Battleground
	country in Region; +1 VP per country controlled that is adjacent to enemy
	superpower; MAY NOT BE HELD!  */
	Score(s, player, Asia)
	s.Cancel(ShuttleDiplomacy)
}

func PlayEuropeScoring(s *State, player Aff) {
	/* Presence: 3; Domination: 7; Control: Automatic Victory; +1 VP per
	controlled Battleground country in Region; +1 VP per country controlled that
	is adjacent to enemy superpower; MAY NOT BE HELD!  */
	Score(s, player, Europe)
}

func PlayMiddleEastScoring(s *State, player Aff) {
	/* Presence: 3; Domination: 5; Control: 7; +1 VP per controlled Battleground
	country in Region; MAY NOT BE HELD!  */
	Score(s, player, MiddleEast)
	s.Cancel(ShuttleDiplomacy)
}

func PlayDuckAndCover(s *State, player Aff) {
	/* Degrade the DEFCON level by 1. The US receives VP equal to 5 minus the
	   current DEFCON level.  */
	s.DegradeDefcon(1)
	s.GainVP(USA, 5-s.Defcon)
}

func PlayFiveYearPlan(s *State, player Aff) {
	/* The USSR must randomly discard a card. If the card has a US associated
	   Event, the Event occurs immediately. If the card has a USSR associated Event
	   or an Event applicable to both players, then the card must be discarded
	   without triggering the Event.  */
	if len(s.Hands[SOV].Cards) == 0 {
		s.Transcribe("USSR has no card to discard.")
		return
	}
	card := SelectRandomCard(s, SOV)
	s.Transcribe(fmt.Sprintf("%s selected.", card))
	s.Hands[SOV].Remove(card)
	if card.Aff == USA {
		PlayEvent(s, USA, card)
	} else {
		s.Transcribe(fmt.Sprintf("%s discarded per Five Year Plan.", card))
		s.Discard.Push(card)
	}
}

func PlaySocialistGovernments(s *State, player Aff) {
	/* Remove a total of 3 US Influence from any countries in Western Europe
	   (removing no more than 2 Influence per country). This Event cannot be used
	   after the “#83 – The Iron Lady” Event has been played.  */
	SelectInfluence(s, player, "Remove 3 US influence (no more than 2 per country)",
		LessInf(USA, 1), 3,
		MaxPerCountry(2), InRegion(WestEurope), HasInfluence(USA))
}

func PlayFidel(s *State, player Aff) {
	/* Remove all US Influence from Cuba. USSR adds sufficient Influence in Cuba
	   for Control.  */
	cuba := s.Countries[Cuba]
	zeroInf(s, cuba, USA)
	setInf(s, cuba, SOV, cuba.Stability)
}

func PlayVietnamRevolts(s *State, player Aff) {
	/* Add 2 USSR Influence to Vietnam. For the remainder of the turn, the USSR
	   receives +1 Operations to the Operations value of a card that uses all its
	   Operations in Southeast Asia.  */
	s.Event(VietnamRevolts, player)
	plusInf(s, s.Countries[Vietnam], SOV, 2)
}

func PlayBlockade(s *State, player Aff) {
	/* Unless the US immediately discards a card with an Operations value of 3 or
	   more, remove all US Influence from West Germany.  */
	enoughOps := ExceedsOps(2, s, USA)
	if hasInHand(s, USA, enoughOps) &&
		"discard" == SelectChoice(s, USA,
			"Discard a card with >=3 Ops, or remove all influence from West Germany?",
			"discard", "remove") {
		card := SelectCard(s, USA, CardBlacklist(TheChinaCard), enoughOps)
		s.Hands[USA].Remove(card)
		s.Transcribe(fmt.Sprintf("US discards %s for Blockade.", card))
		s.Discard.Push(card)
	} else {
		zeroInf(s, s.Countries[WGermany], USA)
	}
}

func PlayKoreanWar(s *State, player Aff) {
	/* North Korea invades South Korea. Roll a die and subtract (-1) from the die
	   roll for every US controlled country adjacent to South Korea. On a
	   modified die roll of 4-6, the USSR receives 2 VP and replaces all US
	   Influence in South Korea with USSR Influence. The USSR adds 2 to its
	   Military Operations Track.  */
	s.MilOps[SOV] += 2
	roll := SelectRoll(s)
	skorea := s.Countries[SKorea]
	mod := skorea.NumControlledNeighbors(USA)
	if mod > 0 {
		s.Transcribe(fmt.Sprintf("%s rolls %d.", player, roll))
	} else {
		s.Transcribe(fmt.Sprintf("%s rolls %d -%d (US controlled adjacent).", player, roll, mod))
	}
	switch roll - mod {
	case 4, 5, 6:
		s.Transcribe("Korean War succeeds.")
		s.GainVP(SOV, 2)
		plusInf(s, skorea, SOV, skorea.Inf[USA])
		zeroInf(s, skorea, USA)
	default:
		s.Transcribe("Korean War fails.")
	}
}

func PlayRomanianAbdication(s *State, player Aff) {
	/* Remove all US Influence from Romania. The USSR adds sufficient Influence
	   to Romania for Control.  */
	romania := s.Countries[Romania]
	zeroInf(s, romania, USA)
	setInf(s, romania, SOV, romania.Stability)
}

func PlayArabIsraeliWar(s *State, player Aff) {
	/* Pan-Arab Coalition invades Israel. Roll a die and subtract (-1) from the
	   die roll for Israel, if it is US controlled, and for every US controlled
	   country adjacent to Israel. On a modified die roll of 4-6, the USSR
	   receives 2 VP and replaces all US Influence in Israel with USSR Influence.
	   The USSR adds 2 to its Military Operations Track. This Event cannot be
	   used after the “#65 – Camp David Accords” Event has been played.  */
	s.MilOps[SOV] += 2
	roll := SelectRoll(s)
	israel := s.Countries[Israel]
	mods := []Mod{
		{-israel.NumControlledNeighbors(USA), "US controlled adjacent"}}
	if israel.Controlled() == USA {
		mods = append(mods, Mod{-1, "US control of Israel"})
	}
	mod := TotalMod(mods)
	if mod > 0 {
		s.Transcribe(fmt.Sprintf("%s rolls %d.", player, roll))
	} else {
		s.Transcribe(fmt.Sprintf("%s rolls %d %s.", player, roll, ModSummary(mods)))
	}
	switch roll + mod {
	case 4, 5, 6:
		s.Transcribe("Arab-Israeli War succeeds.")
		s.GainVP(SOV, 2)
		plusInf(s, israel, SOV, israel.Inf[USA])
		zeroInf(s, israel, USA)
	default:
		s.Transcribe("Arab-Israeli War fails.")
	}
}

func PlayComecon(s *State, player Aff) {
	/* Add 1 USSR Influence to each of 4 non-US controlled countries of Eastern
	   Europe.  */
	SelectInfluence(s, player, "Choose 4 non-US controlled countries",
		PlusInf(SOV, 1), 4,
		MaxPerCountry(1), InRegion(EastEurope), NotControlledBy(USA))
}

func PlayNasser(s *State, player Aff) {
	/* Add 2 USSR Influence to Egypt. The US removes half, rounded up, of its
	   Influence from Egypt.  */
	egypt := s.Countries[Egypt]
	plusInf(s, egypt, SOV, 2)
	loss := egypt.Inf[USA] / 2
	lessInf(s, egypt, USA, loss)
}

func PlayWarsawPactFormed(s *State, player Aff) {
	/* Remove all US influence from 4 countries in Eastern Europe or add 5 USSR
	   Influence to any countries in Eastern Europe (adding no more than 2
	   Influence per country). This Event allows the “#21 – NATO” card to be
	   played as an Event.  */

	s.Event(WarsawPactFormed, player)
	switch SelectChoice(s, player, "Remove US influence or add USSR influence?", "remove", "add") {
	case "remove":
		s.Transcribe("USSR will remove US influence in Eastern Europe.")
		SelectInfluence(s, player, "4 countries to lose all US influence",
			ZeroInf(USA), 4,
			MaxPerCountry(1), InRegion(EastEurope), HasInfluence(USA))
	case "add":
		s.Transcribe("USSR will add USSR influence in Eastern Europe.")
		SelectInfluence(s, player, "5 influence",
			PlusInf(SOV, 1), 5,
			MaxPerCountry(2), InRegion(EastEurope))
	}
}

func PlayDeGaulleLeadsFrance(s *State, player Aff) {
	/* Remove 2 US Influence from France and add 1 USSR Influence to France. This
	   Event cancels the effect(s) of the “#21 – NATO” Event for France only.  */
	s.Event(DeGaulleLeadsFrance, player)
	france := s.Countries[France]
	lessInf(s, france, USA, 2)
	plusInf(s, france, SOV, 1)
}

func PlayCapturedNaziScientist(s *State, player Aff) {
	/* Move the Space Race Marker ahead by 1 space.  */
	box, _ := nextSRBox(s, player)
	box.Enter(s, player)
}

func PlayTrumanDoctrine(s *State, player Aff) {
	/* Remove all USSR Influence from a single uncontrolled country in Europe.  */
	// XXX what if there are no options?
	SelectInfluence(s, player, "Remove all USSR influence from an uncontrolled country in Europe",
		ZeroInf(SOV), 1,
		InRegion(Europe), ControlledBy(NEU), HasInfluence(SOV))
}

func PlayOlympicGames(s *State, player Aff) {
	/* This player sponsors the Olympics. The opponent must either participate or
	   boycott. If the opponent participates, each player rolls a die and the
	   sponsor adds 2 to their roll. The player with the highest modified die
	   roll receives 2 VP (reroll ties). If the opponent boycotts, degrade the
	   DEFCON level by 1 and the sponsor may conduct Operations as if they
	   played a 4 Ops card.  */
	choice := SelectChoice(s, player.Opp(), "Participate or boycott Olympics?", "participate", "boycott")
	switch choice {
	case "participate":
		s.Transcribe(fmt.Sprintf("%s participates in the Olympics.", player.Opp()))
		rolls := [2]int{0, 0}
		tied := true
		for tied {
			rolls[USA] = SelectRoll(s)
			rolls[SOV] = SelectRoll(s)
			rolls[player] += 2
			tied = rolls[USA] == rolls[SOV]
		}
		s.Transcribe(fmt.Sprintf("US rolls %d +2.", rolls[USA]))
		s.Transcribe(fmt.Sprintf("USSR rolls %d.", rolls[SOV]))
		switch {
		case rolls[USA] > rolls[SOV]:
			s.GainVP(USA, 2)
		case rolls[SOV] > rolls[USA]:
			s.GainVP(SOV, 2)
		}
	case "boycott":
		s.Transcribe(fmt.Sprintf("%s boycotts the Olympics.", player.Opp()))
		s.DegradeDefcon(1)
		s.Transcribe(fmt.Sprintf("%s conducts operations.", player))
		ConductOps(s, player, PseudoCard(4))
	}
}

func PlayNATO(s *State, player Aff) {
	/* The USSR cannot make Coup Attempts or Realignment rolls against any US
	   controlled countries in Europe. US controlled countries in Europe cannot
	   be attacked by play of the “#36 – Brush War” Event. This card requires
	   prior play of either the “#16 – Warsaw Pact Formed” or “#23 – Marshall
	   Plan” Event(s) in order to be played as an Event.  */
	s.Event(NATO, player)
}

func PlayIndependentReds(s *State, player Aff) {
	/* Add US Influence to either Yugoslavia, Romania, Bulgaria, Hungary, or
	   Czechoslovakia so that it equals the USSR Influence in that country.  */
	SelectOneInfluence(s, player, "Choose either Yugoslavia, Romania, Bulgaria, Hungary, or Czechoslovakia",
		MatchInf(SOV, USA),
		InCountries(Yugoslavia, Romania, Bulgaria, Hungary, Czechoslovakia))
}

func PlayMarshallPlan(s *State, player Aff) {
	/* Add 1 US Influence to each of any 7 non-USSR controlled countries in
	   Western Europe. This Event allows the “#21 – NATO” card to be played as an
	   Event.  */
	SelectInfluence(s, player, "Choose 7 non-USSR controlled countries",
		PlusInf(USA, 1), 7,
		MaxPerCountry(1), InRegion(WestEurope), NotControlledBy(SOV))
	s.Event(MarshallPlan, player)
}

func PlayIndoPakistaniWar(s *State, player Aff) {
	/* India invades Pakistan or vice versa (player’s choice). Roll a die and
	   subtract (-1) from the die roll for every enemy controlled country
	   adjacent to the target of the invasion (India or Pakistan). On a modified
	   die roll of 4-6, the player receives 2 VP and replaces all the opponent’s
	   Influence in the target country with their Influence. The player adds 2 to
	   its Military Operations Track.  */
	c := SelectCountry(s, player, "Choose who gets invaded", India, Pakistan)
	s.Transcribe(fmt.Sprintf("%s will be invaded.", c))
	s.MilOps[SOV] += 2
	roll := SelectRoll(s)
	mod := c.NumControlledNeighbors(player.Opp())
	if mod > 0 {
		s.Transcribe(fmt.Sprintf("%s rolls %d.", player, roll))
	} else {
		s.Transcribe(fmt.Sprintf("%s rolls %d -%d (%s controlled adjacent).", player, roll, mod, player.Opp()))
	}
	switch roll - mod {
	case 4, 5, 6:
		s.Transcribe("Indo-Pakistani War succeeds.")
		s.GainVP(player, 2)
		plusInf(s, c, player, c.Inf[player.Opp()])
		zeroInf(s, c, player.Opp())
	default:
		s.Transcribe("Indo-Pakistani War fails.")
	}
}

func PlayContainment(s *State, player Aff) {
	/* All Operations cards played by the US, for the remainder of this turn,
	   receive +1 to their Operations value (to a maximum of 4 Operations per
	   card).  */
	s.TurnEvent(Containment, player)
}

func PlayCIACreated(s *State, player Aff) {
	/* The USSR reveals their hand of cards for this turn. The US may use the
	   Operations value of this card to conduct Operations.  */
	ShowHand(s, SOV, USA)
	s.Transcribe(fmt.Sprintf("%s conducts operations with %s.", player, Cards[CIACreated]))
	ConductOps(s, player, PseudoCard(1))
}

func PlayUSJapanMutualDefensePact(s *State, player Aff) {
	/* The US adds sufficient Influence to Japan for Control. The USSR cannot
	   make Coup Attempts or Realignment rolls against Japan.  */
	japan := s.Countries[Japan]
	toControl := japan.Stability + japan.Inf[SOV]
	setInf(s, japan, USA, Max(japan.Inf[USA], toControl))
	s.Event(USJapanMutualDefensePact, player)
}

func PlaySuezCrisis(s *State, player Aff) {
	/* Remove a total of 4 US Influence from France, the United Kingdom and
	   Israel (removing no more than 2 Influence per country).  */
	franceIsraelOrUK := func(c *Country) error {
		switch c.Id {
		case France, Israel, UK:
			return nil
		default:
			return errors.New("Choose only France, UK, Israel")
		}
	}
	SelectInfluence(s, player, "Remove 4 from France, UK, Israel",
		LessInf(USA, 1), 4,
		MaxPerCountry(2), HasInfluence(USA), franceIsraelOrUK)
}

func PlayEastEuropeanUnrest(s *State, player Aff) {
	/* Early or Mid War: Remove 1 USSR Influence from 3 countries in Eastern
	   Europe. Late War: Remove 2 USSR Influence from 3 countries in Eastern
	   Europe.  */
	reduction := 1
	if s.Era() == Late {
		reduction = 2
	}
	SelectInfluence(s, player, "Choose 3 countries in E Europe",
		LessInf(SOV, reduction), 3,
		MaxPerCountry(1), InRegion(EastEurope), HasInfluence(SOV))
}

func PlayDecolonization(s *State, player Aff) {
	/* Add 1 USSR Influence to each of any 4 countries in Africa and/or Southeast
	   Asia.  */
	SelectInfluence(s, player, "Choose 4 countries in Africa or SE Asia",
		PlusInf(SOV, 1), 4,
		MaxPerCountry(1), InRegion(Africa, SoutheastAsia))
}

func PlayRedScarePurge(s *State, player Aff) {
	/* All Operations cards played by the opponent, for the remainder of this
	   turn, receive -1 to their Operations value (to a minimum value of 1
	   Operations point).  */
	s.TurnEvent(RedScarePurge, player)
}

func PlayUNIntervention(s *State, player Aff) {
	/* Play this card simultaneously with a card containing an opponent’s
	   associated Event. The opponent’s associated Event is canceled but you may
	   use the Operations value of the opponent’s card to conduct Operations.
	   This Event cannot be played during the Headline Phase.  */
	opponentEvent := func(c Card) bool {
		return c.Aff == player.Opp()
	}
	card := SelectCard(s, player, opponentEvent, CardBlacklist(TheChinaCard))
	s.Transcribe(fmt.Sprintf("%s conducts operations with %s.", player, card))
	s.Hands[player].Remove(card)
	ConductOps(s, player, card)
	s.Transcribe(fmt.Sprintf("%s to discard.", card))
	s.Discard.Push(card)
	if s.Effect(U2Incident) {
		s.Transcribe("The USSR receives VP due to U2 Incident")
		s.GainVP(SOV, 1)
	}
}

func PlayDeStalinization(s *State, player Aff) {
	/* The USSR may reallocate up to a total of 4 Influence from one or more
	   countries to any non-US controlled countries (adding no more than 2
	   Influence per country).  */
	from := SelectInfluence(s, player, "Choose 4 influence to relocate",
		LessInf(SOV, 1), 4, HasInfluence(SOV))
	if len(from) == 0 {
		return
	}
	// Note that this does permit the SOV player to remove more influence than
	// she is allowed to place; she'd be unable to place all of them and be
	// forced to lose however many were unplacable. Rely on 'undo'.
	SelectInfluence(s, player, fmt.Sprintf("Relocate %d influence to non-US controlled countries", len(from)),
		PlusInf(SOV, 1), len(from),
		MaxPerCountry(2), NotControlledBy(USA))
}

func PlayNuclearTestBan(s *State, player Aff) {
	/* The player receives VP equal to the current DEFCON level minus 2 then
	   improves the DEFCON level by 2.  */
	s.GainVP(player, s.Defcon-2)
	s.ImproveDefcon(2)
}

func PlayFormosanResolution(s *State, player Aff) {
	/* If this card’s Event is in effect, Taiwan will be treated as a
	   Battleground country, for scoring purposes only, if Taiwan is US
	   controlled when the Asia Scoring Card is played. This Event is cancelled
	   after the US has played the “#6 – The China Card” card.  */
	s.Event(FormosanResolution, player)
}

func PlayDefectors(s *State, player Aff) {
	/* The US may play this card during the Headline Phase in order to cancel the
	   USSR Headline Event (including a scoring card). The canceled card is
	   placed into the discard pile. If this card is played by the USSR during
	   its action round, the US gains 1 VP.  */
	// XXX: bug here that would allow US to play defectors and gain 1VP, e.g.
	// if soviets play Grain Sales to Soviets.
	// Can't address this if we're saying the event is always implemented by
	// the card's affiliated player ...
	if s.Phasing == SOV && s.AR > 0 {
		s.Transcribe("US gains VP due to Defectors.")
		s.GainVP(USA, 1)
	}
}

func PlayTheCambridgeFive(s *State, player Aff) {
	/* The US reveals all scoring cards in their hand of cards. The USSR player
	   may add 1 USSR Influence to a single Region named on one of the revealed
	   scoring cards. This card can not be played as an Event during the Late War. */
	scoringCards := []string{}
	regions := []Region{}
	for _, c := range s.Hands[USA].Cards {
		if c.Scoring() {
			scoringCards = append(scoringCards, c.Name)
			regions = append(regions, c.ScoringRegion())
		}
	}
	s.Transcribe(fmt.Sprintf("%s scoring cards: %s\n", USA, strings.Join(scoringCards, ", ")))
	if len(scoringCards) == 0 {
		return
	}
	SelectInfluence(s, player, "Place one influence in one of the regions",
		PlusInf(SOV, 1), 1,
		InRegion(regions...))
}

func PlaySpecialRelationship(s *State, player Aff) {
	/* Add 1 US Influence to a single country adjacent to the U.K. if the U.K.
	   is US-controlled but NATO is not in effect. Add 2 US Influence to a single
	   country in Western Europe, and the US gains 2 VP, if the U.K. is
	   US-controlled and NATO is in effect. */
	ukControlled := s.Countries[UK].Controlled() == USA
	switch {
	case ukControlled && s.Effect(NATO):
		SelectOneInfluence(s, player, "Choose a country in W Europe",
			PlusInf(USA, 2),
			InRegion(WestEurope))
		s.GainVP(USA, 2)
	case ukControlled:
		nextToUK := func(c *Country) error {
			for _, adj := range c.AdjCountries {
				if adj.Id == UK {
					return nil
				}
			}
			return fmt.Errorf("%s not adjacent to UK", c)
		}
		SelectOneInfluence(s, player, "Choose a country adjacent to the UK",
			PlusInf(USA, 1),
			nextToUK)
	}
}

func PlayNORAD(s *State, player Aff) {
	/* Add 1 US Influence to a single country containing US Influence, at the
	   end of each Action Round, if Canada is US-controlled and the DEFCON level
	   moved to 2 during that Action Round. This Event is canceled by the “#42 –
	   Quagmire” Event. */
	s.Event(NORAD, player)
}

/*
 * Mid War
 * -------
 */

func PlayBrushWar(s *State, player Aff) {
	/* The player attacks any country with a stability number of 1 or 2. Roll a
	   die and subtract (-1) from the die roll for every adjacent enemy controlled
	   country. On a modified die roll of 3-6, the player receives 1 VP and
	   replaces all the opponent’s Influence in the target country with their
	   Influence. The player adds 3 to its Military Operations Track. */
	stabLTE := func(c *Country) error {
		if c.Stability > 2 {
			return fmt.Errorf("Country needs to have stability of 1 or 2")
		}
		return nil
	}
	c := SelectOneInfluence(s, player, "Country to attack",
		NoOp,
		stabLTE)
	s.Transcribe(fmt.Sprintf("%s will be attacked.", c))
	s.MilOps[player] += 3
	roll := SelectRoll(s)
	mod := c.NumControlledNeighbors(player.Opp())
	if mod > 0 {
		s.Transcribe(fmt.Sprintf("%s rolls %d.", player, roll))
	} else {
		s.Transcribe(fmt.Sprintf("%s rolls %d -%d (%s controlled adjacent).", player, roll, mod, player.Opp()))
	}
	switch roll - mod {
	case 3, 4, 5, 6:
		s.Transcribe("Brush War succeeds.")
		s.GainVP(player, 1)
		plusInf(s, c, player, c.Inf[player.Opp()])
		zeroInf(s, c, player.Opp())
	default:
		s.Transcribe("Brush War fails.")
	}
}

func PlayCentralAmericaScoring(s *State, player Aff) {
	/* Presence: 1; Domination: 3; Control: 5; +1 VP per controlled Battleground
	   country in Region; +1 VP per country controlled that is adjacent to enemy
	   superpower; MAY NOT BE HELD! */
	// XXX are we enforcing may-not-hold scoring card rule?
	Score(s, player, CentralAmerica)
}

func PlaySoutheastAsiaScoring(s *State, player Aff) {
	/* 1 VP each for Control of Burma, Cambodia/Laos, Vietnam, Malaysia,
	   Indonesia and the Philippines. 2 VP for Control of Thailand; MAY NOT BE
	   HELD! */
	usaMods := []Mod{}
	sovMods := []Mod{}
	for _, c := range SoutheastAsia.Countries {
		points := 1
		if c == Thailand {
			points = 2
		}
		switch s.Countries[c].Controlled() {
		case USA:
			usaMods = append(usaMods, Mod{points, s.Countries[c].Name})
		case SOV:
			sovMods = append(usaMods, Mod{points, s.Countries[c].Name})
		}
	}
	usaScore := TotalMod(usaMods)
	sovScore := TotalMod(sovMods)
	s.Transcribe(fmt.Sprintf("US scores %d: %s.", usaScore, ModSummary(usaMods)))
	s.Transcribe(fmt.Sprintf("USSR scores %d: %s.", sovScore, ModSummary(sovMods)))
	switch {
	case usaScore > sovScore:
		s.GainVP(USA, usaScore-sovScore)
	case sovScore > usaScore:
		s.GainVP(SOV, sovScore-usaScore)
	}
}

func PlayArmsRace(s *State, player Aff) {
	/* Compare each player’s value on the Military Operations Track. If the
	   phasing player has a higher value than their opponent on the Military
	   Operations Track, that player receives 1 VP. If the phasing player has a
	   higher value than their opponent, and has met the “required” amount, on the
	   Military Operations Track, that player receives 3 VP instead. */
	playerMilOps := s.MilOps[player]
	oppMilOps := s.MilOps[player.Opp()]
	switch {
	case playerMilOps > oppMilOps && playerMilOps >= s.Defcon:
		s.Transcribe(fmt.Sprintf("%s gains VP for exceeding %s Military Operations and meeting requirements.", player, player.Opp()))
		s.GainVP(player, 3)
	case playerMilOps > oppMilOps:
		s.Transcribe(fmt.Sprintf("%s gains VP for exceeding %s Military Operations.", player, player.Opp()))
		s.GainVP(player, 1)
	}
}

func PlayCubanMissileCrisis(s *State, player Aff) {
	/* Set the DEFCON level to 2. Any Coup Attempts by your opponent, for the
	   remainder of this turn, will result in Global Thermonuclear War. Your
	   opponent will lose the game. This card’s Event may be canceled, at any time,
	   if the USSR removes 2 Influence from Cuba or the US removes 2 Influence from
	   West Germany or Turkey. */
	s.SetDefcon(2)
	s.TurnEvent(CubanMissileCrisis, player)
}

func PlayNuclearSubs(s *State, player Aff) {
	/* US Operations used for Coup Attempts in Battleground countries, for the
	   remainder of this turn, do not degrade the DEFCON level. This card’s Event
	   does not apply to any Event that would affect the DEFCON level (ex. the
	   “#40 – Cuban Missile Crisis” Event). */
	s.TurnEvent(NuclearSubs, player)
}

func PlayQuagmire(s *State, player Aff) {
	/* On the US’s next action round, it must discard an Operations card with a
	   value of 2 or more and roll 1-4 on a die to cancel this Event. Repeat this
	   Event for each US action round until the US successfully rolls 1-4 on a die.
	   If the US is unable to discard an Operations card, it must play all of its
	   scoring cards and then skip each action round for the rest of the turn. This
	   Event cancels the effect(s) of the “#106 – NORAD” Event (if applicable). */
	s.Event(Quagmire, player)
	s.Cancel(NORAD)
}

func PlaySALTNegotiations(s *State, player Aff) {
	/* Improve the DEFCON level by 2. For the remainder of the turn, both
	   players receive -1 to all Coup Attempt rolls. The player of this card’s
	   Event may look through the discard pile, pick any 1 non-scoring card, reveal
	   it to their opponent and then place the drawn card into their hand. */
	s.ImproveDefcon(2)
	s.TurnEvent(SALTNegotiations, player)
	notScoring := func(c Card) bool {
		return !c.Scoring()
	}
	if "yes" == SelectChoice(s, player,
		"Choose a card to add to your hand?",
		"yes", "no") {
		selected := SelectDiscarded(s, player, notScoring)
		s.Transcribe(fmt.Sprintf("%s adds %s into their hand from the discard pile.", player, selected))
		s.Discard.Remove(selected)
		s.Hands[player].Push(selected)
	}
}

func PlayBearTrap(s *State, player Aff) {
	/* On the USSR’s next action round, it must discard an Operations card
	   with a value of 2 or more and roll 1-4 on a die to cancel this Event. Repeat
	   this Event for each USSR action round until the USSR successfully rolls 1-4
	   on a die. If the USSR is unable to discard an Operations card, it must play
	   all of its scoring cards and then skip each action round for the rest of the
	   turn. */
	s.Event(BearTrap, player)
}

func PlaySummit(s *State, player Aff) {
	/* Both players roll a die. Each player receives +1 to the die roll for each
	   Region (Europe, Asia, etc.) they Dominate or Control. The player with the
	   highest modified die roll receives 2 VP and may degrade or improve the
	   DEFCON level by 1 (do not reroll ties). */
	playerRoll := SelectRoll(s)
	oppRoll := SelectRoll(s)
	playerMods, oppMods := []Mod{}, []Mod{}
	for _, region := range Regions {
		sr := ScoreRegion(s.Game, region)
		switch {
		case sr.Levels[player] == Domination || sr.Levels[player] == Control:
			playerMods = append(playerMods, Mod{1, region.Name})
		case sr.Levels[player.Opp()] == Domination || sr.Levels[player.Opp()] == Control:
			oppMods = append(oppMods, Mod{1, region.Name})
		}
	}
	s.Transcribe(fmt.Sprintf("%s rolls %d%s.", player, playerRoll, ModSummary(playerMods)))
	s.Transcribe(fmt.Sprintf("%s rolls %d%s.", player.Opp(), oppRoll, ModSummary(oppMods)))
	playerRoll += TotalMod(playerMods)
	oppRoll += TotalMod(oppMods)
	var winner Aff
	switch {
	case playerRoll > oppRoll:
		winner = player
	case oppRoll > playerRoll:
		winner = player.Opp()
	default:
		s.Transcribe("Tie.")
		return
	}
	s.Transcribe(fmt.Sprintf("%s wins.", winner))
	s.GainVP(winner, 2)
	switch SelectChoice(s, winner,
		"Degrade or improve DEFCON by one level?",
		"improve", "degrade", "neither") {
	case "improve":
		s.ImproveDefcon(1)
	case "degrade":
		s.DegradeDefcon(1)
	}
}

func PlayHowILearnedToStopWorrying(s *State, player Aff) {
	/* Set the DEFCON level to any level desired (1-5). The player adds 5 to its
	   Military Operations Track. */
	choice := SelectChoice(s, player, "Set DEFCON.", "1", "2", "3", "4", "5")
	newDefcon, _ := strconv.Atoi(choice)
	// XXX transcribe milops
	s.MilOps[player] += 5
	s.SetDefcon(newDefcon)
}

func PlayJunta(s *State, player Aff) {
	/* Add 2 Influence to a single country in Central or South America. The
	   player may make free Coup Attempts or Realignment rolls in either Central or
	   South America using the Operations value of this card. */
	SelectOneInfluence(s, player, "Choose a country in Central or South America",
		PlusInf(player, 2),
		InRegion(SouthAmerica, CentralAmerica))
	switch SelectChoice(s, player,
		"Do you want to coup, realign or do nothing in Central/South America?",
		"coup", "realign", "nothing") {
	case "coup":
		OpCoup(s, player, Cards[Junta], true,
			InRegion(SouthAmerica, CentralAmerica))
	case "realign":
		OpRealign(s, player, Cards[Junta], true)
	}
}

func PlayKitchenDebates(s *State, player Aff) {
	/* If the US controls more Battleground countries than the USSR, the US
	   player uses this Event to poke their opponent in the chest and receive 2 VP!
	*/
	usaBG := 0
	sovBG := 0
	for _, c := range s.Countries {
		if !c.Battleground {
			continue
		}
		switch c.Controlled() {
		case SOV:
			sovBG++
		case USA:
			usaBG++
		}
	}
	s.Transcribe(fmt.Sprintf("US controls %d battlegrounds. USSR controls %d.", usaBG, sovBG))
	if usaBG > sovBG {
		s.GainVP(USA, 2)
	}
}

func PlayMissileEnvy(s *State, player Aff) {
	/* Exchange this card for your opponent’s highest value Operations card.
	   If 2 or more cards are tied, opponent chooses. If the exchanged card
	   contains an Event applicable to yourself or both players, it occurs
	   immediately. If it contains an opponent’s Event, use the Operations value
	   (no Event). The opponent must use this card for Operations during their next
	   action round. */
	maxOps := 0
	for _, c := range s.Hands[player.Opp()].Cards {
		if c.Ops > maxOps {
			maxOps = c.Ops
		}
	}
	isMaxOpCard := func(c Card) bool {
		return c.Ops >= maxOps
	}
	selected := SelectCard(s, player.Opp(),
		isMaxOpCard, CardBlacklist(TheChinaCard))
	s.Transcribe(fmt.Sprintf("%s exchanges %s for Missile Envy.", player.Opp(), selected))
	s.Hands[player.Opp()].Remove(selected)
	switch selected.Aff {
	case player, NEU:
		PlayEvent(s, player, selected)
	default:
		s.Transcribe(fmt.Sprintf("%s conducts operations with %s.", player, selected))
		ConductOps(s, player, selected)
		s.Transcribe(fmt.Sprintf("%s to discard.", selected))
		s.Discard.Push(selected)
	}
	s.TurnEvents[MissileEnvy] = player
}

func PlayWeWillBuryYou(s *State, player Aff) {
	/* Degrade the DEFCON level by 1. Unless the #32 UN Intervention card is
	   played as an Event on the US’s next action round, the USSR receives 3 VP.  */
	s.DegradeDefcon(1)
	// Don't do the generic event message
	s.Transcribe("Unless UN Intervention is played as an Event on the US' next action round, the USSR gains 3 VP.")
	s.Events[WeWillBuryYou] = player
}

func PlayBrezhnevDoctrine(s *State, player Aff) {
	/* All Operations cards played by the USSR, for the remainder of this turn,
	   receive +1 to their Operations value (to a maximum of 4 Operations per
	   card). */
	s.TurnEvent(BrezhnevDoctrine, player)
}

func PlayPortugueseEmpireCrumbles(s *State, player Aff) {
	/* Add 2 USSR Influence to Angola and the SE African States. */
	plusInf(s, s.Countries[Angola], SOV, 2)
	plusInf(s, s.Countries[SEAfricanStates], SOV, 2)
}

func PlaySouthAfricanUnrest(s *State, player Aff) {
	/* The USSR either adds 2 Influence to South Africa or adds 1 Influence to
	   South Africa and 2 Influence to a single country adjacent to South Africa. */

	switch SelectChoice(s, player,
		"Option A: add 2 influence to South Africa, or Option B 1 to South Africa and 2 to an adjacent country?",
		"a", "b") {
	case "a":
		plusInf(s, s.Countries[SouthAfrica], SOV, 2)
	default:
		plusInf(s, s.Countries[SouthAfrica], SOV, 1)
		adjToSA := func(c *Country) error {
			for _, adj := range s.Countries[SouthAfrica].AdjCountries {
				if adj.Id == c.Id {
					return nil
				}
			}
			return fmt.Errorf("%s is not adjacent to South Africa", c)
		}
		SelectOneInfluence(s, player, "Add 2 influence to a country adjacent to South Africa",
			PlusInf(SOV, 2),
			adjToSA)
	}
}

func PlayAllende(s *State, player Aff) {
	/* Add 2 USSR Influence to Chile. */
	plusInf(s, s.Countries[Chile], SOV, 2)
}

func PlayWillyBrandt(s *State, player Aff) {
	/* The USSR receives 1 VP and adds 1 Influence to West Germany. This Event
	   cancels the effect(s) of the “#21 – NATO” Event for West Germany only. This
	   Event is prevented / canceled by the “#96 – Tear Down this Wall” Event. */
	s.GainVP(SOV, 1)
	s.Event(WillyBrandt, player)
	plusInf(s, s.Countries[WGermany], SOV, 1)
}

func PlayMuslimRevolution(s *State, player Aff) {
	/* Remove all US Influence from 2 of the following countries: Sudan, Iran,
	   Iraq, Egypt, Libya, Saudi Arabia, Syria, Jordan. This Event cannot be used
	   after the “#110 – AWACS Sale to Saudis” Event has been played. */
	SelectInfluence(s, player, "2 countries to lose all influence",
		ZeroInf(USA), 2,
		InCountries(Sudan, Iran, Iraq, Egypt, Libya, SaudiArabia, Syria, Jordan),
		MaxPerCountry(1), HasInfluence(USA))
}

func PlayABMTreaty(s *State, player Aff) {
	/* Improve the DEFCON level by 1 and then conduct Operations using the
	   Operations value of this card. */
	s.ImproveDefcon(1)
	s.Transcribe(fmt.Sprintf("%s conducts operations with %s.", player, Cards[ABMTreaty]))
	ConductOps(s, player, PseudoCard(Cards[ABMTreaty].Ops))
}

func PlayCulturalRevolution(s *State, player Aff) {
	/* If the US has the “#6 – The China Card” card, the US must give the card
	   to the USSR (face up and available to be played). If the USSR already has
	   “#6 – The China Card” card, the USSR receives 1 VP. */
	if s.ChinaCardPlayer == USA {
		s.ChinaCardMove(SOV, true)
	} else {
		s.Transcribe("The USSR already has the China Card.")
		s.GainVP(SOV, 1)
	}
}

func PlayFlowerPower(s *State, player Aff) {
	/* The USSR receives 2 VP for every US played “War” card (Arab-Israeli War,
	   Korean War, Brush War, Indo-Pakistani War, Iran-Iraq War), used for
	   Operations or an Event, after this card is played. This Event is prevented /
	   canceled by the “#97 – ‘An Evil Empire’” Event. */
	s.Event(FlowerPower, player)
}

func PlayU2Incident(s *State, player Aff) {
	/* The USSR receives 1 VP. If the “#32 – UN Intervention” Event is played
	   later this turn, either by the US or the USSR, the USSR receives an
	   additional 1 VP. */
	s.GainVP(SOV, 1)
	s.TurnEvent(U2Incident, player)
}

func PlayOPEC(s *State, player Aff) {
	/* The USSR receives 1 VP for Control of each of the following countries:
	   Egypt, Iran, Libya, Saudi Arabia, Iraq, Gulf States, Venezuela. This Event
	   cannot be used after the “#86 – North Sea Oil” Event has been played. */
	mods := []Mod{}
	for _, cid := range []CountryId{Egypt, Iran, Libya, SaudiArabia, GulfStates, Venezuela} {
		if s.Countries[cid].Controlled() == SOV {
			mods = append(mods, Mod{1, s.Countries[cid].Name})
		}
	}
	if len(mods) == 0 {
		s.Transcribe("USSR does not score anything for OPEC.")
		return
	}
	s.Transcribe(fmt.Sprintf("USSR Opec scoring: %s.", ModSummary(mods)))
	s.GainVP(SOV, TotalMod(mods))
}

func PlayLoneGunman(s *State, player Aff) {
	/* The US reveals their hand of cards. The USSR may use the Operations value
	   of this card to conduct Operations. */
	ShowHand(s, USA, SOV)
	s.Transcribe(fmt.Sprintf("%s conducts operations with %s.", player, Cards[LoneGunman]))
	ConductOps(s, player, PseudoCard(Cards[LoneGunman].Ops))
}

func PlayColonialRearGuards(s *State, player Aff) {
	/* Add 1 US Influence to each of any 4 countries in Africa and/or Southeast
	   Asia. */
	SelectInfluence(s, player, "Choose 4 countries in Africa and/or Southeast Asia",
		PlusInf(USA, 1), 4,
		MaxPerCountry(1), InRegion(Africa, SoutheastAsia))
}

func PlayPanamaCanalReturned(s *State, player Aff) {
	/* Add 1 US Influence to Panama, Costa Rica and Venezuela.  */
	plusInf(s, s.Countries[Panama], USA, 1)
	plusInf(s, s.Countries[CostaRica], USA, 1)
	plusInf(s, s.Countries[Venezuela], USA, 1)
}

func PlayCampDavidAccords(s *State, player Aff) {
	/* The US receives 1 VP and adds 1 Influence to Israel, Jordan and Egypt.
	   This Event prevents the “#13 – Arab-Israeli War” card from being played as
	   an Event. */
	plusInf(s, s.Countries[Israel], USA, 1)
	plusInf(s, s.Countries[Jordan], USA, 1)
	plusInf(s, s.Countries[Egypt], USA, 1)
	s.GainVP(USA, 1)
	s.Event(CampDavidAccords, player)
}

func PlayPuppetGovernments(s *State, player Aff) {
	/* The US may add 1 Influence to 3 countries that do not contain Influence
	   from either the US or USSR. */
	SelectInfluence(s, player, "Choose 3 countries with no influence from either power",
		PlusInf(USA, 1), 3,
		MaxPerCountry(1), NoInfluence(USA), NoInfluence(SOV))
}

func PlayGrainSalesToSoviets(s *State, player Aff) {
	/* The US randomly selects 1 card from the USSR’s hand (if available). The
	   US must either play the card or return it to the USSR. If the card is
	   returned, or the USSR has no cards, the US may use the Operations value of
	   this card to conduct Operations. */
	if len(s.Hands[SOV].Cards) == 0 {
		s.Transcribe("The USSR player has no cards.")
		s.Transcribe(fmt.Sprintf("US conducts operations with %s.", Cards[GrainSalesToSoviets]))
		ConductOps(s, player, PseudoCard(Cards[GrainSalesToSoviets].Ops))
	} else {
		card := SelectRandomCard(s, SOV)
		s.Transcribe(fmt.Sprintf("The US selects %s from the USSR hand.", card))
		switch SelectChoice(s, player, "Play this card or return it?",
			"play", "return") {
		case "play":
			s.Transcribe(fmt.Sprintf("US plays %s.", card))
			s.Hands[SOV].Remove(card)
			PlayCard(s, player, card)
		default:
			s.Transcribe(fmt.Sprintf("US returns %s to the USSR.", card))
			s.Transcribe(fmt.Sprintf("US conducts operations with %s.", Cards[GrainSalesToSoviets]))
			ConductOps(s, player, PseudoCard(Cards[GrainSalesToSoviets].Ops))
		}
	}
}

func PlayJohnPaulIIElectedPope(s *State, player Aff) {
	/* Remove 2 USSR Influence from Poland and add 1 US Influence to Poland.
	   This Event allows the “#101 – Solidarity” card to be played as an Event.
	*/
	c := s.Countries[Poland]
	lessInf(s, c, SOV, 2)
	plusInf(s, c, USA, 1)
	s.Event(JohnPaulIIElectedPope, player)
}

func PlayLatinAmericanDeathSquads(s *State, player Aff) {
	/* All of the phasing player’s Coup Attempts in Central and South America,
	   for the remainder of this turn, receive +1 to their die roll. All of the
	   opponent’s Coup Attempts in Central and South America, for the remainder of
	   this turn, receive -1 to their die roll. */
	s.Event(LatinAmericanDeathSquads, player)
}

func PlayOASFounded(s *State, player Aff) {
	/* Add a total of 2 US Influence to any countries in Central or South America. */
	SelectInfluence(s, player, "Add a total of 2 influence to countries in Central or South America",
		PlusInf(USA, 1), 2,
		InRegion(CentralAmerica, SouthAmerica))
}

func PlayNixonPlaysTheChinaCard(s *State, player Aff) {
	/* If the USSR has the “#6 – The China Card” card, the USSR must give the
	   card to the US (face down and unavailable for immediate play). If the US
	   already has the “#6 – The China Card” card, the US receives 2 VP. */
	if s.ChinaCardPlayer == SOV {
		s.ChinaCardMove(USA, false)
	} else {
		s.Transcribe("The US already has the China Card.")
		s.GainVP(USA, 2)
	}
}

func PlaySadatExpelsSoviets(s *State, player Aff) {
	/* Remove all USSR Influence from Egypt and add 1 US Influence to Egypt. */
	c := s.Countries[Egypt]
	zeroInf(s, c, SOV)
	plusInf(s, c, USA, 1)
}

func PlayShuttleDiplomacy(s *State, player Aff) {
	/* If this card’s Event is in effect, subtract (-1) a Battleground country
	   from the USSR total and then discard this card during the next scoring of
	   the Middle East or Asia (which ever comes first). */
	s.Event(ShuttleDiplomacy, player)
}

func PlayTheVoiceOfAmerica(s *State, player Aff) {
	/* Remove 4 USSR Influence from any countries NOT in Europe (removing no
	   more than 2 Influence per country). */
	SelectInfluence(s, player, "Remove a total of 4 USSR influence from countries not in Europe (no more than 2 per country)",
		LessInf(SOV, 1), 4,
		InRegion(Asia, Africa, CentralAmerica, SouthAmerica, MiddleEast),
		MaxPerCountry(2), HasInfluence(SOV))
}

func PlayLiberationTheology(s *State, player Aff) {
	/* Add a total of 3 USSR Influence to any countries in Central America
	   (adding no more than 2 Influence per country). */
	SelectInfluence(s, player, "Add a total of 3 influence (no more than 2 per country) to countries in Central America",
		PlusInf(SOV, 1), 3,
		InRegion(CentralAmerica), MaxPerCountry(2))
}

func PlayUssuriRiverSkirmish(s *State, player Aff) {
	/* If the USSR has the “#6 – The China Card” card, the USSR must give the
	   card to the US (face up and available for play). If the US already has the
	   “#6 – The China Card” card, add a total of 4 US Influence to any countries
	   in Asia (adding no more than 2 Influence per country). */
	if s.ChinaCardPlayer == SOV {
		s.ChinaCardMove(USA, true)
	} else {
		s.Transcribe("The US already has the China Card.")
		SelectInfluence(s, player, "Add a total of 4 influence to countries in Central or South America",
			PlusInf(USA, 1), 4,
			InRegion(Asia), MaxPerCountry(2))
	}
}

func PlayAskNotWhatYourCountry(s *State, player Aff) {
	/* The US may discard up to their entire hand of cards (including scoring
	   cards) to the discard pile and draw replacements from the draw pile. The
	   number of cards to be discarded must be decided before drawing any
	   replacement cards from the draw pile. */
	toDiscard := SelectSomeCards(s, USA,
		"Discard up to entire hand of cards",
		s.Hands[USA].Cards)
	toDraw := len(toDiscard)
	if toDraw == 0 {
		return
	}
	for _, c := range toDiscard {
		s.Transcribe(fmt.Sprintf("US discards %s.", c))
		s.Hands[USA].Remove(c)
	}
	s.Discard.Push(toDiscard...)
	s.Transcribe(fmt.Sprintf("US draws %d cards.", toDraw))
	drawn := s.Deck.Draw(toDraw)
	s.Hands[USA].Push(drawn...)
	ShowHand(s, USA, USA)
}

func PlayAllianceForProgress(s *State, player Aff) {
	/* The US receives 1 VP for each US controlled Battleground country in
	   Central and South America. */
	countries := append(CentralAmerica.Countries, SouthAmerica.Countries...)
	mods := []Mod{}
	for _, cId := range countries {
		c := s.Countries[cId]
		if c.Battleground && c.Controlled() == USA {
			mods = append(mods, Mod{1, c.Name})
		}
	}
	if len(mods) == 0 {
		s.Transcribe("US controls no Central or South American battleground countries.")
		return
	}
	s.Transcribe(fmt.Sprintf("US Alliance for Progress scoring: %s", ModSummary(mods)))
	s.GainVP(USA, TotalMod(mods))
}

func PlayAfricaScoring(s *State, player Aff) {
	/* Presence: 1; Domination: 4; Control: 6; +1 VP per controlled Battleground
	   country in Region; MAY NOT BE HELD! */
	Score(s, player, Africa)
}

func PlayOneSmallStep(s *State, player Aff) {
	/* If you are behind on the Space Race Track, the player uses this Event to
	   move their marker 2 spaces forward on the Space Race Track. The player
	   receives VP only from the final space moved into. */
	if s.SpaceRace[player] >= s.SpaceRace[player.Opp()] {
		s.Transcribe(fmt.Sprintf("%s is not behind on the Space Race.", player))
		return
	}
	srb, _ := nextSRBox(s, player)
	if _, ok := s.SREvents[srb.SideEffect]; ok {
		delete(s.SREvents, srb.SideEffect)
	}
	s.SpaceRace[player] += 2
	s.Transcribe(fmt.Sprintf("%s advances to %d on the Space Race.", player, s.SpaceRace[player]))
	srb, _ = nextSRBox(s, player)
	srb.Enter(s, player)
}

func PlaySouthAmericaScoring(s *State, player Aff) {
	/* Presence: 2; Domination: 5; Control: 6; +1 VP per controlled Battleground
	   country in Region; MAY NOT BE HELD! */
	Score(s, player, SouthAmerica)
}

func PlayChe(s *State, player Aff) {
	/* The USSR may perform a Coup Attempt, using this card’s Operations value,
	   against a non-Battleground country in Central America, South America or
	   Africa. The USSR may perform a second Coup Attempt, against a different
	   non-Battleground country in Central America, South America or Africa, if the
	   first Coup Attempt removed any US Influence from the target country. */
	// They technically don't have to coup. Would always be a "yes" except if
	// under cuban missile crisis, so, we do need to ask.
	if "yes" != SelectChoice(s, player, "Perform a Coup attempt with Che?", "yes", "no") {
		s.Transcribe("USSR elects not to Coup with Che.")
		return
	}
	notBg := func(c *Country) error {
		if c.Battleground {
			return fmt.Errorf("%s is a battleground", c.Name)
		}
		return nil
	}
	couped := OpCoup(s, player, Cards[Che], false,
		InRegion(SouthAmerica, CentralAmerica, Africa),
		notBg)
	if couped {
		OpCoup(s, player, Cards[Che], false,
			InRegion(SouthAmerica, CentralAmerica, Africa),
			notBg)
	}
}

func PlayOurManInTehran(s *State, player Aff) {
	/* If the US controls at least one Middle East country, the US player uses
	   this Event to draw the top 5 cards from the draw pile. The US may discard
	   any or all of the drawn cards, after revealing the discarded card(s) to the
	   USSR player, without triggering the Event(s). Any remaining drawn cards are
	   returned to the draw pile and the draw pile is reshuffled. */
	controlled := false
	for _, c := range MiddleEast.Countries {
		if s.Countries[c].Controlled() == USA {
			controlled = true
			break
		}
	}
	if !controlled {
		s.Transcribe("US controls no Middle East countries.")
		return
	}
	cards := s.Deck.Draw(5)
	// Solicit US player to discard each card
	toDiscard := SelectSomeCards(s, USA,
		"Discard which",
		cards)
	discardedSet := make(map[CardId]bool)
	discarded := []string{}
	for _, c := range toDiscard {
		discarded = append(discarded, c.Name)
		discardedSet[c.Id] = true
	}
	backToDraw := []Card{}
	for _, c := range cards {
		if !discardedSet[c.Id] {
			backToDraw = append(backToDraw, c)
		}
	}
	if len(toDiscard) > 0 {
		for _, c := range toDiscard {
			s.Transcribe(fmt.Sprintf("US discards %s.", c))
		}
	} else {
		s.Transcribe("No cards are discarded.\n")
	}
	s.Discard.Push(toDiscard...)
	// Return other cards to draw pile and reshuffle
	s.Deck.Push(backToDraw...)
	s.Transcribe(fmt.Sprintf("%d cards return to the deck, and it is reshuffled.", len(backToDraw)))
	cards = SelectShuffle(s, s.Deck)
	s.Deck.Reorder(cards)
}

/*
 * Late War
 * --------
 */

func PlayIranianHostageCrisis(s *State, player Aff) {
	/* Remove all US Influence and add 2 USSR Influence to Iran. This card’s
	   Event requires the US to discard 2 cards, instead of 1 card, if the “#92 –
	   Terrorism” Event is played. */
	iran := s.Countries[Iran]
	zeroInf(s, iran, USA)
	plusInf(s, iran, SOV, 2)
	s.Event(IranianHostageCrisis, player)
}

func PlayTheIronLady(s *State, player Aff) {
	/* Add 1 USSR Influence to Argentina and remove all USSR Influence from the
	   United Kingdom. The US receives 1 VP. This Event prevents the “#7 –
	   Socialist Governments” card from being played as an Event. */
	plusInf(s, s.Countries[Argentina], SOV, 1)
	zeroInf(s, s.Countries[UK], SOV)
	s.GainVP(USA, 1)
	s.Event(TheIronLady, player)
}

func PlayReaganBombsLibya(s *State, player Aff) {
	/* The US receives 1 VP for every 2 USSR Influence in Libya. */
	s.Transcribe("The US gains VP for each 2 USSR influence in Libya.")
	s.GainVP(USA, s.Countries[Libya].Inf[SOV]/2)
}

func PlayStarWars(s *State, player Aff) {
	/* If the US is ahead on the Space Race Track, the US player uses this Event
	   to look through the discard pile, pick any 1 non-scoring card and play it
	   immediately as an Event. */
	if s.SpaceRace[USA] <= s.SpaceRace[SOV] {
		s.Transcribe("The US is not ahead on the Space Race.")
		return
	}
	// Limit choice to playable events
	canPlayEvent := func(c Card) bool {
		return !c.Prevented(s.Game)
	}
	card := SelectDiscarded(s, player,
		canPlayEvent,
		CardBlacklist(AsiaScoring, EuropeScoring,
			MiddleEastScoring, CentralAmericaScoring, SouthAmericaScoring,
			SoutheastAsiaScoring, AfricaScoring))
	s.Discard.Remove(card)
	s.Transcribe(fmt.Sprintf("%s picks %s from the discard pile.", player, card))
	PlayEvent(s, player, card)
}

func PlayNorthSeaOil(s *State, player Aff) {
	/* The US may play 8 cards (in 8 action rounds) for this turn only. This
	   Event prevents the “#61 – OPEC” card from being played as an Event. */
	// Turn event handles the 8 action rounds, permanent event handles
	// preventing OPEC
	// Don't message the turn event.
	s.TurnEvents[NorthSeaOil] = player
	s.Event(NorthSeaOil, player)
	s.Transcribe("The US may play 8 cards for this turn only.")
}

func PlayTheReformer(s *State, player Aff) {
	/* Add 4 USSR Influence to Europe (adding no more than 2 Influence per
	   country). If the USSR is ahead of the US in VP, 6 Influence may be added to
	   Europe instead. The USSR may no longer make Coup Attempts in Europe. */
	s.Event(TheReformer, player)
	n := 4
	if s.VP < 0 {
		n = 6
	}
	SelectInfluence(s, player, fmt.Sprintf("Add %d influence in Europe, no more than 2 per country", n),
		PlusInf(SOV, 1), n,
		MaxPerCountry(2), InRegion(Europe))
}

func PlayMarineBarracksBombing(s *State, player Aff) {
	/* Remove all US Influence in Lebanon and remove a total of 2 US Influence
	   from any countries in the Middle East. */
	zeroInf(s, s.Countries[Lebanon], USA)
	SelectInfluence(s, player, "Remove 2 US influence from the Middle East",
		LessInf(USA, 1), 2,
		InRegion(MiddleEast), HasInfluence(USA))
}

func PlaySovietsShootDownKAL007(s *State, player Aff) {
	/* Degrade the DEFCON level by 1 and the US receives 2 VP. The US may place
	   influence or make Realignment rolls, using this card, if South Korea is US
	   controlled. */
	s.DegradeDefcon(1)
	s.GainVP(USA, 2)
	if s.Countries[SKorea].Controlled() == USA {
		s.Transcribe(fmt.Sprintf("The US may place influence or make realignment rolls with %s.", Cards[SovietsShootDownKAL007]))
		ConductOps(s, player, Cards[SovietsShootDownKAL007], INFLUENCE, REALIGN)
	} else {
		s.Transcribe("South Korea is not US controlled.")
	}
}

func PlayGlasnost(s *State, player Aff) {
	/* Improve the DEFCON level by 1 and the USSR receives 2 VP. The USSR may
	   make Realignment rolls or add Influence, using this card, if the “#87 – The
	   Reformer” Event has already been played. */
	s.ImproveDefcon(1)
	s.GainVP(SOV, 2)
	if s.Effect(TheReformer) {
		s.Transcribe(fmt.Sprintf("The USSR may place influence or make realignment rolls with %s.", Cards[Glasnost]))
		ConductOps(s, player, Cards[Glasnost], REALIGN, INFLUENCE)
	} else {
		s.Transcribe("The Reformer Event has not been played.")
	}
}

func PlayOrtegaElectedInNicaragua(s *State, player Aff) {
	/* Remove all US Influence from Nicaragua. The USSR may make a free Coup
	   Attempt, using this card’s Operations value, in a country adjacent to
	   Nicaragua. */
	nicaragua := s.Countries[Nicaragua]
	zeroInf(s, nicaragua, USA)
	adjToNicaragua := func(c *Country) error {
		for _, neighbor := range nicaragua.AdjCountries {
			if neighbor.Id == c.Id {
				return nil
			}
		}
		return fmt.Errorf("%s is not adjacent to Nicaragua", c.Name)
	}
	OpCoup(s, player, Cards[OrtegaElectedInNicaragua], true,
		adjToNicaragua)
}

func PlayTerrorism(s *State, player Aff) {
	/* The player’s opponent must randomly discard 1 card from their hand. If
	   the “#82 – Iranian Hostage Crisis” Event has already been played, a US
	   player (if applicable) must randomly discard 2 cards from their hand. */
	opp := player.Opp()
	if len(s.Hands[opp].Cards) == 0 {
		s.Transcribe(fmt.Sprintf("%s has no cards in their hand.", opp))
		return
	}
	card := SelectRandomCard(s, opp)
	s.Transcribe(fmt.Sprintf("%s is discarded from %s hand.", card, opp))
	s.Hands[opp].Remove(card)
	if opp == USA && s.Effect(IranianHostageCrisis) {
		if len(s.Hands[opp].Cards) == 0 {
			return
		}
		card := SelectRandomCard(s, opp)
		s.Transcribe(fmt.Sprintf("%s is additionally discarded from %s hand (Iranian Hostage Crisis).", card, opp))
		s.Hands[opp].Remove(card)
	}
}

func PlayIranContraScandal(s *State, player Aff) {
	/* All US Realignment rolls, for the remainder of this turn, receive -1 to
	   their die roll. */
	s.TurnEvent(IranContraScandal, player)
}

func PlayChernobyl(s *State, player Aff) {
	/* The US must designate a single Region (Europe, Asia, etc.) that, for the
	   remainder of the turn, the USSR cannot add Influence to using Operations
	   points. */
	region := SelectRegion(s, player, "Choose a region where USSR is blocked from influencing for the turn")
	// Don't use the generic event message here
	s.TurnEvents[Chernobyl] = player
	s.ChernobylRegion = region
	s.Transcribe(fmt.Sprintf("The USSR may not add Influence using Operations points to %s for the remainder of the turn.", region))
}

func PlayLatinAmericanDebtCrisis(s *State, player Aff) {
	/* The US must immediately discard a card with an Operations value of 3 or
	   more or the USSR may double the amount of USSR Influence in 2 countries in
	   South America. */
	enoughOps := ExceedsOps(2, s, USA)
	if hasInHand(s, USA, enoughOps) &&
		"discard" == SelectChoice(s, USA,
			"Discard a card with >=3 Ops, or double USSR influence in two SAM countries?",
			"discard", "whatever") {
		card := SelectCard(s, USA, CardBlacklist(TheChinaCard), enoughOps)
		s.Hands[USA].Remove(card)
		s.Transcribe(fmt.Sprintf("%s discarded for Latin American Debt Crisis.", card))
		s.Discard.Push(card)
	} else {
		SelectInfluence(s, player, "Double USSR influence in 2 countries in South America",
			DoubleInf(SOV), 2,
			InRegion(SouthAmerica), MaxPerCountry(1), HasInfluence(SOV))
	}
}

func PlayTearDownThisWall(s *State, player Aff) {
	/* Add 3 US Influence to East Germany. The US may make free Coup Attempts or
	   Realignment rolls in Europe using the Operations value of this card. This
	   Event prevents / cancels the effect(s) of the “#55 – Willy Brandt” Event. */
	plusInf(s, s.Countries[EGermany], USA, 3)
	s.Transcribe(fmt.Sprintf("The US may make free coup attempts or realignment rolls with %s.", Cards[TearDownThisWall]))
	ConductOpsFree(s, player, Cards[TearDownThisWall], COUP, REALIGN)
	s.Event(TearDownThisWall, player)
	s.Cancel(WillyBrandt)
}

func PlayAnEvilEmpire(s *State, player Aff) {
	/* The US receives 1 VP. This Event prevents / cancels the effect(s) of the
	   “#59 – Flower Power” Event. */
	s.GainVP(USA, 1)
	s.Event(AnEvilEmpire, player)
	s.Cancel(FlowerPower)
}

func PlayAldrichAmesRemix(s *State, player Aff) {
	/* The US reveals their hand of cards, face-up, for the remainder of the
	   turn and the USSR discards a card from the US hand. */
	ShowHand(s, USA, SOV)
	card := selectCardFrom(s, SOV, s.Hands[USA].Cards, false)
	s.Hands[USA].Remove(card)
	s.Transcribe(fmt.Sprintf("%s discarded from US hand for Aldrich Ames Remix.", card))
	s.Discard.Push(card)
}

func PlayPershingIIDeployed(s *State, player Aff) {
	/* The USSR receives 1 VP. Remove 1 US Influence from any 3 countries in
	   Western Europe. */
	s.GainVP(SOV, 1)
	SelectInfluence(s, player, "Remove 1 US Influence from any 3 countries in W Europe",
		LessInf(USA, 1), 3,
		MaxPerCountry(1), InRegion(WestEurope), HasInfluence(USA))
}

func PlayWargames(s *State, player Aff) {
	/* If the DEFCON level is 2, the player may immediately end the game after
	   giving their opponent 6 VP. How about a nice game of chess? */
	if "yes" == SelectChoice(s, player, "Give opponent 6 VP and end the game?", "yes", "no") {
		s.Transcribe(fmt.Sprintf("%s gives their opponent 6 VP and ends the game.", player))
		s.GainVP(player.Opp(), 6)
		// XXX: game end, writeme
	} else {
		s.Transcribe(fmt.Sprintf("%s chooses to not end the game.", player))
	}
}

func PlaySolidarity(s *State, player Aff) {
	/* Add 3 US Influence to Poland. This card requires prior play of the “#68 –
	   John Paul II Elected Pope” Event in order to be played as an Event. */
	plusInf(s, s.Countries[Poland], USA, 3)
}

func PlayIranIraqWar(s *State, player Aff) {
	/* Iran invades Iraq or vice versa (player’s choice). Roll a die and
	   subtract (-1) from the die roll for every enemy controlled country adjacent
	   to the target of the invasion (Iran or Iraq). On a modified die roll of 4-6,
	   the player receives 2 VP and replaces all the opponent’s Influence in the
	   target country with their Influence. The player adds 2 to its Military
	   Operations Track. */
	c := SelectCountry(s, player, "Choose who gets invaded", Iraq, Iran)
	s.Transcribe(fmt.Sprintf("%s will be invaded.", c))
	s.MilOps[player] += 2
	roll := SelectRoll(s)
	mod := c.NumControlledNeighbors(player.Opp())
	if mod > 0 {
		s.Transcribe(fmt.Sprintf("%s rolls %d.", player, roll))
	} else {
		s.Transcribe(fmt.Sprintf("%s rolls %d -%d (%s controlled adjacent).", player, roll, mod, player.Opp()))
	}
	switch roll - mod {
	case 4, 5, 6:
		s.Transcribe("Iran-Iraq War succeeds.")
		s.GainVP(player, 2)
		plusInf(s, c, player, c.Inf[player.Opp()])
		zeroInf(s, c, player.Opp())
	default:
		s.Transcribe("Iran-Iraq War fails.")
	}
}

func PlayYuriAndSamantha(s *State, player Aff) {
	/* The USSR receives 1 VP for each US Coup Attempt performed during the
	   remainder of the Turn. */
	s.TurnEvent(YuriAndSamantha, player)
}

func PlayAWACSSaleToSaudis(s *State, player Aff) {
	/* Add 2 US Influence to Saudi Arabia. This Event prevents the “#56 – Muslim
	   Revolution” card from being played as an Event. */
	plusInf(s, s.Countries[SaudiArabia], USA, 2)
	s.Event(AWACSSaleToSaudis, player)
}
