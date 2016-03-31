package twistr

type CoupLog struct {
	Country *Country
	Roll    int
}

type RealignLog struct {
	Country *Country
	RollUSA int
	RollSOV int
}
