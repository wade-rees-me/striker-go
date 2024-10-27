package table

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	"strconv"
	"strings"

	"github.com/wade-rees-me/striker-go/cmd/sim/cards"
	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

/*
type Strategy struct {
	Playbook     string
	Counts       []int
	Bets         []int
	Insurance    string
	SoftDouble   map[string][]string
	HardDouble   map[string][]string
	PairSplit    map[string][]string
	SoftStand    map[string][]string
	HardStand    map[string][]string
	NumberOfCards int
	JsonResponse []map[string]interface{}
}
*/
type Strategy struct {
    Playbook       string                       `json:"playbook"`
    Counts         []int                        `json:"counts"`
    Bets           []int                        `json:"bets"`
    Insurance      string                       `json:"insurance"`
    SoftDouble     map[string][]string          `json:"soft-double"`
    HardDouble     map[string][]string          `json:"hard-double"`
    PairSplit      map[string][]string          `json:"pair-split"`
    SoftStand      map[string][]string          `json:"soft-stand"`
    HardStand      map[string][]string          `json:"hard-stand"`
    SoftSurrender  map[string][]string          `json:"soft-surrender"`
    HardSurrender  map[string][]string          `json:"hard-surrender"`
	NumberOfCards int
	JsonResponse []map[string]interface{}
}

func NewStrategy(decks, strategy string, numberOfCards int) *Strategy {
	s := &Strategy{NumberOfCards: numberOfCards}
	if strategy != "mimic" {
		err := s.fetchJson("http://localhost:57910/striker/v1/strategy")
		if err != nil {
			log.Fatalf("Error fetching JSON: %v", err)
		}
		fmt.Printf("decks: %s\n", decks)
		fmt.Printf("strategy: %s\n", strategy)
		err = s.fetchTable(decks, strategy)
		if err != nil {
			log.Fatalf("Error fetching table: %v", err)
		}
	}
	return s
}

func (s *Strategy) fetchJson(url string) error {
	fmt.Printf("url: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	err = json.Unmarshal(body, &s.JsonResponse)
	fmt.Printf("err: %v\n", err)
	return err
}

func (s *Strategy) fetchTable(decks, strategy string) error {
	for _, item := range s.JsonResponse {
		if item["playbook"] == decks && item["hand"] == strategy {
			payload, err := json.Marshal(item["payload"])
			if err != nil {
				return err
			}

fmt.Printf("%s\n", payload)

var strategy Strategy
//err = json.Unmarshal([]byte(payload), s)
//err = json.Unmarshal([]byte(payload), &strategy)
//err = json.Unmarshal(payload, &strategy)
if err := json.Unmarshal([]byte(payload), &strategy); err != nil {
	log.Fatal(err)
}
fmt.Printf("%v\n", strategy)


/*
			//var jsonPayload map[string]interface{}
			//err = json.Unmarshal(payload, &jsonPayload)
			err = json.Unmarshal(payload, s)
			if err != nil {
				return fmt.Errorf("Error parsing strategy table payload")
			}
*/

/*
			s.Playbook = jsonPayload["playbook"].(string)
			s.Counts = parseIntSlice(jsonPayload["counts"].([]interface{}))
			s.Bets = parseIntSlice(jsonPayload["bets"].([]interface{}))
			s.Insurance = jsonPayload["insurance"].(string)
			s.SoftDouble = parseStringMap(jsonPayload["soft-double"].(map[string]interface{}))
			s.HardDouble = parseStringMap(jsonPayload["hard-double"].(map[string]interface{}))
			s.PairSplit = parseStringMap(jsonPayload["pair-split"].(map[string]interface{}))
			s.SoftStand = parseStringMap(jsonPayload["soft-stand"].(map[string]interface{}))
			s.HardStand = parseStringMap(jsonPayload["hard-stand"].(map[string]interface{}))
*/
			return nil
		}
	}
	return fmt.Errorf("No matching strategy found")
}

func (s *Strategy) GetBet(seenCards []int) int {
	trueCount := s.getTrueCount(seenCards, s.getRunningCount(seenCards))
	bet := clamp(trueCount*2, constants.MinimumBet, constants.MaximumBet)
	if bet%2 != 0 {
		bet -= 1
	}
	return bet
}

func (s *Strategy) GetInsurance(seenCards *[13]int) bool {
	return false
}
func (s *Strategy) GetDouble(seenCards *[13]int, total int, soft bool, up *cards.Card) bool {
	return false
}
func (s *Strategy) GetSplit(seenCards *[13]int, pair, up *cards.Card) bool {
	return false
}
func (s *Strategy) GetStand(seenCards *[13]int, total int, soft bool, up *cards.Card) bool {
	return true
}

func (s *Strategy) getRunningCount(seenCards []int) int {
	running := 0
	for i, count := range s.Counts {
		running += count * seenCards[i]
	}
	return running
}

func (s *Strategy) getTrueCount(seenCards []int, runningCount int) int {
	unseen := s.NumberOfCards
	for _, card := range seenCards[2:12] {
		unseen -= card
	}
	if unseen > 0 {
		return int(float64(runningCount) / (float64(unseen) / 26.0))
	}
	return 0
}

func (s *Strategy) processValue(value string, trueCount int, missingValue bool) bool {
	if value == "" {
		return missingValue
	}

	switch strings.ToLower(value) {
	case "yes", "y":
		return true
	case "no", "n":
		return false
	}

	if strings.HasPrefix(value, "R") {
		v, _ := strconv.Atoi(value[1:])
		return trueCount <= v
	}
	v, _ := strconv.Atoi(value)
	return trueCount >= v
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}

func parseIntSlice(data []interface{}) []int {
	result := make([]int, len(data))
	for i, v := range data {
		result[i] = int(v.(float64))
	}
	return result
}

func parseStringMap(data map[string]interface{}) map[string][]string {
	result := make(map[string][]string)
	for key, val := range data {
		result[key] = parseStringSlice(val.([]interface{}))
	}
	return result
}

func parseStringSlice(data []interface{}) []string {
	result := make([]string, len(data))
	for i, v := range data {
		result[i] = v.(string)
	}
	return result
}

/*
func main() {
	// Example usage
	decks := "some_decks"
	strategy := "some_strategy"
	numberOfCards := 52
	strat := NewStrategy(decks, strategy, numberOfCards)
	fmt.Println("Playbook:", strat.Playbook)
}

*/

