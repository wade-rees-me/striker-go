package main

import (
	"fmt"
	"log"
	"time"

	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
	"github.com/wade-rees-me/striker-go/cmd/sim/table"
	"github.com/wade-rees-me/striker-go/cmd/sim/simulator"
	"github.com/wade-rees-me/striker-go/cmd/sim/logger"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

//
func main() {
	args := arguments.NewArguments()
	args.ParseArguments()

	name := generateName()
    rules := &table.Rules{}
    if err := rules.LoadTable(args.GetDecks()); err != nil {
        log.Fatalf("Failed to load rules: %v", err)
    }

	logger := logger.NewLogger(constants.StrikerWhoAmI,  args.NumberOfHands < 1000)
    params := arguments.NewParameters(name, args.GetDecks(), args.GetStrategy(), args.GetNumberOfDecks(), args.NumberOfHands, rules, logger)

    logger.Simulation(fmt.Sprintf("Start: %s ...\n\n", name))
    logger.Simulation("  -- arguments -------------------------------------------------------------------\n");
    params.Print();
    rules.Print(logger);
    logger.Simulation("  --------------------------------------------------------------------------------\n");

	simulator := simulator.NewSimulation(params);
	simulator.SimulatorProcess()
    logger.Simulation(fmt.Sprintf("End: %s\n\n", name))
}

//
func generateName() string {
	t := time.Now()

	year := t.Year()
	month := int(t.Month())
	day := t.Day()
	unixTime := t.Unix()

	name := fmt.Sprintf("%s_%04d_%02d_%02d_%012d", constants.StrikerWhoAmI, year, month, day, unixTime)
	return name
}

