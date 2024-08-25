package simulator

import (
	"time"
)

type SimulationDatabaseTable struct {
	Playbook    string `json:"playbook"`
	Guid        string `json:"guid"`
	Simulator   string `json:"simulator"`
	Summary     string `json:"summary"`
	Simulations string `json:"simulations"`
	Rounds      string `json:"rounds"`
	Hands       string `json:"hands"`
	TotalBet    string `json:"totalbet"`
	TotalWon    string `json:"totalwon"`
	Advantage   string `json:"advantage"`
	TotalTime   string `json:"totaltime"`
	AverageTime string `json:"averagetime"`
	Parameters  string `json:"parameters"`
	Payload		string `json:"payload"`
}

type SimulationParameters struct {
	Guid          string
	Processor     string
	Timestamp     string
	Decks         string // single-deck
	Strategy      string // basic
	Playbook      string // single-deck-basic
	BlackjackPays string
	Tables        int64
	Rounds        int64
	NumberOfDecks int
	Penetration   float64
	OptimumTables int64
	TableRules    *RulesTableStruct
}

type SimulationReport struct {
	TotalRounds int64
	TotalHands  int64
	TotalBet    int64
	TotalWon    int64
	Start       time.Time
	End         time.Time
	Duration    time.Duration
}
