package twistr

import "fmt"
import "log"

// All WIP. Maybe obliterate it.

// Realignment
func Realign(s *State, player Aff, c *Country) {
	rollUsa := SelectRoll(s, player)
	rollSov := SelectRoll(s, player)
	s.Transcribe(fmt.Sprintf("%s realigns %s.", player, c))
	realign(s, c, rollUsa, rollSov)
	s.Commit()
}

func realignMods(target Country) (modsUsa []Mod, modsSov []Mod) {
	switch {
	case target.Inf[USA] > target.Inf[SOV]:
		modsUsa = append(modsUsa, Mod{1, "US influence"})
	case target.Inf[SOV] > target.Inf[USA]:
		modsSov = append(modsSov, Mod{1, "USSR influence"})
	}
	usaAdj, sovAdj := 0, 0
	for _, neighbor := range target.AdjCountries {
		control := neighbor.Controlled()
		switch control {
		case USA:
			usaAdj += 1
		case SOV:
			sovAdj += 1
		}
	}
	if usaAdj > 0 {
		modsUsa = append(modsUsa, Mod{usaAdj, "US controlled adjacent"})
	}
	if sovAdj > 0 {
		modsSov = append(modsSov, Mod{sovAdj, "USSR controlled adjacent"})
	}
	return
}

func realign(s *State, target *Country, rollUSA, rollSOV int) {
	modsUsa, modsSov := realignMods(*target)
	if s.Effect(IranContraScandal) {
		modsUsa = append(modsUsa, Mod{-1, "Iran-Contra Scandal"})
	}
	if len(modsUsa) > 0 {
		s.Transcribe(fmt.Sprintf("US rolls %d %s.", rollUSA, ModSummary(modsUsa)))
	} else {
		s.Transcribe(fmt.Sprintf("US rolls %d.", rollUSA))
	}
	if len(modsSov) > 0 {
		s.Transcribe(fmt.Sprintf("USSR rolls %d %s.", rollSOV, ModSummary(modsSov)))
	} else {
		s.Transcribe(fmt.Sprintf("USSR rolls %d.", rollSOV))
	}
	rollUSA += TotalMod(modsUsa)
	rollSOV += TotalMod(modsSov)
	initUSA := target.Inf[USA]
	initSOV := target.Inf[SOV]
	switch {
	case rollUSA > rollSOV:
		target.Inf[SOV] -= Min((rollUSA - rollSOV), target.Inf[SOV])

	case rollSOV > rollUSA:
		target.Inf[USA] -= Min((rollSOV - rollUSA), target.Inf[USA])
	}
	if initUSA > target.Inf[USA] {
		s.Transcribe(fmt.Sprintf("%d US influence removed", initUSA-target.Inf[USA]))
	} else if initSOV > target.Inf[SOV] {
		s.Transcribe(fmt.Sprintf("%d USSR influence removed", initSOV-target.Inf[SOV]))
	} else {
		s.Transcribe("No influence removed")
	}
}

func coupMods(s *State, player Aff, target *Country) (mods []Mod) {
	if s.Effect(SALTNegotiations) {
		mods = append(mods, Mod{-1, "SALT Negotiations"})
	}
	if s.Effect(LatinAmericanDeathSquads, player) {
		mods = append(mods, Mod{1, "LatAm Death Squads"})
	}
	if s.Effect(LatinAmericanDeathSquads, player.Opp()) {
		mods = append(mods, Mod{-1, "LatAm Death Squads"})
	}
	return
}

func ComputeCardOps(s *State, player Aff, card Card, countries []*Country) int {
	return card.Ops + TotalMod(opsMods(s, player, card, countries))
}

func opsMods(s *State, player Aff, card Card, countries []*Country) (mods []Mod) {
	tmpTotal := card.Ops
	if player == SOV && s.Effect(VietnamRevolts) && AllIn(countries, SoutheastAsia) {
		mods = append(mods, Mod{1, "Vietnam Revolts"})
		tmpTotal += 1
	}
	if card.Id == TheChinaCard && AllIn(countries, Asia) {
		mods = append(mods, Mod{1, "The China Card"})
		tmpTotal += 1
	}
	// Brezhnev/containment/redscare computation is surprisingly complicated.
	// The following switch statement is comprehensive.
	brezhnev := player == SOV && s.Effect(BrezhnevDoctrine)
	containment := player == USA && s.Effect(Containment)
	redscare := s.Effect(RedScarePurge, player.Opp())
	switch {
	// Red scare will lower an op total above 4, containment/brezhnev won't help
	case redscare && tmpTotal > 4:
		mods = append(mods, Mod{-1, "Red Scare/Purge"})
	// If the total is <= 4, the two can cancel each other out regardless of total
	case redscare && containment:
		mods = append(mods, Mod{-1, "Red Scare/Purge"}, Mod{1, "Containment"})
	case redscare && brezhnev:
		mods = append(mods, Mod{-1, "Red Scare/Purge"}, Mod{1, "Brezhnev Doctrine"})
	// With no containment/brezhnev, redscare only decrs to 1 ops min
	case redscare && tmpTotal > 1:
		mods = append(mods, Mod{-1, "Red Scare/Purge"})
	case redscare:
		mods = append(mods, Mod{0, "Red Scare/Purge"})
	// Similarly, with no red scare, containment/brezhnev only incrs to 4 ops max
	case brezhnev && tmpTotal < 4:
		mods = append(mods, Mod{1, "Brezhnev Doctrine"})
	case containment && tmpTotal < 4:
		mods = append(mods, Mod{1, "Containment"})
	}
	log.Printf("Computed mods for %s playing %s: %s\n", player, card, ModSummary(mods))
	return
}

// Coup
func Coup(s *State, player Aff, card Card, c *Country, free bool) (success bool) {
	s.Transcribe(fmt.Sprintf("%s coups %s.", player, c))
	if s.Effect(CubanMissileCrisis, player.Opp()) {
		if !CancelCubanMissileCrisis(s, player) {
			s.Transcribe(fmt.Sprintf("%s perturbs the delicate balance of the Cuban missile crisis!", player))
			ThermoNuclearWar(s, player)
		}
	}
	if s.Effect(YuriAndSamantha) && player == USA {
		s.Transcribe("The USSR gains VP for Yuri And Samantha")
		s.GainVP(SOV, 1)
	}
	roll := SelectRoll(s, player)
	mods := opsMods(s, player, card, []*Country{c})
	ops := card.Ops + TotalMod(mods)
	success = coup(s, player, ops, roll, c, free)
	s.Commit()
	return
}

func coup(s *State, player Aff, ops int, roll int, target *Country, free bool) (removedInfluence bool) {
	mods := coupMods(s, player, target)
	delta := roll + TotalMod(mods) + ops - (target.Stability * 2)
	if len(mods) > 0 {
		s.Transcribe(fmt.Sprintf("Result: %d +%d (ops) %s -%d (2x stability).", roll, ops, ModSummary(mods), 2*target.Stability))
	} else {
		s.Transcribe(fmt.Sprintf("Result: %d +%d (ops) -%d (2x stability).", roll, ops, 2*target.Stability))
	}
	removedInfluence = delta > 0
	if removedInfluence {
		oppCurInf := target.Inf[player.Opp()]
		removed := Min(oppCurInf, delta)
		gained := delta - removed
		target.Inf[player] += gained
		target.Inf[player.Opp()] -= removed
		s.Transcribe(fmt.Sprintf("%s %s influence reduced by %d, now %d.", target, player.Opp(), removed, target.Inf[player.Opp()]))
		if gained > 0 {
			s.Transcribe(fmt.Sprintf("%s %s influence increased by %d, now %d.", target, player, gained, target.Inf[player]))
		}
	} else {
		s.Transcribe("No influence removed.")
	}
	if target.Battleground {
		if s.Effect(NuclearSubs) && player == USA {
			s.Transcribe("Defcon is unaffected due to Nuclear Subs")
		} else {
			s.DegradeDefcon(1)
		}
	}
	if !free {
		s.AddMilOps(player, ops)
	}
	return
}

// A country cannot be coup'd if it lacks any of the opponent's influence.
// Some permanent events also impose coup restrictions, e.g. NATO with Europe.
func CanCoup(s *State, player Aff, free bool) countryCheck {
	return func(t *Country) error {
		switch {
		case t.Inf[player.Opp()] < 1:
			return fmt.Errorf("No %s influence in %s", player.Opp(), t.Name)
		case natoProtected(s, player, t):
			return fmt.Errorf("%s protected by NATO", t.Name)
		case japanProtected(s, player, t):
			return fmt.Errorf("%s protected by US/Japan Mutual Defense Pact", t.Name)
		case s.Effect(TheReformer) && player == SOV && t.In(Europe):
			return fmt.Errorf("%s protected by The Reformer", t.Name)
		case defconProtected(s, t) && !free:
			return fmt.Errorf("%s protected by DEFCON", t.Name)
		default:
			return nil
		}
	}
}

func CanRealign(s *State, player Aff, free bool) countryCheck {
	return func(t *Country) error {
		switch {
		case natoProtected(s, player, t):
			return fmt.Errorf("%s protected by NATO", t.Name)
		case japanProtected(s, player, t):
			return fmt.Errorf("%s protected by US/Japan Mutual Defense Pact", t.Name)
		case t.Inf[player.Opp()] < 1:
			return fmt.Errorf("No %s influence in %s", player.Opp(), t.Name)
		case defconProtected(s, t) && !free:
			return fmt.Errorf("%s protected by DEFCON", t.Name)
		default:
			return nil
		}
	}
}

func defconProtected(s *State, t *Country) bool {
	// asia 3, defcon 5, not protected
	// europe 4, defcon 3, protected
	// middle east 2, defcon 2, protected
	return t.Region.Volatility >= s.Defcon
}

func natoProtected(s *State, player Aff, t *Country) bool {
	degaulled := s.Effect(DeGaulleLeadsFrance) && t.Id == France
	willyd := s.Effect(WillyBrandt) && t.Id == WGermany
	natod := (s.Effect(NATO) && player == SOV && t.In(Europe) && t.Controlled() == USA)
	return natod && !degaulled && !willyd
}

func japanProtected(s *State, player Aff, t *Country) bool {
	return s.Effect(USJapanMutualDefensePact) && t.Id == Japan && player == SOV
}

type countryChange func(*State, *Country)

func plusInf(s *State, c *Country, aff Aff, n int) {
	c.Inf[aff] += n
	s.Transcribe(fmt.Sprintf("%s influence in %s +%d, now %d.", aff, c, n, c.Inf[aff]))
}

func setInf(s *State, c *Country, aff Aff, to int) {
	added := to - c.Inf[aff]
	if added <= 0 {
		s.Transcribe(fmt.Sprintf("%s influence in %s is at %d.", aff, c, to))
		return
	}
	s.Transcribe(fmt.Sprintf("%s influence in %s +%d, now %d.", aff, c, added))
	c.Inf[aff] = to
}

func lessInf(s *State, c *Country, aff Aff, n int) {
	if c.Inf[aff] == 0 {
		s.Transcribe(fmt.Sprintf("%s influence in %s is at 0.", aff, c))
		return
	}
	c.Inf[aff] = Max(0, c.Inf[aff]-n)
	s.Transcribe(fmt.Sprintf("%s influence in %s -%d, now %d.", aff, c, n, c.Inf[aff]))
}

func doubleInf(s *State, c *Country, aff Aff) {
	if c.Inf[aff] > 0 {
		s.Transcribe(fmt.Sprintf("%s influence doubled in %s, now %d.", aff, c, c.Inf[aff]))
	}
	c.Inf[aff] *= 2
}

func zeroInf(s *State, c *Country, aff Aff) {
	if c.Inf[aff] > 0 {
		s.Transcribe(fmt.Sprintf("All %s influence in %s removed.", aff, c))
	}
	c.Inf[aff] = 0
}

func matchInf(s *State, c *Country, toMatch, toReceive Aff) {
	if c.Inf[toReceive] >= c.Inf[toMatch] {
		s.Transcribe(fmt.Sprintf("%s already matches %s influence in %s.", toReceive, toMatch, c))
		return
	}
	c.Inf[toReceive] = c.Inf[toMatch]
	s.Transcribe(fmt.Sprintf("%s matches %s influence in %s, now %d.", toReceive, toMatch, c, c.Inf[toReceive]))
}

func PlusInf(aff Aff, n int) countryChange {
	return func(s *State, c *Country) {
		plusInf(s, c, aff, n)
	}
}

func LessInf(aff Aff, n int) countryChange {
	return func(s *State, c *Country) {
		lessInf(s, c, aff, n)
	}
}

func DoubleInf(aff Aff) countryChange {
	return func(s *State, c *Country) {
		doubleInf(s, c, aff)
	}
}

func ZeroInf(aff Aff) countryChange {
	return func(s *State, c *Country) {
		zeroInf(s, c, aff)
	}
}

func MatchInf(toMatch, toReceive Aff) countryChange {
	return func(s *State, c *Country) {
		matchInf(s, c, toMatch, toReceive)
	}
}

func NoOp(s *State, c *Country) {
	return
}

func NormalCost(target *Country) int {
	return 1
}

func OpInfluenceCost(player Aff) func(*Country) int {
	return func(target *Country) int {
		controlled := target.Controlled()
		if controlled == player.Opp() {
			return 2
		}
		return 1
	}
}

func OpsLimit(s *State, player Aff, card Card) func([]*Country) int {
	return func(cs []*Country) int {
		return ComputeCardOps(s, player, card, cs)
	}
}

func LimitN(n int) func([]*Country) int {
	return func(cs []*Country) int {
		return n
	}
}

func SelectInfluence(s *State, player Aff, message string, change countryChange, n int, checks ...countryCheck) []*Country {
	return selectInfluence(s, player, message, change, LimitN(n), false, NormalCost, checks...)
}

func SelectInfluenceExactly(s *State, player Aff, message string, change countryChange, n int, checks ...countryCheck) []*Country {
	return selectInfluence(s, player, message, change, LimitN(n), true, NormalCost, checks...)
}

func SelectOneInfluence(s *State, player Aff, message string, change countryChange, checks ...countryCheck) *Country {
	return selectInfluence(s, player, message, change, LimitN(1), true, NormalCost, checks...)[0]
}

func selectInfluence(s *State, player Aff, message string, change countryChange, nFun func([]*Country) int, exactly bool, costFun func(*Country) int, checks ...countryCheck) []*Country {
	remote := player != s.LocalPlayer
	used := 0
	chosen := []*Country{}
	var c *Country
	var err error
loop:
	if err != nil {
		s.UI.Message(err.Error())
		err = nil
	}
	log.Printf("Reading from %s '%s'\n", player, message)
	if !s.ReadInto(&c, remote) {
		localInput(s, &c, message)
	}
	cost := costFun(c)
	n := nFun(append(chosen, c))
	switch {
	case c == EndSelectCountry && exactly:
		err = fmt.Errorf("Invalid choice")
		goto loop
	case c == EndSelectCountry:
		// We are done!
		if !remote {
			s.Log(c)
		}
		return chosen
	case used+cost > n:
		err = fmt.Errorf("Too much! That would use %d.", used+cost)
		goto loop
	default:
		for _, check := range checks {
			if err = check(c); err != nil {
				goto loop
			}
		}
	}
	// Success!
	used += cost
	chosen = append(chosen, c)
	log.Printf("Added %s, now used %d\n", c.Name, used)
	// Must log the country before applying the countryChange. This allows
	// countryChange implementations to write to the log, which must follow
	// country selection.
	if !remote {
		s.Log(c)
	}
	change(s, c)
	s.Redraw(s.Game)
	if used == n {
		return chosen
	}
	goto loop
}
