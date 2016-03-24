package twistr

import "fmt"

func Headline(s *State) {
	usaHl := SelectCard(s, USA, TheChinaCard)
	sovHl := SelectCard(s, SOV, TheChinaCard)
	// XXX: If space race show headline is in play, deal with that yo
	MessageBoth(s, fmt.Sprintf("USA selects %s, and SOV selects %s", usaHl.Name, sovHl.Name))
	// Check ops
	if usaHl.Ops >= sovHl.Ops {
		PlayEvent(s, USA, usaHl)
		MessageBoth(s, fmt.Sprintf("USA goes first with %s, Soviets go first with %s", usaHl.Name, sovHl.Name))
	} else {
		PlayEvent(s, SOV, sovHl)
		MessageBoth(s, fmt.Sprintf("Soviets go first with %s, USA goes first with %s", sovHl.Name, usaHl.Name))
	}
}
