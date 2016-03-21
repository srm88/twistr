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
	RollUSA int
	RollSOV int
}

type InfluenceLog struct {
	Countries []*Country
}

type CardLog struct {
	Card Card
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
