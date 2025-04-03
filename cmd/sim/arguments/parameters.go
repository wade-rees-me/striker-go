package arguments

import (
	"fmt"
	"log"
	"time"
	"encoding/json"

	"github.com/dustin/go-humanize"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

type Parameters struct {
	Playbook	  	string
	Name		  	string
	Processor	  	string
	Timestamp	  	string
	Decks		  	string
	Strategy	  	string
	NumberOfDecks 	int
	NumberOfThreads	int64
	NumberOfHands 	int64
	NumberOfShares	int64
	Verbose			bool
}

// NewParameters is the constructor for Parameters struct
func NewParameters(decks, strategy string, numDecks int, numberOfHands, numberOfThreads int64) *Parameters {
	params := &Parameters{
		Decks:		  	 decks,
		Strategy:	  	 strategy,
		NumberOfDecks:	 numDecks,
		NumberOfHands:	 numberOfHands,
		NumberOfThreads: numberOfThreads,
		NumberOfShares:  numberOfHands / numberOfThreads,
		Playbook:	  	 fmt.Sprintf("%s-%s", decks, strategy),
		Name:		  	 generateName(),
		Processor:	  	 constants.StrikerWhoAmI,
	}

	params.getCurrentTime()
	return params
}

// Print the parameters
func (p *Parameters) Print() {
	fmt.Printf("    %-26s: %s\n", "Processor", p.Processor)
	fmt.Printf("    %-26s: %d\n", "Threads", p.NumberOfThreads)
	fmt.Printf("    %-26s: %s\n", "Name", p.Name)
	fmt.Printf("    %-26s: %s\n", "Version", constants.StrikerVersion)
	fmt.Printf("    %-26s: %s\n", "Playbook", p.Playbook)
	fmt.Printf("    %-26s: %s\n", "Decks", p.Decks)
	fmt.Printf("    %-26s: %s\n", "Strategy", p.Strategy)
	fmt.Printf("    %-26s: %17s\n", "Number of hands", humanize.Comma(p.NumberOfHands))
	fmt.Printf("    %-26s: %17s\n", "Thread's share of hands", humanize.Comma(p.NumberOfShares))
	fmt.Printf("    %-26s: %s\n", "Timestamp", p.Timestamp)
}

// Get the current timestamp in the desired format
func (p *Parameters) getCurrentTime() {
	p.Timestamp = time.Now().Format(constants.TimeLayout)
}

// Serialize parameters to JSON
func (p *Parameters) Serialize() string {
	data := map[string]interface{}{
		"playbook":			p.Playbook,
		"name":				p.Name,
		"processor":		p.Processor,
		"timestamp":		p.Timestamp,
		"decks":			p.Decks,
		"strategy":			p.Strategy,
		"rounds":			p.NumberOfHands,
		"number_of_decks":	p.NumberOfDecks,
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Failed to serialize parameters: %v", err)
	}

	return string(jsonBytes)
}

//
func generateName() string {
	t := time.Now()

	year := t.Year()
	month := int(t.Month())
	day := t.Day()
	unixTime := t.Unix()

	name := fmt.Sprintf("%s_%04d_%02d_%02d_%012d", constants.StrikerWhoAmI, year, month, day, unixTime)
	return name
}

