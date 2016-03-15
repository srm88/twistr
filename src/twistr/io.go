package twistr

import (
    "bufio"
    "bytes"
    "io"
    "os"
    "strings"
)

type UI interface {
    Solicit(message string, options []string) (reply string)
}

type TerminalUI struct {
    in *bufio.Reader
    buf bytes.Buffer
}

func MakeTerminalUI() *TerminalUI {
    return &TerminalUI{in: bufio.NewReader(os.Stdin)}
}

func (t *TerminalUI) Solicit(message string, options []string) string {
    t.buf.WriteString(message)
    t.buf.WriteString(" [ ")
    t.buf.WriteString(strings.Join(options, " "))
    t.buf.WriteString(" ]\n")
    io.Copy(os.Stdout, &t.buf)
    t.buf.Reset()
    text, err := t.in.ReadString('\n')
    if err != nil {
        panic(err.Error())
    }
    return text
}
