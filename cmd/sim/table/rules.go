package table

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/wade-rees-me/striker-go/cmd/sim/logger"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

type Rules struct {
	Playbook            string  `json:"playbook"`
	HitSoft17           bool    `json:"hitSoft17"`
	Surrender           bool    `json:"surrender"`
	DoubleAnyTwoCards   bool    `json:"doubleAnyTwoCards"`
	DoubleAfterSplit    bool    `json:"doubleAfterSplit"`
	ResplitAces         bool    `json:"resplitAces"`
	HitSplitAces        bool    `json:"hitSplitAces"`
	BlackjackBets       int     `json:"blackjackBets"`
	BlackjackPays       int     `json:"blackjackPays"`
	Penetration         float64 `json:"penetration"`
}

// Logger struct placeholder (assuming you have a Logger defined elsewhere)
type Logger struct{}

// Mock function for logger.Simulation (for the sake of completeness)
func (l *Logger) Simulation(message string) {
	fmt.Print(message)
}

// Function to load rules table by calling fetchRulesTable
func (r *Rules) LoadTable(decks string) error {
	url := "http://" + constants.RulesUrl + "/" + decks
	if err := r.fetchTable(url); err != nil {
		log.Printf("Error fetching rules table: %v\n", err)
		return err
	}
	return nil
}

// Function to fetch rules table using HTTP GET
func (r *Rules) fetchTable(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON
	var result struct {
		Payload string `json:"payload"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("error parsing payload: %w", err)
	}

	// Parse rules from the payload
	if err := json.Unmarshal([]byte(result.Payload), r); err != nil {
		return fmt.Errorf("error parsing rules: %w", err)
	}

	return nil
}

//
func (r *Rules) Print(logger *logger.Logger) {
	logger.Simulation(fmt.Sprintf("    %-24s\n", "Table Rules"))
	logger.Simulation(fmt.Sprintf("      %-24s: %s\n", "Table", r.Playbook))
	logger.Simulation(fmt.Sprintf("      %-24s: %t\n", "Hit soft 17", r.HitSoft17))
	logger.Simulation(fmt.Sprintf("      %-24s: %t\n", "Surrender", r.Surrender))
	logger.Simulation(fmt.Sprintf("      %-24s: %t\n", "Double any two cards", r.DoubleAnyTwoCards))
	logger.Simulation(fmt.Sprintf("      %-24s: %t\n", "Double after split", r.DoubleAfterSplit))
	logger.Simulation(fmt.Sprintf("      %-24s: %t\n", "Resplit aces", r.ResplitAces))
	logger.Simulation(fmt.Sprintf("      %-24s: %t\n", "Hit split aces", r.HitSplitAces))
	logger.Simulation(fmt.Sprintf("      %-24s: %d\n", "Blackjack bets", r.BlackjackBets))
	logger.Simulation(fmt.Sprintf("      %-24s: %d\n", "Blackjack pays", r.BlackjackPays))
	logger.Simulation(fmt.Sprintf("      %-24s: %0.3f %%\n", "Penetration", r.Penetration))
}
