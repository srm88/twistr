package twistr

type Card struct {
	Id   CardId
	Aff  Aff
	Ops  int
	Name string
	Text string
	Star bool
}
