package twistr

import (
	"os"
)

type Game struct {
	UI
	*State
	Master      bool
	LocalPlayer Aff
	Aof         *Aof
	Txn         *TxnLog
}

func NewGame(ui UI, aofPath string, state *State) (*Game, error) {
	in, err := os.Open(aofPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		in, err = os.OpenFile(aofPath, os.O_CREATE|os.O_RDONLY, 0666)
		if err != nil {
			return nil, err
		}
	}

	txn, err := OpenTxnLog(aofPath)
	if err != nil {
		return nil, err
	}

	g := &Game{
		UI:          ui,
		State:       state,
		Master:      false,
		LocalPlayer: USA,
		Aof:         NewAof(in, txn),
		Txn:         txn,
	}
	return g, nil
}

func (g *Game) Close() error {
	g.UI.Close()
	return g.Aof.Close()
}

type State struct {
	VP              int
	Defcon          int
	MilOps          [2]int
	SpaceRace       [2]int
	Turn            int
	AR              int
	Phasing         Aff
	Countries       map[CountryId]*Country
	Events          map[CardId]Aff
	TurnEvents      map[CardId]Aff
	SpaceAttempts   [2]int
	SREvents        map[SpaceId]Aff
	Removed         *Deck
	Discard         *Deck
	Deck            *Deck
	Hands           [2]*Deck
	ChinaCardPlayer Aff
	ChinaCardFaceUp bool
	ChernobylRegion Region
}

func NewState() *State {
	return &State{
		VP:              0,
		Defcon:          5,
		MilOps:          [2]int{0, 0},
		SpaceRace:       [2]int{0, 0},
		Turn:            1,
		AR:              1,
		Phasing:         SOV,
		Countries:       Countries,
		Events:          make(map[CardId]Aff),
		TurnEvents:      make(map[CardId]Aff),
		SpaceAttempts:   [2]int{0, 0},
		SREvents:        make(map[SpaceId]Aff),
		Removed:         NewDeck(),
		Discard:         NewDeck(),
		Deck:            NewDeck(),
		Hands:           [2]*Deck{NewDeck(), NewDeck()},
		ChinaCardPlayer: SOV,
		ChinaCardFaceUp: true,
	}
}

func (s *State) ImproveDefcon(n int) {
	s.Defcon = Min(s.Defcon+n, 5)
}

func (s *State) DegradeDefcon(n int) {
	s.Defcon -= n
	if s.Defcon < 2 {
		// XXX writeme
		panic("Thermonuclear war!")
	}
}

func (s *State) Era() Era {
	switch {
	case s.Turn < 4:
		return Early
	case s.Turn < 8:
		return Mid
	default:
		return Late
	}
}

func (s *State) ActionsPerTurn() int {
	if s.Era() == Early {
		return 6
	}
	return 7
}

func (s *State) Effect(which CardId, player ...Aff) bool {
	aff, ok := s.Events[which]
	if ok && (len(player) == 0 || player[0] == aff) {
		return true
	}
	aff, ok = s.TurnEvents[which]
	return ok && (len(player) == 0 || player[0] == aff)
}

// Cancel ends an event.
func (s *State) Cancel(event CardId) {
	// XXX: this would clobber NorthSeaOil, which registers both a turn-
	// and permanent event.
	delete(s.Events, event)
	delete(s.TurnEvents, event)
}

// CancelTurnEvents cancels all turn-based events currently in effect.
func (s *State) CancelTurnEvents() {
	s.TurnEvents = make(map[CardId]Aff)
}

func (s *State) ChinaCardPlayed() {
	s.ChinaCardPlayer = s.ChinaCardPlayer.Opp()
	s.ChinaCardFaceUp = false
}

func (s *State) GainVP(player Aff, n int) {
	switch player {
	case USA:
		s.VP += n
	case SOV:
		s.VP -= n
	}
}
