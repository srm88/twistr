package twistr

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type TerminalUI struct {
	in *bufio.Reader
	bytes.Buffer
}

func MakeTerminalUI() *TerminalUI {
	return &TerminalUI{in: bufio.NewReader(os.Stdin)}
}

func (t *TerminalUI) Solicit(player Aff, message string, choices []string) string {
	fmt.Fprintf(t, "[%s] %s", player, strings.TrimRight(message, "\n"))
	if len(choices) > 0 {
		fmt.Fprintf(t, " [ %s ]", strings.Join(choices, " "))
	}
	t.WriteString("\n")
	io.Copy(os.Stdout, t)
	t.Reset()
	text, err := t.in.ReadString('\n')
	if err != nil {
		panic(err.Error())
	}
	return strings.ToLower(strings.TrimSpace(text))
}

func (t *TerminalUI) Message(player Aff, message string) {
	fmt.Fprintf(t, "[%s] %s\n", player, strings.TrimRight(message, "\n"))
	io.Copy(os.Stdout, t)
}
