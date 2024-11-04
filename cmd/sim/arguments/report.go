package arguments

import (
	"time"
)

type Report struct {
	TotalRounds int64
	TotalHands  int64
	TotalBet    int64
	TotalWon    int64
	Start       time.Time
	End         time.Time
	Duration    time.Duration
}

