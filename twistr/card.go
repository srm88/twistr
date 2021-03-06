package twistr

type Card struct {
	Id   CardId
	Aff  Aff
	Ops  int
	Name string
	Text string
	Star bool
	Era  Era
	Impl func(*State, Aff)
}

func (c Card) Equal(other Card) bool {
	return c.Id == other.Id
}

func (c Card) String() string {
	return c.Name
}

func (c Card) Scoring() bool {
	return c.Ops == 0
}

func (c Card) IsWar() bool {
	switch c.Id {
	case ArabIsraeliWar, KoreanWar, BrushWar, IndoPakistaniWar, IranIraqWar:
		return true
	default:
		return false
	}
}

func (c Card) ScoringRegion() Region {
	switch c.Id {
	case AsiaScoring:
		return Asia
	case EuropeScoring:
		return Europe
	case MiddleEastScoring:
		return MiddleEast
	case CentralAmericaScoring:
		return CentralAmerica
	case SouthAmericaScoring:
		return SouthAmerica
	case SoutheastAsiaScoring:
		return SoutheastAsia
	case AfricaScoring:
		return Africa
	default:
		return Region{}
	}
}

// Prevented returns whether the card's event is prevented from play. E.g.
// "Tear Down this Wall" prevents play of Willy Brandt as an event.
func (c Card) Prevented(g *Game) bool {
	switch {
	case c.Id == WillyBrandt && g.Effect(TearDownThisWall):
		return true
	case c.Id == FlowerPower && g.Effect(AnEvilEmpire):
		return true
	case c.Id == ArabIsraeliWar && g.Effect(CampDavidAccords):
		return true
	case c.Id == SocialistGovernments && g.Effect(TheIronLady):
		return true
	case c.Id == OPEC && g.Effect(NorthSeaOil):
		return true
	case c.Id == MuslimRevolution && g.Effect(AWACSSaleToSaudis):
		return true
	case c.Id == NATO && !(g.Effect(MarshallPlan) || g.Effect(WarsawPactFormed)):
		return true
	case c.Id == Solidarity && !g.Effect(JohnPaulIIElectedPope):
		return true
	case c.Id == TheCambridgeFive && g.Era() == Late:
		return true
	case c.Id == Wargames && g.Defcon != 2:
		return true
	default:
		return false
	}
}

func (c Card) Ref() string {
	return cardNameLookup[c.Id]
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
			// If the ordering introduced more cards, push them on the end
			d.Cards = append(d.Cards, ordering[i:]...)
			return
		}
		d.Cards[i] = c
	}
}

func (d *Deck) Remove(card Card) {
	for i, c := range d.Cards {
		if c.Equal(card) {
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

// Create a cardFilter that rejects specific cards.
func CardBlacklist(blacklist ...CardId) cardFilter {
	return func(c Card) bool {
		for _, bad := range blacklist {
			if bad == c.Id {
				return false
			}
		}
		return true
	}
}
