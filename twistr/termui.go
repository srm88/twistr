package twistr

import "bufio"
import "fmt"
import "os"
import "strings"

type TerminalUI struct {
	in *bufio.Reader
}

func MakeTerminalUI() *TerminalUI {
	return &TerminalUI{in: bufio.NewReader(os.Stdin)}
}

func (t *TerminalUI) Input() (string, error) {
	text, err := t.in.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.ToLower(strings.TrimSpace(text)), nil
}

func (t *TerminalUI) Message(message string) error {
	_, err := fmt.Fprintf(os.Stdout, "%s\n", strings.TrimRight(message, "\n"))
	return err
}

func (t *TerminalUI) Redraw(s *State) error {
	return nil
}
