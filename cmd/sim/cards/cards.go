package cards

const (
	Spades  = "spades"
	Diamond = "diamond"
	Clubs   = "clubs"
	Hearts  = "hearts"
)

const (
	Two   = "two"
	Three = "three"
	Four  = "four"
	Five  = "five"
	Six   = "six"
	Seven = "seven"
	Eight = "eight"
	Nine  = "nine"
	Ten   = "ten"
	Jack  = "jack"
	Queen = "queen"
	King  = "king"
	Ace   = "ace"
)

type Card struct {
	Suit   string // Suit of the card (e.g., "hearts")
	Rank   string // Rank of the card (e.g., "ace")
	Value  int    // Value of the card for game calculations
	Index  int    // Index of the card in a deck
	Offset int    // Index of the card in a suit
}

var Suits = []string{Spades, Diamond, Clubs, Hearts}

func NewCard(suit, rank string, value, offset int) *Card {
	c := new(Card)
	c.Suit = suit
	c.Rank = rank
	c.Value = value
	c.Offset = offset
	return c
}
