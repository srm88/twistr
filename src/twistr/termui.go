package twistr

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

type TerminalUI struct {
	in  *bufio.Reader
	buf bytes.Buffer
}

func MakeTerminalUI() *TerminalUI {
	return &TerminalUI{in: bufio.NewReader(os.Stdin)}
}

func (t *TerminalUI) Solicit(message string, choices []string) string {
	t.buf.WriteString(message)
	if len(choices) > 0 {
		t.buf.WriteString(" [ ")
		t.buf.WriteString(strings.Join(choices, " "))
		t.buf.WriteString(" ]")
	}
	t.buf.WriteString("\n")
	io.Copy(os.Stdout, &t.buf)
	t.buf.Reset()
	text, err := t.in.ReadString('\n')
	if err != nil {
		panic(err.Error())
	}
	return strings.ToLower(strings.TrimSpace(text))
}
