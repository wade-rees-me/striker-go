package simulator

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
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
	response, err := http.Get(buildUrl("bet", p.SeenCards, nil, 0, p.Parameters.Playbook, p.NumberOfCards, 0))
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
	return ClampInt(aux.Bet, constants.MinimumBet, constants.MaximumBet)
}

func (p *Player) GetInsurance() bool {
	getUrl := buildUrl("insurance", p.SeenCards, nil, 0, p.Parameters.Playbook, p.NumberOfCards, 0)
	response, err := http.Get(getUrl)
	if err != nil {
		return false
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}

	if err = json.Unmarshal([]byte(body), &aux); err != nil {
		return false
	}
	//if aux.Insurance {fmt.Printf("Insurance: %v\n", aux.Insurance)}
	return aux.Insurance
}

func (p *Player) GetSurrender(have *[13]int, up int) bool {
	getUrl := buildUrl("surrender", p.SeenCards, have, 0, p.Parameters.Playbook, p.NumberOfCards, up)
	response, err := http.Get(getUrl)
	if err != nil {
		return false
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}

	if err = json.Unmarshal([]byte(body), &aux); err != nil {
		return false
	}
	//if aux.Surrender {fmt.Printf("Surrender: %v\n", aux.Surrender)}
	return aux.Surrender
}

func (p *Player) GetDouble(have *[13]int, up int) bool {
	getUrl := buildUrl("double", p.SeenCards, have, 0, p.Parameters.Playbook, p.NumberOfCards, up)
	response, err := http.Get(getUrl)
	if err != nil {
		return false
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}

	if err = json.Unmarshal([]byte(body), &aux); err != nil {
		return false
	}
	return aux.Double
}

func (p *Player) GetSplit(pair, up int) bool {
	getUrl := buildUrl("split", p.SeenCards, nil, pair, p.Parameters.Playbook, p.NumberOfCards, up)
	response, err := http.Get(getUrl)
	if err != nil {
		return false
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}

	if err = json.Unmarshal([]byte(body), &aux); err != nil {
		return false
	}
	//if aux.Split {fmt.Printf("Split: %v\n", aux.Split)}
	return aux.Split
}

func (p *Player) GetStand(have *[13]int, up int) bool {
	getUrl := buildUrl("stand", p.SeenCards, have, 0, p.Parameters.Playbook, p.NumberOfCards, up)
	response, err := http.Get(getUrl)
	if err != nil {
		return false
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}

	if err = json.Unmarshal([]byte(body), &aux); err != nil {
		return false
	}
	return aux.Stand
}

func buildUrl(baseUrl string, seenData *[13]int, haveData *[13]int, pair int, playbook string, cards, up int) string {
	params := url.Values{}
	params.Add("playbook", playbook)
	params.Add("cards", fmt.Sprintf("%d", cards))
	params.Add("up", fmt.Sprintf("%d", up))
	params.Add("pair", fmt.Sprintf("%d", pair))

	fullUrl, err := url.Parse(fmt.Sprintf("http://%s/%s", constants.StrategyBasicUrl, baseUrl))
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

func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
