package twistr

type Card struct {
	Id   CardId
	Aff  Aff
	Ops  int
	Name string
	Text string
	Star bool
}

func (c Card) String() string {
	return c.Name
}

type Deck struct {
	cards []Card
}

func NewDeck() *Deck {
	return &Deck{cards: []Card{}}
}

// Shuffle does not modify the deck in place, but rather returns the new order
// of its cards. Use Reorder to change the deck's order.
func (d *Deck) Shuffle() []Card {
	order := make([]Card, len(d.cards))
	for i, j := range rng.Perm(len(d.cards)) {
		order[i] = d.cards[j]
	}
	return order
}

func (d *Deck) Reorder(ordering []Card) {
	curLen := len(d.cards)
	var i int
	var c Card
	// Assign in-place until we reach the current bound of the deck
	for i, c = range ordering {
		if i == curLen {
			break
		}
		d.cards[i] = c
	}
	// If the ordering introduced more cards, push them on the end
	if i < len(ordering) {
		d.cards = append(d.cards, ordering[i:]...)
	}
}

func (d *Deck) Push(cards ...Card) {
	d.cards = append(d.cards, cards...)
}

func (d *Deck) Draw(n int) (draws []Card) {
	draws, d.cards = d.cards[:n], d.cards[n:]
	return
}
