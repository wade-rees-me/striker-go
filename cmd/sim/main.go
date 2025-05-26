package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/wade-rees-me/striker-go/cmd/sim/arguments"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
	"github.com/wade-rees-me/striker-go/cmd/sim/simulator"
	"github.com/wade-rees-me/striker-go/cmd/sim/table"
	"github.com/wade-rees-me/striker-go/cmd/sim/xlog"
)

func main() {
	err := xlog.InitSyslog(xlog.SyslogAddress, xlog.SyslogPort)
	if err != nil {
		log.Fatalf("init_syslog failed: %v", err)
	}
	defer xlog.CloseSyslog()
	var wg sync.WaitGroup
	args := arguments.NewArguments()
	params := arguments.NewParameters(args)
	rules := table.NewRules(args.GetDecks())
	strategy := table.NewStrategy(args.GetDecks(), args.GetStrategy(), args.GetNumberOfDecks())
	simulators := make([]*simulator.Simulator, constants.NumberOfCoresLogical)
	finalReport := new(arguments.Report)

	runtime.GOMAXPROCS(runtime.NumCPU()) // Utilize all available CPUs

	xlog.LogInfo("Simulation started at %v", time.Now())
	xlog.LogError("Encountered an error: %s", "bad luck")
	xlog.LogFatal("Simulation aborted: %s", "fatal error")

	fmt.Printf("Start: %s ...\n", constants.StrikerWhoAmI)
	fmt.Printf("  -- arguments -------------------------------------------------------------------\n")
	params.Print()
	rules.Print()
	fmt.Printf("  --------------------------------------------------------------------------------\n")

	finalReport.InitFinal(params, time.Now())
	for i := 1; i <= int(args.NumberOfThreads); i++ {
		wg.Add(1)
		sim := simulator.NewSimulator(params, rules, strategy)
		simulators[i-1] = sim
		go sim.SimulatorProcess(i, &wg)
	}

	wg.Wait() // Wait for all workers to finish
	for i := 1; i <= int(args.NumberOfThreads); i++ {
		sim := simulators[i-1]
		finalReport.Merge(sim.GetReport())
	}

	finalReport.Finish(time.Now())
	fmt.Printf("  -- results ---------------------------------------------------------------------\n")
	finalReport.Print(args.NumberOfThreads)
	fmt.Printf("  --------------------------------------------------------------------------------\n")
	fmt.Printf("  -- insert ----------------------------------------------------------------------\n")
	finalReport.Insert(params, rules)
	fmt.Printf("  --------------------------------------------------------------------------------\n")
}
