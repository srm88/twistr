package twistr

import (
	"bytes"
	"io"
	"log"
	"os"
)

type State struct {
	UI
	*Game
	History     *History
	Master      bool
	LocalPlayer Aff
	aof         io.WriteCloser
}

// Checkpoint game. User cannot undo past the point this is called.
func (s *State) Commit() {
	if s.History.InReplay() {
		return
	}
	buffered := s.History.Commit()
	if _, err := s.aof.Write([]byte(buffered)); err != nil {
		log.Fatalf("Failed to flush to aof: %s\n", err.Error())
	}
	s.Redraw(s.Game)
}

func (s *State) CanUndo() bool {
	return s.History.CanPop()
}

func (s *State) Log(thing interface{}) (err error) {
	var b []byte
	if b, err = Marshal(thing); err != nil {
		log.Println(err)
		return
	}
	if _, err = s.History.Write(b); err != nil {
		log.Println(err)
	}
	return
}

func (s *State) ReadInto(thing interface{}) bool {
	ok, line := s.History.Next()
	if !ok {
		return false
	}
	if err := Unmarshal(line, thing); err != nil {
		log.Printf("Corrupt log! Tried to parse '%s' into %s\n", line, thing)
		return false
	}
	return true
}

func (s *State) Undo() {
	s.History.Pop()
	// Totally reset all state, and replay history.
	s.Game = NewGame()
	Start(s)
}

func NewState(ui UI, aofPath string, game *Game) (*State, error) {
	in, err := os.Open(aofPath)
	var history *History
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		history = NewHistory(ui)
	} else {
		b := new(bytes.Buffer)
		if _, err := io.Copy(b, in); err != nil {
			return nil, err
		}
		if err := in.Close(); err != nil {
			return nil, err
		}
		if b.Len() > 0 {
			history = NewHistoryBacklog(ui, b.String())
		} else {
			history = NewHistory(ui)
		}
	}

	out, err := os.OpenFile(aofPath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	s := &State{
		UI:          history,
		Game:        game,
		History:     history,
		Master:      false,
		LocalPlayer: USA,
		aof:         out,
	}
	return s, nil
}

func (s *State) Close() error {
	s.UI.Close()
	return s.aof.Close()
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
