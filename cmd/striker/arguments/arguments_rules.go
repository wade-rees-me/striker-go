package arguments

import (
	"flag"
)

type clRulesStruct struct {
	VegasFlag bool
	RenoFlag  bool
}

var CLRules = new(clRulesStruct)

func init() {
	flag.BoolVar(&CLRules.VegasFlag, "rules-vegas", false, "Use Vegas rules (default).")
	flag.BoolVar(&CLRules.RenoFlag, "rules-reno", false, "Use Reno rules.")
}

func (c *clRulesStruct) Get() string {
	if c.RenoFlag {
		return "reno"
	}
	return "vegas"
}
