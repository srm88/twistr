package twistr

import (
	"log"
	"os"
)

type State struct {
	UI
	*Game
	aofPath     string
	History     *History
	Master      bool
	LocalPlayer Aff
	Aof         *Aof
	Txn         *TxnLog
}

// Checkpoint game. User cannot rewind past the point this is called.
func (s *State) Commit() {
	if s.History.InReplay() {
		// Still need to flush the AOF in replay mode, or pending writes will
		// not be written to disk until the first post-replay commit -- this
		// means if the user rewinds again before that next commit, we would
		// lose the *entire* aof.
		// Another thing to address with replay-writes-to-tmp-aof, mv-when-done
		s.Txn.Flush()
		return
	}
	s.Txn.Flush()
	s.History.Commit()
	s.Redraw(s.Game)
}

func (s *State) CanRewind() bool {
	return s.History.CanPop()
}

func (s *State) ReadInto(thing interface{}) bool {
	// Sorta gross.
	// Read first from history (replay), then aof (initial load)
	// If reading from history, *write* to aof (this should write to a temp aof).
	// If neither history nor aof have any buffered commands, return false.
	var ok bool
	var line string
	ok, line = s.History.Next()
	// Re-log into the AOF!
	if ok {
		log.Printf("Re-logging '%s'\n", line)
		s.Aof.Write(append([]byte(line), '\n'))
	} else {
		ok, line = s.Aof.Next()
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

func (s *State) Rewind() {
	// XXX: durability: should rewrite aof during replay, `mv` in place when done
	s.History.Pop()
	// WOwwwwwwWWW this is nuts!!
	// Drop any buffered aof writes on the floor, truncate the aof on disk,
	// totally reset all state, and REPLAY HISTORY.
	s.Txn.Reset()
	os.Truncate(s.aofPath, 0)
	s.Game = NewGame()
	Start(s)
}

func NewState(ui UI, aofPath string, game *Game) (*State, error) {
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
	s := &State{
		UI:          history,
		Game:        game,
		aofPath:     aofPath,
		History:     history,
		Master:      false,
		LocalPlayer: USA,
		Aof:         NewAof(in, txn, history),
		Txn:         txn,
	}
	return s, nil
}

func (s *State) Close() error {
	s.UI.Close()
	return s.Aof.Close()
}

type Game struct {
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

func NewGame() *Game {
	resetCountries()
	return &Game{
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

func (s *Game) ImproveDefcon(n int) {
	s.Defcon = Min(s.Defcon+n, 5)
}

func (s *Game) DegradeDefcon(n int) {
	s.Defcon -= n
	if s.Defcon < 2 {
		// XXX writeme
		panic("Thermonuclear war!")
	}
}

func (s *Game) Era() Era {
	switch {
	case s.Turn < 4:
		return Early
	case s.Turn < 8:
		return Mid
	default:
		return Late
	}
}

func (s *Game) ActionsPerTurn() int {
	if s.Era() == Early {
		return 6
	}
	return 7
}

func (s *Game) Effect(which CardId, player ...Aff) bool {
	aff, ok := s.Events[which]
	if ok && (len(player) == 0 || player[0] == aff) {
		return true
	}
	aff, ok = s.TurnEvents[which]
	return ok && (len(player) == 0 || player[0] == aff)
}

// Cancel ends an event.
func (s *Game) Cancel(event CardId) {
	// XXX: this would clobber NorthSeaOil, which registers both a turn-
	// and permanent event.
	delete(s.Events, event)
	delete(s.TurnEvents, event)
}

// CancelTurnEvents cancels all turn-based events currently in effect.
func (s *Game) CancelTurnEvents() {
	s.TurnEvents = make(map[CardId]Aff)
}

func (s *Game) ChinaCardPlayed() {
	s.ChinaCardPlayer = s.ChinaCardPlayer.Opp()
	s.ChinaCardFaceUp = false
}

func (s *Game) GainVP(player Aff, n int) {
	switch player {
	case USA:
		s.VP += n
	case SOV:
		s.VP -= n
	}
}
