package twistr

import "fmt"

func Headline(s *Game) {
	// XXX DEFECTORS
	var usaHl, sovHl Card
	headlineSecond, ok := s.SREvents[OppHeadlineFirst]
	switch {
	case ok && headlineSecond == USA:
		MessageBoth(s, "USSR must choose the headline first")
		sovHl = SelectCard(s, SOV, CardBlacklist(TheChinaCard))
		MessageBoth(s, fmt.Sprintf("USSR selects %s", sovHl.Name))
		s.Txn.Flush()
		usaHl = SelectCard(s, USA, CardBlacklist(TheChinaCard))
		MessageBoth(s, fmt.Sprintf("USA selects %s", sovHl.Name))
		s.Txn.Flush()
	case ok && headlineSecond == SOV:
		MessageBoth(s, "US must choose the headline first")
		usaHl = SelectCard(s, USA, CardBlacklist(TheChinaCard))
		MessageBoth(s, fmt.Sprintf("USA selects %s", sovHl.Name))
		s.Txn.Flush()
		sovHl = SelectCard(s, SOV, CardBlacklist(TheChinaCard))
		MessageBoth(s, fmt.Sprintf("USSR selects %s", sovHl.Name))
		s.Txn.Flush()
	default:
		sovHl = SelectCard(s, SOV, CardBlacklist(TheChinaCard))
		s.Txn.Flush()
		usaHl = SelectCard(s, USA, CardBlacklist(TheChinaCard))
		s.Txn.Flush()
		MessageBoth(s, fmt.Sprintf("USA selects %s, and USSR selects %s", usaHl.Name, sovHl.Name))
	}
	s.Hands[USA].Remove(usaHl)
	s.Hands[SOV].Remove(sovHl)
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
	s.Txn.Flush()
}
