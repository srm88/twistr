package twistr

import (
	"log"
	"strings"
)

type History struct {
	wrapped   UI
	inputs    []string
	index     int
	watermark int
	replaying bool
}

func NewHistoryBacklog(ui UI, backlog string) *History {
	inputs := strings.Split(backlog, "\n")
	if inputs[len(inputs)-1] == "" {
		inputs = inputs[:len(inputs)-1]
	}
	return &History{
		wrapped:   ui,
		inputs:    inputs,
		index:     0,
		watermark: len(inputs),
		replaying: true,
	}
}

func NewHistory(ui UI) *History {
	return &History{
		wrapped:   ui,
		inputs:    []string{},
		index:     0,
		watermark: 0,
		replaying: false,
	}
}

func (r *History) Dump() {
	log.Printf(">>> DUMP\nindex:     %d\nwatermark: %d\n", r.index, r.watermark)
	for i, l := range r.inputs {
		log.Printf("%3d: %s\n", i, l)
	}
}

func (r *History) Solicit(player Aff, message string, choices []string) (reply string) {
	if !r.InReplay() {
		if r.replaying {
			r.replaying = false
		}
		// Passthru
		reply = r.wrapped.Solicit(player, message, choices)
		return
	}
	panic("Can't solicit in replay mode")
}

func (r *History) Next() (bool, string) {
	if !r.InReplay() {
		return false, ""
	}
	reply := r.inputs[r.index]
	log.Printf("Replay %3d: %s\n", r.index, reply)
	r.index++
	return true, reply
}

func (r *History) Message(player Aff, message string) {
	if !r.InReplay() && !r.replaying {
		r.wrapped.Message(player, message)
		return
	}
	log.Printf("Suppress: %s\n", message)
}

func (r *History) ShowMessages(ms []string) {
	// Future, fix this to get messages from this object
	if !r.InReplay() {
		r.wrapped.ShowMessages(ms)
	}
}

func (r *History) ShowCards(cs []Card) {
	if !r.InReplay() {
		r.wrapped.ShowCards(cs)
	}
}

func (r *History) Redraw(g *Game) {
	// XXX this is preventing redraw of anything when starting game up
	if !r.InReplay() {
		r.wrapped.Redraw(g)
	}
}

func (r *History) Close() error {
	return r.wrapped.Close()
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
	r.replaying = true
}

func (r *History) Write(input []byte) (n int, err error) {
	// Never do when replaying ...
	lines := strings.Split(string(input), "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) == 0 {
		return
	}
	r.inputs = append(r.inputs, lines...)
	r.index = len(r.inputs)
	n = len(input)
	return
}
