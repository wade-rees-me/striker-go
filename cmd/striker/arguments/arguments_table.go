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
	flag.BoolVar(&CLTable.DeckSingleFlag, "table-deck-single", false, "Use a single deck table (default).")
	flag.BoolVar(&CLTable.DeckDoubleFlag, "table-deck-double", false, "Use a double deck table. (This option is currently unimplemented)")
	flag.BoolVar(&CLTable.DeckMultiFlag, "table-deck-multi", false, "Use a multi deck table (6 deck shoe). (This option is currently unimplemented)")
}

func (c *clTableStruct) Get() string {
	if c.DeckDoubleFlag { // (This option is currently unimplemented)
		return "single" //"double"
	}
	if c.DeckMultiFlag { // (This option is currently unimplemented)
		return "single" //"multi"
	}
	return "single"
}
