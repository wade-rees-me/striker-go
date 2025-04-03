package arguments

import (
	"fmt"
	"time"
	"bytes"
	"encoding/json"
	//"sync"
	"io"
	"log"
	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/wade-rees-me/striker-go/cmd/sim/table"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

type Report struct {
	TotalRounds 	int64
	TotalHands  	int64
	TotalBet		int64
	TotalWon		int64
 	TotalBlackjacks	int64
	TotalDoubles	int64
	TotalSplits		int64
	TotalWins		int64
	TotalLoses		int64
	TotalPushes		int64
	TotalThreads	int64
	Advantage		float64
	Start			time.Time
	End				time.Time
	Duration		time.Duration
}

func (r *Report) Merge(b *Report) {
	r.TotalRounds += b.TotalRounds
	r.TotalHands += b.TotalHands
	r.TotalBet += b.TotalBet
	r.TotalWon += b.TotalWon
	r.TotalBlackjacks += b.TotalBlackjacks
	r.TotalDoubles += b.TotalDoubles
	r.TotalSplits += b.TotalSplits
	r.TotalWins += b.TotalWins
	r.TotalLoses += b.TotalLoses
	r.TotalPushes += b.TotalPushes
}

func (report *Report) Print(numberOfThreads int64) {
	report.TotalThreads = numberOfThreads
    report.Advantage = float64(report.TotalWon) / float64(report.TotalBet) * float64(100)

	fmt.Printf("\n")
	fmt.Printf("  -- results ---------------------------------------------------------------------\n")
	fmt.Printf("    %-26s: %17s\n", "Number of hands", humanize.Comma(report.TotalHands))
	fmt.Printf("    %-26s: %17s\n", "Number of rounds",  humanize.Comma(report.TotalRounds))
	fmt.Printf("    %-26s: %17s %+08.3f average bet per hand\n", "Total bet", humanize.Comma(report.TotalBet), (float64(report.TotalBet) / float64(report.TotalHands)))
	fmt.Printf("    %-26s: %17s %+08.3f average win per hand\n", "Total won", humanize.Comma(report.TotalWon), (float64(report.TotalWon) / float64(report.TotalHands)))
	fmt.Printf("    %-26s: %17s %+08.3f %% of total hands\n", "Total blackjacks", humanize.Comma(report.TotalBlackjacks), (float64(report.TotalBlackjacks) / float64(report.TotalHands) * 100.0))
	fmt.Printf("    %-26s: %17s %+08.3f %% of total hands\n", "Total doubles", humanize.Comma(report.TotalDoubles), (float64(report.TotalDoubles) / float64(report.TotalHands) * 100.0))
	fmt.Printf("    %-26s: %17s %+08.3f %% of total hands\n", "Total split", humanize.Comma(report.TotalSplits), (float64(report.TotalSplits) / float64(report.TotalHands) * 100.0))
	fmt.Printf("    %-26s: %17s %+08.3f %% of total hands\n", "Total wins", humanize.Comma(report.TotalWins), (float64(report.TotalWins) / float64(report.TotalHands) * 100.0))
	fmt.Printf("    %-26s: %17s %+08.3f %% of total hands\n", "Total pushes", humanize.Comma(report.TotalPushes), (float64(report.TotalPushes) / float64(report.TotalHands) * 100.0))
	fmt.Printf("    %-26s: %17s %+08.3f %% of total hands\n", "Total loses", humanize.Comma(report.TotalLoses), (float64(report.TotalLoses) / float64(report.TotalHands) * 100.0))
	fmt.Printf("    %-26s: %17s seconds\n", "Total time", humanize.Comma(int64(report.Duration.Seconds())))
	fmt.Printf("    %-26s: %17s threads\n", "Number of threads", humanize.Comma(int64(report.TotalThreads)))
	fmt.Printf("    %-26s: %17s %s\n", "Average time", humanize.Comma(int64(float64(report.Duration.Seconds()) * float64(1000000000) / float64(report.TotalHands))), "seconds per 1,000,000,000 hands")
	fmt.Printf("    %-26s: %17s %+08.3f %%\n", "Player advantage", "", report.Advantage) // House Edge (%)=(Total Loss/Total Bet)×100
	fmt.Printf("  --------------------------------------------------------------------------------\n\n")
	fmt.Printf("\n")
}

//
func (r *Report) Insert(p *Parameters, l *table.Rules) error {
	url := fmt.Sprintf("http://%s/%s/%s/%s", constants.SimulationsUrl, p.Processor, p.Playbook, p.Name)

	// Convert data to JSON
	jsonData, err := json.Marshal(r)
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
