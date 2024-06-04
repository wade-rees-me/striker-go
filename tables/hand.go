package tables

import (
	"fmt"

	"github.com/wade-rees-me/striker-go/cards"
)

type Hand struct {
	Cards   []cards.Card // Cards in the hand
	total   int          // Total value of the hand
	softAce int          // Number of aces valued as 11
}

type HandReport struct {
	HandsDealt  int64
	Blackjacks  int64
	HandsPlayed int64
	HandsBusted int64
}

func NewHand() *Hand {
	return new(Hand)
}

func (h *Hand) Reset() {
	h.Cards = h.Cards[:0]
	h.total = 0
	h.softAce = 0
}

func (h *Hand) Draw(c *cards.Card) *cards.Card {
	h.Cards = append(h.Cards, *c)
	h.totalHand()
	return c
}

func (h *Hand) Blackjack() bool {
	return len(h.Cards) == 2 && h.total == 21
}

func (h *Hand) Pair() bool {
	return len(h.Cards) == 2 && h.Cards[0].Rank == h.Cards[1].Rank
}

func (h *Hand) PairOf(r string) bool {
	return h.Pair() && (r == h.Cards[0].Rank)
}

func (h *Hand) Busted() bool {
	return h.total > 21
}

func (h *Hand) Soft() bool {
	return h.softAce > 0
}

func (h *Hand) Total() int {
	return h.total
}

func (h *Hand) FirstTwoCardTotal() int {
	return h.total
}

func (h *Hand) GetCard(i int) *cards.Card {
	return &h.Cards[i]
}

func (h *Hand) Soft17() bool {
	return h.Total() == 17 && h.Soft()
}

func (h *Hand) SplitPair() *cards.Card {
	if h.Pair() {
		card := h.Cards[1]
		h.Cards = h.Cards[:1]
		h.totalHand()
		return &card
	}
	panic("Trying to split a non-pair")
}

func (h *Hand) totalHand() {
	h.total = 0
	h.softAce = 0
	for i := 0; i < len(h.Cards); i++ {
		h.total += h.Cards[i].Value
		if h.Cards[i].Value == 11 {
			h.softAce++
		}
	}
	for h.total > 21 && h.softAce > 0 {
		h.total -= 10
		h.softAce--
	}
}

func (h *Hand) ToString() string {
	if h.Soft() {
		return fmt.Sprintf("soft %d", h.Total())
	}
	return fmt.Sprintf("hard %d", h.Total())
}
