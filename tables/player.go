package tables

import (
	"strings"

	"github.com/wade-rees-me/striker-go/cards"
	"github.com/wade-rees-me/striker-go/database"
)

const MaxSplitHands = 3

type PlayStrategy struct {
	HardDouble map[int][]string
	SoftDouble map[int][]string
	PairSplit  map[int][]string
	HardStand  map[int][]string
	SoftStand  map[int][]string
}

type Player struct {
	TableRules   *database.TableRules
	PlayStrategy PlayStrategy
	Wager        Wager
	Splits       [MaxSplitHands]Wager
	splitCount   int
	PlayerReport PlayerReport
}

type WagerReport struct {
	TotalBet  int64
	TotalWon  int64
	Advantage float64
}

type PlayerReport struct {
	Wager WagerReport
	Hand  HandReport
}

func NewPlayer(tr *database.TableRules, playStrategy *PlayStrategy) *Player {
	p := new(Player)
	p.TableRules = tr
	p.PlayStrategy = *playStrategy
	return p
}

func NewPlayerReport() *PlayerReport {
	p := new(PlayerReport)
	return p
}

func (p *Player) PlaceBet(b int) {
	p.Wager.Reset()
	p.Wager.Bet(b)
	p.splitCount = 0
	p.PlayerReport.Hand.HandsDealt++
}

func (p *Player) Play(s *cards.Shoe, up *cards.Card) {
	p.PlayerReport.Hand.HandsPlayed++
	if p.Wager.Hand.Blackjack() {
		return
	}
	if p.double(&p.Wager, up) {
		if p.TableRules.DoubleAnyTwoCards || (p.Wager.Hand.Total() == 10 || p.Wager.Hand.Total() == 11) {
			p.Wager.Double()
			p.Wager.Hand.Draw(s.Draw())
			return
		}
	}
	if p.split(&p.Wager, up) {
		split := p.Splits[p.splitCount]
		p.splitCount++
		if !p.TableRules.ResplitAces && !p.TableRules.HitSplitAces && p.Wager.Hand.PairOf(cards.Ace) {
			p.Wager.SplitWager(&split)
			p.Wager.Hand.Draw(s.Draw())
			split.Hand.Draw(s.Draw())
			return
		}
		p.Wager.SplitWager(&split)
		p.Wager.Hand.Draw(s.Draw())
		p.PlaySplit(&p.Wager, s, up)
		split.Hand.Draw(s.Draw())
		p.PlaySplit(&split, s, up)
		return
	}
	for !p.stand(&p.Wager, up) {
		p.Wager.Hand.Draw(s.Draw())
	}
	if p.Wager.Hand.Busted() {
		p.PlayerReport.Hand.HandsBusted++
	}
}

func (p *Player) PlaySplit(h *Wager, s *cards.Shoe, up *cards.Card) {
	if p.TableRules.DoubleAfterSplit && (p.TableRules.HitSplitAces || !p.Wager.Hand.PairOf(cards.Ace)) && p.double(h, up) {
		p.Wager.Double()
		p.Wager.Hand.Draw(s.Draw())
		return
	}
	if (p.TableRules.ResplitAces || !p.Wager.Hand.PairOf(cards.Ace)) && p.split(h, up) {
		split := p.Splits[p.splitCount]
		p.splitCount++
		h.Hand.Draw(s.Draw())
		p.PlaySplit(h, s, up)
		split.Hand.Draw(s.Draw())
		p.PlaySplit(&split, s, up)
		return
	}
	if p.TableRules.HitSplitAces || !p.Wager.Hand.PairOf(cards.Ace) {
		for !p.stand(h, up) {
			h.Hand.Draw(s.Draw())
		}
	}
	if p.Wager.Hand.Busted() {
		p.PlayerReport.Hand.HandsBusted++
	}
}

func (p *Player) Draw(c *cards.Card) *cards.Card {
	return p.Wager.Hand.Draw(c)
}

func (p *Player) BustedOrBlackjack() bool {
	return p.Wager.Hand.Busted() || p.Wager.Hand.Blackjack()
}

func (p *Player) Payoff(blackjack bool, busted bool, total int) {
	if p.Wager.Hand.Blackjack() {
		p.PlayerReport.Hand.Blackjacks++
	}
	p.payoffHand(&p.Wager, blackjack, busted, total)
	for i := 0; i < p.splitCount; i++ {
		p.payoffHand(&p.Splits[i], blackjack, busted, total)
	}
}

func (p *Player) payoffHand(w *Wager, blackjack bool, busted bool, total int) {
	if w.Hand.Blackjack() {
		if !blackjack {
			w.WonBlackjack()
		}
	} else if w.Hand.Busted() {
		w.Lost()
		p.PlayerReport.Hand.HandsBusted++
	} else if busted || (w.Hand.Total() > total) {
		w.Won()
	} else if total > w.Hand.Total() {
		w.Lost()
	}

	p.PlayerReport.Wager.TotalBet += int64(w.AmountBet + w.DoubleBet)
	p.PlayerReport.Wager.TotalWon += int64(w.AmountWon)
}

func (p *Player) double(h *Wager, up *cards.Card) bool {
	if h.Hand.Soft() {
		return p.strategyHelper2(p.PlayStrategy.SoftDouble, h.Hand.Total(), up, false)
	}
	return p.strategyHelper2(p.PlayStrategy.HardDouble, h.Hand.Total(), up, false)
}

func (p *Player) split(h *Wager, up *cards.Card) bool {
	if p.splitCount < MaxSplitHands && h.Hand.Pair() {
		return p.strategyHelper2(p.PlayStrategy.PairSplit, h.Hand.Total(), up, true)
	}
	return false
}

func (p *Player) stand(h *Wager, up *cards.Card) bool {
	if h.Hand.Soft() {
		return p.strategyHelper2(p.PlayStrategy.SoftStand, h.Hand.Total(), up, true)
	}
	return p.strategyHelper2(p.PlayStrategy.HardStand, h.Hand.Total(), up, true)
}

func (p *Player) strategyHelper2(strategyMap map[int][]string, total int, up *cards.Card, defaultValue bool) bool {
	row, ok := strategyMap[total]
	if ok {
		return strings.ToLower(row[up.Value]) == "yes"
	}
	return defaultValue
}

func (p *Player) GetReport() *PlayerReport {
	p.PlayerReport.Wager.Advantage = float64(p.PlayerReport.Wager.TotalBet) / float64(p.PlayerReport.Wager.TotalWon) / 100.0
	return &p.PlayerReport
}
