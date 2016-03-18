package twistr

// Complete ordering of draw deck (shuffling)
type DeckShuffleLog struct {
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
