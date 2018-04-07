package twistr

import "bytes"
import "fmt"
import "strings"

type Mode interface {
	Display(UI) Mode
	Command(string) bool
}

type LogMode struct {
	lines   []string
	columns int
	rows    int
	start   int
}

func NewLogMode(messages []string) *LogMode {
	m := &LogMode{
		lines:   []string{},
		columns: 100,
		rows:    30,
	}
	for _, msg := range messages {
		wrapped := wordWrap(msg, m.columns)
		rows := len(wrapped)
		if rows > 0 && wrapped[rows-1] == "" {
			wrapped = wrapped[:rows-1]
		}
		m.lines = append(m.lines, wrapped...)
	}
	m.start = Max(0, len(m.lines)-m.rows)
	return m
}

func (m *LogMode) Display(ui UI) Mode {
	ui.ShowMessages(m.lines[m.start:Min(len(m.lines), m.start+m.rows)])
	return m
}

func (m *LogMode) Command(raw string) bool {
	switch raw {
	case "next":
		if m.start+m.rows < len(m.lines) {
			m.start += m.rows
		}
		return true
	case "prev":
		m.start = Max(0, m.start-m.rows)
		return true
	}
	return false
}

type SpaceMode struct {
	spaceRace [2]int
}

func NewSpaceMode(sr [2]int) *SpaceMode {
	return &SpaceMode{sr}
}

func (m *SpaceMode) Display(ui UI) Mode {
	ui.ShowSpaceRace(m.spaceRace)
	return m
}

func (m *SpaceMode) Command(raw string) bool {
	return false
}

type CardMode struct {
	cards []Card
	start int
}

func NewCardMode(cards []Card) *CardMode {
	return &CardMode{cards, 0}
}

func (m *CardMode) Display(ui UI) Mode {
	ui.ShowCards(m.cards[m.start:])
	return m
}

func (m *CardMode) Command(raw string) bool {
	switch raw {
	case "next":
		if m.start+6 < len(m.cards) {
			m.start += 6
		}
		return true
	case "prev":
		m.start = Max(0, m.start-6)
		return true
	}
	return false
}

func parseCommand(raw string) (cmd string, args []string) {
	tokens := strings.Split(raw, " ")
	if len(tokens) == 0 {
		return "", nil
	}
	return tokens[0], tokens[1:]
}

func modal(s *State, command string) bool {
	cmd, args := parseCommand(command)
	switch cmd {
	case "help":
		s.UI.Message("Commands: 'undo' 'hand' 'log' 'spacerace' 'board' 'card <card>'")
	case "hand":
		ShowHand(s, s.LocalPlayer, s.LocalPlayer, true)
	case "log":
		s.Enter(NewLogMode(s.Game.Transcript))
		s.Redraw(s.Game)
	case "spacerace":
		s.Enter(NewSpaceMode(s.Game.SpaceRace))
		s.Redraw(s.Game)
	case "board":
		s.Enter(nil)
		s.Redraw(s.Game)
	case "deck":
		s.Enter(NewCardMode(s.Deck.Cards))
		s.Redraw(s.Game)
	case "opponent":
		if !s.Ability(ViewOpponentHand, s.LocalPlayer) {
			s.UI.Message("Cannot view opponent's hand.")
			return true
		}
		ShowHand(s, s.LocalPlayer.Opp(), s.LocalPlayer, true)
	case "discard":
		if !s.Ability(ViewDiscard, s.LocalPlayer) {
			s.UI.Message("Cannot view discard pile.")
			return true
		}
		s.Enter(NewCardMode(s.Discard.Cards))
		s.Redraw(s.Game)
	case "card":
		if len(args) != 1 {
			break
		}
		card, err := lookupCard(args[0])
		if err != nil {
			s.UI.Message(err.Error())
		}
		s.Enter(NewCardMode([]Card{card}))
		s.Redraw(s.Game)
	case "barf":
		s.History.Dump()
	case "undo":
		if !s.CanUndo() {
			s.UI.Message("Cannot undo the last action.")
			return true
		}
		s.Undo()
		panic("Should never get here!")
	default:
		ret := false
		if s.Mode != nil {
			ret = s.Mode.Command(cmd)
			s.Redraw(s.Game)
		}
		return ret
	}
	return true
}

func localInput(s *State, inp interface{}, message string, choices ...string) {
	var err error
	validChoice := func(in string) bool {
		if len(choices) == 0 {
			return true
		}
		for _, choice := range choices {
			if choice == in {
				return true
			}
		}
		return false
	}
retry:
	inputStr := Solicit(s.UI, message, choices)
	if ok := modal(s, inputStr); ok {
		goto retry
	}
	if len(choices) > 0 && !validChoice(inputStr) {
		err = fmt.Errorf("'%s' is not a valid choice", inputStr)
	} else {
		err = Unmarshal(inputStr, inp)
	}
	if err != nil {
		message = err.Error() + ". Try again?"
		goto retry
	}
}

func input(ui UI, inp interface{}, message string, choices ...string) {
	var err error
	validChoice := func(in string) bool {
		if len(choices) == 0 {
			return true
		}
		for _, choice := range choices {
			if choice == in {
				return true
			}
		}
		return false
	}
retry:
	inputStr := Solicit(ui, message, choices)
	if len(choices) > 0 && !validChoice(inputStr) {
		err = fmt.Errorf("'%s' is not a valid choice", inputStr)
	} else {
		err = Unmarshal(inputStr, inp)
	}
	if err != nil {
		message = err.Error() + ". Try again?"
		goto retry
	}
}

func Solicit(ui UI, message string, choices []string) (reply string) {
	buf := bytes.NewBufferString(strings.TrimRight(message, "\n"))
	if len(choices) > 0 {
		fmt.Fprintf(buf, " [ %s ]", strings.Join(choices, " "))
	}
	ui.Message(buf.String())
	reply, _ = ui.Input()
	return
}

type UI interface {
	Input() (string, error)
	Message(message string) error
	ShowMessages([]string)
	ShowCards([]Card)
	ShowSpaceRace([2]int)
	// Inconsistent. Doesn't always use game -- other modes have their own
	// state instead of relying on parameter.
	Redraw(*Game)
	Close() error
}
