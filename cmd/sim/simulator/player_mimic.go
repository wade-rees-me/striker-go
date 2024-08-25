package simulator

import (
	"github.com/wade-rees-me/striker-go/cmd/sim/cards"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

func (p *Player) PlaceMimicBet() {
	p.Wager.Reset()
	p.Wager.Bet(int64(constants.MinimumBet))
}

func (p *Player) MimicDealer(s *cards.Shoe) {
	for !p.MimicStand() {
		p.Draw(s.Draw())
	}
}

func (p *Player) MimicStand() bool {
	if p.Wager.Hand.Soft17() {
		return false
	}
	return p.Wager.Hand.Total() >= 17
}
