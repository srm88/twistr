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
