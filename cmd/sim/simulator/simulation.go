package simulator

type Simulation struct {
	Playbook    string `json:"playbook"`
	Guid        string `json:"guid"`
	Simulator   string `json:"simulator"`
	Summary     string `json:"summary"`
	Simulations string `json:"simulations"`
	Rounds      string `json:"rounds"`
	Hands       string `json:"hands"`
	TotalBet    string `json:"total_bet"`
	TotalWon    string `json:"total_won"`
	Advantage   string `json:"advantage"`
	TotalTime   string `json:"total_time"`
	AverageTime string `json:"average_time"`
	Parameters  string `json:"parameters"`
	Rules		string `json:"rules"`
	Payload     string `json:"payload"`
}

