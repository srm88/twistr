package twistr

import "fmt"

func Headline(s *State) {
	usaHl := SelectCard(s, USA, CardBlacklist(TheChinaCard))
	s.Txn.Flush()
	sovHl := SelectCard(s, SOV, CardBlacklist(TheChinaCard))
	s.Txn.Flush()
	s.Hands[USA].Remove(usaHl)
	s.Hands[SOV].Remove(sovHl)
	// XXX: If space race show headline is in play, deal with that yo
	MessageBoth(s, fmt.Sprintf("USA selects %s, and USSR selects %s", usaHl.Name, sovHl.Name))
	// Check ops
	if usaHl.Ops >= sovHl.Ops {
		s.Phasing = USA
		PlayEvent(s, USA, usaHl)
		s.Phasing = Sov
		PlayEvent(s, SOV, sovHl)
	} else {
		s.Phasing = Sov
		PlayEvent(s, SOV, sovHl)
		s.Phasing = USA
		PlayEvent(s, USA, usaHl)
	}
	s.Txn.Flush()
}
