package constants

import (
	"os"
)

// General constants
const (
	StrikerVersion = "v01.00.01"
	TimeLayout = "2006-01-02 15:04:05 -0700 MST"
)

// Simulation constants
const (
	MaxNumberOfRounds     = int64(1000000000)
	MinNumberOfRounds     = int64(10000)
	DefaultNumberOfRounds = int64(1000000)
	MaxNumberOfTables     = 4
	MinNumberOfTables     = 1
	MaxSplitHands         = 3
	StrikerWhoAmI         = "striker-go"

	MinimumBet = 2
	MaximumBet = 98
)

var RulesUrl = os.Getenv("STRIKER_URL_RULES")
var SimulationUrl = os.Getenv("STRIKER_URL_SIMULATION")
var StrategyUrl = os.Getenv("STRIKER_URL_STRATEGY")
var StrategyBasicUrl = os.Getenv("STRIKER_URL_BASIC")
var MachineCounterUrl = os.Getenv("STRIKER_URL_MACHINE")

