package twistr

import (
	"fmt"
	gc "github.com/rthornton128/goncurses"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	world = `          CAN 4 ----.                     NOR 4 -- SWE 4 -- FIN 4 --.
        /   :        '-._               /   :    /   :        :      '-.
       /                 '-._          /   DNK 3   EDE 3 - POL 3 --.    '                        Turn:
                             '---- UK  5     :   \   :   \   :      '---  * U.S.S.R. *           VP:
   * U.S.A. *                         :  \        \ /   | \  /                                   US Mil:
                                      \   BLX 3 - WDE 4 | CZE 3          /       \               USSR Mil:
   /       \                           \    :   /   :   |   :           /         \              DEFCON:
  /        |                            \      /       /    \          /           \             US Space:
MEX 2     CUB 3                         FRA 3 '   AUT 4 - HUN 3 - ROU 3             \            USSR Space:
  :         :                             :   \     :       :   /   :                \
   \        \   \                       /  \   \   /        / .'       \              \                   NKR 3
   GTM 1     \   HTI 1 - DOM 1      ESP 2 -+- ITA 2 ---- YUG 3  BGR 3 - TUR 2          \                    :  
     :   \    \    :       :          :    |    :   \      :   /  :   /   :             \                   |
 /        \    \                      |     \        '-._  \  / _.---'      \            \                SKR 3
SLV 1 - HND 2 - NIC 1               MAR 3 - ALG 2-TUN 2   GRC 2    LBN 1 - SYR 2         AFG 2              :   \
  :       :   /   :                   :       :     :       :    /   :   /   :         /   :                |    \
          |  /                        |       |      \        .-ISR 4-+-'IRQ 3 - IRN 2     |                |     JPN 4
        CRI 3 - PAN 2   VEN 2       WAS 2   SHS 1  LBY 2-EGY 2    :   |    :           - PAK 2              |   /   :  
          :       :       :           :       :      :     :       \  |   /|\              :              TWN 3     |
                   \    /  |          |        \           |        JOR 2  |  GST 3         \               :       |
           ECU 2 - COL 1   |        CIV 2 --- NGA 1      SDN 1        :    |    :            IND 3                  |
             :       :     |          :         :          :   \          \|/                  :                    |
             |             |                    |               ETH 1    SAU 3                  \                   |
           PER 2 - BOL 2   |                  CMR 1 - ZIR 1       :   \    :                     BUR 2 - LAO 1    PHL 2
             :       :   BRA 2                  :       :       KEN 2  \                           :       :        :  
            /        /     :                          /   \       :   - SOM 2                           /   \       |
          CHL 3   PRY 2    /                         /     \      |       :                         THA 2 - VNM 1   |
            :       :     /                         /   ZWE 1 - SEA 1                                 :       :    /
             \    /   \  /                     AGO 1      :       :                                   \       _ IDN 1
             ARG 2 -- URY 2                      :        |                                          MYS 2 --'    :  
               :        :                        \    _ BWA 2                                          :   \ 
                                                ZAF 3     :                                                 \ AUS 4
       Action Round:                              :                                                             :  `
)

type Pos struct {
	X int
	Y int
}

var (
	turnPos      Pos = Pos{2, 109}
	vpPos        Pos = Pos{3, 109}
	usaMilOpsPos Pos = Pos{4, 109}
	sovMilOpsPos Pos = Pos{5, 109}
	defconPos    Pos = Pos{6, 109}
	usaSpacePos  Pos = Pos{7, 109}
	sovSpacePos  Pos = Pos{8, 109}
	arPos        Pos = Pos{32, 21}
)

const (
	C_BattleName = 1 + iota
	C_BattleStab
	C_NormalName
	C_NormalStab
	C_SovControl
	C_UsaControl
	C_SovInfluence
	C_UsaInfluence
	C_Cam
	C_Sam
	C_Weu
	C_Eeu
	C_Mde
	C_Afr
	C_Asi
	C_Sea
)

var (
	regionColors = map[string][2]int16{
		"CentralAmerica": [2]int16{C_Cam, C_Cam},
		"SouthAmerica":   [2]int16{C_Sam, C_Sam},
		"Europe":         [2]int16{C_Weu, C_Eeu},
		"WestEurope":     [2]int16{C_Weu, C_Weu},
		"EastEurope":     [2]int16{C_Eeu, C_Eeu},
		"MiddleEast":     [2]int16{C_Mde, C_Mde},
		"Africa":         [2]int16{C_Afr, C_Afr},
		"Asia":           [2]int16{C_Asi, C_Asi},
		"SoutheastAsia":  [2]int16{C_Asi, C_Sea},
	}
)

var (
	data = map[CountryId]struct {
		Code string
		Y    int
		X    int
	}{
		Afghanistan:       {"AFG", 14, 89},
		Algeria:           {"ALG", 14, 44},
		Angola:            {"AGO", 28, 47},
		Argentina:         {"ARG", 29, 13},
		SEAfricanStates:   {"SEA", 27, 64},
		Austria:           {"AUT", 8, 50},
		Australia:         {"AUS", 31, 110},
		Bulgaria:          {"BGR", 11, 64},
		Bolivia:           {"BOL", 23, 19},
		Brazil:            {"BRA", 24, 25},
		Burma:             {"BUR", 23, 97},
		Botswana:          {"BWA", 30, 56},
		Benelux:           {"BLX", 5, 42},
		Canada:            {"CAN", 0, 10},
		IvoryCoast:        {"CIV", 20, 36},
		Cameroon:          {"CMR", 23, 46},
		Chile:             {"CHL", 26, 10},
		Colombia:          {"COL", 20, 19},
		CostaRica:         {"CRI", 17, 8},
		Cuba:              {"CUB", 8, 10},
		Czechoslovakia:    {"CZE", 5, 58},
		Denmark:           {"DNK", 2, 43},
		DominicanRep:      {"DOM", 11, 25},
		Ecuador:           {"ECU", 20, 11},
		EGermany:          {"EDE", 2, 51},
		Egypt:             {"EGY", 17, 57},
		Ethiopia:          {"ETH", 22, 64},
		Finland:           {"FIN", 0, 60},
		France:            {"FRA", 8, 40},
		Greece:            {"GRC", 14, 58},
		GulfStates:        {"GST", 19, 78},
		Guatemala:         {"GTM", 11, 3},
		Honduras:          {"HND", 14, 8},
		Haiti:             {"HTI", 11, 17},
		Hungary:           {"HUN", 8, 58},
		Indonesia:         {"IDN", 28, 112},
		Israel:            {"ISR", 16, 64},
		India:             {"IND", 20, 93},
		Iraq:              {"IRQ", 16, 73},
		Iran:              {"IRN", 16, 81},
		Italy:             {"ITA", 11, 46},
		Jordan:            {"JOR", 19, 68},
		Japan:             {"JPN", 16, 114},
		Kenya:             {"KEN", 24, 64},
		LaosCambodia:      {"LAO", 23, 105},
		Lebanon:           {"LBN", 14, 67},
		Libya:             {"LBY", 17, 51},
		Morocco:           {"MAR", 14, 36},
		Mexico:            {"MEX", 8, 0},
		Malaysia:          {"MYS", 29, 101},
		Nigeria:           {"NGA", 20, 46},
		Nicaragua:         {"NIC", 14, 16},
		NKorea:            {"NKR", 10, 106},
		Norway:            {"NOR", 0, 42},
		Panama:            {"PAN", 17, 16},
		Peru:              {"PER", 23, 11},
		Philippines:       {"PHL", 23, 114},
		Pakistan:          {"PAK", 17, 89},
		Poland:            {"POL", 2, 59},
		Paraguay:          {"PRY", 26, 18},
		Romania:           {"ROU", 8, 66},
		SaudiArabia:       {"SAU", 22, 73},
		Sudan:             {"SDN", 20, 57},
		Sweden:            {"SWE", 0, 51},
		SKorea:            {"SKR", 13, 106},
		Somalia:           {"SOM", 25, 72},
		SpainPortugal:     {"ESP", 11, 36},
		SaharanStates:     {"SHS", 17, 44},
		ElSalvador:        {"SLV", 14, 0},
		Syria:             {"SYR", 14, 75},
		Thailand:          {"THA", 26, 100},
		Tunisia:           {"TUN", 14, 50},
		Turkey:            {"TUR", 11, 72},
		Taiwan:            {"TWN", 18, 106},
		UK:                {"UK ", 3, 35},
		Uruguay:           {"URY", 29, 22},
		Venezuela:         {"VEN", 17, 24},
		Vietnam:           {"VNM", 26, 108},
		WestAfricanStates: {"WAS", 17, 36},
		WGermany:          {"WDE", 5, 50},
		Yugoslavia:        {"YUG", 11, 57},
		SouthAfrica:       {"ZAF", 31, 48},
		Zaire:             {"ZIR", 23, 54},
		Zimbabwe:          {"ZWE", 27, 56},
	}
)

func initColors() {
	gc.InitPair(C_BattleName, gc.C_WHITE, 61)
	gc.InitPair(C_BattleStab, gc.C_WHITE, 1)
	gc.InitPair(C_NormalName, gc.C_BLACK, 230)
	gc.InitPair(C_NormalStab, gc.C_BLACK, 227)

	gc.InitPair(C_SovControl, 227, 9)
	gc.InitPair(C_UsaControl, gc.C_WHITE, 32)
	gc.InitPair(C_SovInfluence, 9, gc.C_WHITE)
	gc.InitPair(C_UsaInfluence, 32, gc.C_WHITE)

	gc.InitPair(C_Cam, gc.C_BLACK, 187)
	gc.InitPair(C_Sam, gc.C_BLACK, 150)
	gc.InitPair(C_Weu, gc.C_BLACK, 140)
	gc.InitPair(C_Eeu, gc.C_BLACK, 182)
	gc.InitPair(C_Mde, gc.C_BLACK, 195)
	gc.InitPair(C_Afr, gc.C_BLACK, 229)
	gc.InitPair(C_Asi, gc.C_BLACK, 214)
	gc.InitPair(C_Sea, gc.C_BLACK, 220)
}

type NCursesUI struct {
	*gc.Window
}

func MakeNCursesUI() *NCursesUI {
	scr, err := gc.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if !gc.HasColors() {
		log.Fatal("No colors")
	}
	if err := gc.StartColor(); err != nil {
		log.Fatal(err)
	}
	gc.Echo(false)
	initColors()
	return &NCursesUI{scr}
}

func (nc *NCursesUI) Input() (string, error) {
	nc.Move(37, 0)
	gc.Echo(true)
	nc.Refresh()
	text, err := nc.GetString(100)
	if err != nil {
		return "", err
	}
	gc.Echo(false)
	nc.Move(37, 0)
	nc.ClearToEOL()
	return strings.ToLower(strings.TrimSpace(text)), nil
}

func (nc *NCursesUI) Message(message string) error {
	nc.Move(36, 0)
	nc.ClearToEOL()
	nc.MovePrint(36, 0, strings.TrimRight(message, "\n"))
	nc.MoveTo(37, 0)
	return nil
}

func (nc *NCursesUI) Close() error {
	gc.End()
	return nil
}

func (nc *NCursesUI) Redraw(s *State) {
	var name, stab, infUsa, infSov int16
	nc.MovePrint(0, 0, world)
	for id, extra := range data {
		country := s.Countries[id]
		if country.Battleground {
			name, stab = C_BattleName, C_BattleStab
		} else {
			name, stab = C_NormalName, C_NormalStab
		}
		nc.ColorOn(name)
		nc.MovePrint(extra.Y, extra.X, extra.Code)
		nc.ColorOff(name)
		nc.ColorOn(stab)
		nc.MovePrint(extra.Y, extra.X+3, fmt.Sprintf(" %d", country.Stability))
		nc.ColorOff(stab)
		// Influence
		infUsa, infSov = infColors(country)
		nc.AttrOn(gc.A_BOLD)
		nc.ColorOn(infUsa)
		nc.MovePrint(extra.Y+1, extra.X, infString(country.Inf[USA]))
		nc.ColorOff(infUsa)
		// Separator space should have the country's background color
		colors := countryColors(country)
		nc.ColorOn(colors[USA])
		nc.MovePrint(extra.Y+1, extra.X+2, " ")
		nc.ColorOff(colors[USA])
		nc.ColorOn(infSov)
		nc.MovePrint(extra.Y+1, extra.X+3, infString(country.Inf[SOV]))
		nc.ColorOff(infSov)
		nc.AttrOff(gc.A_BOLD)
	}
	// Draw game metadata
	nc.MovePrint(turnPos.X, turnPos.Y, strconv.Itoa(s.Turn))
	var vp string
	switch {
	case s.VP > 0:
		vp = fmt.Sprintf("US +%d", s.VP)
	case s.VP < 0:
		vp = fmt.Sprintf("USSR +%d", -s.VP)
	default:
		vp = "0"
	}
	nc.MovePrint(vpPos.X, vpPos.Y, vp)
	nc.MovePrint(usaMilOpsPos.X, usaMilOpsPos.Y, strconv.Itoa(s.MilOps[USA]))
	nc.MovePrint(sovMilOpsPos.X, sovMilOpsPos.Y, strconv.Itoa(s.MilOps[SOV]))
	nc.MovePrint(defconPos.X, defconPos.Y, strconv.Itoa(s.Defcon))
	nc.MovePrint(usaSpacePos.X, usaSpacePos.Y, strconv.Itoa(s.SpaceRace[USA]))
	nc.MovePrint(sovSpacePos.X, sovSpacePos.Y, strconv.Itoa(s.SpaceRace[SOV]))
	var phasingColor int16
	if s.Phasing == USA {
		phasingColor = C_UsaControl
	} else {
		phasingColor = C_SovControl
	}
	nc.ColorOn(phasingColor)
	nc.MovePrint(arPos.X, arPos.Y, fmt.Sprintf("%2d", s.AR))
	nc.ColorOff(phasingColor)
	nc.Refresh()
	nc.Move(37, 0)
	return nil
}

// countryColors returns the default coloring for a country, notwithstanding
// control by either superpower.
// Western european countries are simply western european, but Austria
// straddles both, so we treat it as european.
func countryColors(country *Country) [2]int16 {
	var regionName string
	switch country.Region.Name {
	case "Europe":
		switch {
		case country.In(WestEurope) && country.In(EastEurope):
			regionName = Europe.Name
		case country.In(WestEurope):
			regionName = WestEurope.Name
		default:
			regionName = EastEurope.Name
		}
	case "Asia":
		if country.In(SoutheastAsia) {
			regionName = SoutheastAsia.Name
		} else {
			regionName = Asia.Name
		}
	default:
		regionName = country.Region.Name
	}
	return regionColors[regionName]
}

func infColors(country *Country) (int16, int16) {
	colors := countryColors(country)
	influenceColor := func(aff Aff) int16 {
		if country.Inf[aff] > 0 {
			switch aff {
			case SOV:
				return C_SovInfluence
			default:
				return C_UsaInfluence
			}
		}
		return colors[aff]
	}
	switch country.Controlled() {
	case USA:
		return C_UsaControl, influenceColor(SOV)
	case SOV:
		return influenceColor(USA), C_SovControl
	default:
		return influenceColor(USA), influenceColor(SOV)
	}
}

func infString(inf int) string {
	switch inf {
	case 0:
		return "  "
	default:
		return fmt.Sprintf("%2d", inf)
	}
}
