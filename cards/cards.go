package cards

const (
	Ace   = "ace"
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
)

const (
	Spades  = "spades"
	Diamond = "diamond"
	Clubs   = "clubs"
	Hearts  = "hearts"
)

type Card struct {
	Suit  string // Suit of the card (e.g., "hearts")
	Rank  string // Rank of the card (e.g., "ace")
	Value int    // Value of the card for game calculations
	Index int    // Index of the card in a deck
}

var Suits = []string{Spades, Diamond, Clubs, Hearts}

func NewCard(suit, rank string, value int) *Card {
	c := new(Card)
	c.Suit = suit
	c.Rank = rank
	c.Value = value
	return c
}
