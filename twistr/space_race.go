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

func (s SpaceId) String() string {
	switch s {
	case NoAbility:
		return ""
	case TwoSpace:
		return "play 2 Space Race cards per turn"
	case OppHeadlineFirst:
		return "choose & show headline card after opponent"
	case DiscardHeld:
		return "discard 1 held card"
	case ExtraAR:
		return "take 8 Action Rounds"
	default:
		return "?"
	}
}

type SRBox struct {
	MaxRoll    int
	OpsNeeded  int
	FirstVP    int
	SecondVP   int
	SideEffect SpaceId
	Name       string
}

var SRTrack []SRBox = []SRBox{
	SRBox{0, 0, 0, 0, NoAbility, "Start"},
	SRBox{3, 2, 2, 1, NoAbility, "Earth Satellite"},
	SRBox{4, 2, 0, 0, TwoSpace, "Animal in Space"},
	SRBox{3, 2, 2, 0, NoAbility, "Man in Space"},
	SRBox{4, 2, 0, 0, OppHeadlineFirst, "Man in Earth Orbit"},
	SRBox{3, 3, 3, 1, NoAbility, "Lunar Orbit"},
	SRBox{4, 3, 0, 0, DiscardHeld, "Eagle/Bear has Landed"},
	SRBox{3, 3, 4, 2, NoAbility, "Space Shuttle"},
	SRBox{2, 4, 2, 0, ExtraAR, "Space Station"},
}

func CanAdvance(s *State, player Aff, ops int) bool {
	pos := s.SpaceRace[player]
	maxAttempts := 1
	if s.SREffect(TwoSpace, player) {
		maxAttempts = 2
	}
	switch {
	case s.SpaceAttempts[player] >= maxAttempts:
		return false
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
	pos := 1 + s.SpaceRace[player]
	oppPos := s.SpaceRace[player.Opp()]

	s.Transcribe(fmt.Sprintf("%s advances to %d on the Space Race.", player, pos))
	if pos > oppPos {
		s.GainVP(player, srb.FirstVP)
	} else {
		s.GainVP(player, srb.SecondVP)
	}

	if _, ok := s.SREvents[srb.SideEffect]; ok {
		if srb.SideEffect != NoAbility {
			s.Transcribe(fmt.Sprintf("%s may no longer %s.", player.Opp(), srb.SideEffect.String()))
		}
		delete(s.SREvents, srb.SideEffect)
	} else {
		if srb.SideEffect != NoAbility {
			s.Transcribe(fmt.Sprintf("%s may now %s.", player, srb.SideEffect.String()))
		}
		s.SREvents[srb.SideEffect] = player
	}
	s.SpaceRace[player]++
}
