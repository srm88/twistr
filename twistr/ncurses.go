package twistr

import (
	"bytes"
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
   * U.S.A. *                         :  \        \ /   | \  /                                   US Space:
                                      \   BLX 3 - WDE 4 | CZE 3          /       \               USSR Space:
   /       \                           \    :   /   :   |   :           /         \
  /        |                            \      /       /    \          /           \
MEX 2     CUB 3                         FRA 3 '   AUT 4 - HUN 3 - ROU 3             \
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
               :        :                        \    _ BWA 2    DEFCON  | 5 | 4 | 3 | 2 | X |         :   \
                                                ZAF 3     :      MilOps  | 5 | 4 | 3 | 2 | 1 | 0 |          \ AUS 4
       Action Round:                              :                      | 5 | 4 | 3 | 2 | 1 | 0 |              :  `
)

type Pos struct {
	X int
	Y int
}

var (
	turnPos      Pos = Pos{109, 2}
	vpPos        Pos = Pos{109, 3}
	usaSpacePos  Pos = Pos{109, 4}
	sovSpacePos  Pos = Pos{109, 5}
	arPos        Pos = Pos{21, 32}
	defconPos    Pos = Pos{74, 30}
	usaMilOpsPos Pos = Pos{74, 31}
	sovMilOpsPos Pos = Pos{74, 32}
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
	// Card drawing
	C_Early
	C_Mid
	C_Late
	C_CardText
	// Spacerace
	C_SpaceDefault
	C_SpaceDetails
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

	gc.InitPair(C_Early, gc.C_WHITE, 38)
	gc.InitPair(C_Mid, gc.C_WHITE, 25)
	gc.InitPair(C_Late, gc.C_WHITE, 236)
	gc.InitPair(C_CardText, gc.C_BLACK, 153)

	gc.InitPair(C_SpaceDefault, 1, gc.C_BLACK)
	gc.InitPair(C_SpaceDetails, gc.C_WHITE, 0)
}

type NCursesUI struct {
	*gc.Window
}

func MakeNCursesUI() *NCursesUI {
	scr, err := gc.Init()
	if err != nil {
		log.Println(err)
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

func (nc *NCursesUI) Solicit(player Aff, message string, choices []string) string {
	nc.Move(36, 0)
	nc.ClearToEOL()
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "[%s] %s", player, strings.TrimRight(message, "\n"))
	if len(choices) > 0 {
		fmt.Fprintf(buf, " [ %s ]", strings.Join(choices, " "))
	}
	nc.MovePrint(36, 0, buf.String())
	nc.Move(37, 0)
	gc.Echo(true)
	nc.Refresh()
	text, err := nc.GetString(100)
	if err != nil {
		panic(err.Error())
	}
	gc.Echo(false)
	nc.Move(37, 0)
	nc.ClearToEOL()
	return strings.ToLower(strings.TrimSpace(text))
}

func (nc *NCursesUI) Message(player Aff, message string) {
	nc.Move(36, 0)
	nc.ClearToEOL()
	nc.MovePrint(36, 0, fmt.Sprintf("[%s] %s", player, strings.TrimRight(message, "\n")))
	nc.Move(37, 0)
	nc.GetChar()
}

func (nc *NCursesUI) Close() error {
	gc.End()
	return nil
}

func (nc *NCursesUI) Redraw(g *Game) {
	nc.clear()
	var name, stab, infUsa, infSov int16
	nc.MovePrint(0, 0, world)
	for id, extra := range data {
		country := g.Countries[id]
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
	nc.MovePrint(turnPos.Y, turnPos.X, strconv.Itoa(g.Turn))
	var vp string
	switch {
	case g.VP > 0:
		vp = fmt.Sprintf("US +%d", g.VP)
	case g.VP < 0:
		vp = fmt.Sprintf("USSR +%d", -g.VP)
	default:
		vp = "0"
	}
	nc.MovePrint(vpPos.Y, vpPos.X, vp)
	nc.MovePrint(usaSpacePos.Y, usaSpacePos.X, strconv.Itoa(g.SpaceRace[USA]))
	nc.MovePrint(sovSpacePos.Y, sovSpacePos.X, strconv.Itoa(g.SpaceRace[SOV]))

	nc.ColorOn(C_SpaceDefault)
	nc.MovePrint(defconPos.Y, defconPos.X+(4*(5-g.Defcon)), " @ ")
	nc.ColorOff(C_SpaceDefault)
	nc.ColorOn(C_UsaControl)
	nc.MovePrint(usaMilOpsPos.Y, usaMilOpsPos.X+(4*(5-g.MilOps[USA])), "USA")
	nc.ColorOff(C_UsaControl)
	nc.ColorOn(C_SovControl)
	nc.MovePrint(sovMilOpsPos.Y, sovMilOpsPos.X+(4*(5-g.MilOps[SOV])), "SOV")
	nc.ColorOff(C_SovControl)

	var phasingColor int16
	if g.Phasing == USA {
		phasingColor = C_UsaControl
	} else {
		phasingColor = C_SovControl
	}
	nc.ColorOn(phasingColor)
	nc.MovePrint(arPos.Y, arPos.X, fmt.Sprintf("%2d", g.AR))
	nc.ColorOff(phasingColor)
	nc.Refresh()
	nc.Move(37, 0)
}

func (nc *NCursesUI) clear() {
	for i := 0; i < maxHeight; i++ {
		nc.Move(i, 0)
		nc.ClearToEOL()
	}
	nc.Move(37, 0)
	nc.Refresh()
}

func (nc *NCursesUI) ShowMessages(messages []string) {
	nc.clear()
	start := 0
	if len(messages) > maxHeight {
		start = len(messages) - maxHeight
	}
	for i, msg := range messages[start:] {
		nc.MovePrint(i, 0, msg)
	}
	nc.Move(37, 0)
	nc.Refresh()
}

const (
	spaceBox = `.----------.
|          |
|          |
|          |
|          |
|          |
|          |
'----------'`
	spaceWidth  = 12
	spaceHeight = 8
	spaceNameY  = spaceHeight + 1
	spaceOpsY   = spaceNameY + 3
)

func (nc *NCursesUI) ShowSpaceRace(positions [2]int) {
	nc.clear()
	x := 5
	y := 2
	for i, box := range SRTrack {
		offsetX := i * (spaceWidth + 1)
		usaHere := positions[USA] == i
		sovHere := positions[SOV] == i
		nc.drawSpaceBox(box, usaHere, sovHere, Pos{x + offsetX, y})
	}
}

func (nc *NCursesUI) drawSpaceBox(box SRBox, usaHere, sovHere bool, start Pos) {
	nc.ColorOn(C_SpaceDefault)
	for i, line := range strings.Split(spaceBox, "\n") {
		nc.MovePrint(start.Y+i, start.X, line)
	}
	nc.ColorOff(C_SpaceDefault)
	if usaHere {
		nc.ColorOn(C_UsaControl)
		nc.MovePrint(start.Y+2, start.X+(spaceWidth-2)/2, "US")
		nc.ColorOff(C_UsaControl)
	}
	if sovHere {
		nc.ColorOn(C_SovControl)
		nc.MovePrint(start.Y+4, start.X+(spaceWidth-4)/2, "USSR")
		nc.ColorOff(C_SovControl)
	}
	if box.FirstVP > 0 || box.SecondVP > 0 {
		nc.ColorOn(C_SpaceDefault)
		nc.MovePrint(start.Y+(spaceHeight-2), start.X+2, fmt.Sprintf("%d/%d", box.FirstVP, box.SecondVP))
		nc.ColorOff(C_SpaceDefault)
	}

	nc.ColorOn(C_SpaceDefault)
	offsetY := start.Y + spaceNameY
	for _, line := range wordWrap(box.Name, spaceWidth) {
		nc.MovePrint(offsetY, start.X, line)
		offsetY++
	}
	nc.ColorOff(C_SpaceDefault)
	if box.OpsNeeded > 0 {
		offsetY := start.Y + spaceOpsY
		nc.ColorOn(C_SpaceDetails)
		nc.MovePrint(offsetY, start.X, fmt.Sprintf("%d Ops: 1-%d", box.OpsNeeded, box.MaxRoll))
		offsetY++
		for _, line := range wordWrap(box.SideEffect.String(), spaceWidth) {
			nc.MovePrint(offsetY, start.X, line)
			offsetY++
		}
		nc.ColorOff(C_SpaceDetails)
	}
}

const (
	cardWidth   = 35
	cardsPerRow = 3
	maxHeight   = 36
)

type cardRow []Card

func (cr cardRow) Height() int {
	h := 0
	for _, c := range cr {
		height := cardHeight(c)
		if height > h {
			h = height
		}
	}
	return h
}

func cardRows(cards []Card) []cardRow {
	rows := []cardRow{}
	for i := 0; i < len(cards); i += cardsPerRow {
		end := Min(i+cardsPerRow, len(cards))
		rows = append(rows, cardRow(cards[i:end]))
	}
	return rows
}

func cardHeight(card Card) int {
	lines := wordWrap(card.Text, cardWidth)
	return 2 + len(lines)
}

func (nc *NCursesUI) ShowCards(cards []Card) {
	nc.clear()
	x := 5
	y := 2
	offsetY := 0
	rows := cardRows(cards)
	for _, row := range rows {
		if offsetY+row.Height() > maxHeight {
			break
		}

		for i, c := range row {
			offsetX := i * (cardWidth + 1)
			nc.drawCard(c, Pos{x + offsetX, y + offsetY})
		}

		offsetY += row.Height() + 1

	}
	nc.Refresh()
	nc.Move(37, 0)
}

func (nc *NCursesUI) drawCard(card Card, start Pos) {
	// Ops
	var affColor int16
	switch card.Aff {
	case USA:
		affColor = C_UsaControl
	case SOV:
		affColor = C_SovControl
	default:
		affColor = C_SovInfluence
	}
	var opsStr string
	switch card.Ops {
	case 0:
		opsStr = " - "
	default:
		opsStr = " " + fmt.Sprintf("%1d", card.Ops) + " "
	}
	nc.ColorOn(affColor)
	nc.MovePrint(start.Y, start.X, opsStr)
	nc.ColorOff(affColor)

	// War heading
	var eraColor int16
	switch card.Era {
	case Early:
		eraColor = C_Early
	case Mid:
		eraColor = C_Mid
	default:
		eraColor = C_Late
	}
	nc.ColorOn(eraColor)
	lineFormat := "%-" + strconv.Itoa(cardWidth-len(opsStr)) + "s"
	nc.MovePrint(start.Y, start.X+len(opsStr), fmt.Sprintf(lineFormat, card.Era.String()))
	nc.ColorOff(eraColor)

	// Name
	lineFormat = "%-" + strconv.Itoa(cardWidth) + "s"
	nc.AttrOn(gc.A_BOLD)
	nc.ColorOn(C_CardText)
	nc.MovePrint(start.Y+1, start.X, fmt.Sprintf(lineFormat, card.Name))
	nc.AttrOff(gc.A_BOLD)
	lines := wordWrap(card.Text, cardWidth)
	for i, line := range lines {
		nc.MovePrint(start.Y+2+i, start.X, fmt.Sprintf(lineFormat, line))
	}
	nc.ColorOff(C_CardText)
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
