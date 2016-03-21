package twistr

import "fmt"

type SpaceId int

const (
	NoAbility SpaceId = iota
	TwoSpace
	OppHeadlineFirst
	DiscardHeld
	ExtraAR
)

type SRBox struct {
	MaxRoll    int
	OpsNeeded  int
	FirstVP    int
	SecondVP   int
	SideEffect SpaceId
}

var SRTrack []SRBox = []SRBox{
	SRBox{0, 0, 0, 0, NoAbility},
	SRBox{3, 2, 2, 1, NoAbility},
	SRBox{4, 2, 0, 0, TwoSpace},
	SRBox{3, 2, 2, 0, NoAbility},
	SRBox{4, 2, 0, 0, OppHeadlineFirst},
	SRBox{3, 3, 3, 1, NoAbility},
	SRBox{4, 3, 0, 0, DiscardHeld},
	SRBox{3, 3, 4, 2, NoAbility},
	SRBox{2, 4, 2, 0, ExtraAR},
}

func CanAdvance(s *State, player Aff, ops int) bool {
	pos := s.SpaceRace[player]
	switch {
	case pos >= len(SRTrack)-1:
		return false
	case SRTrack[pos+1].OpsNeeded > ops:
		return false
	default:
		return true
	}
}

func nextSRBox(s *State, player Aff) (srb SRBox, err error) {
	pos := s.SpaceRace[player]

	if pos >= len(SRTrack)-1 {
		return SRBox{}, fmt.Errorf("Player is already at the last Space Race position")
	}

	return SRTrack[pos+1], err
}

func (srb SRBox) Enter(s *State, player Aff) {
	pos := s.SpaceRace[player]
	oppPos := s.SpaceRace[player.Opp()]

	if pos+1 > oppPos {
		s.GainVP(player, srb.FirstVP)
	} else {
		s.GainVP(player, srb.SecondVP)
	}

	if _, ok := s.SREvents[srb.SideEffect]; ok {
		delete(s.SREvents, srb.SideEffect)
	} else {
		s.SREvents[srb.SideEffect] = player
	}

	s.SpaceRace[player]++
}
