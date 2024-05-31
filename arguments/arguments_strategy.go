package arguments

import (
	"flag"
)

type clStrategyStruct struct {
	BasicFlag      bool
	LinearFlag     bool
	PolynomialFlag bool
}

var CLStrategy = new(clStrategyStruct)

func init() {
	flag.BoolVar(&CLStrategy.BasicFlag, "strategy-basic", false, "Use the basic strategy tables (default).")
	flag.BoolVar(&CLStrategy.LinearFlag, "strategy-linear", false, "Use the linear regression strategy tables.")
	flag.BoolVar(&CLStrategy.PolynomialFlag, "strategy-polynomial", false, "Use the polynomial regression strategy tables.")
}

func (c *clStrategyStruct) Get() string {
	if c.PolynomialFlag { // (This option is currently unimplemented)
		return "polynomial"
	}
	if c.LinearFlag { // (This option is currently unimplemented)
		return "linear"
	}
	return "basic"
}
