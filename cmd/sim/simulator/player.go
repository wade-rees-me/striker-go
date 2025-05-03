package simulator

import (
	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
	"github.com/wade-rees-me/striker-go/cmd/sim/cards"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
	"github.com/wade-rees-me/striker-go/cmd/sim/table"
)

type Player struct {
	Wager         cards.Wager
	Splits        [constants.MaxSplitHands]cards.Wager
	SplitCount    int
	Rules         *table.Rules
	Strategy      *table.Strategy
	Report        arguments.Report
	NumberOfCards int
	SeenCards     *[cards.MAXIMUM_CARD_VALUE + 1]int
}

func NewPlayer(rules *table.Rules, strategy *table.Strategy, numberOfCards int) *Player {
	p := new(Player)
	p.Rules = rules
	p.Strategy = strategy
	p.NumberOfCards = numberOfCards
	p.Wager.InitWager(constants.MinimumBet, constants.MaximumBet)
	for i := 0; i < len(p.Splits); i++ {
		p.Splits[i].InitWager(constants.MinimumBet, constants.MaximumBet)
	}
	return p
}

func (p *Player) Shuffle() {
	p.SeenCards = new([cards.MAXIMUM_CARD_VALUE + 1]int)
}

func (p *Player) PlaceBet(mimic bool) {
	p.Wager.Reset()
	for i := 0; i < len(p.Splits); i++ {
		p.Splits[i].Reset()
	}
	p.SplitCount = 0
	if mimic {
		p.Wager.Bet(int64(constants.MinimumBet))
	} else {
		p.Wager.Bet(int64(p.Strategy.GetBet(p.SeenCards)))
	}
}

func (p *Player) Insurance() {
	if p.Strategy.GetInsurance(p.SeenCards) {
		p.Wager.InsuranceBet = p.Wager.AmountBet / 2
	}
}

func (p *Player) Play(s *cards.Shoe, up *cards.Card, mimic bool) {
	if p.Wager.Hand.Blackjack() {
		p.Report.TotalBlackjacks++
		return
	}

	if mimic {
		for !p.MimicStand() {
			p.Wager.Hand.Draw(s.Draw())
		}
		return
	}

	if p.Strategy.GetDouble(p.SeenCards, p.Wager.Hand.Total(), p.Wager.Hand.Soft(), up) {
		p.Wager.Double()
		p.Draw(&p.Wager.Hand, s)
		p.Report.TotalDoubles++
		return
	}

	if p.Wager.Hand.Pair() && p.Strategy.GetSplit(p.SeenCards, &p.Wager.Hand.Cards[0], up) {
		split := &p.Splits[p.SplitCount]
		p.SplitCount++
		p.Report.TotalSplits++
		if p.Wager.Hand.PairOfAces() {
			p.Report.TotalSplitsAce++
			p.Wager.SplitWager(split)
			p.Draw(&p.Wager.Hand, s)
			p.Draw(&split.Hand, s)
			return
		}
		p.Wager.SplitWager(split)
		p.Draw(&p.Wager.Hand, s)
		p.PlaySplit(&p.Wager, s, up)
		p.Draw(&split.Hand, s)
		p.PlaySplit(split, s, up)
		return
	}

	doStand := p.Strategy.GetStand(p.SeenCards, p.Wager.Hand.Total(), p.Wager.Hand.Soft(), up)
	for !p.Wager.Hand.Busted() && !doStand {
		p.Draw(&p.Wager.Hand, s)
		if !p.Wager.Hand.Busted() {
			doStand = p.Strategy.GetStand(p.SeenCards, p.Wager.Hand.Total(), p.Wager.Hand.Soft(), up)
		}
	}
}

func (p *Player) PlaySplit(w *cards.Wager, shoe *cards.Shoe, up *cards.Card) {
	if w.Hand.Pair() && p.SplitCount < constants.MaxSplitHands {
		if p.Strategy.GetSplit(p.SeenCards, &w.Hand.Cards[0], up) {
			split := &p.Splits[p.SplitCount]
			p.SplitCount++
			p.Report.TotalSplits++
			w.SplitWager(split)
			p.Draw(&w.Hand, shoe)
			p.PlaySplit(w, shoe, up)
			p.Draw(&split.Hand, shoe)
			p.PlaySplit(split, shoe, up)
			return
		}
	}

	doStand := p.Strategy.GetStand(p.SeenCards, w.Hand.Total(), w.Hand.Soft(), up)
	for !w.Hand.Busted() && !doStand {
		p.Draw(&w.Hand, shoe)
		if !w.Hand.Busted() {
			doStand = p.Strategy.GetStand(p.SeenCards, w.Hand.Total(), w.Hand.Soft(), up)
		}
	}
}

func (p *Player) Draw(h *cards.Hand, s *cards.Shoe) *cards.Card {
	card := s.Draw()
	p.Show(card)
	h.Draw(card)
	return card
}

func (p *Player) Show(c *cards.Card) {
	p.SeenCards[c.Value]++
}

func (p *Player) BustedOrBlackjack() bool {
	if p.SplitCount == 0 {
		return p.Wager.Hand.Busted() || p.Wager.Hand.Blackjack()
	}
	if !p.Wager.Hand.Busted() {
		return false
	}
	for i := 0; i < p.SplitCount; i++ {
		if !p.Splits[i].Hand.Busted() {
			return false
		}
	}
	return true
}

func (p *Player) Payoff(dealerBlackjack bool, dealerBusted bool, dealerTotal int) {
	if p.SplitCount == 0 {
		p.payoffHand(&p.Wager, dealerBlackjack, dealerBusted, dealerTotal)
		return
	}

	p.payoffSplit(&p.Wager, dealerBusted, dealerTotal)
	for i := 0; i < p.SplitCount; i++ {
		p.payoffSplit(&p.Splits[i], dealerBusted, dealerTotal)
	}
}

func (p *Player) payoffHand(w *cards.Wager, dealerBlackjack bool, dealerBusted bool, dealerTotal int) {
	if dealerBlackjack {
		w.WonInsurance()
	} else {
		w.LostInsurance()
	}

	if dealerBlackjack {
		if w.Hand.Blackjack() {
			w.Push()
			p.Report.TotalPushes++
		} else {
			w.Lost()
			p.Report.TotalLoses++
		}
	} else if w.Hand.Blackjack() {
		w.WonBlackjack(int64(p.Rules.BlackjackPays), int64(p.Rules.BlackjackBets))
	} else if w.Hand.Busted() {
		w.Lost()
		p.Report.TotalLoses++
	} else if dealerBusted || (w.Hand.Total() > dealerTotal) {
		w.Won()
		p.Report.TotalWins++
	} else if dealerTotal > w.Hand.Total() {
		w.Lost()
		p.Report.TotalLoses++
	} else {
		w.Push()
		p.Report.TotalPushes++
	}
	p.Report.TotalWon += w.AmountWon
	p.Report.TotalBet += w.AmountBet + w.InsuranceBet
}

func (p *Player) payoffSplit(w *cards.Wager, dealerBusted bool, dealerTotal int) {
	if w.Hand.Busted() {
		w.Lost()
		p.Report.TotalLoses++
	} else if dealerBusted || (w.Hand.Total() > dealerTotal) {
		w.Won()
		p.Report.TotalWins++
	} else if dealerTotal > w.Hand.Total() {
		w.Lost()
		p.Report.TotalLoses++
	} else {
		w.Push()
		p.Report.TotalPushes++
	}
	p.Report.TotalWon += w.AmountWon
	p.Report.TotalBet += w.AmountBet
}

func (p *Player) MimicStand() bool {
	if p.Wager.Hand.Soft17() {
		return false
	}
	return p.Wager.Hand.Total() >= 17
}
