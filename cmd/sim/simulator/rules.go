package simulator

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type RulesDatabaseTable struct {
	Playbook string `json:"playbook"`
	Payload   string `json:"payload"`
}

type RulesTableStruct struct {
	Playbook          string  `json:"playbook"`
	HitSoft17         bool    `json:"hitSoft17"`
	Surrender         bool    `json:"surrender"`
	DoubleAnyTwoCards bool    `json:"doubleAnyTwoCards"`
	DoubleAfterSplit  bool    `json:"doubleAfterSplit"`
	ResplitAces       bool    `json:"resplitAces"`
	HitSplitAces      bool    `json:"hitSplitAces"`
	BlackjackPays     string  `json:"blackjackPays"`
	Penetration       float64 `json:"penetration"`
}

var RulesUrl = os.Getenv("STRIKER_URL_RULES")
var TableRules RulesTableStruct

//
func LoadTableRules(decks string) {
	if err := FetchRulesTable("http://" + RulesUrl + "/" + decks); err != nil {
		panic(err.Error())
	}
}

//
func FetchRulesTable(url string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(body), &TableRules); err != nil {
		return err
	}
	return nil
}

