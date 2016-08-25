package twistr

import (
	"bytes"
	"fmt"
	"go/doc"
	"math/rand"
	"strings"
	"time"
)

var (
	rng *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Mod struct {
	Diff int
	Name string
}

func (m Mod) String() string {
	if m.Diff >= 0 {
		return fmt.Sprintf("+%d (%s)", m.Diff, m.Name)
	} else {
		return fmt.Sprintf("%d (%s)", m.Diff, m.Name)
	}
}

func ModSummary(ms []Mod) string {
	if len(ms) == 0 {
		return ""
	}
	ss := make([]string, len(ms))
	for i, m := range ms {
		ss[i] = m.String()
	}
	return " " + strings.Join(ss, " ")
}

func TotalMod(ms []Mod) (total int) {
	for _, m := range ms {
		total += m.Diff
	}
	return total
}

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
