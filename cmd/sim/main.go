package main

import (
	"fmt"

	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
	"github.com/wade-rees-me/striker-go/cmd/sim/table"
	"github.com/wade-rees-me/striker-go/cmd/sim/simulator"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

//
func main() {
    fmt.Printf("Start: %s ...\n\n", constants.StrikerWhoAmI)
	args := arguments.NewArguments()
    params := arguments.NewParameters(args.GetDecks(), args.GetStrategy(), args.GetNumberOfDecks(), args.NumberOfHands)
    rules := table.NewRules(args.GetDecks())
	strategy := table.NewStrategy(args.GetDecks(), args.GetStrategy(), args.GetNumberOfDecks() * 52)
	sim := simulator.NewSimulator(params, rules, strategy)

    fmt.Printf("  -- arguments -------------------------------------------------------------------\n");
    params.Print();
    rules.Print();
    fmt.Printf("  --------------------------------------------------------------------------------\n");

	sim.SimulatorProcess()
    fmt.Printf("End: %s\n\n", constants.StrikerWhoAmI)
}

