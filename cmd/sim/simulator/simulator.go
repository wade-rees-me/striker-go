package simulator

import (
	"encoding/json"
	"fmt"
	"log"

	"bytes"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
	//"github.com/wade-rees-me/striker-go/cmd/sim/utilities"
)

type Simulation struct {
	Name       string
	Guid       string
	Simulator  string
	Playbook   string
	Year       int
	Month      int
	Day        int
	Parameters *SimulationParameters
	Report     SimulationReport
	TableList  []Table
}

func NewSimulation(parameters *SimulationParameters) *Simulation {
	s := new(Simulation)
	t := time.Now()
	s.Year = t.Year()
	s.Month = int(t.Month())
	s.Day = t.Day()
	s.Name = fmt.Sprintf("go-striker-%4d_%02d_%02d_%012d", s.Year, s.Month, s.Day, t.Unix())
	s.Guid = uuid.New().String()
	s.Parameters = parameters

	for tableNumber := int64(1); tableNumber <= parameters.Tables; tableNumber++ {
		table := NewTable(tableNumber, parameters)
		player := NewPlayer(parameters, table.Shoe.ShoeReport.NumberOfCards)
		table.AddPlayer(player)
		s.TableList = append(s.TableList, *table)
	}

	return s
}

func RunOnce(parameters *SimulationParameters) {
	log.Printf("Starting striker simulation: ...\n")
	if err := SimulatorProcess(NewSimulation(parameters)); err != nil {
		log.Printf("Simulation failed: %s", err)
	}
}

func SimulatorProcess(s *Simulation) error {
	s.RunSimulation()

	tbs := new(SimulationDatabaseTable)

	jsonData, err := json.Marshal(s.Parameters)
	if err == nil {
		tbs.Parameters = string(jsonData)
	}

	tbs.Guid = s.Parameters.Guid
	tbs.Playbook = s.Parameters.Playbook
	tbs.Simulator = constants.StrikerWhoAmI
	tbs.Summary = "no"
	tbs.Simulations = "1"
	tbs.Rounds = fmt.Sprintf("%d", s.Report.TotalRounds)
	tbs.Hands = fmt.Sprintf("%d", s.Report.TotalHands)
	tbs.TotalBet = fmt.Sprintf("%d", s.Report.TotalBet)
	tbs.TotalWon = fmt.Sprintf("%d", s.Report.TotalWon)
	tbs.TotalTime = fmt.Sprintf("%d", int64(s.Report.Duration.Seconds()))
	tbs.AverageTime = fmt.Sprintf("%06.2f seconds", s.Report.Duration.Seconds()*float64(1000000)/float64(s.Report.TotalHands))
	tbs.Advantage = fmt.Sprintf("%+04.3f%%", (float64(s.Report.TotalWon) / float64(s.Report.TotalBet) * float64(100)))

	fmt.Printf("\n")
	fmt.Printf("Number of rounds:  %s\n", tbs.Rounds)
	fmt.Printf("Number of hands:   %d\n", s.Report.TotalHands)
	fmt.Printf("Total bet:         %d, average bet per hand: %.2f\n", s.Report.TotalBet, (float64(s.Report.TotalBet) / float64(s.Report.TotalHands)))
	fmt.Printf("Total won:         %d, average win per hand: %.2f\n", s.Report.TotalWon, (float64(s.Report.TotalWon) / float64(s.Report.TotalHands)))
	fmt.Printf("Total time:        %s seconds\n", tbs.TotalTime)
	fmt.Printf("Average time:      %s per 1,000,000 hands\n", tbs.AverageTime)
	fmt.Printf("Player advantage:  %s\n", tbs.Advantage) /* House Edge (%)=(Total Loss/Total Bet)×100 */
	fmt.Printf("\n")

	if err := InsertSimulationTable(tbs, s.Playbook); err != nil {
		log.Printf("Failed to insert into Simulation table: %s", err)
		return err
	}

	return nil
}

func (s *Simulation) RunSimulation() {
	var wg sync.WaitGroup

	log.Printf("Simulation %v, started at %v", s.Name, time.Now())
	wg.Add(len(s.TableList))
	for i := range s.TableList {
		t := &s.TableList[i]
		if "mimic" == s.Parameters.Strategy {
			go t.SessionMimic(&wg)
		} else {
			go t.Session(&wg)
		}
	}
	wg.Wait()

	// Merge tables into one report
	for i := range s.TableList {
		t := &s.TableList[i]

		s.Report.TotalRounds += t.Report.TotalRounds
		s.Report.TotalHands += t.Report.TotalHands
		s.Report.TotalBet += t.Player.Report.TotalBet
		s.Report.TotalWon += t.Player.Report.TotalWon
		s.Report.Duration += t.Report.Duration
	}
}

func InsertSimulationTable(s *SimulationDatabaseTable, playbook string) error {
	url := fmt.Sprintf("http://%s/%s/%s/%s", constants.SimulationUrl, s.Simulator, playbook, s.Guid)
	//log.Printf("Insert Simulation: %s\n", url)

	// Convert data to JSON
	jsonData, err := json.Marshal(s)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil
	}
	log.Printf("Insert Simulation: %v\n", string(jsonData))

	// Create a new POST request with JSON data
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Send the request using http.DefaultClient
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Print the response status and body
	log.Printf("Response Status: %v\n", resp.Status)
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return nil
		}
		log.Printf("Response Body: %v\n", string(body))
	}

	return nil
}
