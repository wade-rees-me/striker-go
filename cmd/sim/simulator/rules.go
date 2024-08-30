package simulator

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

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

var TableRules RulesTableStruct

func LoadTableRules(decks string) {
	if err := FetchRulesTable("http://" + constants.RulesUrl + "/" + decks); err != nil {
		panic(err.Error())
	}
}

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
