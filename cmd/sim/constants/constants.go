package constants

import (
	"os"
)

// General constants
const (
	StrikerVersion = "v02.02.00"
	TimeLayout = "2006-01-02 15:04:05 -0700"
)

// Simulation constants
const (
	MaximumNumberOfHands	= int64(25000000000)
	MinimumNumberOfHands	= int64(100)
	DefaultNumberOfHands	= int64(250000000)
	DatabaseNumberOfHands	= int64(250000000)
	MaxSplitHands			= 18
	StrikerWhoAmI			= "striker-go"

	MinimumBet = 2
	MaximumBet = 98
	TrueCountBet = 2
	TrueCounTMultiplier = 26
)

var StrategyUrl = os.Getenv("STRIKER_URL_STRATEGY")
var RulesUrl = os.Getenv("STRIKER_URL_RULES")
var SimulationUrl = os.Getenv("STRIKER_URL_SIMULATION")
