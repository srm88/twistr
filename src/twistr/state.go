package twistr

type State struct {
	VP int

	Defcon int

	MilOps [2]int

	SpaceRace [2]int

	Turn int
	AR   int

	Countries map[CountryId]*Country

	Events map[CardId]*Card

	Removed []*Card

	Discard []*Card

	Hands [2]map[CardId]*Card

	ChinaCardPlayer Aff
	ChinaCardFaceUp bool
}
