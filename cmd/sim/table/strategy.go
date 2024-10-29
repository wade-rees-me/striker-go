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
	NumberOfCards int
	JsonResponse []map[string]interface{}
	JsonPayload []map[string]interface{}
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
	return err
}

func (s *Strategy) fetchTable(decks, strategy string) error {
	for _, item := range s.JsonResponse {
		if item["playbook"] == decks && item["hand"] == strategy {
			payload, err := json.Marshal(item["payload"])
			if err != nil {
				return err
			}

			payString := string(payload)
			newPay := payString[1 : len(payString)-1]
			jsonStr := strings.ReplaceAll(newPay, "\\", "")
//fmt.Printf("jsonStr:::: %v\n", jsonStr)

			var result map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
				fmt.Println("Error parsing JSON:", err)
				log.Fatalf("Error parsing JSON string: %v", err)
			}
//fmt.Printf("result:::: %v\n", result)

            s.Playbook = result["playbook"].(string)
            s.Counts = parseIntSlice(result["counts"].([]interface{}))
            s.Bets = parseIntSlice(result["bets"].([]interface{}))
            s.Insurance = result["insurance"].(string)
            s.SoftDouble = parseStringMap(result["soft-double"].(map[string]interface{}))
            s.HardDouble = parseStringMap(result["hard-double"].(map[string]interface{}))
            s.PairSplit = parseStringMap(result["pair-split"].(map[string]interface{}))
            s.SoftStand = parseStringMap(result["soft-stand"].(map[string]interface{}))
            s.HardStand = parseStringMap(result["hard-stand"].(map[string]interface{}))
//fmt.Printf("strategy:::: %v\n", s)

			return nil
		}
	}
	return fmt.Errorf("No matching strategy found")
}

func (s *Strategy) GetBet(seenCards *[13]int) int {
	trueCount := s.getTrueCount(seenCards, s.getRunningCount(seenCards))
	bet := clamp(trueCount*2, constants.MinimumBet, constants.MaximumBet)
	if bet%2 != 0 {
		bet -= 1
	}
	return bet
}

func (s *Strategy) GetInsurance(seenCards *[13]int) bool {
    trueCount := s.getTrueCount(seenCards, s.getRunningCount(seenCards))
    return s.processValue(s.Insurance, trueCount, false)
}

func (s *Strategy) GetDouble(seenCards *[13]int, total int, soft bool, up *cards.Card) bool {
//fmt.Printf("getDouble\n")
    trueCount := s.getTrueCount(seenCards, s.getRunningCount(seenCards))
    if (soft) {
        return s.processValue(s.SoftDouble[strconv.Itoa(total)][up.Offset], trueCount, false)
    }
    return s.processValue(s.HardDouble[strconv.Itoa(total)][up.Offset], trueCount, false)
}

func (s *Strategy) GetSplit(seenCards *[13]int, pair, up *cards.Card) bool {
    trueCount := s.getTrueCount(seenCards, s.getRunningCount(seenCards))
    return s.processValue(s.PairSplit[strconv.Itoa(pair.Value)][up.Offset], trueCount, false)
}
func (s *Strategy) GetStand(seenCards *[13]int, total int, soft bool, up *cards.Card) bool {
    trueCount := s.getTrueCount(seenCards, s.getRunningCount(seenCards))
    if (soft) {
        return s.processValue(s.SoftStand[strconv.Itoa(total)][up.Offset], trueCount, false)
    }
    return s.processValue(s.HardStand[strconv.Itoa(total)][up.Offset], trueCount, false)
}

func (s *Strategy) getRunningCount(seenCards *[13]int) int {
	running := 0
	for i, count := range s.Counts {
		running += count * seenCards[i]
	}
	return running
}

func (s *Strategy) getTrueCount(seenCards *[13]int, runningCount int) int {
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

