package simulator

import (
	//"fmt"
	"log"
	"sync"
	"time"

	"github.com/wade-rees-me/striker-go/cmd/sim/cards"
)

type Table struct {
	Index	   int64
	Dealer     *cards.Dealer
	Player     *Player
	Shoe       cards.Shoe
	Parameters *SimulationParameters
	Report     SimulationReport
}

func NewTable(index int64, parameters *SimulationParameters) *Table {
	t := new(Table)
	t.Index = index
	t.Parameters = parameters
	t.Dealer = cards.NewDealer(parameters.TableRules.HitSoft17)
	deck := cards.NewDeck(cards.Suits, cards.Blackjack, 1)
	t.Shoe = *cards.NewShoe(*deck, parameters.NumberOfDecks, parameters.Penetration)
	return t
}

func (t *Table) AddPlayer(player *Player) {
	t.Player = player
}

func (t *Table) Session(wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("  Beg table %02d: rounds: %d", t.Index, t.Parameters.Rounds)
	t.Report.Start = time.Now()
	t.Report.TotalRounds = t.Parameters.Rounds
	for i := int64(0); i < t.Parameters.Rounds; i++ {
		t.Shoe.Shuffle()
		t.Player.Shuffle()
		for !t.Shoe.ShouldShuffle() { // Until the cut card is passed
			t.Report.TotalHands++
			t.Dealer.Reset()
			t.Player.PlaceBet()
			up := t.dealCards()
			if up.BlackjackAce() {
				t.Player.Insurance()
			}
			if !t.Dealer.Hand.Blackjack() { // Dealer does not have 21
				t.Player.Play(&t.Shoe, up)
				if !t.Player.BustedOrBlackjack() { // If the player busted or has blackjack the dealer does not play
					t.Dealer.Play(&t.Shoe)
				}
			}
			t.Player.Payoff(t.Dealer.Hand.Blackjack(), t.Dealer.Hand.Busted(), t.Dealer.Hand.Total())
			t.show(up)
		}
	}

	t.Report.End = time.Now()
	t.Report.Duration = time.Since(t.Report.Start).Round(time.Second)
	log.Printf("  End table %02d: ended at %v, total elapsed time: %s", t.Index, t.Report.End, t.Report.Duration)
}

func (t *Table) SessionMimic(wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("  Beg table %02d (mimic): rounds: %d", t.Index, t.Parameters.Rounds)
	t.Report.Start = time.Now()
	t.Report.TotalRounds = t.Parameters.Rounds
	for i := int64(0); i < t.Parameters.Rounds; i++ {
		t.Shoe.Shuffle()
		t.Player.Shuffle()
		for !t.Shoe.ShouldShuffle() { // Until the cut card is passed
			t.Report.TotalHands++
			t.Dealer.Reset()
			t.Player.PlaceMimicBet()
			t.dealCards()

			if !t.Dealer.Hand.Blackjack() {
				t.Player.MimicDealer(&t.Shoe)
				if !t.Player.BustedOrBlackjack() { // If the player busted or has blackjack the dealer does not play
					t.Dealer.Play(&t.Shoe)
				}
			}
			t.Player.Payoff(t.Dealer.Hand.Blackjack(), t.Dealer.Hand.Busted(), t.Dealer.Hand.Total())
		}
	}

	t.Report.End = time.Now()
	t.Report.Duration = time.Since(t.Report.Start).Round(time.Second)
	log.Printf("  End table %02d (mimic): ended at %v, total elapsed time: %s", t.Index, t.Report.End, t.Report.Duration)
}

func (t *Table) dealCards() *cards.Card {
	t.Player.Draw(t.Shoe.Draw())
	up := t.Dealer.Draw(t.Shoe.Draw())
	t.Player.Show(up)
	t.Player.Draw(t.Shoe.Draw())
	t.Dealer.Draw(t.Shoe.Draw())
	return up
}

func (t *Table) show(up *cards.Card) {
	for _, card := range t.Dealer.Hand.Cards {
		if up.Index != card.Index {
			t.Player.Show(&card)
		}
	}
}
