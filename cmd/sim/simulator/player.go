package simulator

import (
	//"fmt"

	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
	"github.com/wade-rees-me/striker-go/cmd/sim/cards"
	//"github.com/wade-rees-me/striker-go/cmd/sim/table"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

type Player struct {
	Wager         cards.Wager
	Splits        [constants.MaxSplitHands]cards.Wager
	SplitCount    int
	Parameters    *arguments.Parameters
	Report        arguments.Report
	NumberOfCards int
	SeenCards     *[13]int
}

func NewPlayer(parameters *arguments.Parameters, numberOfCards int) *Player {
	p := new(Player)
	p.Parameters = parameters
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
		p.Wager.Bet(int64(p.GetBet()))
	}
}

func (p *Player) Insurance() {
	if p.GetInsurance() {
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

	haveCards := getHave(&p.Wager.Hand)
	//doSurrender := p.GetSurrender(haveCards, up.Offset)
	//if doSurrender {
	if p.GetSurrender(haveCards, up.Offset) {
		p.Wager.Hand.Surrender = true
		return
	}

	doDouble := p.GetDouble(haveCards, up.Offset)
	if doDouble && (p.Parameters.Rules.DoubleAnyTwoCards || (p.Wager.Hand.Total() == 10 || p.Wager.Hand.Total() == 11)) {
		p.Wager.Double()
		p.Wager.Hand.Draw(s.Draw())
		return
	}

	if p.Wager.Hand.Pair() && p.GetSplit(p.Wager.Hand.Cards[0].Value, up.Offset) {
		split := &p.Splits[p.SplitCount]
		p.SplitCount++
		if p.Wager.Hand.PairOfAces() {
			if !p.Parameters.Rules.ResplitAces && !p.Parameters.Rules.HitSplitAces {
				p.Wager.SplitWager(split)
				p.Wager.Hand.Draw(s.Draw())
				split.Hand.Draw(s.Draw())
				return
			}
		}
		p.Wager.SplitWager(split)
		p.Wager.Hand.Draw(s.Draw())
		p.PlaySplit(&p.Wager, s, up)
		split.Hand.Draw(s.Draw())
		p.PlaySplit(split, s, up)
		return
	}

	doStand := p.GetStand(haveCards, up.Offset)
	for !p.Wager.Hand.Busted() && !doStand {
		p.Wager.Hand.Draw(s.Draw())
		doStand = p.GetStand(getHave(&p.Wager.Hand), up.Offset)
	}
}

func (p *Player) PlaySplit(w *cards.Wager, s *cards.Shoe, up *cards.Card) {
	haveCards := getHave(&w.Hand)
	if p.Parameters.Rules.DoubleAfterSplit && p.GetDouble(haveCards, up.Offset) {
		w.Double()
		w.Hand.Draw(s.Draw())
		return
	}

	if w.Hand.Pair() && p.SplitCount < constants.MaxSplitHands {
		if p.GetSplit(w.Hand.Cards[0].Value, up.Offset) {
			if !w.Hand.PairOfAces() || (w.Hand.PairOfAces() && p.Parameters.Rules.ResplitAces) {
				split := &p.Splits[p.SplitCount]
				p.SplitCount++
				w.SplitWager(split)
				w.Hand.Draw(s.Draw())
				p.PlaySplit(w, s, up)
				split.Hand.Draw(s.Draw())
				p.PlaySplit(split, s, up)
				return
			}
		}
	}

	if w.Hand.Cards[0].BlackjackAce() && !p.Parameters.Rules.HitSplitAces {
		return
	}

	doStand := p.GetStand(haveCards, up.Offset)
	for !w.Hand.Busted() && !doStand {
		w.Hand.Draw(s.Draw())
		doStand = p.GetStand(getHave(&w.Hand), up.Offset)
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

	if w.Hand.Surrender {
		p.Report.TotalWon -= w.AmountBet / 2
	} else {
		if dealerBlackjack {
			if w.Hand.Blackjack() {
				w.Push()
			} else {
				w.Lost()
			}
		} else if w.Hand.Blackjack() {
			w.WonBlackjack(int64(p.Parameters.Rules.BlackjackPays), int64(p.Parameters.Rules.BlackjackBets))
		} else if w.Hand.Busted() {
			w.Lost()
		} else if dealerBusted || (w.Hand.Total() > dealerTotal) {
			w.Won()
		} else if dealerTotal > w.Hand.Total() {
			w.Lost()
		} else {
			w.Push()
		}
		p.Report.TotalWon += (w.AmountWon + w.DoubleWon)
	}
	p.Report.TotalBet += w.AmountBet + w.DoubleBet + w.InsuranceBet
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
	p.Report.TotalWon += (w.AmountWon + w.DoubleWon)
	p.Report.TotalBet += w.AmountBet + w.DoubleBet
	// fmt.Printf("  Payoff.Splits(%d, %d) [%v] %v:%d\n", w.AmountBet + w.DoubleBet, (w.AmountWon + w.DoubleWon), w.Hand.Cards, dealerBusted, dealerTotal)
}

func getHave(hand *cards.Hand) *[13]int {
	haveCards := new([13]int)
	for _, card := range hand.Cards {
		haveCards[card.Offset]++
	}
	return haveCards
}

func (p *Player) MimicStand() bool {
	if p.Wager.Hand.Soft17() {
		return false
	}
	return p.Wager.Hand.Total() >= 17
}
