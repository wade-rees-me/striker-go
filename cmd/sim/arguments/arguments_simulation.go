package arguments

import (
	"flag"
	"fmt"

	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

type clSimulationStruct struct {
	Tables        int64
	Rounds        int64
	BlackjackPays string
	Penetration   float64
}

var CLSimulation = new(clSimulationStruct)

func init() {
	flag.Int64Var(&CLSimulation.Tables, "number-of-tables", constants.MinNumberOfTables, fmt.Sprintf("Number of tables (minimum %d; maximum %d).", constants.MinNumberOfTables, constants.MaxNumberOfTables))
	flag.Int64Var(&CLSimulation.Rounds, "number-of-rounds", constants.MinNumberOfRounds, fmt.Sprintf("Number of rounds (minimum %d; maximum %d).", constants.MinNumberOfRounds, constants.MaxNumberOfRounds))
	flag.StringVar(&CLSimulation.BlackjackPays, "table-blackjack-pays", "3:2", "Set the payout for blackjack pays.")
	flag.Float64Var(&CLSimulation.Penetration, "table-penetration", float64(0.75), "Set the deck penetration before shuffle.")
}
