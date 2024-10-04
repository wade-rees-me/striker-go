package arguments

import (
	"fmt"
	"os"
	"strconv"

	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

type Arguments struct {
	MimicFlag       bool
	BasicFlag       bool
	LinearFlag      bool
	PolynomialFlag  bool
	HighLowFlag     bool
	WongFlag        bool
	StrikerFlag     bool
	SingleDeckFlag  bool
	DoubleDeckFlag  bool
	SixShoeFlag     bool
	NumberOfHands   int64
}

func NewArguments() *Arguments {
	return &Arguments{
		NumberOfHands: constants.DefaultNumberOfHands,
	}
}

// Parse command-line arguments
func (args *Arguments) ParseArguments() {
	argv := os.Args[1:]  // Skip the first element (program name)
	for i := 0; i < len(argv); i++ {
		switch argv[i] {
		case "-h", "--number-of-hands":
			if i+1 < len(argv) {
				i++
				hands, err := strconv.ParseInt(argv[i], 10, 64)
				if err != nil || hands < constants.MinimumNumberOfHands || hands > constants.MaximumNumberOfHands {
					fmt.Fprintf(os.Stderr, "Number of hands must be between %d and %d\n", constants.MinimumNumberOfHands, constants.MaximumNumberOfHands)
					os.Exit(1)
				}
				args.NumberOfHands = hands
			}
		case "-M", "--mimic":
			args.MimicFlag = true
		case "-B", "--basic":
			args.BasicFlag = true
		case "-L", "--linear":
			args.LinearFlag = true
		case "-P", "--polynomial":
			args.PolynomialFlag = true
		case "-H", "--high-low":
			args.HighLowFlag = true
		case "-W", "--wong":
			args.WongFlag = true
		case "-S", "--striker":
			args.StrikerFlag = true
		case "-1", "--single-deck":
			args.SingleDeckFlag = true
		case "-2", "--double-deck":
			args.DoubleDeckFlag = true
		case "-6", "--six-shoe":
			args.SixShoeFlag = true
		case "--help":
			args.PrintHelpMessage()
			os.Exit(0)
		case "--version":
			args.PrintVersion()
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "Error: Invalid argument: %s\n", argv[i])
			os.Exit(2)
		}
	}
}

func (args *Arguments) PrintVersion() {
	fmt.Printf("%s: version: %s\n", constants.StrikerWhoAmI, constants.StrikerVersion)
}

func (args *Arguments) PrintHelpMessage() {
	fmt.Println(`Usage: strikerGO [options]
Options:
  -h, --help                               Show this help message
  -v, --version                            Display the program version
  -h, --number-of-hands <number of hands>  The number of hands to play in this simulation
  -M, --mimic                              Use the mimic dealer player strategy
  -B, --basic                              Use the basic player strategy
  -L, --linear                             Use the liner regression player strategy
  -P, --polynomial                         Use the polynomial regression player strategy
  -H, --high-low                           Use the high low count player strategy
  -W, --wong                               Use the Wong count player strategy
  -S, --striker                            Use the Striker machine learning player strategy
  -1, --single-deck                        Use a single deck of cards and rules
  -2, --double-deck                        Use a double deck of cards and rules
  -6, --six-shoe                           Use a six deck shoe of cards and rules`)
}

// Get current strategy as a string
func (args *Arguments) GetStrategy() string {
	switch {
	case args.MimicFlag:
		return "mimic"
	case args.PolynomialFlag:
		return "polynomial"
	case args.LinearFlag:
		return "linear"
	case args.HighLowFlag:
		return "high-low"
	case args.WongFlag:
		return "wong"
	case args.StrikerFlag:
		return "striker"
	default:
		return "basic"
	}
}

// Get deck type as a string
func (args *Arguments) GetDecks() string {
	switch {
	case args.DoubleDeckFlag:
		return "double-deck"
	case args.SixShoeFlag:
		return "six-shoe"
	default:
		return "single-deck"
	}
}

// Get the number of decks
func (args *Arguments) GetNumberOfDecks() int {
	switch {
	case args.DoubleDeckFlag:
		return 2
	case args.SixShoeFlag:
		return 6
	default:
		return 1
	}
}

