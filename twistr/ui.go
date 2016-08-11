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

func (m BoardMode) Display(ui UI) Mode {
	ui.Redraw(m.s)
	return m
}

func (m BoardMode) Command(raw string) Mode {
	return m
}

type LogMode struct {
	messages []string
}

func (m LogMode) Display(ui UI) Mode {
	ui.ShowMessages(m.messages)
	return m
}

func (m LogMode) Command(raw string) Mode {
	return m
}

type CardMode struct {
	cards []Card
}

func (m CardMode) Display(ui UI) Mode {
	ui.ShowCards(m.cards)
	return m
}

func (m CardMode) Command(raw string) Mode {
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
		s.Mode = LogMode{s.Messages}
		s.Redraw(s)
		return
	case "board":
		s.Mode = BoardMode{s}
		s.Redraw(s)
		return
	case "card":
		if len(args) != 1 {
			break
		}
		card, err := lookupCard(args[0])
		if err != nil {
			s.UI.Message(who, err.Error())
			return
		}
		s.Mode = CardMode{[]Card{card}}
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
