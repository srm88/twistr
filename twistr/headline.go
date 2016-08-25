package twistr

import "fmt"

func SelectHeadline(s *State, player Aff) Card {
	return SelectCard(s, player, CardBlacklist(TheChinaCard, UNIntervention))
}

func Headline(s *State) {
	var usaHl, sovHl Card
	headlineSecond, ok := s.SREvents[OppHeadlineFirst]
	switch {
	case ok && headlineSecond == USA:
		s.Transcribe("USSR must choose the headline first")
		sovHl = SelectHeadline(s, SOV)
		s.Transcribe(fmt.Sprintf("USSR selects %s", sovHl.Name))
		s.Commit()
		usaHl = SelectHeadline(s, USA)
		s.Transcribe(fmt.Sprintf("USA selects %s", sovHl.Name))
		s.Commit()
	case ok && headlineSecond == SOV:
		s.Transcribe("US must choose the headline first")
		usaHl = SelectHeadline(s, USA)
		s.Transcribe(fmt.Sprintf("USA selects %s", sovHl.Name))
		s.Commit()
		sovHl = SelectHeadline(s, SOV)
		s.Transcribe(fmt.Sprintf("USSR selects %s", sovHl.Name))
		s.Commit()
	default:
		sovHl = SelectHeadline(s, SOV)
		s.Commit()
		usaHl = SelectHeadline(s, USA)
		s.Commit()
		s.Transcribe(fmt.Sprintf("USA selects %s, and USSR selects %s", usaHl.Name, sovHl.Name))
	}
	s.Hands[USA].Remove(usaHl)
	s.Hands[SOV].Remove(sovHl)
	if usaHl.Id == Defectors {
		s.Transcribe("USSR event canceled by Defectors.")
		s.Discard.Push(usaHl)
		s.Discard.Push(sovHl)
		s.Commit()
		return
	}
	// Check ops
	if usaHl.Ops >= sovHl.Ops {
		s.Phasing = USA
		PlayEvent(s, USA, usaHl)
		s.Phasing = SOV
		PlayEvent(s, SOV, sovHl)
	} else {
		s.Phasing = SOV
		PlayEvent(s, SOV, sovHl)
		s.Phasing = USA
		PlayEvent(s, USA, usaHl)
	}
	// XXX: should probably flush in PlayEvent. Need a more consistent
	// convention for when to flush.
	s.Commit()
}
