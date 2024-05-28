package tables

import (
	"fmt"
	"sync"
	"time"

	"github.com/wade-rees-me/striker-go/arguments"
	"github.com/wade-rees-me/striker-go/cards"
	"github.com/wade-rees-me/striker-go/database"
	"github.com/wade-rees-me/striker-go/logger"
)

type Table struct {
	Number           int
	Dealer           *Dealer
	Player           *Player
	Shoe             cards.Shoe
	TableRules       *database.TableRules
	SimulationReport SimulationReport
}

type TableReport struct {
	NumberOfRounds     int64
	NumberOfHandsDealt int64
	TableRules         *database.TableRules
	BlackjackPays      string
	Penatration        float64
	ShoeReport         *cards.ShoeReport
	DealerReport       *DealerReport
	PlayerReport       *PlayerReport
}

type SimulationReport struct {
	Name        string
	Code        string
	Start       time.Time
	End         time.Time
	Duration    time.Duration
	TableReport TableReport
}

func NewTable(r *database.TableRules, p *Player, tableNumber, numberOfDecks int, penatration float64, name, code string) *Table {
	t := new(Table)
	t.Number = tableNumber
	t.Dealer = NewDealer(r.HitSoft17)
	t.Player = p
	t.TableRules = r
	deck := cards.NewDeck(cards.Suits, cards.Blackjack, 1)
	t.Shoe = *cards.NewShoe(*deck, numberOfDecks, penatration)

	t.SimulationReport.Name = fmt.Sprintf("%s_table_%02d", name, tableNumber)
	t.SimulationReport.Code = code
	t.SimulationReport.TableReport.TableRules = r

	return t
}

func (t *Table) Session(wg *sync.WaitGroup, numberOfRounds int) {
	defer wg.Done()

	logger.Log.Info(fmt.Sprintf("  Beg table: %v, rounds: %v\n", t.Number, numberOfRounds))
	t.SimulationReport.Start = time.Now()
	t.SimulationReport.TableReport.NumberOfRounds = int64(numberOfRounds)
	for i := 0; i < numberOfRounds; i++ {
		t.Shoe.Shuffle()
		for !t.Shoe.ShouldShuffle() { // Until the cut card is passed
			t.SimulationReport.TableReport.NumberOfHandsDealt++
			t.Dealer.Reset()
			t.Player.PlaceBet(2)
			up := t.dealCards()
			if !t.Dealer.Hand.Blackjack() { // Dealer does not have 21
				t.Player.Play(&t.Shoe, up)
			}
			if !t.Player.BustedOrBlackjack() { // If the player busted or has blackjack the dealer does not play
				t.Dealer.Play(&t.Shoe)
			}
			t.Player.Payoff(t.Dealer.Hand.Blackjack(), t.Dealer.Hand.Busted(), t.Dealer.Hand.Total())
			t.Dealer.Statistics()
		}
	}

	t.SimulationReport.End = time.Now()
	t.SimulationReport.Duration = time.Since(t.SimulationReport.Start).Round(time.Second)
	logger.Log.Info(fmt.Sprintf("  End table: %v, ended at %v, total elapsed time: %v", t.Number, t.SimulationReport.End, t.SimulationReport.Duration))
}

func (t *Table) dealCards() *cards.Card {
	t.Player.Draw(t.Shoe.Draw())
	up := t.Dealer.Draw(t.Shoe.Draw())
	t.Player.Draw(t.Shoe.Draw())
	t.Dealer.Draw(t.Shoe.Draw())
	return up
}

func (t *Table) GetReport() *SimulationReport {
	t.SimulationReport.TableReport.ShoeReport = t.Shoe.GetReport()
	t.SimulationReport.TableReport.DealerReport = t.Dealer.GetReport()
	t.SimulationReport.TableReport.PlayerReport = t.Player.GetReport()
	t.SimulationReport.TableReport.BlackjackPays = arguments.CLSimulation.BlackjackPays
	t.SimulationReport.TableReport.Penatration = arguments.CLSimulation.Penatration
	return &t.SimulationReport
}
