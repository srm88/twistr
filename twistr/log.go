package twistr

// OurManInTehran: draw top 5 cards, discard any or all

// AskNotWhatYourCountry: discard N from hand, draw replacements

// Versatile log for any card change. Meaning is entirely context-driven.
// Memo field is purely for log human-readability.
// Uses include:
// Cards removed from deck (hand draw)
// Added cards to deck (start of midwar)
// Complete ordering of draw deck (shuffling)
// Drawing a card from discard (SALTNegotiations)
type CardLog struct {
	Memo  string
	Cards []Card
}

type CoupLog struct {
	Country *Country
	Roll    int
}

type RealignLog struct {
	Country *Country
	RollUS  int
	RollSov int
}

type InfluenceLog struct {
	Countries []*Country
}

type CardPlayLog struct {
	Player Aff
	Card   Card
	Kind   PlayKind
}

type SpaceLog struct {
	Roll int
}

type OpsLog struct {
	Kind OpsKind
}

type OpponentOpsLog struct {
	First Aff
}
