package cards

type Card struct {
	Suit   string // Suit of the card (e.g., "hearts")
	Rank   string // Rank of the card (e.g., "ace")
	Key	   string
	Value  int // Value of the card for game calculations
	Index  int // Index of the card in a deck
}

func NewCard(suit, rank, key string, value int) *Card {
	c := new(Card)
	c.Suit = suit
	c.Rank = rank
	c.Key = key
	c.Value = value
	return c
}

func (c *Card) BlackjackAce() bool {
	return c.Value == 11
}
