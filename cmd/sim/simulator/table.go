package simulator

import (
	"fmt"
	"time"
	"os"

	"github.com/dustin/go-humanize"

	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
	"github.com/wade-rees-me/striker-go/cmd/sim/table"
	"github.com/wade-rees-me/striker-go/cmd/sim/cards"
)

const STATUS_DOT = 25000;
const STATUS_LINE = 1000000;

type Table struct {
	Index      int64
	Dealer     *cards.Dealer
	Player     *Player
	Shoe       cards.Shoe
	Parameters *arguments.Parameters
	Report     arguments.Report
	Up		   *cards.Card
	Down	   *cards.Card
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
	fmt.Printf("    Start: %s table session\n", t.Parameters.Strategy);
	fmt.Printf("      Start: table playing %s hands\n", humanize.Comma(t.Parameters.NumberOfHands))
	t.Report.Start = time.Now()
	for t.Report.TotalHands < t.Parameters.NumberOfHands {
        t.Status(t.Report.TotalRounds, t.Report.TotalHands)
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

	t.Report.End = time.Now()
	t.Report.Duration = time.Since(t.Report.Start).Round(time.Second)
	fmt.Printf("\n      End: table\n")
	fmt.Printf("    End: table session\n");
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
	if round == 0 {
		fmt.Printf("        ")
	}
	if (round+1)%STATUS_DOT == 0 {
		fmt.Printf(".")
	}
	if (round+1)%STATUS_LINE == 0 {
		fmt.Printf(" : %s (rounds), %s (hands)\n", humanize.Comma(round+1), humanize.Comma(hand))
		fmt.Printf("        ")
	}
	os.Stdout.Sync()
}
