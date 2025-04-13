package arguments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
	"github.com/wade-rees-me/striker-go/cmd/sim/table"
)

type Report struct {
	Name            string  `json:"guid"`
	Version         string  `json:"version"`
	Playbook        string  `json:"playbook"`
	Simulator       string  `json:"simulator"`
	Strategy        string  `json:"strategy"`
	Decks           string  `json:"decks"`
	Epoch           string  `json:"epoch"`
	TotalRounds     int64   `json:"rounds"`
	TotalHands      int64   `json:"hands"`
	TotalBet        int64   `json:"total_bet"`
	TotalWon        int64   `json:"total_won"`
	TotalBlackjacks int64   `json:"total_blackjacks"`
	TotalDoubles    int64   `json:"total_doubles"`
	TotalSplits     int64   `json:"total_splits"`
	TotalWins       int64   `json:"total_wins"`
	TotalLoses      int64   `json:"total_loses"`
	TotalPushes     int64   `json:"total_pushes"`
	TotalThreads    int64   `json:"threads"`
	Start           int64   `json:"start"`
	End             int64   `json:"end"`
	Duration        int64   `json:"duration"`
	Advantage       float64 `json:"advantage"`
	PerBillion      float64 `json:"per_billion"`
}

func (r *Report) InitFinal(parameters *Parameters, start time.Time) {
	r.Init()
	r.Name = parameters.Name
	r.Version = constants.StrikerVersion
	r.Playbook = parameters.Playbook
	r.Simulator = parameters.Processor
	r.Strategy = parameters.Strategy
	r.Decks = parameters.Decks
	r.TotalThreads = parameters.NumberOfThreads
	r.Epoch = parameters.Epoch
	r.Start = start.Unix()
	r.End = int64(0)
	r.Duration = int64(0)
	r.Advantage = float64(0.0)
	r.PerBillion = float64(0.0)
}

func (r *Report) Init() {
	r.TotalRounds = int64(0)
	r.TotalHands = int64(0)
	r.TotalBet = int64(0)
	r.TotalWon = int64(0)
	r.TotalBlackjacks = int64(0)
	r.TotalDoubles = int64(0)
	r.TotalSplits = int64(0)
	r.TotalWins = int64(0)
	r.TotalLoses = int64(0)
	r.TotalPushes = int64(0)
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

func (r *Report) Finish(end time.Time) {
	r.End = end.Unix()

	r.Duration = r.End - r.Start
	if r.TotalBet > 0 {
		r.Advantage = float64(r.TotalWon) / float64(r.TotalBet) * 100.0
	}
	if r.TotalHands > 0 {
		r.PerBillion = (float64(r.Duration) / float64(r.TotalHands)) * float64(constants.Billion)
	}
}

func (report *Report) Print(numberOfThreads int64) {
	fmt.Printf("    %-26s: %17s\n", "Number of hands", humanize.Comma(report.TotalHands))
	fmt.Printf("    %-26s: %17s\n", "Number of rounds", humanize.Comma(report.TotalRounds))
	fmt.Printf("    %-26s: %17s %+08.3f average bet per hand\n", "Total bet", humanize.Comma(report.TotalBet),
		(float64(report.TotalBet) / float64(report.TotalHands)))
	fmt.Printf("    %-26s: %17s %+08.3f average win per hand\n", "Total won", humanize.Comma(report.TotalWon),
		(float64(report.TotalWon) / float64(report.TotalHands)))
	fmt.Printf("    %-26s: %17s %+08.3f %% of total hands\n", "Total blackjacks", humanize.Comma(report.TotalBlackjacks),
		(float64(report.TotalBlackjacks) / float64(report.TotalHands) * 100.0))
	fmt.Printf("    %-26s: %17s %+08.3f %% of total hands\n", "Total doubles", humanize.Comma(report.TotalDoubles),
		(float64(report.TotalDoubles) / float64(report.TotalHands) * 100.0))
	fmt.Printf("    %-26s: %17s %+08.3f %% of total hands\n", "Total split", humanize.Comma(report.TotalSplits),
		(float64(report.TotalSplits) / float64(report.TotalHands) * 100.0))
	fmt.Printf("    %-26s: %17s %+08.3f %% of total hands\n", "Total wins", humanize.Comma(report.TotalWins),
		(float64(report.TotalWins) / float64(report.TotalHands) * 100.0))
	fmt.Printf("    %-26s: %17s %+08.3f %% of total hands\n", "Total pushes", humanize.Comma(report.TotalPushes),
		(float64(report.TotalPushes) / float64(report.TotalHands) * 100.0))
	fmt.Printf("    %-26s: %17s %+08.3f %% of total hands\n", "Total loses", humanize.Comma(report.TotalLoses),
		(float64(report.TotalLoses) / float64(report.TotalHands) * 100.0))
	fmt.Printf("    %-26s: %17s seconds\n", "Total time", humanize.Comma(int64(report.Duration)))
	fmt.Printf("    %-26s: %17s threads\n", "Number of threads", humanize.Comma(int64(report.TotalThreads)))
	fmt.Printf("    %-26s: %17s %s\n", "Average time",
		humanize.Comma(int64(float64(report.Duration)*float64(1000000000)/float64(report.TotalHands))),
		"seconds per 1,000,000,000 hands")
	// House Edge (%)=(Total Loss/Total Bet)×100
	fmt.Printf("    %-26s: %17s %+08.3f %%\n", "Player advantage", "", report.Advantage)
}

func (r *Report) Insert(p *Parameters, l *table.Rules) {
	if r.TotalHands < constants.DatabaseNumberOfHands {
		fmt.Printf("    Error: Not enough hands played (%d). Minimum required is %d\n", r.TotalHands,
			constants.DatabaseNumberOfHands)
		return
	}
	url := fmt.Sprintf("http://%s/%s/%s/%s", constants.SimulationsUrl, p.Processor, p.Playbook, p.Name)

	// Convert data to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		fmt.Println("    Error marshalling JSON: %v\n", err)
		return
	}

	// Create a new POST request with JSON data
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("    Error creating request: %v\n", err)
		return
	}

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Send the request using http.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("    Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Print the response status and body
	if resp.StatusCode == 200 {
		fmt.Printf("    Insert successful\n")
		return
	}
	fmt.Printf("    Failed to insert into Simulation table: %s\n", err)
}
