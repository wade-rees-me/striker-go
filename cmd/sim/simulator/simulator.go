package simulator

import (
	"fmt"
	"sync"
	"time"

	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
	"github.com/wade-rees-me/striker-go/cmd/sim/table"
)

type Simulator struct {
	Name       string
	Simulator  string
	Playbook   string
	Year       int
	Month      int
	Day        int
	Parameters *arguments.Parameters
	Rules      *table.Rules
	Report     arguments.Report
	TableList  []Table
}

func NewSimulator(parameters *arguments.Parameters, rules *table.Rules, strategy *table.Strategy) *Simulator {
	s := new(Simulator)
	t := time.Now()
	s.Year = t.Year()
	s.Month = int(t.Month())
	s.Day = t.Day()
	s.Name = fmt.Sprintf("striker-go-%4d_%02d_%02d_%012d", s.Year, s.Month, s.Day, t.Unix())
	s.Parameters = parameters
	s.Rules = rules

	table := NewTable(1, parameters, rules)
	player := NewPlayer(rules, strategy, table.Shoe.NumberOfCards)
	table.AddPlayer(player)
	s.TableList = append(s.TableList, *table)

	return s
}

func (s *Simulator) GetReport() *arguments.Report {
	return &s.Report
}

func (s *Simulator) SimulatorProcess(id int, wg *sync.WaitGroup) error {
	defer wg.Done() // Mark this goroutine as done when it finishes

	//fmt.Printf("\n  Start: simulation %s\n", s.Name)
	s.RunSimulation()
	//fmt.Printf("  End: simulation\n")

	return nil
}

func (s *Simulator) RunSimulation() {
	for i := range s.TableList {
		t := &s.TableList[i]
		t.Session("mimic" == s.Parameters.Strategy)
	}

	// Merge tables into one report
	for i := range s.TableList {
		t := &s.TableList[i]

		s.Report.TotalRounds += t.Report.TotalRounds
		s.Report.TotalHands += t.Report.TotalHands
		s.Report.TotalBlackjacks += t.Player.Report.TotalBlackjacks
		s.Report.TotalDoubles += t.Player.Report.TotalDoubles
		s.Report.TotalSplits += t.Player.Report.TotalSplits
		s.Report.TotalWins += t.Player.Report.TotalWins
		s.Report.TotalPushes += t.Player.Report.TotalPushes
		s.Report.TotalLoses += t.Player.Report.TotalLoses
		s.Report.TotalBet += t.Player.Report.TotalBet
		s.Report.TotalWon += t.Player.Report.TotalWon
		s.Report.Duration += t.Report.Duration
	}
}
