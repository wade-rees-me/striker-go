package cards

import (
	"fmt"
)

type Hand struct {
	Cards     []Card // Cards in the hand
	HandTotal int    // Total value of the hand
	SoftAce   int    // Number of aces valued as 11
	Surrender bool
}

func NewHand() *Hand {
	return new(Hand)
}

func (h *Hand) Reset() {
	h.Cards = h.Cards[:0]
	h.HandTotal = 0
	h.SoftAce = 0
	h.Surrender = false
}

func (h *Hand) Draw(c *Card) *Card {
	h.Cards = append(h.Cards, *c)
	h.totalHand()
	return c
}

func (h *Hand) Blackjack() bool {
	return len(h.Cards) == 2 && h.HandTotal == 21
}

func (h *Hand) Pair() bool {
	return len(h.Cards) == 2 && h.Cards[0].Rank == h.Cards[1].Rank
}

func (h *Hand) PairOf(r string) bool {
	return h.Pair() && (r == h.Cards[0].Rank)
}

func (h *Hand) Busted() bool {
	return h.HandTotal > 21
}

func (h *Hand) Soft() bool {
	return h.SoftAce > 0
}

func (h *Hand) Total() int {
	return h.HandTotal
}

func (h *Hand) FirstTwoCardTotal() int {
	return h.HandTotal
}

func (h *Hand) GetCard(i int) *Card {
	return &h.Cards[i]
}

func (h *Hand) Soft17() bool {
	return h.Total() == 17 && h.Soft()
}

func (h *Hand) SplitPair() *Card {
	if h.Pair() {
		card := h.Cards[1]
		h.Cards = h.Cards[:1]
		h.totalHand()
		return &card
	}
	panic("Trying to split a non-pair")
}

func (h *Hand) totalHand() {
	h.HandTotal = 0
	h.SoftAce = 0
	for i := 0; i < len(h.Cards); i++ {
		h.HandTotal += h.Cards[i].Value
		if h.Cards[i].Value == 11 {
			h.SoftAce++
		}
	}
	for h.HandTotal > 21 && h.SoftAce > 0 {
		h.HandTotal -= 10
		h.SoftAce--
	}
}

func (h *Hand) ToString() string {
	if h.Soft() {
		return fmt.Sprintf("soft %d", h.Total())
	}
	return fmt.Sprintf("hard %d", h.Total())
}
