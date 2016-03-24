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

// Prevented returns whether the card's event is prevented from play. E.g.
// "Tear Down this Wall" prevents play of Willy Brandt as an event.
func (c Card) Prevented(s *State) bool {
	switch {
	case c.Id == WillyBrandt && s.Effect(TearDownThisWall):
		return true
	case c.Id == FlowerPower && s.Effect(AnEvilEmpire):
		return true
	case c.Id == ArabIsraeliWar && s.Effect(CampDavidAccords):
		return true
	case c.Id == SocialistGovernments && s.Effect(TheIronLady):
		return true
	case c.Id == OPEC && s.Effect(NorthSeaOil):
		return true
	case c.Id == MuslimRevolution && s.Effect(AWACSSaleToSaudis):
		return true
	case c.Id == NATO && !(s.Effect(MarshallPlan) || s.Effect(WarsawPactFormed)):
		return true
	case c.Id == Solidarity && !s.Effect(JohnPaulIIElectedPope):
		return true
	case c.Id == TheCambridgeFive && s.Era() == Late:
		return true
	default:
		return false
	}
}

type Deck struct {
	Cards []Card
}

func NewDeck() *Deck {
	return &Deck{Cards: []Card{}}
}

// Shuffle does not modify the deck in place, but rather returns the new order
// of its cards. Use Reorder to change the deck's order.
func (d *Deck) Shuffle() []Card {
	order := make([]Card, len(d.Cards))
	for i, j := range rng.Perm(len(d.Cards)) {
		order[i] = d.Cards[j]
	}
	return order
}

func (d *Deck) Reorder(ordering []Card) {
	curLen := len(d.Cards)
	var i int
	var c Card
	// Assign in-place until we reach the current bound of the deck
	for i, c = range ordering {
		if i == curLen {
			break
		}
		d.Cards[i] = c
	}
	// If the ordering introduced more cards, push them on the end
	if i < len(ordering) {
		d.Cards = append(d.Cards, ordering[i:]...)
	}
}

func (d *Deck) Remove(card Card) {
	for i, c := range d.Cards {
		if c == card {
			d.Cards = append(d.Cards[:i], d.Cards[i+1:]...)
		}
	}
}

func (d *Deck) Push(cards ...Card) {
	d.Cards = append(d.Cards, cards...)
}

func (d *Deck) Draw(n int) (draws []Card) {
	draws, d.Cards = d.Cards[:n], d.Cards[n:]
	return
}

func (d *Deck) Names() []string {
	names := make([]string, len(d.Cards))
	for i, card := range d.Cards {
		names[i] = card.Name
	}
	return names
}
