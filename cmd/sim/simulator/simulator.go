package simulator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	//"github.com/google/uuid"
	"github.com/dustin/go-humanize"

	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

type Simulation struct {
	Name       string
	Simulator  string
	Playbook   string
	Year       int
	Month      int
	Day        int
	Parameters *arguments.Parameters
	Report     arguments.Report
	TableList  []Table
}

type SimulationDatabaseTable struct {
	Playbook    string `json:"playbook"`
	Guid        string `json:"guid"`
	Simulator   string `json:"simulator"`
	Summary     string `json:"summary"`
	Simulations string `json:"simulations"`
	Rounds      string `json:"rounds"`
	Hands       string `json:"hands"`
	TotalBet    string `json:"bet"`
	TotalWon    string `json:"won"`
	Advantage   string `json:"advantage"`
	TotalTime   string `json:"time"`
	AverageTime string `json:"average"`
	Parameters  string `json:"parameters"`
	Payload     string `json:"payload"`
}

func NewSimulation(parameters *arguments.Parameters) *Simulation {
	s := new(Simulation)
	t := time.Now()
	s.Year = t.Year()
	s.Month = int(t.Month())
	s.Day = t.Day()
	s.Name = fmt.Sprintf("striker-go--%4d_%02d_%02d_%012d", s.Year, s.Month, s.Day, t.Unix())
	s.Parameters = parameters

		table := NewTable(1, parameters)
		player := NewPlayer(parameters, table.Shoe.NumberOfCards)
		table.AddPlayer(player)
		s.TableList = append(s.TableList, *table)

	return s
}

//
func (s *Simulation) SimulatorProcess() error {
	s.Parameters.Logger.Simulation(fmt.Sprintf("\n  Start: simulation %s\n", s.Name))
	s.RunSimulation()
	s.Parameters.Logger.Simulation(fmt.Sprintf("  End: simulation\n"))

	tbs := new(SimulationDatabaseTable)

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

	fmt.Printf("\n")
    s.Parameters.Logger.Simulation("  -- results ---------------------------------------------------------------------\n");
	s.Parameters.Logger.Simulation(fmt.Sprintf("    %-24s: %s\n", "Number of hands", humanize.Comma(s.Report.TotalHands)))
	s.Parameters.Logger.Simulation(fmt.Sprintf("    %-24s: %s\n", "Number of rounds",  humanize.Comma(s.Report.TotalRounds)))
	s.Parameters.Logger.Simulation(fmt.Sprintf("    %-24s: %s, %+04.3f average bet per hand\n", "Total bet", humanize.Comma(s.Report.TotalBet), (float64(s.Report.TotalBet) / float64(s.Report.TotalHands))))
	s.Parameters.Logger.Simulation(fmt.Sprintf("    %-24s: %s, %+04.3f average win per hand\n", "Total won", humanize.Comma(s.Report.TotalWon), (float64(s.Report.TotalWon) / float64(s.Report.TotalHands))))
	s.Parameters.Logger.Simulation(fmt.Sprintf("    %-24s: %s seconds\n", "Total time", humanize.Comma(int64(s.Report.Duration.Seconds()))))
	s.Parameters.Logger.Simulation(fmt.Sprintf("    %-24s: %s per 1,000,000 hands\n", "Average time", tbs.AverageTime))
	s.Parameters.Logger.Simulation(fmt.Sprintf("    %-24s: %s\n", "Player advantage", tbs.Advantage)) /* House Edge (%)=(Total Loss/Total Bet)×100 */
    s.Parameters.Logger.Simulation("  --------------------------------------------------------------------------------\n\n");
	fmt.Printf("\n")

    s.Parameters.Logger.Simulation("  -- insert ----------------------------------------------------------------------\n");
	if err := InsertSimulationTable(tbs, tbs.Playbook); err != nil {
		log.Printf("Failed to insert into Simulation table: %s", err)
		return err
	}
    s.Parameters.Logger.Simulation("  --------------------------------------------------------------------------------\n");

	return nil
}

func (s *Simulation) RunSimulation() {
	var wg sync.WaitGroup
	status := make(chan string) // Channel to pass status updates

	wg.Add(len(s.TableList))
	for i := range s.TableList {
		t := &s.TableList[i]
		//s.Parameters.Logger.Simulation(fmt.Sprintf("    Start: %s table session\n", s.Parameters.Strategy));
		go t.Session(&wg, "mimic" == s.Parameters.Strategy, status)
		//s.Parameters.Logger.Simulation(fmt.Sprintf("    End: table session\n"));
	}
	//wg.Wait()

	// Goroutine to close the status channel after all processes finish
	go func() {
		wg.Wait() // Wait for all goroutines to complete
		close(status) // Close the status channel when done
	}()

	// Continuously receive and print status updates from the channel
	for update := range status {
		fmt.Print(update)
	}
/*
	fmt.Println("All processes completed.")
*/

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



/*
func process(id int, wg *sync.WaitGroup, status chan string) {
	defer wg.Done() // Notify the WaitGroup when done

	for i := 1; i <= 10; i++ {
		time.Sleep(time.Millisecond * 500) // Simulate work
		// Send status update to the channel
		status <- fmt.Sprintf("Process %d is at step %d\n", id, i)
	}
}

func main() {
	var wg sync.WaitGroup   // WaitGroup to wait for all goroutines to complete
	status := make(chan string) // Channel to pass status updates

	// Start two concurrent processes
	wg.Add(2) // We have two goroutines to wait for
	go process(1, &wg, status) // Start process 1
	go process(2, &wg, status) // Start process 2

	// Goroutine to close the status channel after all processes finish
	go func() {
		wg.Wait() // Wait for all goroutines to complete
		close(status) // Close the status channel when done
	}()

	// Continuously receive and print status updates from the channel
	for update := range status {
		fmt.Print(update)
	}

	fmt.Println("All processes completed.")
}
*/


func InsertSimulationTable(s *SimulationDatabaseTable, playbook string) error {
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
