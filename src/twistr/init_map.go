package twistr

import (
    "bytes"
)

var (
    countries map[CountryId]*Country
)

// Temp:
func ByName(name string) *Country {
    for _, c := range countries {
        if c.Name == name {
            return c
        }
    }
    return nil
}

func CountryNames(cs []*Country) string {
    var b bytes.Buffer
    for _, c := range cs {
        b.WriteString(c.Name)
        b.WriteString(" ")
    }
    return b.String()
}

func init() {
    countries = make(map[CountryId]*Country)
    for _, c := range countryTable {
        countries[c.Id] = &Country{
            Id: c.Id,
            Name: c.Name,
            Inf: Influence{c.USInf, c.SovInf},
            Stability: c.Stability,
            Battleground: c.Battleground,
            AdjSuper: c.AdjSuper,
        }
    }
    for _, link := range countryLinks {
        foo := countries[link[0]]
        bar := countries[link[1]]
        foo.AdjCountries = append(foo.AdjCountries, bar)
        bar.AdjCountries = append(bar.AdjCountries, foo)
    }
}

var countryTable = []struct{
    Id CountryId
    Name string
    USInf int
    SovInf int
    Stability int
    Battleground bool
    AdjSuper Aff
    Region RegionId
} {
    { Mexico, "Mexico", 0, 0, 2, true, US, CentralAmerica },
    { Guatemala, "Guatemala", 0, 0, 1, false, Neu, CentralAmerica },
    { ElSalvador, "ElSalvador", 0, 0, 1, false, Neu, CentralAmerica },
    { Honduras, "Honduras", 0, 0, 2, false, Neu, CentralAmerica },
    { CostaRica, "CostaRica", 0, 0, 3, false, Neu, CentralAmerica },
    { Cuba, "Cuba", 0, 0, 3, true, US, CentralAmerica },
    { Nicaragua, "Nicaragua", 0, 0, 1, false, Neu, CentralAmerica },
    { Panama, "Panama", 1, 0, 2, true, Neu, CentralAmerica },
    { Haiti, "Haiti", 0, 0, 1, false, Neu, CentralAmerica },
    { DominicanRep, "DominicanRep", 0, 0, 1, false, Neu, CentralAmerica },
    { Ecuador, "Ecuador", 0, 0, 2, false, Neu, SouthAmerica },
    { Peru, "Peru", 0, 0, 2, false, Neu, SouthAmerica },
    { Colombia, "Colombia", 0, 0, 1, false, Neu, SouthAmerica },
    { Chile, "Chile", 0, 0, 3, true, Neu, SouthAmerica },
    { Venezuela, "Venezuela", 0, 0, 2, true, Neu, SouthAmerica },
    { Argentina, "Argentina", 0, 0, 2, true, Neu, SouthAmerica },
    { Bolivia, "Bolivia", 0, 0, 2, false, Neu, SouthAmerica },
    { Paraguay, "Paraguay", 0, 0, 2, false, Neu, SouthAmerica },
    { Uruguay, "Uruguay", 0, 0, 2, false, Neu, SouthAmerica },
    { Brazil, "Brazil", 0, 0, 2, true, Neu, SouthAmerica },
    { Canada, "Canada", 0, 0, 4, false, US, Europe },
    { UK, "UK", 5, 0, 5, false, Neu, Europe },
    { SpainPortugal, "SpainPortugal", 0, 0, 2, false, Neu, Europe },
    { France, "France", 0, 0, 3, true, Neu, Europe },
    { Benelux, "Benelux", 0, 0, 3, false, Neu, Europe },
    { Norway, "Norway", 0, 0, 4, false, Neu, Europe },
    { Denmark, "Denmark", 0, 0, 3, false, Neu, Europe },
    { WGermany, "WGermany", 0, 0, 4, true, Neu, Europe },
    { EGermany, "EGermany", 3, 0, 3, true, Neu, Europe },
    { Italy, "Italy", 0, 0, 2, true, Neu, Europe },
    { Austria, "Austria", 0, 0, 4, false, Neu, Europe },
    { Sweden, "Sweden", 0, 0, 4, false, Neu, Europe },
    { Czechoslovakia, "Czechoslovakia", 0, 0, 3, false, Neu, Europe },
    { Yugoslavia, "Yugoslavia", 0, 0, 3, false, Neu, Europe },
    { Poland, "Poland", 0, 0, 3, true, Sov, Europe },
    { Greece, "Greece", 0, 0, 2, false, Neu, Europe },
    { Hungary, "Hungary", 0, 0, 3, false, Neu, Europe },
    { Finland, "Finland", 1, 0, 4, false, Sov, Europe },
    { Romania, "Romania", 0, 0, 3, false, Sov, Europe },
    { Bulgaria, "Bulgaria", 0, 0, 3, false, Neu, Europe },
    { Turkey, "Turkey", 0, 0, 2, false, Neu, Europe },
    { Morocco, "Morocco", 0, 0, 3, false, Neu, Africa },
    { WestAfricanStates, "WestAfricanStates", 0, 0, 2, false, Neu, Africa },
    { IvoryCoast, "IvoryCoast", 0, 0, 2, false, Neu, Africa },
    { Algeria, "Algeria", 0, 0, 2, true, Neu, Africa },
    { SaharanStates, "SaharanStates", 0, 0, 1, false, Neu, Africa },
    { Nigeria, "Nigeria", 0, 0, 1, true, Neu, Africa },
    { Tunisia, "Tunisia", 0, 0, 2, false, Neu, Africa },
    { Cameroon, "Cameroon", 0, 0, 1, false, Neu, Africa },
    { Angola, "Angola", 0, 0, 1, true, Neu, Africa },
    { SouthAfrica, "SouthAfrica", 1, 0, 3, true, Neu, Africa },
    { Zaire, "Zaire", 0, 0, 1, true, Neu, Africa },
    { Botswana, "Botswana", 0, 0, 2, false, Neu, Africa },
    { Zimbabwe, "Zimbabwe", 0, 0, 1, false, Neu, Africa },
    { Sudan, "Sudan", 0, 0, 1, false, Neu, Africa },
    { Ethiopia, "Ethiopia", 0, 0, 1, false, Neu, Africa },
    { Kenya, "Kenya", 0, 0, 2, false, Neu, Africa },
    { SEAfricanStates, "SEAfricanStates", 0, 0, 1, false, Neu, Africa },
    { Somalia, "Somalia", 0, 0, 2, false, Neu, Africa },
    { Libya, "Libya", 0, 0, 2, true, Neu, MiddleEast },
    { Egypt, "Egypt", 0, 0, 2, true, Neu, MiddleEast },
    { Israel, "Israel", 1, 0, 4, true, Neu, MiddleEast },
    { Lebanon, "Lebanon", 0, 0, 1, false, Neu, MiddleEast },
    { Jordan, "Jordan", 0, 0, 2, false, Neu, MiddleEast },
    { Syria, "Syria", 1, 0, 2, false, Neu, MiddleEast },
    { Iraq, "Iraq", 1, 0, 3, true, Neu, MiddleEast },
    { SaudiArabia, "SaudiArabia", 0, 0, 3, true, Neu, MiddleEast },
    { GulfStates, "GulfStates", 0, 0, 3, false, Neu, MiddleEast },
    { Iran, "Iran", 1, 0, 2, true, Neu, MiddleEast },
    { Afghanistan, "Afghanistan", 0, 0, 2, false, Sov, Asia },
    { Pakistan, "Pakistan", 0, 0, 2, true, Neu, Asia },
    { India, "India", 0, 0, 3, true, Neu, Asia },
    { Burma, "Burma", 0, 0, 2, false, Neu, Asia },
    { Thailand, "Thailand", 0, 0, 2, true, Neu, Asia },
    { LaosCambodia, "LaosCambodia", 0, 0, 1, false, Neu, Asia },
    { Vietnam, "Vietnam", 0, 0, 1, false, Neu, Asia },
    { Malaysia, "Malaysia", 0, 0, 2, false, Neu, Asia },
    { Indonesia, "Indonesia", 0, 0, 1, false, Neu, Asia },
    { Australia, "Australia", 4, 0, 4, false, Neu, Asia },
    { Taiwan, "Taiwan", 0, 0, 3, false, Neu, Asia },
    { NKorea, "NKorea", 3, 0, 3, true, Sov, Asia },
    { SKorea, "SKorea", 1, 0, 3, true, Neu, Asia },
    { Philippines, "Philippines", 1, 0, 2, false, Neu, Asia },
    { Japan, "Japan", 1, 0, 4, true, US, Asia },
}

var countryLinks = [][2]CountryId{
    // Central America
    { Mexico, Guatemala },
    { Guatemala, ElSalvador },
    { Guatemala, Honduras },
    { ElSalvador, Honduras },
    { Honduras, CostaRica },
    { Honduras, Nicaragua },
    { CostaRica, Nicaragua },
    { Cuba, Nicaragua },
    { Cuba, Haiti },
    { Haiti, DominicanRep },
    { CostaRica, Panama },
    // South America
    { Panama, Colombia },
    { Colombia, Ecuador },
    { Ecuador, Peru },
    { Peru, Chile },
    { Peru, Bolivia },
    { Chile, Argentina },
    { Bolivia, Paraguay },
    { Paraguay, Argentina },
    { Paraguay, Uruguay },
    { Uruguay, Brazil },
    { Brazil, Venezuela },
    { Venezuela, Colombia },
    // Europe
    { Canada, UK },
    { UK, France },
    { UK, Norway },
    { SpainPortugal, France },
    { SpainPortugal, Italy },
    { France, Italy },
    { France, WGermany },
    { Benelux, WGermany },
    { Norway, Sweden },
    { Sweden, Denmark },
    { Sweden, Finland },
    { Denmark, WGermany },
    { WGermany, Austria },
    { Austria, Italy },
    { Austria, EGermany },
    { Austria, Hungary },
    { EGermany, WGermany },
    { EGermany, Poland },
    { EGermany, Czechoslovakia },
    { Poland, Czechoslovakia },
    { Czechoslovakia, Hungary },
    { Hungary, Yugoslavia },
    { Hungary, Romania },
    { Romania, Turkey },
    { Romania, Yugoslavia },
    { Yugoslavia, Italy },
    { Yugoslavia, Greece },
    { Greece, Italy },
    { Greece, Bulgaria },
    { Greece, Turkey },
    // Middle east
    { Turkey, Syria },
    { Syria, Lebanon },
    { Syria, Israel },
    { Lebanon, Jordan },
    { Lebanon, Israel },
    { Israel, Egypt },
    { Israel, Jordan },
    { Egypt, Libya },
    { Jordan, Iraq },
    { Jordan, SaudiArabia },
    { Iraq, SaudiArabia },
    { Iraq, GulfStates },
    { Iraq, Iran },
    { GulfStates, SaudiArabia },
    // Africa
    { Egypt, Sudan },
    { Libya, Tunisia },
    { Algeria, France },
    { Morocco, SpainPortugal },
    { Algeria, Tunisia },
    { Morocco, Algeria },
    { Morocco, WestAfricanStates },
    { WestAfricanStates, IvoryCoast },
    { IvoryCoast, Nigeria },
    { Nigeria, SaharanStates },
    { SaharanStates, Algeria },
    { Nigeria, Cameroon },
    { Cameroon, Zaire },
    { Zaire, Angola },
    { Zaire, Zimbabwe },
    { Angola, Botswana },
    { Angola, SouthAfrica },
    { Botswana, SouthAfrica },
    { Botswana, Zimbabwe },
    { Zimbabwe, SEAfricanStates },
    { SEAfricanStates, Kenya },
    { Kenya, Somalia },
    { Somalia, Ethiopia },
    { Ethiopia, Sudan },
    // Asia
    { Iran, Afghanistan },
    { Iran, Pakistan },
    { Pakistan, Afghanistan },
    { Pakistan, India },
    { India, Burma },
    { Burma, LaosCambodia },
    { LaosCambodia, Thailand },
    { LaosCambodia, Vietnam },
    { Vietnam, Thailand },
    { Thailand, Malaysia },
    { Malaysia, Australia },
    { Malaysia, Indonesia },
    { Indonesia, Philippines },
    { Philippines, Japan },
    { Japan, Taiwan },
    { Japan, SKorea },
    { Taiwan, SKorea },
    { SKorea, NKorea },
}
