package twistr

import "log"

type History struct {
	wrapped   UI
	inputs    []string
	index     int
	watermark int
	replaying bool
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

func (r *History) Redraw(s *State) {
	if !r.InReplay() {
		r.wrapped.Redraw(s)
	}
}

func (r *History) Close() error {
	return r.wrapped.Close()
}

// Custom
func (r *History) InReplay() bool {
	return r.index < len(r.inputs)
}

func (r *History) Commit() {
	r.watermark = r.index
	log.Printf("Watermarked at %d\n", r.index)
}

func (r *History) CanPop() bool {
	return len(r.inputs) > r.watermark
}

func (r *History) Pop() {
	// Called to begin replay. Discard the most recent input, everything else
	// will be replayed.
	if !r.CanPop() {
		panic("Rewound beyond a checkpoint in the game")
	}
	r.inputs = r.inputs[:len(r.inputs)-1]
	r.index = 0
	r.replaying = true
}

func (r *History) Write(input []byte) (n int, err error) {
	// Never do when replaying ...
	r.inputs = append(r.inputs, string(input))
	r.index = len(r.inputs)
	n = len(input)
	return
}
