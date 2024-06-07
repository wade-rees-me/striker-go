package cards

import (
	"math/rand"
)

type Shoe struct {
	cards        []Card // Cards currently in the shoe
	inplay       []Card // Cards that are in play
	discards     []Card // Cards that have been discarded
	ForceShuffle bool   // Flag to force a shuffle
	ShoeReport   ShoeReport
}

type ShoeReport struct {
	NumberOfDecks    int
	NumberOfCards    int
	CutCard          int
	NumberOfShuffles int64
	NumberOutOfCards int64
}

func NewShoe(deck Deck, numberOfDecks int, penetration float64) *Shoe {
	shoe := new(Shoe)
	shoe.ShoeReport.NumberOfDecks = numberOfDecks

	for i := 0; i < shoe.ShoeReport.NumberOfDecks; i++ {
		shoe.discards = append(shoe.discards, deck.Cards...)
	}

	shoe.ShoeReport.NumberOfCards = len(shoe.discards)
	shoe.ShoeReport.CutCard = int(float64(shoe.ShoeReport.NumberOfCards) * penetration)

	for i := 0; i < shoe.ShoeReport.NumberOfCards; i++ {
		shoe.discards[i].Index = i
	}

	return shoe
}

func (s *Shoe) Shuffle() (err error) {
	s.discards = append(s.discards, s.cards...)
	s.discards = append(s.discards, s.inplay...)
	s.cards = nil
	s.inplay = nil
	s.ForceShuffle = false
	return s.ShuffleDiscardsFisherYates()
}

func (s *Shoe) ShuffleDiscards() (err error) {
	s.ForceShuffle = true
	s.ShoeReport.NumberOutOfCards++
	return s.ShuffleDiscardsFisherYates()
}

func (s *Shoe) ShuffleDiscardsFisherYates() (err error) {
	rand.Shuffle(len(s.discards), func(i, j int) { s.discards[i], s.discards[j] = s.discards[j], s.discards[i] })
	s.cards = append(s.cards, s.discards...)
	s.discards = nil
	s.Discard(s.Draw()) // Burn a card
	s.ShoeReport.NumberOfShuffles++
	return nil
}

func (s *Shoe) ShouldShuffle() bool {
	s.discards = append(s.discards, s.inplay...)
	s.inplay = nil
	return (len(s.cards) < (s.ShoeReport.NumberOfCards - s.ShoeReport.CutCard)) || s.ForceShuffle
}

func (s *Shoe) Draw() *Card {
	if len(s.cards) == 0 {
		err := s.ShuffleDiscards()
		if err != nil {
			panic(err)
		}
		if len(s.cards) == 0 {
			panic("shuffle discards")
		}
	}

	card := s.cards[0]
	s.cards = s.cards[1:]
	s.inplay = append(s.inplay, card)
	return &card
}

func (s *Shoe) Discard(card *Card) {
	s.discards = append(s.discards, *card)
	for i, c := range s.inplay {
		if *card == c {
			s.inplay = append(s.inplay[:i], s.inplay[i+1:]...)
		}
	}
}

func (s *Shoe) CardsUndealt() int {
	return len(s.cards)
}

func (s *Shoe) CardsInplay() int {
	return len(s.inplay)
}

func (s *Shoe) CardsDiscarded() int {
	return len(s.discards)
}

func (s *Shoe) GetReport() *ShoeReport {
	return &s.ShoeReport
}
