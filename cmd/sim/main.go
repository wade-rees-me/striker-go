package main

import (
	"fmt"
	//"log"
	//"time"

	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
	"github.com/wade-rees-me/striker-go/cmd/sim/table"
	"github.com/wade-rees-me/striker-go/cmd/sim/simulator"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

//
func main() {
    fmt.Printf("Start: %s ...\n\n", constants.StrikerWhoAmI)
	argumentsX := arguments.NewArguments()
    parametersX := arguments.NewParameters(argumentsX.GetDecks(), argumentsX.GetStrategy(), argumentsX.GetNumberOfDecks(), argumentsX.NumberOfHands)
    rulesX := table.NewRules(argumentsX.GetDecks())
	strategyX := table.NewStrategy(argumentsX.GetDecks(), argumentsX.GetStrategy(), argumentsX.GetNumberOfDecks() * 52)
	simulatorX := simulator.NewSimulator(parametersX, rulesX, strategyX)

    fmt.Printf("  -- arguments -------------------------------------------------------------------\n");
    parametersX.Print();
    rulesX.Print();
    fmt.Printf("  --------------------------------------------------------------------------------\n");

	simulatorX.SimulatorProcess()
    fmt.Printf("End: %s\n\n", constants.StrikerWhoAmI)
}

