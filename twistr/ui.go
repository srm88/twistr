package twistr

import (
	"fmt"
	"strings"
)

const (
	MetaRune = ';'
)

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
		CardMode(s, s.Hands[s.Phasing].Cards)
		return
	case "card":
		if len(args) != 1 {
			break
		}
		card, err := lookupCard(args[0])
		if err != nil {
			s.Message(who, err.Error())
			return
		}
		CardMode(s, []Card{card})
		return
	}
	s.Message(s.Phasing, "Unknown command")
}

func CardMode(s *State, cards []Card) {
	s.ShowCards(s, cards)
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
	ShowCards(*State, []Card)
	Redraw(*State)
	Close() error
}
