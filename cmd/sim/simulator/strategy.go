package simulator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

const (
	BET = "bet"
	INSURANCE = "insurance"
	SURRENDER = "surrender"
	DOUBLE = "double"
	SPLIT = "split"
	STAND = "stand"
	PLAY = "play"
)

var aux struct {
	Bet       int  `json:"bet"`
	Insurance bool `json:"insurance"`
	Double    bool `json:"double"`
	Split     bool `json:"split"`
	Surrender bool `json:"surrender"`
	Stand     bool `json:"stand"`
}

func (p *Player) GetBet() int {
	response, err := http.Get(buildUrl(BET, p.SeenCards, nil, 0, p.Parameters.Playbook, p.NumberOfCards, 0))
	if err != nil || response.StatusCode != 200 {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal([]byte(body), &aux); err != nil {
		panic(err)
	}
	return aux.Bet
}

func (p *Player) GetInsurance() bool {
	p.GetHTTP(buildUrl(INSURANCE, p.SeenCards, nil, 0, p.Parameters.Playbook, p.NumberOfCards, 0))
	return aux.Insurance
}

func (p *Player) GetSurrender(have *[13]int, up int) bool {
	p.GetHTTP(buildUrl(SURRENDER, p.SeenCards, have, 0, p.Parameters.Playbook, p.NumberOfCards, up))
	return aux.Surrender
}

func (p *Player) GetDouble(have *[13]int, up int) bool {
	p.GetHTTP(buildUrl(DOUBLE, p.SeenCards, have, 0, p.Parameters.Playbook, p.NumberOfCards, up))
	return aux.Double
}

func (p *Player) GetSplit(pair, up int) bool {
	p.GetHTTP(buildUrl(SPLIT, p.SeenCards, nil, pair, p.Parameters.Playbook, p.NumberOfCards, up))
	return aux.Split
}

func (p *Player) GetStand(have *[13]int, up int) bool {
	p.GetHTTP(buildUrl(STAND, p.SeenCards, have, 0, p.Parameters.Playbook, p.NumberOfCards, up))
	return aux.Stand
}

//
func (p *Player) GetPlay(have *[13]int, pair, up int) bool {
	p.GetHTTP(buildUrl(PLAY, p.SeenCards, have, pair, p.Parameters.Playbook, p.NumberOfCards, up))
	return aux.Stand
}

//
func (p *Player) GetHTTP(url string) {
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal([]byte(body), &aux); err != nil {
		panic(err)
	}
}

//
func buildUrl(baseUrl string, seenData *[13]int, haveData *[13]int, pair int, playbook string, cards, up int) string {
	params := url.Values{}
	params.Add("playbook", playbook)
	params.Add("cards", fmt.Sprintf("%d", cards))
	params.Add("up", fmt.Sprintf("%d", up))
	params.Add("pair", fmt.Sprintf("%d", pair))

	fullUrl, err := url.Parse(fmt.Sprintf("http://%s/%s", constants.StrategyUrl, baseUrl))
	if err != nil {
		return baseUrl
	}

	if seenData != nil {
		seenPayload, err := json.Marshal(seenData)
		if err == nil {
			params.Add("seen", string(seenPayload))
		}
	}
	if haveData != nil {
		havePayload, err := json.Marshal(haveData)
		if err == nil {
			params.Add("have", string(havePayload))
		}
	}
	fullUrl.RawQuery = params.Encode()

	return fullUrl.String()
}

