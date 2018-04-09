package twistr

import "bufio"
import "io"
import "log"
import "strings"

type CmdOut struct {
	inputs []string
	w      io.Writer
}

func NewCmdOut(w io.Writer) *CmdOut {
	return &CmdOut{
		inputs: []string{},
		w:      w,
	}
}

func (c *CmdOut) Pop() {
	if len(c.inputs) == 0 {
		return
	}
	c.inputs = c.inputs[:len(c.inputs)-1]
}

func (c *CmdOut) Write(input []byte) (n int, err error) {
	lines := inputLines(string(input))
	if len(lines) == 0 {
		return
	}
	c.inputs = append(c.inputs, lines...)
	n = len(input)
	return
}

func (c *CmdOut) Commit() string {
	contents := strings.Join(c.inputs, "\n")
	if len(contents) == 0 {
		return ""
	}
	log.Printf("Committing to remote: %s\n", contents)
	_, err := c.w.Write([]byte(contents + "\n"))
	if err != nil {
		log.Println(err)
	}
	c.inputs = []string{}
	return contents
}

type CmdIn struct {
	*bufio.Scanner
	Inputs     chan string
	KillSwitch chan bool
}

func NewCmdIn(scanner *bufio.Scanner) *CmdIn {
	ci := &CmdIn{
		Scanner: scanner,
		Inputs:  make(chan string, 1000),
	}
	go ci.consume()
	return ci
}

func (ci *CmdIn) consume() {
	defer close(ci.Inputs)
	for ci.Scan() {
		select {
		case <-ci.KillSwitch:
			log.Println("LinkIn received the done signal")
			return
		default:
			ci.Inputs <- ci.Text()
		}
	}
	if err := ci.Err(); err != nil {
		log.Printf("Error after exhausting the CmdIn: %s\n", err.Error())
	}
}

type History struct {
	UI        UI
	inputs    []string
	index     int
	watermark int
	Replaying bool
}

func NewHistoryBacklog(ui UI, backlog []string) *History {
	end := len(backlog)
	for ; end > 1; end-- {
		if backlog[end-1] != "" {
			break
		}
	}
	return &History{
		UI:        ui,
		inputs:    backlog[:end],
		index:     0,
		watermark: end,
		Replaying: true,
	}
}

func NewHistory(ui UI) *History {
	return &History{
		UI:        ui,
		inputs:    []string{},
		index:     0,
		watermark: 0,
		Replaying: false,
	}
}

func (r *History) Dump() {
	log.Printf(">>> DUMP\nindex:     %d\nwatermark: %d\n", r.index, r.watermark)
	for i, l := range r.inputs {
		log.Printf("%3d: %s\n", i, l)
	}
}

func (r *History) Input() (reply string, err error) {
	if !r.InReplay() {
		if r.Replaying {
			r.Replaying = false
		}
		// Passthru
		reply, err = r.UI.Input()
		return
	}
	panic("Can't solicit in replay mode")
}

func (r *History) Next() (bool, string) {
	if !r.InReplay() {
		return false, ""
	}
	reply := r.inputs[r.index]
	r.index++
	return true, reply
}

func (r *History) Message(message string) error {
	if !r.InReplay() {
		return r.UI.Message(message)
	}
	return nil
}

func (r *History) ShowMessages(ms []string) {
	// Future, fix this to get messages from this object
	if !r.InReplay() {
		r.UI.ShowMessages(ms)
	}
}

func (r *History) ShowCards(cs []Card) {
	if !r.InReplay() {
		r.UI.ShowCards(cs)
	}
}

func (r *History) ShowSpaceRace(sr [2]int) {
	if !r.InReplay() {
		r.UI.ShowSpaceRace(sr)
	}
}

func (r *History) Redraw(g *Game) {
	// XXX this is preventing redraw of anything when starting game up
	if !r.InReplay() {
		r.UI.Redraw(g)
	}
}

func (r *History) Close() error {
	return r.UI.Close()
}

// Custom
func (r *History) InReplay() bool {
	return r.index < len(r.inputs)
}

// Checkpoint input at the current point; undo-ing past this point will not
// be allowed.
// Return inputs buffered since last commit as newline-joined string.
func (r *History) Commit() string {
	buffered := strings.Join(r.inputs[r.watermark:], "\n")
	r.watermark = r.index
	log.Printf("Watermarked at %d\n", r.index)
	return buffered
}

func (r *History) CanPop() bool {
	return len(r.inputs) > r.watermark
}

func (r *History) Pop() {
	// Called to begin replay. Discard the most recent input, everything else
	// will be replayed.
	if !r.CanPop() {
		panic("Undid beyond a checkpoint in the game")
	}
	r.inputs = r.inputs[:len(r.inputs)-1]
	r.index = 0
	r.Replaying = true
}

func (r *History) Write(input []byte) (n int, err error) {
	// Never do when replaying. This means s.Log is safe to call on replayed
	// input.
	if r.InReplay() {
		log.Printf("History not writing %s, in replay (InReplay %v, Replaying %v\n", input, r.InReplay(), r.Replaying)
		return
	}
	lines := inputLines(string(input))
	log.Printf("History writing %s\n", lines)
	if len(lines) == 0 {
		return
	}
	r.inputs = append(r.inputs, lines...)
	r.index = len(r.inputs)
	n = len(input)
	return
}

func inputLines(input string) []string {
	lines := strings.Split(input, "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}
