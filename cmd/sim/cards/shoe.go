package cards

import (
	"math/rand"
	"time"
)

const (
	MINIMUM_CARD_VALUE = 2
	MAXIMUM_CARD_VALUE = 11
)

// Shoe represents a collection of cards
type Shoe struct {
	cards		  []*Card
	forceShuffle  bool
	NumberOfCards int
	cutCard		  int
	burnCard	  int
	nextCard	  int
	lastDiscard   int
	Random		  *rand.Rand
}

// Suits and card names/constants
var suits = []string{"spades", "diamonds", "clubs", "hearts"}
var cardNames = []string{"two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "jack", "queen", "king", "ace"}
var cardValues = []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10, 11}
var cardKeys = []string{"2", "3", "4", "5", "6", "7", "8", "9", "X", "X", "X", "X", "A"}

// NewShoe creates a new shoe of cards
func NewShoe(numberOfDecks int, penetration float64) *Shoe {
	cards := []*Card{}

	for i := 0; i < numberOfDecks; i++ {
		for _, suit := range suits {
			for j, name := range cardNames {
				card := NewCard(suit, name, cardKeys[j], cardValues[j])
				cards = append(cards, card)
			}
		}
	}

	numberOfCards := len(cards)
	cutCard := int(float64(numberOfCards) * penetration)

	shoe := &Shoe{
		cards:			cards,
		forceShuffle:	false,
		NumberOfCards:	numberOfCards,
		cutCard:		cutCard,
		burnCard:		1,
		nextCard:		numberOfCards,
		lastDiscard:	numberOfCards,
		Random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	//rand.Seed(time.Now().UnixNano())
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shoe.Shuffle()
	return shoe
}

// Shuffle shuffles the cards in the shoe
func (s *Shoe) Shuffle() {
	s.lastDiscard = s.NumberOfCards
	s.forceShuffle = false
	s.shuffleRandom()
}

// shuffleRandom shuffles the cards using the Fisher-Yates algorithm
func (s *Shoe) shuffleRandom() {
	for i := len(s.cards) - 1; i > 0; i-- {
		j := s.Random.Intn(i + 1)
		s.cards[i], s.cards[j] = s.cards[j], s.cards[i]
	}
	s.nextCard = s.burnCard
}

// DrawCard draws the next card from the shoe
func (s *Shoe) Draw() *Card {
	if s.nextCard >= s.NumberOfCards {
		s.forceShuffle = true
		s.shuffleRandom()
	}
	card := s.cards[s.nextCard]
	s.nextCard++
	return card
}

// ShouldShuffle checks if the shoe should be shuffled
func (s *Shoe) ShouldShuffle() bool {
	s.lastDiscard = s.nextCard
	return (s.nextCard >= s.cutCard) || s.forceShuffle
}

// IsAce checks if a card is an Ace
func (c *Card) IsAce() bool {
	return c.Rank == "ace"
}

