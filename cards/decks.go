package cards

type Deck struct {
	Cards []Card
}

func NewDeck(s []string, r map[string]int, copies int) *Deck {
	deck := new(Deck)
	for i := 0; i < copies; i++ {
		for _, suit := range s {
			for key, value := range r {
				card := NewCard(suit, key, value)
				deck.Cards = append(deck.Cards, *card)
			}
		}
	}
	return deck
}

func (deck *Deck) Search(suit, rank string) []Card {
	var cards []Card
	for _, card := range deck.Cards {
		if (suit == "*" || suit == card.Suit) && (rank == "*" || rank == card.Rank) {
			cards = append(cards, card)
		}
	}
	return cards
}
