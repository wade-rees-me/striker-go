package arguments

import (
	"flag"
)

type clStrategyStruct struct {
	MimicFlag      bool
	BasicFlag      bool
	LinearFlag     bool
	PolynomialFlag bool
	HighLowFlag    bool
	WongFlag       bool
	StrikerFlag    bool
}

var CLStrategy = new(clStrategyStruct)

func init() {
	flag.BoolVar(&CLStrategy.MimicFlag, "strategy-mimic", false, "Use the mimic the dealer strategy tables.")
	flag.BoolVar(&CLStrategy.BasicFlag, "strategy-basic", false, "Use the basic strategy tables (default).")
	flag.BoolVar(&CLStrategy.LinearFlag, "strategy-linear", false, "Use the linear regression strategy tables.")
	flag.BoolVar(&CLStrategy.PolynomialFlag, "strategy-polynomial", false, "Use the polynomial regression strategy tables.")
	flag.BoolVar(&CLStrategy.HighLowFlag, "strategy-high-low", false, "Use the High Low strategy tables.")
	flag.BoolVar(&CLStrategy.WongFlag, "strategy-wong", false, "Use the Wong strategy tables.")
	flag.BoolVar(&CLStrategy.StrikerFlag, "strategy-striker", false, "Use the Striker strategy tables.")
}

func (c *clStrategyStruct) Get() string {
	if c.MimicFlag {
		return "mimic"
	}
	if c.PolynomialFlag {
		return "polynomial"
	}
	if c.LinearFlag {
		return "linear"
	}
	if c.HighLowFlag {
		return "high-low"
	}
	if c.WongFlag {
		return "wong"
	}
	if c.StrikerFlag {
		return "striker"
	}
	return "basic"
}
