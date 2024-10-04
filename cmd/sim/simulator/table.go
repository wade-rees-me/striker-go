package simulator

import (
	"fmt"
	"sync"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/wade-rees-me/striker-go/cmd/sim/logger"
	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
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
}

func NewTable(index int64, parameters *arguments.Parameters) *Table {
	t := new(Table)
	t.Index = index
	t.Parameters = parameters
	t.Dealer = cards.NewDealer(parameters.Rules.HitSoft17)
	t.Shoe = *cards.NewShoe(cards.DeckOfPokerCards, parameters.NumberOfDecks, parameters.Rules.Penetration)
	return t
}

func (t *Table) AddPlayer(player *Player) {
	t.Player = player
}

func (t *Table) Session(wg *sync.WaitGroup, mimic bool, status chan string) {
	defer wg.Done()

	t.Parameters.Logger.Simulation(fmt.Sprintf("    Start: %s table session\n", t.Parameters.Strategy));
	t.Parameters.Logger.Simulation(fmt.Sprintf("      Start: table playing %s hands\n", humanize.Comma(t.Parameters.NumberOfHands)))
	t.Report.Start = time.Now()
	for t.Report.TotalHands < t.Parameters.NumberOfHands {
        t.Status(t.Report.TotalRounds, t.Report.TotalHands, t.Parameters.Logger, status)
		t.Report.TotalRounds++
		t.Shoe.Shuffle()
		t.Player.Shuffle()

		for !t.Shoe.ShouldShuffle() { // Until the cut card is passed
			t.Report.TotalHands++
			t.Dealer.Reset()
			t.Player.PlaceBet(mimic)
			up := t.dealCards()

			if !mimic && up.BlackjackAce() {
				t.Player.Insurance()
			}

			if !t.Dealer.Hand.Blackjack() { // Dealer does not have 21
				t.Player.Play(&t.Shoe, up, mimic)
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
	t.Parameters.Logger.Simulation(fmt.Sprintf("\n      End: table\n"))
	t.Parameters.Logger.Simulation(fmt.Sprintf("    End: table session\n"));
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

func (t *Table) Status(round int64, hand int64, logger *logger.Logger, status chan string) {
	if round == 0 {
		//logger.Simulation("        ")
		status <- fmt.Sprintf("        ")
	}
	if (round+1)%STATUS_DOT == 0 {
		//logger.Simulation(".")
		status <- fmt.Sprintf(".")
	}
	if (round+1)%STATUS_LINE == 0 {
		// Format the round and hand count with commas
		buffer := fmt.Sprintf(" : %s (rounds), %s (hands)\n", humanize.Comma(round+1), humanize.Comma(hand))
		//logger.Simulation(buffer)
		status <- fmt.Sprintf(buffer)
		//logger.Simulation("        ")
		status <- fmt.Sprintf("        ")
	}
}
