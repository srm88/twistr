package twistr

import (
	"bytes"
	"fmt"
	"go/doc"
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

func wordWrap(body string, columns int) []string {
	b := new(bytes.Buffer)
	doc.ToText(b, body, "", "", columns)
	return strings.Split(b.String(), "\n")
}

func Roll() int {
	return rng.Intn(6) + 1
}

func MessageBoth(ui UI, message string) {
	ui.Message(USA, message)
	ui.Message(SOV, message)
}

func Debug(message string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, message, a...)
}
