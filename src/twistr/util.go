package twistr

import (
	"math/rand"
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Roll() int {
	return rand.Intn(6) + 1
}

func Opp(player Aff) Aff {
	// Relies on neutral being last in the const, i.e. US and Sov are 0 and 1.
	return player ^ 1
}
