package twistr

import (
	"fmt"
	"io"
	"log"
)

type State struct {
	UI
	*Game
	Mode        Mode
	Master      bool
	LocalPlayer Aff
	History     *History
	LinkIn      *CmdIn
	LinkOut     *CmdOut
	Aof         io.Writer
}

// Checkpoint game. User cannot undo past the point this is called.
func (s *State) Commit() {
	if s.History.InReplay() {
		log.Println("Not committing, in replay")
		return
	}
	log.Printf("Committing...\n")
	s.LinkOut.Commit()
	buffered := s.History.Commit()
	if len(buffered) > 0 {
		log.Printf("Writing buffered to aof\n")
		if _, err := s.Aof.Write(append([]byte(buffered), '\n')); err != nil {
			log.Fatalf("Failed to flush to aof: %s\n", err.Error())
		}
	} else {
		log.Printf("nothing buffered to aof\n")
	}
	s.Redraw(s.Game)
}

func (s *State) CanUndo() bool {
	return s.History.CanPop()
}

func (s *State) Log(thing interface{}) (err error) {
	if s.History.InReplay() || s.History.Replaying {
		log.Printf("Not logging, replay %v replaying %v\n", s.History.InReplay(), s.History.Replaying)
		return nil
	}
	var b []byte
	if b, err = Marshal(thing); err != nil {
		log.Println(err)
		return
	}
	log.Printf("Logging %s to history+linkout\n", string(b))
	if _, err = s.History.Write(b); err != nil {
		log.Println(err)
	}
	if _, err = s.LinkOut.Write(append(b, '\n')); err != nil {
		log.Println(err)
	}
	return
}

func (s *State) ReadInto(thing interface{}, fromRemote bool) bool {
	var ok bool
	var line string
	ok, line = s.History.Next()
	if !ok {
		if !fromRemote {
			return false
		}
		// Autocommit to preclude deadlock
		s.Commit()
		ok, line = s.LinkIn.Next()
		if !ok {
			return false
		}
		log.Printf("Read %s in from remote. Writing to history\n", line)
		if _, err := s.History.Write([]byte(line)); err != nil {
			log.Println(err)
		}
	} else {
		log.Printf("Read %s in from history\n", line)
	}

	if err := Unmarshal(line, thing); err != nil {
		log.Printf("Corrupt log! Tried to parse '%s' into %s\n", line, thing)
		return false
	}
	return true
}

func (s *State) Undo() {
	s.History.Pop()
	s.LinkOut.Pop()
	// Totally reset all state, and replay history.
	s.Game = NewGame()
	Start(s)
}

func (s *State) Close() error {
	return s.UI.Close()
}

func (s *State) Redraw(g *Game) {
	// Careful ...
	if s.Mode != nil {
		s.Mode = s.Mode.Display(s.UI)
	} else {
		s.UI.Redraw(g)
	}
}

func (s *State) Message(msg string) {
	s.UI.Message(msg)
}

// Transcribe is used for player-independent, objective, happenings in the
// game.
func (s *State) Transcribe(msg string) {
	s.Game.Transcript = append(s.Game.Transcript, msg)
	s.UI.Message(msg)
}

func (s *State) Enter(o Mode) {
	if s.History.InReplay() {
		return
	}
	s.Mode = o
}

func NewState(history *History, game *Game, isMaster bool, localPlayer Aff, aof io.Writer) *State {
	return &State{
		UI:          history,
		Mode:        nil,
		Game:        game,
		History:     history,
		Master:      isMaster,
		LocalPlayer: localPlayer,
		Aof:         aof,
	}
}

func (s *State) SetDefcon(n int) {
	switch {
	case n > s.Defcon:
		s.ImproveDefcon(n - s.Defcon)
	case n < s.Defcon:
		s.DegradeDefcon(s.Defcon - n)
	default:
		s.Transcribe(fmt.Sprintf("Defcon remains %d.", s.Defcon))
	}
}

func (s *State) ImproveDefcon(n int) {
	newDefcon := Min(s.Defcon+n, 5)
	s.Transcribe(fmt.Sprintf("Defcon improves by %d, now at %d.", newDefcon-s.Defcon, newDefcon))
	s.Defcon = newDefcon
}

func (s *State) DegradeDefcon(n int) {
	s.Defcon -= n
	if s.Defcon < 2 {
		ThermoNuclearWar(s, s.Phasing)
	} else {
		s.Transcribe(fmt.Sprintf("Defcon degrades by %d, now at %d.", n, s.Defcon))
	}

}

func (s *State) TurnEvent(event CardId, player Aff) {
	s.Transcribe(fmt.Sprintf("%s is in effect for the remainder of the turn.", Cards[event]))
	s.TurnEvents[event] = player
}

func (s *State) Event(event CardId, player Aff) {
	s.Transcribe(fmt.Sprintf("%s is now in effect.", Cards[event]))
	s.Events[event] = player
}

// Cancel ends an event.
func (s *State) Cancel(event CardId) {
	if s.Effect(event) {
		s.Transcribe(fmt.Sprintf("%s is canceled.", Cards[event]))
	}
	delete(s.Events, event)
	delete(s.TurnEvents, event)
}

func (s *State) ChinaCardPlayed() {
	if s.ChinaCardPlayer == USA {
		s.Cancel(FormosanResolution)
	}
	s.ChinaCardMove(s.ChinaCardPlayer.Opp(), false)
}

func (s *State) ChinaCardMove(to Aff, faceUp bool) {
	s.ChinaCardPlayer = to
	s.ChinaCardFaceUp = faceUp
	if faceUp {
		s.Transcribe(fmt.Sprintf("%s receives The China Card, face down.", to))
	} else {
		s.Transcribe(fmt.Sprintf("%s receives The China Card, face up.", to))
	}
}

func (s *State) GainVP(player Aff, n int) {
	switch player {
	case USA:
		s.VP += n
		s.Transcribe(fmt.Sprintf("USA gains %d VP, now at %d.", n, s.VP))
	case SOV:
		s.VP -= n
		s.Transcribe(fmt.Sprintf("USSR gains %d VP, now at %d.", n, s.VP))
	}
	if s.VP == 20 {
		AutoWin(s, USA, "20 VP")
	} else if s.VP == -20 {
		AutoWin(s, SOV, "20 VP")
	}
}

func (s *State) AddMilOps(player Aff, n int) {
	s.MilOps[player] += n
	s.Transcribe(fmt.Sprintf("%s adds %d to its Military Operations track.", player, n))
}

// XXX remove
func (s *State) MessageOne(player Aff, message string) error {
	if player != s.LocalPlayer {
		return nil
	}
	return s.UI.Message(message)
}

func (s *State) EnablePlayer(which Ability, player Aff) {
	if player != s.LocalPlayer {
		return
	}
	s.TurnAbilities[player][which] = true
	s.UI.Message(which.Message())
}

func (s *State) CancelAbility(which Ability, player Aff) {
	s.TurnAbilities[player][which] = false
}

type Game struct {
	Transcript      []string
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
	TurnAbilities   [2]map[Ability]bool
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
		Transcript:      []string{},
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
		TurnAbilities:   [2]map[Ability]bool{make(map[Ability]bool), make(map[Ability]bool)},
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

func (s *Game) SREffect(which SpaceId, player Aff) bool {
	aff, ok := s.SREvents[which]
	return ok && player == aff
}

func (s *Game) Effect(which CardId, player ...Aff) bool {
	aff, ok := s.Events[which]
	if ok && (len(player) == 0 || player[0] == aff) {
		return true
	}
	aff, ok = s.TurnEvents[which]
	return ok && (len(player) == 0 || player[0] == aff)
}

func (s *Game) Ability(which Ability, player Aff) bool {
	return s.TurnAbilities[player][which]
}

func (s *Game) CancelTurnAbilities() {
	s.TurnAbilities[USA] = make(map[Ability]bool)
	s.TurnAbilities[SOV] = make(map[Ability]bool)
}

// CancelTurnEvents cancels all turn-based events currently in effect.
func (s *Game) CancelTurnEvents() {
	s.TurnEvents = make(map[CardId]Aff)
}
