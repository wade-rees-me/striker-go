package cards

var Blackjack = map[string][]int{
	Two:   {2, 0},
	Three: {3, 1},
	Four:  {4, 2},
	Five:  {5, 3},
	Six:   {6, 4},
	Seven: {7, 5},
	Eight: {8, 6},
	Nine:  {9, 7},
	Ten:   {10, 8},
	Jack:  {10, 9},
	Queen: {10, 10},
	King:  {10, 11},
	Ace:   {11, 12},
}

func (c *Card) BlackjackAce() bool {
	return c.Value == 11
}

