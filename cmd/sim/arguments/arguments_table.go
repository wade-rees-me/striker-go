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
	flag.BoolVar(&CLTable.DeckSingleFlag, "table-single-deck", false, "Use a single deck table (default).")
	flag.BoolVar(&CLTable.DeckDoubleFlag, "table-double-deck", false, "Use a double deck table.")
	flag.BoolVar(&CLTable.DeckMultiFlag, "table-six-shoe", false, "Use a six shoe table (6 deck shoe).")
}

func (c *clTableStruct) Get() (string, int) {
	if c.DeckDoubleFlag {
		return "double-deck", 2
	}
	if c.DeckMultiFlag {
		return "six-shoe", 6
	}
	return "single-deck", 1
}
