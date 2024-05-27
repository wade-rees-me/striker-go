package arguments

import (
	"flag"
	"fmt"

	"github.com/wade-rees-me/striker-go/cmd/striker/constants"
	"github.com/wade-rees-me/striker-go/cmd/striker/utilities"
)

type clSimulationStruct struct {
	Tables        int
	Rounds        int
	BlackjackPays string
	Penatration   float64
}

var CLSimulation = new(clSimulationStruct)

func init() {
	flag.IntVar(&CLSimulation.Tables, "number-of-tables", constants.MinNumberOfTables, fmt.Sprintf("Number of tables (minimum %s; maximum %s).", utilities.Format(constants.MinNumberOfTables), utilities.Format(constants.MaxNumberOfTables)))
	flag.IntVar(&CLSimulation.Rounds, "number-of-rounds", constants.MinNumberOfRounds, fmt.Sprintf("Number of rounds (minimum %s; maximum %s).", utilities.Format(constants.MinNumberOfRounds), utilities.Format(constants.MaxNumberOfRounds)))
	flag.StringVar(&CLSimulation.BlackjackPays, "table-blackjack-pays", "3:2", "Blackjack pays.")
	flag.Float64Var(&CLSimulation.Penatration, "table-penatration", float64(0.75), "Deck penatrations before shuffle.")
}
