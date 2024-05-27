package main

import (
	"flag"
	"fmt"

	screen "github.com/aditya43/clear-shell-screen-golang"
	"github.com/wade-rees-me/striker-go/arguments"
	"github.com/wade-rees-me/striker-go/constants"
	"github.com/wade-rees-me/striker-go/simulators"
	"github.com/wade-rees-me/striker-go/utilities"
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
