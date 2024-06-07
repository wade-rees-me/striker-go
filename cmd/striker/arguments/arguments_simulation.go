package arguments

import (
	"flag"
	"fmt"

	"github.com/wade-rees-me/striker-go/constants"
	"github.com/wade-rees-me/striker-go/utilities"
)

type clSimulationStruct struct {
	Tables        int
	Rounds        int
	BlackjackPays string
	Penetration   float64
}

var CLSimulation = new(clSimulationStruct)

func init() {
	flag.IntVar(&CLSimulation.Tables, "number-of-tables", constants.MinNumberOfTables, fmt.Sprintf("Number of tables (minimum %s; maximum %s).", utilities.Format(constants.MinNumberOfTables), utilities.Format(constants.MaxNumberOfTables)))
	flag.IntVar(&CLSimulation.Rounds, "number-of-rounds", constants.MinNumberOfRounds, fmt.Sprintf("Number of rounds (minimum %s; maximum %s).", utilities.Format(constants.MinNumberOfRounds), utilities.Format(constants.MaxNumberOfRounds)))
	flag.StringVar(&CLSimulation.BlackjackPays, "table-blackjack-pays", "3:2", "Set the payout for blackjack pays.")
	flag.Float64Var(&CLSimulation.Penetration, "table-penetration", float64(0.75), "Set the deck penetration before shuffl.")
}
