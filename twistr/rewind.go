package twistr

type History struct {
	wrapped   UI
	Inputs    []string
	index     int
	watermark int
}

func NewHistory(ui UI) *History {
	return &History{
		wrapped:   ui,
		Inputs:    []string{},
		index:     0,
		watermark: 0,
	}
}

func (r *History) Solicit(player Aff, message string, choices []string) (reply string) {
	if !r.InReplay() {
		// Passthru
		reply = r.wrapped.Solicit(player, message, choices)
		return
	}
	reply = r.Inputs[r.index]
	r.index++
	return
}

func (r *History) Message(player Aff, message string) {
	if !r.InReplay() {
		r.wrapped.Message(player, message)
	}
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
	return r.index < len(r.Inputs)
}

func (r *History) Commit() {
	r.watermark = r.index
}

func (r *History) CanPop() bool {
	return len(r.Inputs) > r.watermark
}

func (r *History) Pop() {
	// Called to begin replay. Discard the most recent input, everything else
	// will be replayed.
	if !r.CanPop() {
		panic("Rewound beyond a checkpoint in the game")
	}
	r.Inputs = r.Inputs[:len(r.Inputs)-1]
	r.index = 0
}

func (r *History) Write(input []byte) (n int, err error) {
	// Never do when replaying ...
	r.Inputs = append(r.Inputs, string(input))
	r.index = len(r.Inputs)
	n = len(input)
	return
}
