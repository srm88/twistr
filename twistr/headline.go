package twistr

import "fmt"

func Headline(s *State) {
	usaHl := SelectCard(s, USA, TheChinaCard)
	sovHl := SelectCard(s, SOV, TheChinaCard)
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
}
