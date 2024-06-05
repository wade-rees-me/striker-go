package tables

import (
	"fmt"
	"strings"

	"github.com/wade-rees-me/striker-go/cards"
	"github.com/wade-rees-me/striker-go/database"
	"github.com/wade-rees-me/striker-go/logger"
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
	TableRules    *database.DBRulesPayload
	PlayStrategy  PlayStrategy
	Wager         Wager
	Splits        [MaxSplitHands]Wager
	splitCount    int
	PlayerReport  PlayerReport
	blackjackPays int
	blackjackBets int
}

type WagerReport struct {
	Advantage       float64
	TotalWon        int64
	TotalBet        int64
	TotalBetWon     int64
	TotalBetLost    int64
	TotalDoubleBet  int64
	TotalDoubleWon  int64
	TotalDoubleLost int64
}

type PlayerReport struct {
	Wager WagerReport
	Hand  HandReport
}

func NewPlayer(tr *database.DBRulesPayload, playStrategy *PlayStrategy, blackjackPays string) *Player {
	p := new(Player)
	p.TableRules = tr
	p.PlayStrategy = *playStrategy
	_, err := fmt.Sscanf(blackjackPays, "%d:%d", &p.blackjackPays, &p.blackjackBets)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse blackjack pays: %v", err))
	}
	return p
}

func NewPlayerReport() *PlayerReport {
	p := new(PlayerReport)
	return p
}

func (p *Player) PlaceBet(b int64) {
	p.Wager.Reset()
	for i := 0; i < len(p.Splits); i++ {
		p.Splits[i].Reset()
	}
	p.splitCount = 0
	p.Wager.Bet(b)
	p.PlayerReport.Hand.HandsDealt++
}

func (p *Player) Play(s *cards.Shoe, up *cards.Card) {
	p.PlayerReport.Hand.HandsPlayed++
	if p.Wager.Hand.Blackjack() {
		return
	}
	if p.double(&p.Wager, up) && (p.TableRules.DoubleAnyTwoCards || (p.Wager.Hand.Total() == 10 || p.Wager.Hand.Total() == 11)) {
		p.Wager.Double()
		p.Wager.Hand.Draw(s.Draw())
		logger.Log.Debug(fmt.Sprintf("  double: %d (%s, %s draw %s) vs %s", p.Wager.Hand.Total(), (p.Wager.Hand.GetCard(0)).Rank, (p.Wager.Hand.GetCard(1)).Rank, (p.Wager.Hand.GetCard(2)).Rank, up.Rank))
		return
	}
	if p.split(&p.Wager, up) {
		split := &p.Splits[p.splitCount]
		p.splitCount++
		if !p.TableRules.ResplitAces && !p.TableRules.HitSplitAces && p.Wager.Hand.PairOf(cards.Ace) {
			logger.Log.Debug(fmt.Sprintf("  split Aces: vs %s", up.Rank))
			p.Wager.SplitWager(split)
			p.Wager.Hand.Draw(s.Draw())
			split.Hand.Draw(s.Draw())
			return
		}
		logger.Log.Debug(fmt.Sprintf("  split pair: %ss vs %s", (p.Wager.Hand.GetCard(0)).Rank, up.Rank))
		p.Wager.SplitWager(split)
		p.Wager.Hand.Draw(s.Draw())
		p.PlaySplit(&p.Wager, s, up)
		split.Hand.Draw(s.Draw())
		p.PlaySplit(split, s, up)
		return
	}
	for !p.Wager.Hand.Busted() && !p.stand(&p.Wager, up) {
		logger.Log.Debug(fmt.Sprintf("  hit: %s vs %s", printHand(&p.Wager), up.Rank))
		p.Wager.Hand.Draw(s.Draw())
	}
	logger.Log.Debug(fmt.Sprintf("  stand: %s vs %s", printHand(&p.Wager), up.Rank))
	if p.Wager.Hand.Busted() {
		p.PlayerReport.Hand.HandsBusted++
	}
}

func (p *Player) PlaySplit(w *Wager, s *cards.Shoe, up *cards.Card) {
	if p.TableRules.DoubleAfterSplit && (p.TableRules.HitSplitAces || !w.Hand.PairOf(cards.Ace)) && p.double(w, up) {
		logger.Log.Debug(fmt.Sprintf("  double after split: %d vs %s", w.Hand.Total(), up.Rank))
		w.Double()
		w.Hand.Draw(s.Draw())
		return
	}
	if (p.TableRules.ResplitAces || !w.Hand.PairOf(cards.Ace)) && p.split(w, up) {
		logger.Log.Debug(fmt.Sprintf("  resplit pair: %s vs %s", (w.Hand.GetCard(0)).Rank, up.Rank))
		split := &p.Splits[p.splitCount]
		p.splitCount++
		w.SplitWager(split)
		w.Hand.Draw(s.Draw())
		p.PlaySplit(w, s, up)
		split.Hand.Draw(s.Draw())
		p.PlaySplit(split, s, up)
		return
	}
	if p.TableRules.HitSplitAces || !w.Hand.PairOf(cards.Ace) {
		for !w.Hand.Busted() && !p.stand(w, up) {
			w.Hand.Draw(s.Draw())
			logger.Log.Debug(fmt.Sprintf("  split Hit: %d vs %s", w.Hand.Total(), up.Rank))
		}
	}
	logger.Log.Debug(fmt.Sprintf("  split Stand: %d vs %s", w.Hand.Total(), up.Rank))
	if w.Hand.Busted() {
		p.PlayerReport.Hand.HandsBusted++
	}
}

func (p *Player) Draw(c *cards.Card) *cards.Card {
	return p.Wager.Hand.Draw(c)
}

func (p *Player) BustedOrBlackjack() bool {
	all := p.Wager.Hand.Busted() || p.Wager.Hand.Blackjack()
	for i := 0; i < p.splitCount; i++ {
		all = all || p.Splits[i].Hand.Busted()
	}
	return all
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
	if blackjack { // Dealer blackjack
		if w.Hand.Blackjack() {
			w.Push()
		} else {
			w.Lost()
		}
	} else if w.Hand.Blackjack() {
		w.WonBlackjack(int64(p.blackjackPays), int64(p.blackjackBets))
	} else if w.Hand.Busted() {
		w.Lost()
		p.PlayerReport.Hand.HandsBusted++
	} else if busted || (w.Hand.Total() > total) {
		w.Won()
	} else if total > w.Hand.Total() {
		w.Lost()
	} else {
		w.Push()
	}
	logger.Log.Debug(fmt.Sprintf("    payoff.wager: %d (bet) -> (%s) vs (%s) = %s ", (w.AmountBet + w.DoubleBet), printHand(w), printDealerHand(blackjack, busted, total), printResults(w.AmountWon, w.DoubleWon)))

	p.PlayerReport.Wager.TotalBet += w.AmountBet + w.DoubleBet
	p.PlayerReport.Wager.TotalWon += w.AmountWon
	p.PlayerReport.Wager.TotalWon += w.DoubleWon

	p.PlayerReport.Wager.TotalDoubleBet += w.DoubleBet
	if w.DoubleWon > 0 {
		p.PlayerReport.Wager.TotalDoubleWon += w.DoubleBet
	}
	if w.DoubleWon < 0 {
		p.PlayerReport.Wager.TotalDoubleLost += w.DoubleWon
	}
	if w.AmountWon > 0 {
		p.PlayerReport.Wager.TotalBetWon += w.AmountWon
	}
	if w.AmountWon < 0 {
		p.PlayerReport.Wager.TotalBetLost += w.AmountWon
	}
}

func printDealerHand(blackjack bool, busted bool, total int) string {
	if blackjack {
		return "blackjack"
	}
	if busted {
		return "busted"
	}
	return fmt.Sprintf("%d", total)
}
func printHand(w *Wager) string {
	if w.Hand.Busted() {
		return fmt.Sprintf("busted %d", w.Hand.Total())
	}
	if w.Hand.Soft() {
		return fmt.Sprintf("soft %d", w.Hand.Total())
	}
	return fmt.Sprintf("hard %d", w.Hand.Total())
}
func printResults(a int64, d int64) string {
	if a > 0 {
		return fmt.Sprintf("won %d", a+d)
	}
	if a < 0 {
		return fmt.Sprintf("lost %d", a+d)
	}
	return fmt.Sprintf("pushed %d", a+d)
}

func (p *Player) double(h *Wager, up *cards.Card) bool {
	if h.Hand.Soft() {
		return p.strategyHelper2(p.PlayStrategy.SoftDouble, h.Hand.Total(), up, false)
	}
	return p.strategyHelper2(p.PlayStrategy.HardDouble, h.Hand.Total(), up, false)
}

func (p *Player) split(h *Wager, up *cards.Card) bool {
	if p.splitCount <= MaxSplitHands && h.Hand.Pair() {
		return p.strategyHelper2(p.PlayStrategy.PairSplit, h.Hand.GetCard(0).Value, up, false)
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
	logger.Log.Debug(fmt.Sprintf("Strategy.table missing row : %d", total))
	return defaultValue
}

func (p *Player) GetReport() *PlayerReport {
	p.PlayerReport.Wager.Advantage = float64(p.PlayerReport.Wager.TotalWon) / float64(p.PlayerReport.Wager.TotalBet) * 100.0
	return &p.PlayerReport
}
