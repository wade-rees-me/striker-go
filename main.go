package main

import (
	"flag"
	"fmt"

	screen "github.com/aditya43/clear-shell-screen-golang"
	"github.com/wade-rees-me/striker-go/cmd/striker/arguments"
	"github.com/wade-rees-me/striker-go/cmd/striker/constants"
	"github.com/wade-rees-me/striker-go/cmd/striker/simulators"
	"github.com/wade-rees-me/striker-go/cmd/striker/utilities"
)

func main() {
	flag.Parse()

	if arguments.CLFlags.HelpFlag || len(flag.Args()) > 0 {
		flag.PrintDefaults()
		return
	}
	if arguments.CLFlags.VersionFlag {
		fmt.Println("Version: ", constants.StrikerVersion)
		return
	}

	if !arguments.CLFlags.QueueFlag {
		screen.Clear()
		screen.MoveTopLeft()
		utilities.Banner()
		simulators.SimulatorRunOnce()
		return
	}
	utilities.Banner()
	simulators.SimulatorRunQueue()
}
