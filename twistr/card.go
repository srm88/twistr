package twistr

import (
	"math/rand"
)

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
	cards []*Card
}

func (d *Deck) Shuffle() {
	deckLen := len(d.cards)
	for i := 0; i < 2*deckLen; i++ {
		x := rand.Intn(deckLen)
		d.cards[i], d.cards[x] = d.cards[x], d.cards[i]
	}
}

func (d *Deck) ShuffleIn(cards []*Card) {
	d.cards = append(d.cards, cards...)
	d.Shuffle()
}

func (d *Deck) Push(card *Card) {
	d.cards = append(d.cards, card)
}

func (d *Deck) Draw(n int) (draws []*Card) {
	draws, d.cards = d.cards[:n], d.cards[n:]
	return
}
