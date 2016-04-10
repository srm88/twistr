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
	score(s, player, Asia)
}

func PlayEuropeScoring(s *State, player Aff) {
	/* Presence: 3; Domination: 7; Control: Automatic Victory; +1 VP per
	controlled Battleground country in Region; +1 VP per country controlled that
	is adjacent to enemy superpower; MAY NOT BE HELD!  */
	score(s, player, Europe)
}

func PlayMiddleEastScoring(s *State, player Aff) {
	/* Presence: 3; Domination: 5; Control: 7; +1 VP per controlled Battleground
	country in Region; MAY NOT BE HELD!  */
	score(s, player, MiddleEast)
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
	card := SelectRandomCard(s, SOV)
	s.Hands[SOV].Remove(card)
	if card.Aff == USA {
		PlayEvent(s, USA, card)
	} else {
		s.Discard.Push(card)
	}
}

func PlaySocialistGovernments(s *State, player Aff) {
	/* Remove a total of 3 US Influence from any countries in Western Europe
	   (removing no more than 2 Influence per country). This Event cannot be used
	   after the “#83 – The Iron Lady” Event has been played.  */
	// XXX: doesn't check that the user isn't removing more influence than
	// a country has
	cs := SelectInfluenceForce(s, player, func() ([]*Country, error) {
		return SelectNInfluenceCheck(s, player,
			"Remove 3 US influence (no more than 2 per country)", 3,
			MaxPerCountry(2), InRegion(WestEurope))
	})
	RemoveInfluence(s, USA, cs)
}

func PlayFidel(s *State, player Aff) {
	/* Remove all US Influence from Cuba. USSR adds sufficient Influence in Cuba
	   for Control.  */
	cuba := s.Countries[Cuba]
	cuba.Inf[USA] = 0
	cuba.Inf[SOV] = cuba.Stability
}

func PlayVietnamRevolts(s *State, player Aff) {
	/* Add 2 USSR Influence to Vietnam. For the remainder of the turn, the USSR
	   receives +1 Operations to the Operations value of a card that uses all its
	   Operations in Southeast Asia.  */
	s.Events[VietnamRevolts] = player
	s.Countries[Vietnam].Inf[SOV] += 2
}

func PlayBlockade(s *State, player Aff) {
	/* Unless the US immediately discards a card with an Operations value of 3 or
	   more, remove all US Influence from West Germany.  */
	choice := SelectChoice(s, USA, "Discard a card with >=3 Ops, or remove all influence from West Germany?", "discard", "remove")
	switch choice {
	case "discard":
		card := SelectCard(s, USA, CardBlacklist(TheChinaCard), ExceedsOps(2))
		s.Discard.Push(card)
	case "remove":
		s.Countries[WGermany].Inf[USA] = 0
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
	if (roll - mod) > 3 {
		s.GainVP(SOV, 2)
		skorea.Inf[SOV] += skorea.Inf[USA]
		skorea.Inf[USA] = 0
	}
}

func PlayRomanianAbdication(s *State, player Aff) {
	/* Remove all US Influence from Romania. The USSR adds sufficient Influence
	   to Romania for Control.  */
	romania := s.Countries[Romania]
	romania.Inf[USA] = 0
	romania.Inf[SOV] = romania.Stability
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
	mod := israel.NumControlledNeighbors(USA)
	if israel.Controlled() == USA {
		mod += 1
	}
	if (roll - mod) > 3 {
		s.GainVP(SOV, 2)
		israel.Inf[SOV] += israel.Inf[USA]
		israel.Inf[USA] = 0
	}
}

func PlayComecon(s *State, player Aff) {
	/* Add 1 USSR Influence to each of 4 non-US controlled countries of Eastern
	   Europe.  */
	cs := SelectInfluenceForce(s, player, func() ([]*Country, error) {
		return SelectNInfluenceCheck(s, player,
			"Choose 4 non-US controlled countries", 4,
			MaxPerCountry(1), InRegion(EastEurope), NotControlledBy(USA))
	})
	PlaceInfluence(s, SOV, cs)
}

func PlayNasser(s *State, player Aff) {
	/* Add 2 USSR Influence to Egypt. The US removes half, rounded up, of its
	   Influence from Egypt.  */
	egypt := s.Countries[Egypt]
	egypt.Inf[SOV] += 2
	loss := egypt.Inf[USA] / 2
	egypt.Inf[USA] -= loss
}

func PlayWarsawPactFormed(s *State, player Aff) {
	/* Remove all US influence from 4 countries in Eastern Europe or add 5 USSR
	   Influence to any countries in Eastern Europe (adding no more than 2
	   Influence per country). This Event allows the “#21 – NATO” card to be
	   played as an Event.  */

	s.Events[WarsawPactFormed] = player
	choice := SelectChoice(s, player, "Remove US influence or add USSR influence?", "remove", "add")
	switch choice {
	case "remove":
		cs := SelectInfluenceForce(s, player, func() ([]*Country, error) {
			return SelectNInfluenceCheck(s, player,
				"4 countries to lose all US influence", 4,
				MaxPerCountry(1), InRegion(EastEurope))
		})
		for _, c := range cs {
			c.Inf[USA] = 0
		}
	case "add":
		cs := SelectInfluenceForce(s, player, func() ([]*Country, error) {
			return SelectNInfluenceCheck(s, player, "5 influence", 5,
				MaxPerCountry(2), InRegion(EastEurope))
		})
		PlaceInfluence(s, SOV, cs)
	}
}

func PlayDeGaulleLeadsFrance(s *State, player Aff) {
	/* Remove 2 US Influence from France and add 1 USSR Influence to France. This
	   Event cancels the effect(s) of the “#21 – NATO” Event for France only.  */
	s.Events[DeGaulleLeadsFrance] = player
	france := s.Countries[France]
	france.Inf[USA] = Max(0, france.Inf[USA]-2)
	france.Inf[SOV] += 1
}

func PlayCapturedNaziScientist(s *State, player Aff) {
	/* Move the Space Race Marker ahead by 1 space.  */
	box, _ := nextSRBox(s, player)
	box.Enter(s, player)
}

func PlayTrumanDoctrine(s *State, player Aff) {
	/* Remove all USSR Influence from a single uncontrolled country in Europe.  */
	cs := SelectInfluenceForce(s, player, func() ([]*Country, error) {
		return SelectNInfluenceCheck(s, player, "1 country", 1,
			InRegion(Europe), ControlledBy(NEU))
	})
	cs[0].Inf[SOV] = 0
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
		rolls := [2]int{0, 0}
		tied := true
		for tied {
			rolls[USA] = SelectRoll(s)
			rolls[SOV] = SelectRoll(s)
			rolls[player] += 2
			tied = rolls[USA] == rolls[SOV]
		}
		switch {
		case rolls[USA] > rolls[SOV]:
			s.GainVP(USA, 2)
		case rolls[SOV] > rolls[USA]:
			s.GainVP(SOV, 2)
		}
	case "boycott":
		s.DegradeDefcon(1)
		ConductOps(s, player, PseudoCard(4))
	}
}

func PlayNATO(s *State, player Aff) {
	/* The USSR cannot make Coup Attempts or Realignment rolls against any US
	   controlled countries in Europe. US controlled countries in Europe cannot
	   be attacked by play of the “#36 – Brush War” Event. This card requires
	   prior play of either the “#16 – Warsaw Pact Formed” or “#23 – Marshall
	   Plan” Event(s) in order to be played as an Event.  */
	s.Events[NATO] = player
}

func PlayIndependentReds(s *State, player Aff) {
	/* Add US Influence to either Yugoslavia, Romania, Bulgaria, Hungary, or
	   Czechoslovakia so that it equals the USSR Influence in that country.  */
	c := SelectCountry(s, player, "Choose a country to match USSR influence",
		Yugoslavia, Romania, Bulgaria,
		Hungary, Czechoslovakia)
	c.Inf[USA] = Max(c.Inf[USA], c.Inf[SOV])
}

func PlayMarshallPlan(s *State, player Aff) {
	/* Add 1 US Influence to each of any 7 non-USSR controlled countries in
	   Western Europe. This Event allows the “#21 – NATO” card to be played as an
	   Event.  */
	cs := SelectInfluenceForce(s, player, func() ([]*Country, error) {
		return SelectNInfluenceCheck(s, player,
			"Choose 7 non-USSR controlled countries", 7,
			MaxPerCountry(1), InRegion(WestEurope), NotControlledBy(SOV))
	})
	PlaceInfluence(s, USA, cs)
	s.Events[MarshallPlan] = player
}

func PlayIndoPakistaniWar(s *State, player Aff) {
	/* India invades Pakistan or vice versa (player’s choice). Roll a die and
	   subtract (-1) from the die roll for every enemy controlled country
	   adjacent to the target of the invasion (India or Pakistan). On a modified
	   die roll of 4-6, the player receives 2 VP and replaces all the opponent’s
	   Influence in the target country with their Influence. The player adds 2 to
	   its Military Operations Track.  */
	c := SelectCountry(s, player, "Choose who gets invaded", India, Pakistan)
	s.MilOps[SOV] += 2
	roll := SelectRoll(s)
	mod := c.NumControlledNeighbors(player.Opp())
	if (roll - mod) > 3 {
		s.GainVP(player, 2)
		c.Inf[player] += c.Inf[player.Opp()]
		c.Inf[player.Opp()] = 0
	}
}

func PlayContainment(s *State, player Aff) {
	/* All Operations cards played by the US, for the remainder of this turn,
	   receive +1 to their Operations value (to a maximum of 4 Operations per
	   card).  */
	// XXX turn-duration events #14
	s.Events[Containment] = player
}

func PlayCIACreated(s *State, player Aff) {
	/* The USSR reveals their hand of cards for this turn. The US may use the
	   Operations value of this card to conduct Operations.  */
	ShowHand(s, SOV, USA)
	ConductOps(s, player, PseudoCard(1))
}

func PlayUSJapanMutualDefensePact(s *State, player Aff) {
	/* The US adds sufficient Influence to Japan for Control. The USSR cannot
	   make Coup Attempts or Realignment rolls against Japan.  */
	japan := s.Countries[Japan]
	toControl := japan.Stability + japan.Inf[SOV]
	japan.Inf[USA] = Max(japan.Inf[USA], toControl)
	s.Events[USJapanMutualDefensePact] = player
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
	cs := SelectInfluenceForce(s, player, func() ([]*Country, error) {
		return SelectNInfluenceCheck(s, player,
			"Remove 4 from France, UK, Israel", 4,
			MaxPerCountry(2), franceIsraelOrUK)
	})
	RemoveInfluence(s, USA, cs)
}

func PlayEastEuropeanUnrest(s *State, player Aff) {
	/* Early or Mid War: Remove 1 USSR Influence from 3 countries in Eastern
	   Europe. Late War: Remove 2 USSR Influence from 3 countries in Eastern
	   Europe.  */
	cs := SelectInfluenceForce(s, player, func() ([]*Country, error) {
		return SelectNInfluenceCheck(s, player,
			"Choose 3 countries in E Europe", 3,
			MaxPerCountry(1), InRegion(EastEurope))
	})
	RemoveInfluence(s, SOV, cs)
	if s.Era() == Late {
		RemoveInfluence(s, SOV, cs)
	}
}

func PlayDecolonization(s *State, player Aff) {
	/* Add 1 USSR Influence to each of any 4 countries in Africa and/or Southeast
	   Asia.  */
	cs := SelectInfluenceForce(s, player, func() ([]*Country, error) {
		return SelectNInfluenceCheck(s, player,
			"Choose 4 countries in Africa or SE Asia", 4,
			MaxPerCountry(1), InRegion(Africa, SoutheastAsia))
	})
	PlaceInfluence(s, SOV, cs)
}

func PlayRedScarePurge(s *State, player Aff) {
	/* All Operations cards played by the opponent, for the remainder of this
	   turn, receive -1 to their Operations value (to a minimum value of 1
	   Operations point).  */
	// XXX turn-duration events #14
	s.Events[RedScarePurge] = player
}

func PlayUNIntervention(s *State, player Aff) {
	/* Play this card simultaneously with a card containing an opponent’s
	   associated Event. The opponent’s associated Event is canceled but you may
	   use the Operations value of the opponent’s card to conduct Operations.
	   This Event cannot be played during the Headline Phase.  */
	// XXX: opponent's event
	card := SelectCard(s, player)
	ConductOps(s, player, card)
	s.Discard.Push(card)
}

func PlayDeStalinization(s *State, player Aff) {
	/* The USSR may reallocate up to a total of 4 Influence from one or more
	   countries to any non-US controlled countries (adding no more than 2
	   Influence per country).  */
	removed := make(map[CountryId]int)
	enoughSovInf := func(c *Country) error {
		removed[c.Id] += 1
		if removed[c.Id] > c.Inf[SOV] {
			return fmt.Errorf("You only have %d influence in %s", c.Inf[SOV], c)
		}
		return nil
	}
	from := SelectInfluenceForce(s, player, func() ([]*Country, error) {
		return SelectNInfluenceCheck(s, player,
			"Choose 4 influence to relocate", 4,
			enoughSovInf)
	})
	to := SelectInfluenceForce(s, player, func() ([]*Country, error) {
		return SelectNInfluenceCheck(s, player,
			"Relocate 4 influence to non-US controlled countries", 4,
			MaxPerCountry(2), NotControlledBy(USA))
	})
	RemoveInfluence(s, SOV, from)
	PlaceInfluence(s, SOV, to)
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
	s.Events[FormosanResolution] = player
}

func PlayDefectors(s *State, player Aff) {
	/* The US may play this card during the Headline Phase in order to cancel the
	   USSR Headline Event (including a scoring card). The canceled card is
	   placed into the discard pile. If this card is played by the USSR during
	   its action round, the US gains 1 VP.  */
	// XXX headline special case!
	s.GainVP(USA, 1)
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
	s.Message(player, fmt.Sprintf("%s scoring cards: %s\n", USA, strings.Join(scoringCards, ", ")))
	if len(scoringCards) == 0 {
		return
	}
	cs := SelectInfluenceForce(s, player, func() ([]*Country, error) {
		return SelectNInfluenceCheck(s, player,
			"Place one influence in one of the regions", 1,
			InRegion(regions...))
	})
	PlaceInfluence(s, SOV, cs)
}

func PlaySpecialRelationship(s *State, player Aff) {
	/* Add 1 US Influence to a single country adjacent to the U.K. if the U.K.
	   is US-controlled but NATO is not in effect. Add 2 US Influence to a single
	   country in Western Europe, and the US gains 2 VP, if the U.K. is
	   US-controlled and NATO is in effect. */
	ukControlled := s.Countries[UK].Controlled() == USA
	switch {
	case ukControlled && s.Effect(NATO):
		cs := SelectInfluenceForce(s, player, func() ([]*Country, error) {
			return SelectNInfluenceCheck(s, player,
				"Choose a country in W Europe", 1,
				InRegion(WestEurope))
		})
		s.Countries[cs[0].Id].Inf[USA] += 2
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
		cs := SelectInfluenceForce(s, player, func() ([]*Country, error) {
			return SelectNInfluenceCheck(s, player,
				"Choose a country adjacent to the UK", 1,
				nextToUK)
		})
		PlaceInfluence(s, USA, cs)
	}
}

func PlayNORAD(s *State, player Aff) {
	/* Add 1 US Influence to a single country containing US Influence, at the
	   end of each Action Round, if Canada is US-controlled and the DEFCON level
	   moved to 2 during that Action Round. This Event is canceled by the “#42 –
	   Quagmire” Event. */
	s.Events[NORAD] = player
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

	cl := SelectInfluenceForce(s, player, func() ([]*Country, error) {
		return SelectNInfluenceCheck(s, player, "1 country", 1,
			stabLTE)
	})
	c := cl[0]
	s.MilOps[player] += 3
	roll := SelectRoll(s)
	mod := c.NumControlledNeighbors(player.Opp())
	if (roll - mod) > 2 {
		s.GainVP(player, 1)
		c.Inf[player] += c.Inf[player.Opp()]
		c.Inf[player.Opp()] = 0
	}
}

func PlayCentralAmericaScoring(s *State, player Aff) {
	/* Presence: 1; Domination: 3; Control: 5; +1 VP per controlled Battleground
	   country in Region; +1 VP per country controlled that is adjacent to enemy
	   superpower; MAY NOT BE HELD! */
	score(s, player, CentralAmerica)
}

func PlaySoutheastAsiaScoring(s *State, player Aff) {
	/* 1 VP each for Control of Burma, Cambodia/Laos, Vietnam, Malaysia,
	   Indonesia and the Philippines. 2 VP for Control of Thailand; MAY NOT BE
	   HELD! */
	vp := 0
	for _, c := range []CountryId{Burma, LaosCambodia, Vietnam, Malaysia, Indonesia, Philippines} {
		if s.Countries[c].Controlled() == player {
			vp++
		}
		if s.Countries[Thailand].Controlled() == player {
			vp += 2
		}
	}
	s.GainVP(player, vp)
}

func PlayArmsRace(s *State, player Aff) {
	/* Compare each player’s value on the Military Operations Track. If the
	   phasing player has a higher value than their opponent on the Military
	   Operations Track, that player receives 1 VP. If the phasing player has a
	   higher value than their opponent, and has met the “required” amount, on the
	   Military Operations Track, that player receives 3 VP instead. */
	playerMilOps := s.MilOps[player]
	oppMilOps := s.MilOps[player.Opp()]
	// mil ops has to be greater than or equal to defcon to meet "required" amount
	switch {
	case playerMilOps > oppMilOps && playerMilOps >= s.Defcon:
		s.GainVP(player, 3)
	case playerMilOps > oppMilOps:
		s.GainVP(player, 1)
	}
}

func PlayCubanMissileCrisis(s *State, player Aff) {
	/* Set the DEFCON level to 2. Any Coup Attempts by your opponent, for the
	   remainder of this turn, will result in Global Thermonuclear War. Your
	   opponent will lose the game. This card’s Event may be canceled, at any time,
	   if the USSR removes 2 Influence from Cuba or the US removes 2 Influence from
	   West Germany or Turkey. */
	s.Defcon = 2
	s.TurnEvents[CubanMissileCrisis] = player
}

func PlayNuclearSubs(s *State, player Aff) {
	/* US Operations used for Coup Attempts in Battleground countries, for the
	   remainder of this turn, do not degrade the DEFCON level. This card’s Event
	   does not apply to any Event that would affect the DEFCON level (ex. the
	   “#40 – Cuban Missile Crisis” Event). */
	s.TurnEvents[NuclearSubs] = player
}

func PlayQuagmire(s *State, player Aff) {
	/* On the US’s next action round, it must discard an Operations card with a
	   value of 2 or more and roll 1-4 on a die to cancel this Event. Repeat this
	   Event for each US action round until the US successfully rolls 1-4 on a die.
	   If the US is unable to discard an Operations card, it must play all of its
	   scoring cards and then skip each action round for the rest of the turn. This
	   Event cancels the effect(s) of the “#106 – NORAD” Event (if applicable). */
	s.TurnEvents[Quagmire] = player
	s.Cancel(NORAD)
}

func PlaySALTNegotiations(s *State, player Aff) {
	/* Improve the DEFCON level by 2. For the remainder of the turn, both
	   players receive -1 to all Coup Attempt rolls. The player of this card’s
	   Event may look through the discard pile, pick any 1 non-scoring card, reveal
	   it to their opponent and then place the drawn card into their hand. */
	s.ImproveDefcon(2)
	s.TurnEvents[SALTNegotiations] = player
	choice := s.Solicit(player, "Choose a card to show to your opponent and add to your hand?", []string{"yes", "no"})
	notScoring := func(c Card) bool {
		return !c.Scoring()
	}
	switch choice {
	case "yes":
		selected := SelectDiscarded(s, player, notScoring)
		ShowCard(s, selected, player.Opp())
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
	s.TurnEvents[BearTrap] = player
}

func PlaySummit(s *State, player Aff) {
	/* Both players roll a die. Each player receives +1 to the die roll for each
	   Region (Europe, Asia, etc.) they Dominate or Control. The player with the
	   highest modified die roll receives 2 VP and may degrade or improve the
	   DEFCON level by 1 (do not reroll ties). */
	regions := []Region{CentralAmerica, SouthAmerica, Europe, MiddleEast, Africa, Asia}
	playerRoll := SelectRoll(s)
	oppRoll := SelectRoll(s)
	for _, region := range regions {
		sr := ScoreRegion(s, region)
		switch {
		case sr.Levels[player] == Domination || sr.Levels[player] == Control:
			playerRoll++
		case sr.Levels[player.Opp()] == Domination || sr.Levels[player.Opp()] == Control:
			oppRoll++
		}
	}
	if playerRoll == oppRoll {
		return
	}

	var winner Aff
	if playerRoll > oppRoll {
		winner = player
	} else {
		winner = player.Opp()
	}

	s.GainVP(winner, 2)
	choice := s.Solicit(winner, "Degrade or improve DEFCON by one level?", []string{"improve", "degrade", "neither"})
	switch choice {
	case "leave it":
		return
	case "improve":
		s.ImproveDefcon(1)
	case "degrade":
		s.DegradeDefcon(1)
	}
}

func PlayHowILearnedToStopWorrying(s *State, player Aff) {
	/* Set the DEFCON level to any level desired (1-5). The player adds 5 to its
	   Military Operations Track. */
	choice := s.Solicit(player, "Set DEFCON.", []string{"1", "2", "3", "4", "5"})
	newDefcon, _ := strconv.Atoi(choice)
	s.MilOps[player] = 5
	switch {
	case newDefcon == s.Defcon:
		return
	case newDefcon > s.Defcon:
		s.ImproveDefcon(newDefcon - s.Defcon)
	case newDefcon < s.Defcon:
		s.DegradeDefcon(s.Defcon - newDefcon)
	}
}

func PlayJunta(s *State, player Aff) {
	/* Add 2 Influence to a single country in Central or South America. The
	   player may make free Coup Attempts or Realignment rolls in either Central or
	   South America using the Operations value of this card. */
	centralOrSouth := append(SouthAmerica.Countries, CentralAmerica.Countries...)
	c := SelectCountry(s, player, "Choose a country in Central or South America",
		centralOrSouth...)
	c.Inf[player] += 2
	choice := s.Solicit(player, "Do you want to coup, realign or do nothing in Central/South America?", []string{"coup, realign, nothing"})
	switch choice {
	case "nothing":
		return
	case "coup":
		couped := DoFreeCoup(s, player, Cards[Junta], centralOrSouth)
		fmt.Println("coup isn't implemented, but like, ya did it. ya couped %s", couped)
	case "realign":
        fmt.Println("This will be realign")
	}
}

func PlayKitchenDebates(s *State, player Aff) {
	/* If the US controls more Battleground countries than the USSR, the US
	   player uses this Event to poke their opponent in the chest and receive 2 VP!
	*/
	usaBG := 0
	sovBG := 0
	for _, c := range Countries {
		if c.Battleground {
			if c.Controlled() == SOV {
				sovBG++
			} else if c.Controlled() == USA {
				usaBG++
			}
		}
	}
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
	s.TurnEvents[MissileEnvy] = player
	maxOps := 0
	for _, c := range s.Hands[player.Opp()].Cards {
		if c.Ops > maxOps {
			maxOps = c.Ops
		}
	}

	selected := SelectCard(s, player.Opp(), ExceedsOps(maxOps-1))
	if selected.Aff == player || selected.Aff == NEU {
		PlayEvent(s, player, selected)
	} else {
		ConductOps(s, player, selected)
	}

	s.Hands[player.Opp()].Push(Cards[MissileEnvy])
}

func PlayWeWillBuryYou(s *State, player Aff) {
	/* Degrade the DEFCON level by 1. Unless the #32 UN Intervention card is
	   played as an Event on the US’s next action round, the USSR receives 3 VP.  */
	s.DegradeDefcon(1)
	s.Events[WeWillBuryYou] = player
}

func PlayBrezhnevDoctrine(s *State, player Aff) {
	/* All Operations cards played by the USSR, for the remainder of this turn,
	   receive +1 to their Operations value (to a maximum of 4 Operations per
	   card). */
	s.TurnEvents[BrezhnevDoctrine] = player
}

func PlayPortugueseEmpireCrumbles(s *State, player Aff) {
	/* Add 2 USSR Influence to Angola and the SE African States. */
	Countries[Angola].Inf[SOV] += 2
	Countries[SEAfricanStates].Inf[SOV] += 2
}

func PlaySouthAfricanUnrest(s *State, player Aff) {
	/* The USSR either adds 2 Influence to South Africa or adds 1 Influence to
	   South Africa and 2 Influence to a single country adjacent to South Africa. */
	choice := s.Solicit(player, "Soviets add 2 influence to South Africa, or 1 to South Africa and two to an Adjacent country?", []string{"option 1", "option 2"})
	switch choice {
	case "option 1":
		Countries[SouthAfrica].Inf[SOV] += 2
	case "option 2":
		Countries[SouthAfrica].Inf[SOV] += 1
		adjIds := []CountryId{}
		adj := Countries[SouthAfrica].AdjCountries
		for _, c := range adj {
			adjIds = append(adjIds, c.Id)
		}
		selected := SelectCountry(s, SOV, "Which country would you like to add to?", adjIds...)
		selected.Inf[SOV] += 2
	}
}

func PlayAllende(s *State, player Aff) {
	/* Add 2 USSR Influence to Chile. */
	Countries[Chile].Inf[SOV] += 2
}

func PlayWillyBrandt(s *State, player Aff) {
	/* The USSR receives 1 VP and adds 1 Influence to West Germany. This Event
	   cancels the effect(s) of the “#21 – NATO” Event for West Germany only. This
	   Event is prevented / canceled by the “#96 – Tear Down this Wall” Event. */
	s.GainVP(SOV, 1)
	s.Events[WillyBrandt] = player
	Countries[WGermany].Inf[SOV] += 1
}

func PlayMuslimRevolution(s *State, player Aff) {
	/* Remove all US Influence from 2 of the following countries: Sudan, Iran,
	   Iraq, Egypt, Libya, Saudi Arabia, Syria, Jordan. This Event cannot be used
	   after the “#110 – AWACS Sale to Saudis” Event has been played. */
}

func PlayABMTreaty(s *State, player Aff) {
	/* Improve the DEFCON level by 1 and then conduct Operations using the
	   Operations value of this card. */
	s.ImproveDefcon(1)
	ConductOps(s, player, PseudoCard(Cards[ABMTreaty].Ops))
}

func PlayCulturalRevolution(s *State, player Aff) {
	/* If the US has the “#6 – The China Card” card, the US must give the card
	   to the USSR (face up and available to be played). If the USSR already has
	   “#6 – The China Card” card, the USSR receives 1 VP. */
	if s.ChinaCardPlayer == USA {
		s.ChinaCardPlayer = SOV
		s.ChinaCardFaceUp = true
	} else {
		s.GainVP(SOV, 1)
	}
}

func PlayFlowerPower(s *State, player Aff) {
	/* The USSR receives 2 VP for every US played “War” card (Arab-Israeli War,
	   Korean War, Brush War, Indo-Pakistani War, Iran-Iraq War), used for
	   Operations or an Event, after this card is played. This Event is prevented /
	   canceled by the “#97 – ‘An Evil Empire’” Event. */
	s.Events[FlowerPower] = player
}

func PlayU2Incident(s *State, player Aff) {
	/* The USSR receives 1 VP. If the “#32 – UN Intervention” Event is played
	   later this turn, either by the US or the USSR, the USSR receives an
	   additional 1 VP. */
	s.GainVP(SOV, 1)
	s.TurnEvents[U2Incident] = player
}

func PlayOPEC(s *State, player Aff) {
	/* The USSR receives 1 VP for Control of each of the following countries:
	   Egypt, Iran, Libya, Saudi Arabia, Iraq, Gulf States, Venezuela. This Event
	   cannot be used after the “#86 – North Sea Oil” Event has been played. */
	for _, cid := range []CountryId{Egypt, Iran, Libya, SaudiArabia, GulfStates, Venezuela} {
		if Countries[cid].Controlled() == SOV {
			s.GainVP(SOV, 1)
		}
	}
}

func PlayLoneGunman(s *State, player Aff) {
	/* The US reveals their hand of cards. The USSR may use the Operations value
	   of this card to conduct Operations. */
	ShowHand(s, USA, SOV)
	ConductOps(s, player, PseudoCard(Cards[LoneGunman].Ops))
}

func PlayColonialRearGuards(s *State, player Aff) {
	/* Add 1 US Influence to each of any 4 countries in Africa and/or Southeast
	   Asia. */
}

func PlayPanamaCanalReturned(s *State, player Aff) {
	/* Add 1 US Influence to Panama, Costa Rica and Venezuela.  */
}

func PlayCampDavidAccords(s *State, player Aff) {
	/* The US receives 1 VP and adds 1 Influence to Israel, Jordan and Egypt.
	   This Event prevents the “#13 – Arab-Israeli War” card from being played as
	   an Event. */
}

func PlayPuppetGovernments(s *State, player Aff) {
	/* The US may add 1 Influence to 3 countries that do not contain Influence
	   from either the US or USSR. */
}

func PlayGrainSalesToSoviets(s *State, player Aff) {
	/* The US randomly selects 1 card from the USSR’s hand (if available). The
	   US must either play the card or return it to the USSR. If the card is
	   returned, or the USSR has no cards, the US may use the Operations value of
	   this card to conduct Operations. */
}

func PlayJohnPaulIIElectedPope(s *State, player Aff) {
	/* Remove 2 USSR Influence from Poland and add 1 US Influence to Poland.
	   This Event allows the “#101 – Solidarity” card to be played as an Event.
	*/
}

func PlayLatinAmericanDeathSquads(s *State, player Aff) {
	/* All of the phasing player’s Coup Attempts in Central and South America,
	   for the remainder of this turn, receive +1 to their die roll. All of the
	   opponent’s Coup Attempts in Central and South America, for the remainder of
	   this turn, receive -1 to their die roll. */
}

func PlayOASFounded(s *State, player Aff) {
	/* Add a total of 2 US Influence to any countries in Central or South America. */
}

func PlayNixonPlaysTheChinaCard(s *State, player Aff) {
	/* If the USSR has the “#6 – The China Card” card, the USSR must give the
	   card to the US (face down and unavailable for immediate play). If the US
	   already has the “#6 – The China Card” card, the US receives 2 VP. */
}

func PlaySadatExpelsSoviets(s *State, player Aff) {
	/* Remove all USSR Influence from Egypt and add 1 US Influence to Egypt. */
}

func PlayShuttleDiplomacy(s *State, player Aff) {
	/* If this card’s Event is in effect, subtract (-1) a Battleground country
	   from the USSR total and then discard this card during the next scoring of
	   the Middle East or Asia (which ever comes first). */
}

func PlayTheVoiceOfAmerica(s *State, player Aff) {
	/* Remove 4 USSR Influence from any countries NOT in Europe (removing no
	   more than 2 Influence per country). */
}

func PlayLiberationTheology(s *State, player Aff) {
	/* Add a total of 3 USSR Influence to any countries in Central America
	   (adding no more than 2 Influence per country). */
}

func PlayUssuriRiverSkirmish(s *State, player Aff) {
	/* If the USSR has the “#6 – The China Card” card, the USSR must give the
	   card to the US (face up and available for play). If the US already has the
	   “#6 – The China Card” card, add a total of 4 US Influence to any countries
	   in Asia (adding no more than 2 Influence per country). */
}

func PlayAskNotWhatYourCountry(s *State, player Aff) {
	/* The US may discard up to their entire hand of cards (including scoring
	   cards) to the discard pile and draw replacements from the draw pile. The
	   number of cards to be discarded must be decided before drawing any
	   replacement cards from the draw pile. */
}

func PlayAllianceForProgress(s *State, player Aff) {
	/* The US receives 1 VP for each US controlled Battleground country in
	   Central and South America. */
}

func PlayAfricaScoring(s *State, player Aff) {
	/* Presence: 1; Domination: 4; Control: 6; +1 VP per controlled Battleground
	   country in Region; MAY NOT BE HELD! */
}

func PlayOneSmallStep(s *State, player Aff) {
	/* If you are behind on the Space Race Track, the player uses this Event to
	   move their marker 2 spaces forward on the Space Race Track. The player
	   receives VP only from the final space moved into. */
}

func PlaySouthAmericaScoring(s *State, player Aff) {
	/* Presence: 2; Domination: 5; Control: 6; +1 VP per controlled Battleground
	   country in Region; MAY NOT BE HELD! */
	score(s, player, SouthAmerica)
}

func PlayChe(s *State, player Aff) {
	/* The USSR may perform a Coup Attempt, using this card’s Operations value,
	   against a non-Battleground country in Central America, South America or
	   Africa. The USSR may perform a second Coup Attempt, against a different
	   non-Battleground country in Central America, South America or Africa, if the
	   first Coup Attempt removed any US Influence from the target country. */
	// Free coup
	targets := []CountryId{}
	targets = append(targets, SouthAmerica.Countries...)
	targets = append(targets, CentralAmerica.Countries...)
	targets = append(targets, Africa.Countries...)
	couped := DoFreeCoup(s, player, Cards[Che], targets)
	fmt.Println("placeholder %s", couped)
}

func PlayOurManInTehran(s *State, player Aff) {
	/* If the US controls at least one Middle East country, the US player uses
	   this Event to draw the top 5 cards from the draw pile. The US may discard
	   any or all of the drawn cards, after revealing the discarded card(s) to the
	   USSR player, without triggering the Event(s). Any remaining drawn cards are
	   returned to the draw pile and the draw pile is reshuffled. */
}

/*
 * Late War
 * --------
 */

func PlayIranianHostageCrisis(s *State, player Aff) {
	/* Remove all US Influence and add 2 USSR Influence to Iran. This card’s
	   Event requires the US to discard 2 cards, instead of 1 card, if the “#92 –
	   Terrorism” Event is played. */
}

func PlayTheIronLady(s *State, player Aff) {
	/* Add 1 USSR Influence to Argentina and remove all USSR Influence from the
	   United Kingdom. The US receives 1 VP. This Event prevents the “#7 –
	   Socialist Governments” card from being played as an Event. */
}

func PlayReaganBombsLibya(s *State, player Aff) {
	/* The US receives 1 VP for every 2 USSR Influence in Libya. */
}

func PlayStarWars(s *State, player Aff) {
	/* If the US is ahead on the Space Race Track, the US player uses this Event
	   to look through the discard pile, pick any 1 non-scoring card and play it
	   immediately as an Event. */
}

func PlayNorthSeaOil(s *State, player Aff) {
	/* The US may play 8 cards (in 8 action rounds) for this turn only. This
	   Event prevents the “#61 – OPEC” card from being played as an Event. */
	// Turn event handles the 8 action rounds, permanent event handles
	// preventing OPEC
}

func PlayTheReformer(s *State, player Aff) {
	/* Add 4 USSR Influence to Europe (adding no more than 2 Influence per
	   country). If the USSR is ahead of the US in VP, 6 Influence may be added to
	   Europe instead. The USSR may no longer make Coup Attempts in Europe. */
}

func PlayMarineBarracksBombing(s *State, player Aff) {
	/* Remove all US Influence in Lebanon and remove a total of 2 US Influence
	   from any countries in the Middle East. */
}

func PlaySovietsShootDownKAL007(s *State, player Aff) {
	/* Degrade the DEFCON level by 1 and the US receives 2 VP. The US may place
	   influence or make Realignment rolls, using this card, if South Korea is US
	   controlled. */
}

func PlayGlasnost(s *State, player Aff) {
	/* Improve the DEFCON level by 1 and the USSR receives 2 VP. The USSR may
	   make Realignment rolls or add Influence, using this card, if the “#87 – The
	   Reformer” Event has already been played. */
}

func PlayOrtegaElectedInNicaragua(s *State, player Aff) {
	/* Remove all US Influence from Nicaragua. The USSR may make a free Coup
	   Attempt, using this card’s Operations value, in a country adjacent to
	   Nicaragua. */
}

func PlayTerrorism(s *State, player Aff) {
	/* The player’s opponent must randomly discard 1 card from their hand. If
	   the “#82 – Iranian Hostage Crisis” Event has already been played, a US
	   player (if applicable) must randomly discard 2 cards from their hand. */
}

func PlayIranContraScandal(s *State, player Aff) {
	/* All US Realignment rolls, for the remainder of this turn, receive -1 to
	   their die roll. */
}

func PlayChernobyl(s *State, player Aff) {
	/* The US must designate a single Region (Europe, Asia, etc.) that, for the
	   remainder of the turn, the USSR cannot add Influence to using Operations
	   points. */
}

func PlayLatinAmericanDebtCrisis(s *State, player Aff) {
	/* The US must immediately discard a card with an Operations value of 3 or
	   more or the USSR may double the amount of USSR Influence in 2 countries in
	   South America. */
}

func PlayTearDownThisWall(s *State, player Aff) {
	/* Add 3 US Influence to East Germany. The US may make free Coup Attempts or
	   Realignment rolls in Europe using the Operations value of this card. This
	   Event prevents / cancels the effect(s) of the “#55 – Willy Brandt” Event. */
}

func PlayAnEvilEmpire(s *State, player Aff) {
	/* The US receives 1 VP. This Event prevents / cancels the effect(s) of the
	   “#59 – Flower Power” Event. */
}

func PlayAldrichAmesRemix(s *State, player Aff) {
	/* The US reveals their hand of cards, face-up, for the remainder of the
	   turn and the USSR discards a card from the US hand. */
}

func PlayPershingIIDeployed(s *State, player Aff) {
	/* The USSR receives 1 VP. Remove 1 US Influence from any 3 countries in
	   Western Europe. */
}

func PlayWargames(s *State, player Aff) {
	/* If the DEFCON level is 2, the player may immediately end the game after
	   giving their opponent 6 VP. How about a nice game of chess? */
}

func PlaySolidarity(s *State, player Aff) {
	/* Add 3 US Influence to Poland. This card requires prior play of the “#68 –
	   John Paul II Elected Pope” Event in order to be played as an Event. */
}

func PlayIranIraqWar(s *State, player Aff) {
	/* Iran invades Iraq or vice versa (player’s choice). Roll a die and
	   subtract (-1) from the die roll for every enemy controlled country adjacent
	   to the target of the invasion (Iran or Iraq). On a modified die roll of 4-6,
	   the player receives 2 VP and replaces all the opponent’s Influence in the
	   target country with their Influence. The player adds 2 to its Military
	   Operations Track. */
}

func PlayYuriAndSamantha(s *State, player Aff) {
	/* The USSR receives 1 VP for each US Coup Attempt performed during the
	   remainder of the Turn. */
}

func PlayAWACSSaleToSaudis(s *State, player Aff) {
	/* Add 2 US Influence to Saudi Arabia. This Event prevents the “#56 – Muslim
	   Revolution” card from being played as an Event. */
}
