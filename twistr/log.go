package twistr

// Logs must be referenced in this set, or we will not marshal them correctly.
var logTypes map[string]bool = map[string]bool{
	"CoupLog":    true,
	"RealignLog": true,
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
