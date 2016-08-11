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

type BoardMode struct {
	s *State
}

func NewBoardMode(s *State) *BoardMode {
	return &BoardMode{s}
}

func (m *BoardMode) Display(ui UI) Mode {
	ui.Redraw(m.s)
	return m
}

func (m *BoardMode) Command(raw string) Mode {
	return m
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
		m.lines = append(m.lines, wordWrap(msg, m.columns)...)
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
		if m.start+6 <= len(m.cards) {
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
		s.Mode = NewLogMode(s.Messages)
		s.Redraw(s)
		return
	case "board":
		s.Mode = NewBoardMode(s)
		s.Redraw(s)
		return
	case "deck":
		s.Mode = NewCardMode(s.Deck.Cards)
		s.Redraw(s)
	case "card":
		if len(args) != 1 {
			break
		}
		card, err := lookupCard(args[0])
		if err != nil {
			s.UI.Message(who, err.Error())
			return
		}
		s.Mode = NewCardMode([]Card{card})
		s.Redraw(s)
		return
	default:
		s.Mode = s.Mode.Command(command)
		s.Redraw(s)
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
	Redraw(*State)
	Close() error
}
