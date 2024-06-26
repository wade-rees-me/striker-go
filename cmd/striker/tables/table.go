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
	TableRules       *database.DBRulesPayload
	SimulationReport SimulationReport
}

type TableReport struct {
	NumberOfRounds     int64
	NumberOfHandsDealt int64
	TableRules         *database.DBRulesPayload
	BlackjackPays      string
	Penetration        float64
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
	ElapsedTime string
	TableReport TableReport
}

func NewTable(r *database.DBRulesPayload, p *Player, tableNumber, numberOfDecks int, penetration float64, name, code string) *Table {
	t := new(Table)
	t.Number = tableNumber
	t.Dealer = NewDealer(r.HitSoft17)
	t.Player = p
	t.TableRules = r
	deck := cards.NewDeck(cards.Suits, cards.Blackjack, 1)
	t.Shoe = *cards.NewShoe(*deck, numberOfDecks, penetration)

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
			logger.Log.Debug(fmt.Sprintf("Play: %s (%s, %s) vs %s", printHand(&t.Player.Wager), (t.Player.Wager.Hand.GetCard(0)).Rank, (t.Player.Wager.Hand.GetCard(1)).Rank, up.Rank))
			if !t.Dealer.Hand.Blackjack() { // Dealer does not have 21
				t.Player.Play(&t.Shoe, up)
				if !t.Player.BustedOrBlackjack() { // If the player busted or has blackjack the dealer does not play
					t.Dealer.Play(&t.Shoe)
				}
			}
			t.Player.Payoff(t.Dealer.Hand.Blackjack(), t.Dealer.Hand.Busted(), t.Dealer.Hand.Total())
			t.Dealer.Statistics()
		}
	}

	t.SimulationReport.End = time.Now()
	t.SimulationReport.Duration = time.Since(t.SimulationReport.Start).Round(time.Second)
	t.SimulationReport.ElapsedTime = t.SimulationReport.Duration.String()
	logger.Log.Info(fmt.Sprintf("  End table: %v, ended at %v, total elapsed time: %s", t.Number, t.SimulationReport.End, t.SimulationReport.Duration))
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
	t.SimulationReport.TableReport.Penetration = arguments.CLSimulation.Penetration
	return &t.SimulationReport
}
