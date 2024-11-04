package simulator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	//"sync"
	"time"

	//"github.com/google/uuid"
	"github.com/dustin/go-humanize"

	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
	"github.com/wade-rees-me/striker-go/cmd/sim/table"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

type Simulator struct {
	Name       string
	Simulator  string
	Playbook   string
	Year       int
	Month      int
	Day        int
	Parameters *arguments.Parameters
	Rules *table.Rules
	Report     arguments.Report
	TableList  []Table
}

func NewSimulator(parameters *arguments.Parameters, rules *table.Rules, strategy *table.Strategy) *Simulator {
	s := new(Simulator)
	t := time.Now()
	s.Year = t.Year()
	s.Month = int(t.Month())
	s.Day = t.Day()
	s.Name = fmt.Sprintf("striker-go--%4d_%02d_%02d_%012d", s.Year, s.Month, s.Day, t.Unix())
	s.Parameters = parameters
	s.Rules = rules

	table := NewTable(1, parameters, rules)
	player := NewPlayer(rules, strategy, table.Shoe.NumberOfCards)
	table.AddPlayer(player)
	s.TableList = append(s.TableList, *table)

	return s
}

//
func (s *Simulator) SimulatorProcess() error {
	fmt.Printf("\n  Start: simulation %s\n", s.Name)
	s.RunSimulation()
	fmt.Printf("  End: simulation\n")

	tbs := new(Simulation)

	jsonData, err := json.Marshal(s.Parameters)
	if err == nil {
		tbs.Parameters = string(jsonData)
	}

	tbs.Guid = s.Parameters.Name
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
	tbs.Advantage = fmt.Sprintf("%+04.3f %%", (float64(s.Report.TotalWon) / float64(s.Report.TotalBet) * float64(100)))
	tbs.Parameters = s.Parameters.Serialize()
	tbs.Rules = s.Rules.Serialize()
	tbs.Payload = "n/a"

	fmt.Printf("\n")
    fmt.Printf("  -- results ---------------------------------------------------------------------\n");
	fmt.Printf("    %-24s: %s\n", "Number of hands", humanize.Comma(s.Report.TotalHands))
	fmt.Printf("    %-24s: %s\n", "Number of rounds",  humanize.Comma(s.Report.TotalRounds))
	fmt.Printf("    %-24s: %s, %+04.3f average bet per hand\n", "Total bet", humanize.Comma(s.Report.TotalBet), (float64(s.Report.TotalBet) / float64(s.Report.TotalHands)))
	fmt.Printf("    %-24s: %s, %+04.3f average win per hand\n", "Total won", humanize.Comma(s.Report.TotalWon), (float64(s.Report.TotalWon) / float64(s.Report.TotalHands)))
	fmt.Printf("    %-24s: %s seconds\n", "Total time", humanize.Comma(int64(s.Report.Duration.Seconds())))
	fmt.Printf("    %-24s: %s per 1,000,000 hands\n", "Average time", tbs.AverageTime)
	fmt.Printf("    %-24s: %s\n", "Player advantage", tbs.Advantage) /* House Edge (%)=(Total Loss/Total Bet)×100 */
    fmt.Printf("  --------------------------------------------------------------------------------\n\n");
	fmt.Printf("\n")

	if(s.Report.TotalHands >= constants.DatabaseNumberOfHands) {
    	fmt.Printf("  -- insert ----------------------------------------------------------------------\n");
		if err := InsertSimulationTable(tbs, tbs.Playbook); err != nil {
			log.Printf("Failed to insert into Simulation table: %s", err)
			return err
		}
    	fmt.Printf("  --------------------------------------------------------------------------------\n");
    }

	return nil
}

func (s *Simulator) RunSimulation() {
	//for i := range s.TableList {
		//t := &s.TableList[i]
		//t.Session("mimic" == s.Parameters.Strategy)
	//}
t := &s.TableList[0]
t.Session("mimic" == s.Parameters.Strategy)

	// Merge tables into one report
	//for i := range s.TableList {
		t = &s.TableList[0]

		s.Report.TotalRounds += t.Report.TotalRounds
		s.Report.TotalHands += t.Report.TotalHands
		s.Report.TotalBet += t.Player.Report.TotalBet
		s.Report.TotalWon += t.Player.Report.TotalWon
		s.Report.Duration += t.Report.Duration
	//}
}

//
func InsertSimulationTable(s *Simulation, playbook string) error {
	url := fmt.Sprintf("http://%s/%s/%s/%s", constants.SimulationUrl, s.Simulator, playbook, s.Guid)
	//log.Printf("Insert Simulation: %s\n", url)

	// Convert data to JSON
	jsonData, err := json.Marshal(s)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil
	}
	//log.Printf("Insert Simulation: %v\n", string(jsonData))

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
