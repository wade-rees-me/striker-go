package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	screen "github.com/aditya43/clear-shell-screen-golang"
	"github.com/google/uuid"

	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
	"github.com/wade-rees-me/striker-go/cmd/sim/simulator"
)

//
func main() {
	flag.Parse()

	if arguments.CLFlags.HelpFlag || len(flag.Args()) > 0 {
		flag.PrintDefaults()
		return
	}
	if arguments.CLFlags.VersionFlag {
		log.Printf("Version: %s\n", constants.StrikerVersion)
		return
	}

	screen.Clear()
	screen.MoveTopLeft()

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
	parameters.OptimumTables = constants.MinimumBet
	if parameters.Tables == 1 {
		parameters.Tables = parameters.OptimumTables
		parameters.Rounds /= parameters.OptimumTables
	}

	simulator.LoadTableRules(parameters.Decks)
	parameters.TableRules = &simulator.TableRules

	simulator.RunOnce(parameters)
}
