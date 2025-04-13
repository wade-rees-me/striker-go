package arguments

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

type Parameters struct {
	Playbook        string
	Name            string
	Processor       string
	Epoch           string
	Decks           string
	Strategy        string
	NumberOfDecks   int
	NumberOfThreads int64
	NumberOfHands   int64
	NumberOfShares  int64
	Verbose         bool
}

// constructor for Parameters
func NewParameters(arguments *Arguments) *Parameters {
	params := &Parameters{
		Decks:           arguments.GetDecks(),
		Strategy:        arguments.GetStrategy(),
		NumberOfDecks:   arguments.GetNumberOfDecks(),
		NumberOfHands:   arguments.NumberOfHands,
		NumberOfThreads: arguments.NumberOfThreads,
		NumberOfShares:  (arguments.NumberOfHands / arguments.NumberOfThreads) + 1,
		Playbook:        fmt.Sprintf("%s-%s", arguments.GetDecks(), arguments.GetStrategy()),
		Name:            generateName(),
		Processor:       constants.StrikerWhoAmI,
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
	fmt.Printf("    %-26s: %s\n", "Epoch", p.Epoch)
}

// Get the current epoch in the desired format
func (p *Parameters) getCurrentTime() {
	p.Epoch = time.Now().Format(constants.TimeLayout)
}

func generateName() string {
	t := time.Now()

	year := t.Year()
	month := int(t.Month())
	day := t.Day()
	unixTime := t.Unix()

	name := fmt.Sprintf("%s_%04d_%02d_%02d_%012d", constants.StrikerWhoAmI, year, month, day, unixTime)
	return name
}
