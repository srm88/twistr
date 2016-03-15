package twistr

// Shortcut. US = 0 and Sov = 1, so we can index into any [2] array with a
// player constant.
type Influence [2]int8

type Country struct {
    Id CountryId
    Name string
    Inf Influence
    Stab uint8
    Battleground bool
    AdjSuper Aff
    AdjCountries []*Country
    Region *Region
}

func (c Country) Controlled() Aff {
    switch {
    case (c.Inf[US] - c.Inf[Sov]) >= c.Stab:
        return US
    case (c.Inf[Sov] - c.Inf[US]) >= c.Stab:
        return Sov
    default:
        return Neu
    }
}

type Region struct {
    Id RegionId
    Name string
    Countries []*Country
    Volatility uint8
}
