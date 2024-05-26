package arguments

import (
	"flag"
)

type clFlagsStruct struct {
	QueueFlag   bool
	HelpFlag    bool
	VersionFlag bool
}

var CLFlags = new(clFlagsStruct)

func init() {
	flag.BoolVar(&CLFlags.QueueFlag, "queue", false, "Start in Queue mode reading simulation requests from the queue.")
	flag.BoolVar(&CLFlags.HelpFlag, "help", false, "Print a help message and exit.")
	flag.BoolVar(&CLFlags.VersionFlag, "help-version", false, "Display version information and exit.")
}
