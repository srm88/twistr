package twistr

import "strings"

// Shortcut. USA = 0 and SOV = 1, so we can index into any [2] array with a
// player constant.
type Influence [2]int

type Country struct {
	Id           CountryId
	Name         string
	Inf          Influence
	Stability    int
	Battleground bool
	AdjSuper     Aff
	AdjCountries []*Country
	Region       Region
}

func (c Country) String() string {
	return c.Name
}

func (c Country) Ref() string {
	return strings.ToLower(c.Name)
}

func (c Country) Controlled() Aff {
	switch {
	case (c.Inf[USA] - c.Inf[SOV]) >= c.Stability:
		return USA
	case (c.Inf[SOV] - c.Inf[USA]) >= c.Stability:
		return SOV
	default:
		return NEU
	}
}

func (c Country) NumControlledNeighbors(aff Aff) int {
	n := 0
	for _, c := range c.AdjCountries {
		if c.Controlled() == aff {
			n += 1
		}
	}
	return n
}

func (c Country) In(r Region) bool {
	for _, cid := range r.Countries {
		if cid == c.Id {
			return true
		}
	}
	return false
}

func cloneCountries(cMap map[CountryId]*Country) map[CountryId]*Country {
	clone := make(map[CountryId]*Country)
	for cid, oldC := range cMap {
		// Fill each new country pointer's value with the old pointer's value
		clone[cid] = new(Country)
		*clone[cid] = *oldC
	}
	// Each clone's AdjCountries pointers point at originals. Point to clones.
	for _, newC := range clone {
		for i, oldC := range newC.AdjCountries {
			newC.AdjCountries[i] = clone[oldC.Id]
		}
	}
	return clone
}

func AllIn(cs []*Country, r Region) bool {
	for _, c := range cs {
		if !c.In(r) {
			return false
		}
	}
	return true
}

type Region struct {
	Name       string
	Countries  []CountryId
	Volatility int
}

func (r Region) String() string {
	return r.Name
}

func (r Region) Ref() string {
	return strings.ToLower(r.Name)
}
