package simulator

import (
	"fmt"
	"os"

	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
	"github.com/wade-rees-me/striker-go/cmd/sim/cards"
	"github.com/wade-rees-me/striker-go/cmd/sim/table"
)

type Table struct {
	Index      int64
	Dealer     *cards.Dealer
	Player     *Player
	Shoe       cards.Shoe
	Parameters *arguments.Parameters
	Report     arguments.Report
	Up         *cards.Card
	Down       *cards.Card
}

func NewTable(index int64, parameters *arguments.Parameters, rules *table.Rules) *Table {
	t := new(Table)
	t.Index = index
	t.Parameters = parameters
	t.Dealer = cards.NewDealer(rules.HitSoft17)
	t.Shoe = *cards.NewShoe(parameters.NumberOfDecks, rules.Penetration)
	return t
}

func (t *Table) AddPlayer(player *Player) {
	t.Player = player
}

func (t *Table) Session(mimic bool) {
	for t.Report.TotalHands < t.Parameters.NumberOfShares {
		if t.Parameters.NumberOfThreads == 1 {
			t.Status(t.Report.TotalRounds, t.Report.TotalHands)
		}
		t.Report.TotalRounds++
		t.Shoe.Shuffle()
		t.Player.Shuffle()

		for !t.Shoe.ShouldShuffle() { // Until the cut card is passed
			t.Report.TotalHands++
			t.Dealer.Reset()
			t.Player.PlaceBet(mimic)
			t.dealCards()

			if !mimic && t.Up.BlackjackAce() {
				t.Player.Insurance()
			}

			if !t.Dealer.Hand.Blackjack() { // Dealer does not have 21
				t.Player.Play(&t.Shoe, t.Up, mimic)
				if !t.Player.BustedOrBlackjack() { // If the player busted or has blackjack the dealer does not play
					for !t.Dealer.Stand() {
						card := t.Shoe.Draw()
						t.Dealer.Draw(card)
						t.Player.Show(card)
					}
				}
			}

			t.Player.Show(t.Down)
			t.Player.Payoff(t.Dealer.Hand.Blackjack(), t.Dealer.Hand.Busted(), t.Dealer.Hand.Total())
		}
	}

	if t.Parameters.NumberOfThreads == 1 {
		fmt.Printf("\r")
	}
}

func (t *Table) dealCards() {
	t.Player.Draw(&t.Player.Wager.Hand, &t.Shoe)
	t.Up = t.Shoe.Draw()
	t.Dealer.Draw(t.Up)
	t.Player.Show(t.Up)

	t.Player.Draw(&t.Player.Wager.Hand, &t.Shoe)
	t.Down = t.Shoe.Draw()
	t.Dealer.Draw(t.Down)
}

func (t *Table) Status(round int64, hand int64) {
	spinner := []rune{'|', '/', '-', '\\'}

	fmt.Printf("\r%c Simulating...", spinner[round%int64(len(spinner))])
	os.Stdout.Sync()
}
