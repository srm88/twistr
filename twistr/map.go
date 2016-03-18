package twistr

// Shortcut. US = 0 and Sov = 1, so we can index into any [2] array with a
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

func (c Country) Controlled() Aff {
	switch {
	case (c.Inf[US] - c.Inf[Sov]) >= c.Stability:
		return US
	case (c.Inf[Sov] - c.Inf[US]) >= c.Stability:
		return Sov
	default:
		return Neu
	}
}

func (c Country) In(r Region) bool {
	for _, cid := range r.Countries {
		if cid == c.Id {
			return true
		}
	}
	return false
}

type Region struct {
	Name       string
	Countries  []CountryId
	Volatility int
}

func (r Region) String() string {
	return r.Name
}
