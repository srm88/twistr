package twistr

import "fmt"

func Headline(s *State) {
	usaHl := SelectCard(s, USA, CardBlacklist(TheChinaCard))
    s.Txn.Flush()
	sovHl := SelectCard(s, SOV, CardBlacklist(TheChinaCard))
    s.Txn.Flush()
	// XXX: If space race show headline is in play, deal with that yo
	MessageBoth(s, fmt.Sprintf("USA selects %s, and USSR selects %s", usaHl.Name, sovHl.Name))
	// Check ops
	if usaHl.Ops >= sovHl.Ops {
		PlayEvent(s, USA, usaHl)
		PlayEvent(s, SOV, sovHl)
	} else {
		PlayEvent(s, SOV, sovHl)
		PlayEvent(s, USA, usaHl)
	}
    s.Txn.Flush()
}
