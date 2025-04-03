package main

import (
	"fmt"
	"log"
	"sync"
	"time"
	"runtime"

	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
	"github.com/wade-rees-me/striker-go/cmd/sim/table"
	"github.com/wade-rees-me/striker-go/cmd/sim/simulator"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

//
func main() {
	var wg sync.WaitGroup
	args := arguments.NewArguments()
	params := arguments.NewParameters(args.GetDecks(), args.GetStrategy(), args.GetNumberOfDecks(), args.NumberOfHands, args.NumberOfThreads)
	rules := table.NewRules(args.GetDecks())
	strategy := table.NewStrategy(args.GetDecks(), args.GetStrategy(), args.GetNumberOfDecks() * 52)
	simulators := make([]*simulator.Simulator, constants.NumberOfCoresLogical)
	finalReport := new(arguments.Report)

    runtime.GOMAXPROCS(runtime.NumCPU()) // Utilize all available CPUs

	fmt.Printf("Start: %s ...\n", constants.StrikerWhoAmI)
	fmt.Printf("  -- arguments -------------------------------------------------------------------\n");
	params.Print();
	rules.Print();
	fmt.Printf("  --------------------------------------------------------------------------------\n");
    fmt.Printf("  Start: simulation(%s) on %d cores\n", params.Name, args.NumberOfThreads);

	finalReport.Start = time.Now()
	for i := 1; i <= int(args.NumberOfThreads); i++ {
		wg.Add(1)
		sim := simulator.NewSimulator(params, rules, strategy)
		simulators[i - 1] = sim
		go sim.SimulatorProcess(i, &wg)
	}

	wg.Wait() // Wait for all workers to finish
    fmt.Printf("  End: simulation(%s) on %d cores\n", params.Name, args.NumberOfThreads);
	fmt.Printf("End: %s\n", constants.StrikerWhoAmI)

	finalReport.End = time.Now()
	finalReport.Duration = time.Since(finalReport.Start).Round(time.Second)
	for i := 1; i <= int(args.NumberOfThreads); i++ {
		sim := simulators[i - 1]
		finalReport.Merge(sim.GetReport())
	}
	finalReport.Print(args.NumberOfThreads)

	if(finalReport.TotalHands >= constants.DatabaseNumberOfHands) {
		fmt.Printf("  -- insert ----------------------------------------------------------------------\n");
		if err := finalReport.Insert(params, rules); err != nil {
			log.Printf("Failed to insert into Simulation table: %s", err)
		}
		fmt.Printf("  --------------------------------------------------------------------------------\n");
	}
}

