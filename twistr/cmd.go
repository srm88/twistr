package twistr

type CoupCommand struct {
	Country *Country
	Roll    int
}

type RealignCommand struct {
	Country *Country
	RollUS  int
	RollSov int
}

type InfluenceCommand struct {
	Countries []*Country
}

type SpaceCommand struct {
	Roll int
}

type EventCommand struct {
	Player Aff
	Card   Card
	// XXX
}
