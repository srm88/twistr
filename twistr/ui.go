package twistr

import (
	"fmt"
	"strings"
)

const (
	MetaRune = ';'
)

type Mode interface {
	Display(UI) Mode
	Command(string) Mode
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
		start:   0,
	}
	for _, msg := range messages {
		wrapped := wordWrap(msg, m.columns)
		rows := len(wrapped)
		if rows > 0 && wrapped[rows-1] == "" {
			wrapped = wrapped[:rows-1]
		}
		m.lines = append(m.lines, wrapped...)
	}
	return m
}

func (m *LogMode) Display(ui UI) Mode {
	ui.ShowMessages(m.lines[m.start:Min(len(m.lines), m.start+m.rows)])
	return m
}

func (m *LogMode) Command(raw string) Mode {
	switch raw {
	case "next":
		if m.start+m.rows <= len(m.lines) {
			m.start += m.rows
		}
	case "prev":
		m.start = Max(0, m.start-m.rows)
	}
	return m
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

func (m *CardMode) Command(raw string) Mode {
	switch raw {
	case "next":
		if m.start+6 < len(m.cards) {
			m.start += 6
		}
	case "prev":
		m.start = Max(0, m.start-6)
	}
	return m
}

func parseMeta(input string) (bool, string) {
	if len(input) > 0 && input[0] == MetaRune {
		return true, strings.ToLower(input[1:])
	}
	return false, ""
}

func parseCommand(raw string) (cmd string, args []string) {
	tokens := strings.Split(raw, " ")
	if len(tokens) == 0 {
		return "", nil
	}
	return tokens[0], tokens[1:]
}

func modal(s *State, command string) {
	who := s.Phasing
	cmd, args := parseCommand(command)
	switch cmd {
	case "hand":
		ShowHand(s, s.Phasing, s.Phasing, true)
		return
	case "log":
		s.Enter(NewLogMode(s.Game.Transcript))
		s.Redraw(s.Game)
		return
	case "board":
		s.Enter(nil)
		s.Redraw(s.Game)
		return
	case "deck":
		s.Enter(NewCardMode(s.Deck.Cards))
		s.Redraw(s.Game)
	case "card":
		if len(args) != 1 {
			break
		}
		card, err := lookupCard(args[0])
		if err != nil {
			s.UI.Message(who, err.Error())
			return
		}
		s.Enter(NewCardMode([]Card{card}))
		s.Redraw(s.Game)
		return
	case "barf":
		s.History.Dump()
		return
	case "undo":
		if !s.CanUndo() {
			// XXX message = "Cannot undo."
			return
		}
		s.Undo()
		panic("Should never get here!")
	default:
		if s.Mode != nil {
			s.Enter(s.Mode.Command(cmd))
		}
		s.Redraw(s.Game)
		return
	}
	s.UI.Message(s.Phasing, "Unknown command")
}

func GetInput(s *State, player Aff, inp interface{}, message string, choices ...string) {
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
	inputStr := s.Solicit(player, message, choices)
	if ok, cmd := parseMeta(inputStr); ok {
		modal(s, cmd)
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

type UI interface {
	Solicit(player Aff, message string, choices []string) (reply string)
	Message(player Aff, message string)
	ShowMessages([]string)
	ShowCards([]Card)
	// Inconsistent. Doesn't always use game -- other modes have their own
	// state instead of relying on parameter.
	Redraw(*Game)
	Close() error
}
