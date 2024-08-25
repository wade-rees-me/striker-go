package arguments

import (
	"flag"
)

type clFlagsStruct struct {
	HelpFlag    bool
	VersionFlag bool
}

var CLFlags = new(clFlagsStruct)

func init() {
	flag.BoolVar(&CLFlags.HelpFlag, "help", false, "Print a help message and exit.")
	flag.BoolVar(&CLFlags.VersionFlag, "help-version", false, "Display version information and exit.")
}
