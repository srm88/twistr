package twistr

type CoupInput struct {
	Country *Country
	Roll    int
}

type RealignInput struct {
	Country *Country
	RollUS  int
	RollSov int
}

type InfluenceInput struct {
	Countries []*Country
}

type CardPlayInput struct {
	Player Aff
	Card   Card
	Kind   ActionKind
}

type SpaceInput struct {
	Roll int
}

type OpsInput struct {
	Kind OpsKind
}

type OpponentOpsInput struct {
	First Aff
}
