package arguments

import (
	"flag"
)

type clTableStruct struct {
	DeckSingleFlag bool
	DeckDoubleFlag bool
	DeckMultiFlag  bool
}

var CLTable = new(clTableStruct)

func init() {
	flag.BoolVar(&CLTable.DeckSingleFlag, "table-deck-single", false, "Single deck table (default).")
	flag.BoolVar(&CLTable.DeckDoubleFlag, "table-deck-double", false, "Double deck table.")
	flag.BoolVar(&CLTable.DeckMultiFlag, "table-deck-multi", false, "Multi deck table (6 deck shoe).")
}

func (c *clTableStruct) Get() string {
	if c.DeckDoubleFlag {
		return "double"
	}
	if c.DeckMultiFlag {
		return "multi"
	}
	return "single"
}
