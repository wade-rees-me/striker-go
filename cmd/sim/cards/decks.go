package cards

type Deck struct {
	Cards []Card
}

var DeckOfPokerCards *Deck

var pokerCards = map[string][]int{
	"two":   {2, 0},
	"three": {3, 1},
	"four":  {4, 2},
	"five":  {5, 3},
	"six":   {6, 4},
	"seven": {7, 5},
	"eight": {8, 6},
	"nine":  {9, 7},
	"ten":   {10, 8},
	"jack":  {10, 9},
	"queen": {10, 10},
	"king":  {10, 11},
	"ace":   {11, 12},
}
var suits = []string{"spades", "diamond", "clubs", "hearts"}

func init() {
	DeckOfPokerCards = newDeck(suits, pokerCards, 1)
}

func newDeck(s []string, r map[string][]int, copies int) *Deck {
	deck := new(Deck)
	for i := 0; i < copies; i++ {
		for _, suit := range s {
			for key, value := range r {
				card := NewCard(suit, key, value[0], value[1])
				deck.Cards = append(deck.Cards, *card)
			}
		}
	}
	return deck
}
