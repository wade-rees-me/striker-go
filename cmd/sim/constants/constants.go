package constants

import (
	"fmt"
	"os"
	"strings"
)

// General constants
const (
	StrikerVersion = "v3.00.00"
	TimeLayout     = "2006-01-02 15:04:05 -0700"
)

// Simulation constants
const (
	Million               = int64(1000000)
	Billion               = int64(Million * 1000)
	MaximumNumberOfHands  = int64(Billion * 10)
	MinimumNumberOfHands  = int64(1000)
	DefaultNumberOfHands  = int64(Million * 500)
	DatabaseNumberOfHands = int64(Million * 500)
	MaxSplitHands         = 18
	StrikerWhoAmI         = "striker-go"
	StatusRounds          = int64(1000000)
	MyHostname            = "Striker"

	NumberOfCardsInDeck   = 52
	NumberOfCoresPhysical = 24
	NumberOfCoresLogical  = 32
	NumberOfCoresDefault  = 16

	MinimumBet          = 2
	MaximumBet          = 20
	TrueCountBet        = 2
	TrueCounTMultiplier = 26
)

// Get hostname and check if it matches
func IsMyComputer() bool {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname:", err)
		return false
	}

	myHostname := MyHostname
	return myHostname == hostname
}

var RulesUrl = os.Getenv("STRIKER_URL_RULES")
var ChartsUrl = os.Getenv("STRIKER_URL_CHARTS")
var SimulationsUrl = os.Getenv("STRIKER_URL_SIMULATIONS")

// removes escape sequences like \n, \" and \\
func UnescapeJSON(str string) string {
	var result strings.Builder
	i := 0

	for i < len(str) {
		if str[i] == '\\' {
			i++ // Skip the backslash

			// Handle known escape sequences
			if i < len(str) {
				switch str[i] {
				case 'n':
					result.WriteRune('\n') // Convert \n to newline
				case '"':
					result.WriteRune('"') // Convert \" to "
				case '\\':
					result.WriteRune('\\') // Convert \\ to \
				default:
					result.WriteByte(str[i]) // Copy other characters
				}
			}
		} else {
			result.WriteByte(str[i]) // Copy normal characters
		}
		i++
	}

	return result.String()
}

// removes the leading and trailing quotes from the string if they exist
func StripQuotes(str string) string {
	if len(str) > 1 && str[0] == '"' && str[len(str)-1] == '"' {
		return str[1 : len(str)-1] // Remove first and last quote
	}
	return str
}
