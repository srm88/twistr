package twistr

import (
	"log"
	"os"
)

type Game struct {
	UI
	*State
	aofPath     string
	History     *History
	Master      bool
	LocalPlayer Aff
	Aof         *Aof
	Txn         *TxnLog
}

// Checkpoint game. User cannot rewind past the point this is called.
func (g *Game) Commit() {
	if g.History.InReplay() {
		// Still need to flush the AOF in replay mode, or pending writes will
		// not be written to disk until the first post-replay commit -- this
		// means if the user rewinds again before that next commit, we would
		// lose the *entire* aof.
		// Another thing to address with replay-writes-to-tmp-aof, mv-when-done
		g.Txn.Flush()
		return
	}
	g.Txn.Flush()
	g.History.Commit()
	g.Redraw(g.State)
}

func (g *Game) CanRewind() bool {
	return g.History.CanPop()
}

func (g *Game) ReadInto(thing interface{}) bool {
	// Sorta gross.
	// Read first from history (replay), then aof (initial load)
	// If reading from history, *write* to aof (this should write to a temp aof).
	// If neither history nor aof have any buffered commands, return false.
	var ok bool
	var line string
	ok, line = g.History.Next()
	// Re-log into the AOF!
	if ok {
		log.Printf("Re-logging '%s'\n", line)
		g.Aof.Write(append([]byte(line), '\n'))
	} else {
		ok, line = g.Aof.Next()
		if !ok {
			return false
		}
	}
	if err := Unmarshal(line, thing); err != nil {
		log.Printf("Corrupt log! Tried to parse '%s' into %s\n", line, thing)
		return false
	}
	return true
}

func (g *Game) Rewind() {
	// XXX: durability: should rewrite aof during replay, `mv` in place when done
	g.History.Pop()
	// WOwwwwwwWWW this is nuts!!
	// Drop any buffered aof writes on the floor, truncate the aof on disk,
	// totally reset all state, and REPLAY HISTORY.
	g.Txn.Reset()
	os.Truncate(g.aofPath, 0)
	g.State = NewState()
	Start(g)
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

	history := NewHistory(ui)
	g := &Game{
		UI:          history,
		State:       state,
		aofPath:     aofPath,
		History:     history,
		Master:      false,
		LocalPlayer: USA,
		Aof:         NewAof(in, txn, history),
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
	resetCountries()
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

func (s *State) Clone() *State {
	clone := &State{
		VP:              s.VP,
		Defcon:          s.Defcon,
		MilOps:          [2]int{s.MilOps[0], s.MilOps[1]},
		SpaceRace:       [2]int{s.SpaceRace[0], s.SpaceRace[1]},
		Turn:            s.Turn,
		AR:              s.AR,
		Phasing:         s.Phasing,
		Countries:       cloneCountries(s.Countries),
		Events:          make(map[CardId]Aff),
		TurnEvents:      make(map[CardId]Aff),
		SpaceAttempts:   [2]int{s.SpaceAttempts[0], s.SpaceAttempts[1]},
		SREvents:        make(map[SpaceId]Aff),
		Removed:         s.Removed.Clone(),
		Discard:         s.Discard.Clone(),
		Deck:            s.Deck.Clone(),
		Hands:           [2]*Deck{s.Hands[0].Clone(), s.Hands[1].Clone()},
		ChinaCardPlayer: s.ChinaCardPlayer,
		ChinaCardFaceUp: s.ChinaCardFaceUp,
	}
	for c, a := range s.Events {
		clone.Events[c] = a
	}
	for c, a := range s.TurnEvents {
		clone.TurnEvents[c] = a
	}
	for s, a := range s.SREvents {
		clone.SREvents[s] = a
	}
	return clone
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
