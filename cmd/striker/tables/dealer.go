package tables

import (
	"github.com/wade-rees-me/striker-go/cards"
)

type Dealer struct {
	Hand           Hand
	DealerReport   DealerReport
	TimesHitSoft17 int64
	HitSoft17      bool
}

type DealerReport struct {
	Hand           HandReport
	TimesHitSoft17 int64
}

func NewDealer(hitSoft17 bool) *Dealer {
	d := new(Dealer)
	d.HitSoft17 = hitSoft17
	return d
}

func (d *Dealer) Reset() {
	d.Hand.Reset()
	d.DealerReport.Hand.HandsDealt++
}

func (d *Dealer) Play(s *cards.Shoe) {
	for !d.Stand() {
		d.Draw(s.Draw())
	}
	d.DealerReport.Hand.HandsPlayed++
	if d.Hand.Busted() {
		d.DealerReport.Hand.HandsBusted++
	}
}

func (d *Dealer) Stand() bool {
	if d.HitSoft17 && d.Hand.Soft17() {
		d.DealerReport.TimesHitSoft17++
		return false
	}
	return d.Hand.Total() >= 17
}

func (d *Dealer) Draw(c *cards.Card) *cards.Card {
	return d.Hand.Draw(c)
}

func (d *Dealer) Statistics() {
	if d.Hand.Blackjack() {
		d.DealerReport.Hand.Blackjacks++
	}
}

func (d *Dealer) GetReport() *DealerReport {
	return &d.DealerReport
}
