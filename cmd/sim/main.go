package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
	"github.com/wade-rees-me/striker-go/cmd/sim/simulator"
)

func main() {
	flag.Parse()

	if arguments.CLStrategy.StrikerFlag {
		constants.StrategyUrl = constants.StrategyMlbUrl
	}

	parameters := new(simulator.SimulationParameters)
	parameters.Guid = uuid.New().String()
	parameters.Processor = constants.StrikerWhoAmI
	parameters.Timestamp = (time.Now()).Format(constants.TimeLayout)
	parameters.Decks, parameters.NumberOfDecks = arguments.CLTable.Get()
	parameters.Strategy = arguments.CLStrategy.Get()
	parameters.Playbook = fmt.Sprintf("%s-%s", parameters.Decks, parameters.Strategy)
	parameters.Tables = arguments.CLSimulation.Tables
	parameters.Rounds = arguments.CLSimulation.Rounds
	parameters.BlackjackPays = arguments.CLSimulation.BlackjackPays
	parameters.Penetration = arguments.CLSimulation.Penetration
	parameters.TableRules = &simulator.TableRules

	simulator.LoadTableRules(parameters.Decks)
	simulator.RunOnce(parameters)
}
