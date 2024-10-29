package simulator

import (
	//"fmt"

	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
	"github.com/wade-rees-me/striker-go/cmd/sim/cards"
	"github.com/wade-rees-me/striker-go/cmd/sim/table"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

type Player struct {
	Wager         cards.Wager
	Splits        [constants.MaxSplitHands]cards.Wager
	SplitCount    int
	Rules 		  *table.Rules
	Strategy 	  *table.Strategy
	Report        arguments.Report
	NumberOfCards int
	SeenCards     *[13]int
}

func NewPlayer(rules *table.Rules, strategy *table.Strategy, numberOfCards int) *Player {
	p := new(Player)
	p.Rules = rules
	p.Strategy = strategy
	p.NumberOfCards = numberOfCards
	return p
}

func (p *Player) Shuffle() {
	p.SeenCards = new([13]int)
}

func (p *Player) PlaceBet(mimic bool) {
	p.Wager.Reset()
	for i := 0; i < len(p.Splits); i++ {
		p.Splits[i].Reset()
	}
	p.SplitCount = 0
	if(mimic) {
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
		return
	}

    if(mimic) {
		for !p.MimicStand() {
			p.Wager.Hand.Draw(s.Draw())
		}
        return;
	}

	if p.Strategy.GetDouble(p.SeenCards, p.Wager.Hand.Total(), p.Wager.Hand.Soft(), up) {
		p.Wager.Double()
		p.Wager.Hand.Draw(s.Draw())
		return
	}

	if p.Wager.Hand.Pair() && p.Strategy.GetSplit(p.SeenCards, &p.Wager.Hand.Cards[0], up) {
		split := &p.Splits[p.SplitCount]
		p.SplitCount++
		p.Wager.SplitWager(split)
		if p.Wager.Hand.PairOfAces() {
			p.Wager.Hand.Draw(s.Draw())
			split.Hand.Draw(s.Draw())
			return
		}
		p.Wager.Hand.Draw(s.Draw())
		p.PlaySplit(&p.Wager, s, up)
		split.Hand.Draw(s.Draw())
		p.PlaySplit(split, s, up)
		return
	}

	doStand := p.Strategy.GetStand(p.SeenCards, p.Wager.Hand.Total(), p.Wager.Hand.Soft(), up)
	for !p.Wager.Hand.Busted() && !doStand {
		p.Wager.Hand.Draw(s.Draw())
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
			w.SplitWager(split)
			w.Hand.Draw(shoe.Draw())
			p.PlaySplit(w, shoe, up)
			split.Hand.Draw(shoe.Draw())
			p.PlaySplit(split, shoe, up)
			return
		}
	}

	doStand := p.Strategy.GetStand(p.SeenCards, w.Hand.Total(), w.Hand.Soft(), up)
	for !w.Hand.Busted() && !doStand {
		w.Hand.Draw(shoe.Draw())
		if !w.Hand.Busted() {
			doStand = p.Strategy.GetStand(p.SeenCards, w.Hand.Total(), w.Hand.Soft(), up)
		}
	}
}

func (p *Player) Draw(c *cards.Card) *cards.Card {
	p.Show(c)
	return p.Wager.Hand.Draw(c)
}

func (p *Player) Show(c *cards.Card) {
	p.SeenCards[c.Offset]++
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
		} else {
			w.Lost()
		}
	} else if w.Hand.Blackjack() {
		w.WonBlackjack(int64(p.Rules.BlackjackPays), int64(p.Rules.BlackjackBets))
	} else if w.Hand.Busted() {
		w.Lost()
	} else if dealerBusted || (w.Hand.Total() > dealerTotal) {
		w.Won()
	} else if dealerTotal > w.Hand.Total() {
		w.Lost()
	} else {
		w.Push()
	}
	p.Report.TotalWon += w.AmountWon
	p.Report.TotalBet += w.AmountBet + w.InsuranceBet
}

func (p *Player) payoffSplit(w *cards.Wager, dealerBusted bool, dealerTotal int) {
	if w.Hand.Busted() {
		w.Lost()
	} else if dealerBusted || (w.Hand.Total() > dealerTotal) {
		w.Won()
	} else if dealerTotal > w.Hand.Total() {
		w.Lost()
	} else {
		w.Push()
	}
	p.Report.TotalWon += w.AmountWon
	p.Report.TotalBet += w.AmountBet
	// fmt.Printf("  Payoff.Splits(%d, %d) [%v] %v:%d\n", w.AmountBet, w.AmountWon, w.Hand.Cards, dealerBusted, dealerTotal)
}

func (p *Player) MimicStand() bool {
	if p.Wager.Hand.Soft17() {
		return false
	}
	return p.Wager.Hand.Total() >= 17
}
