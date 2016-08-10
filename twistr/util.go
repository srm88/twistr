package twistr

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
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
	ui.Message(message)
}

func Debug(message string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, strings.TrimRight(message, "\n")+"\n", a...)
}
