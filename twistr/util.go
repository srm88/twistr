package twistr

import (
	"math/rand"
	"time"
)

var (
	rng *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
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
	return rng.Intn(6) + 1
}

func MessageBoth(ui UI, message string) {
	ui.Message(USA, message)
	ui.Message(SOV, message)
}
