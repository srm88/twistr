package twistr

import "fmt"
import "log"

// All WIP. Maybe obliterate it.

// Realignment
func realignMods(target Country) (mods Influence) {
	switch {
	case target.Inf[USA] > target.Inf[SOV]:
		mods[USA] += 1
	case target.Inf[SOV] > target.Inf[USA]:
		mods[SOV] += 1
	}
	for _, neighbor := range target.AdjCountries {
		control := neighbor.Controlled()
		if control != NEU {
			mods[control] += 1
		}
	}
	return
}

func realign(s *State, target *Country, rollUSA, rollSOV int) {
	mods := realignMods(*target)
	rollUSA += mods[USA]
	rollSOV += mods[SOV]
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
		s.Transcribe(fmt.Sprintf("%d soviet influence removed", initSOV-target.Inf[SOV]))
	} else {
		s.Transcribe("No influence removed")
	}
}

func coupBonus(s *State, player Aff, target *Country) (bonus int) {
	if s.Effect(SALTNegotiations) {
		bonus -= 1
	}
	if s.Effect(LatinAmericanDeathSquads, player) {
		bonus += 1
	}
	if s.Effect(LatinAmericanDeathSquads, player.Opp()) {
		bonus -= 1
	}
	return
}

func opsMod(s *State, player Aff, countries []*Country) (mod int) {
	if player == SOV && s.Effect(VietnamRevolts) {
		if AllIn(countries, SoutheastAsia) {
			mod += 1
		}
	}
	return
}

// Coup
func coup(s *State, player Aff, ops int, roll int, target *Country, free bool) bool {
	bonus := coupBonus(s, player, target)
	delta := roll + bonus + ops - (target.Stability * 2)
	if delta <= 0 {
		return false
	}
	oppCurInf := target.Inf[player.Opp()]
	removed := Min(oppCurInf, delta)
	gained := delta - removed
	if target.Battleground {
		// XXX: CubanMissileCrisis, NuclearSubs
		s.DegradeDefcon(1)
	}
	if !free {
		s.MilOps[player] += ops
	}
	target.Inf[player] += gained
	target.Inf[player.Opp()] -= removed
	return true
}

// A country cannot be coup'd if it lacks any of the opponent's influence.
// Some permanent events also impose coup restrictions, e.g. NATO with Europe.
func canCoup(s *State, player Aff, t *Country, free bool) bool {
	switch {
	case t.Inf[player.Opp()] < 1:
		return false
	case natoProtected(s, player, t):
		return false
	case japanProtected(s, player, t):
		return false
	case s.Effect(TheReformer) && player == SOV && t.In(Europe):
		return false
	case defconProtected(s, t) && !free:
		return false
	default:
		return true
	}
}

func canRealign(s *State, player Aff, t *Country, free bool) bool {
	switch {
	case natoProtected(s, player, t):
		return false
	case japanProtected(s, player, t):
		return false
	case t.Inf[player.Opp()] < 1:
		return false
	case defconProtected(s, t) && !free:
		return false
	default:
		return true
	}
}

func defconProtected(s *State, t *Country) bool {
	// asia 3, defcon 5, not protected
	// europe 4, defcon 3, protected
	// middle east 2, defcon 2, protected
	return t.Region.Volatility >= s.Defcon
}

func natoProtected(s *State, player Aff, t *Country) bool {
	return s.Effect(NATO) && player == SOV && t.In(Europe) && t.Controlled() == USA
}

func japanProtected(s *State, player Aff, t *Country) bool {
	return s.Effect(USJapanMutualDefensePact) && t.Id == Japan && player == SOV
}

type influenceChange func(*Country) error

func PlusInf(aff Aff, n int) influenceChange {
	return func(c *Country) error {
		c.Inf[aff] += n
		return nil
	}
}

func LessInf(aff Aff, n int) influenceChange {
	return func(c *Country) error {
		if c.Inf[aff] == 0 {
			return fmt.Errorf("No %s influence in %s", aff, c.Name)
		}
		c.Inf[aff] = Max(0, c.Inf[aff]-n)
		return nil
	}
}

func DoubleInf(aff Aff) influenceChange {
	return func(c *Country) error {
		if c.Inf[aff] == 0 {
			return fmt.Errorf("No %s influence in %s", aff, c.Name)
		}
		c.Inf[aff] *= 2
		return nil
	}
}

func ZeroInf(aff Aff) influenceChange {
	return func(c *Country) error {
		if c.Inf[aff] == 0 {
			return fmt.Errorf("No %s influence in %s", aff, c.Name)
		}
		c.Inf[aff] = 0
		return nil
	}
}

func MatchInf(toMatch, toReceive Aff) influenceChange {
	return func(c *Country) error {
		if c.Inf[toReceive] >= c.Inf[toMatch] {
			return fmt.Errorf("Already match %s influence in %s", toMatch, c.Name)
		}
		c.Inf[toReceive] = c.Inf[toMatch]
		return nil
	}
}

func NoOp(c *Country) error {
	return nil
}

func NormalCost(player Aff, target *Country) int {
	return 1
}

func OpInfluenceCost(player Aff, target *Country) int {
	controlled := target.Controlled()
	if controlled == player.Opp() {
		return 2
	}
	return 1
}

func SelectInfluence(s *State, player Aff, message string, change influenceChange, n int, checks ...countryCheck) []*Country {
	return selectInfluence(s, player, message, change, n, false, NormalCost, checks...)
}

func SelectInfluenceExactly(s *State, player Aff, message string, change influenceChange, n int, checks ...countryCheck) []*Country {
	return selectInfluence(s, player, message, change, n, true, NormalCost, checks...)
}

func SelectOneInfluence(s *State, player Aff, message string, change influenceChange, checks ...countryCheck) *Country {
	return selectInfluence(s, player, message, change, 1, true, NormalCost, checks...)[0]
}

func selectInfluence(s *State, player Aff, message string, change influenceChange, n int, exactly bool, costFun func(Aff, *Country) int, checks ...countryCheck) []*Country {
	used := 0
	chosen := []*Country{}
	var c *Country
	var err error
loop:
	if err != nil {
		s.UI.Message(player, err.Error())
		err = nil
	}
	if !s.ReadInto(&c) {
		GetInput(s, player, &c, message)
	}
	cost := costFun(player, c)
	switch {
	case c == EndSelectCountry && exactly:
		err = fmt.Errorf("Invalid choice")
		goto loop
	case c == EndSelectCountry:
		// We are done!
		s.Log(c)
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
	// User managed to select a country; let's see if the story checks out.
	if err = change(c); err != nil {
		goto loop
	}
	// Success!
	used += cost
	chosen = append(chosen, c)
	log.Printf("Added %s, now used %d\n", c.Name, used)
	s.Log(c)
	s.Redraw(s.Game)
	if used == n {
		return chosen
	}
	goto loop
}
