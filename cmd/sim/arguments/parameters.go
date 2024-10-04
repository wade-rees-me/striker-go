package arguments

import (
	"fmt"
	"log"
	"time"
	"encoding/json"

	//"github.com/google/uuid"
	"github.com/dustin/go-humanize"

	"github.com/wade-rees-me/striker-go/cmd/sim/logger"
	"github.com/wade-rees-me/striker-go/cmd/sim/table"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

type Parameters struct {
	Rules         *table.Rules
	Logger        *logger.Logger
	Playbook      string
	Name          string
	Processor     string
	Timestamp     string
	Decks         string
	Strategy      string
	NumberOfDecks int
	NumberOfHands int64
}

// NewParameters is the constructor for Parameters struct
func NewParameters(name, decks, strategy string, numDecks int, numberOfHands int64, rules *table.Rules, logger *logger.Logger) *Parameters {
	params := &Parameters{
		Rules:         rules,
		Logger:        logger,
		Decks:         decks,
		Strategy:      strategy,
		NumberOfDecks: numDecks,
		NumberOfHands: numberOfHands,
		Playbook:      fmt.Sprintf("%s-%s", decks, strategy),
		Name:	       name,
		Processor:     constants.StrikerWhoAmI,
	}

	params.getCurrentTime()
	return params
}

// Print the parameters
func (p *Parameters) Print() {
	p.Logger.Simulation(fmt.Sprintf("    %-24s: %s\n", "Name", p.Name))
	p.Logger.Simulation(fmt.Sprintf("    %-24s: %s\n", "Playbook", p.Playbook))
	p.Logger.Simulation(fmt.Sprintf("    %-24s: %s\n", "Processor", p.Processor))
	p.Logger.Simulation(fmt.Sprintf("    %-24s: %s\n", "Version", constants.StrikerVersion))
	p.Logger.Simulation(fmt.Sprintf("    %-24s: %s\n", "Number of hands", humanize.Comma(p.NumberOfHands)))
	p.Logger.Simulation(fmt.Sprintf("    %-24s: %s\n", "Timestamp", p.Timestamp))
}

// Get the current timestamp in the desired format
func (p *Parameters) getCurrentTime() {
	p.Timestamp = time.Now().Format(constants.TimeLayout)
}

// Serialize parameters to JSON
func (p *Parameters) Serialize() string {
	data := map[string]interface{}{
		"playbook":          p.Playbook,
		"name":              p.Name,
		"processor":         p.Processor,
		"timestamp":         p.Timestamp,
		"decks":             p.Decks,
		"strategy":          p.Strategy,
		"rounds":            p.NumberOfHands,
		"number_of_decks":   p.NumberOfDecks,
		"hit_soft_17":       p.Rules.HitSoft17,
		"surrender":         p.Rules.Surrender,
		"double_any_two_cards": p.Rules.DoubleAnyTwoCards,
		"double_after_split":   p.Rules.DoubleAfterSplit,
		"resplit_aces":      p.Rules.ResplitAces,
		"hit_split_aces":    p.Rules.HitSplitAces,
		"blackjack_bets":    p.Rules.BlackjackBets,
		"blackjack_pays":    p.Rules.BlackjackPays,
		"penetration":       p.Rules.Penetration,
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Failed to serialize parameters: %v", err)
	}

	return string(jsonBytes)
}

