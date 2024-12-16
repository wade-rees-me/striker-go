package table

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
    SoftDouble     *Chart
    HardDouble     *Chart
    PairSplit      *Chart
    SoftStand      *Chart
    HardStand      *Chart

	NumberOfCards int
	JsonResponse []map[string]interface{}
	JsonPayload []map[string]interface{}
}

func NewStrategy(decks, strategy string, numberOfCards int) *Strategy {
	s := &Strategy{NumberOfCards: numberOfCards}

	s.SoftDouble = NewChart("Soft Double")
	s.HardDouble = NewChart("Hard Double")
	s.PairSplit = NewChart("Pair Split")
	s.SoftStand = NewChart("Soft Stand")
	s.HardStand = NewChart("Hard Stand")

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

		s.SoftDouble.Print()
		s.HardDouble.Print()
		s.PairSplit.Print()
		s.SoftStand.Print()
		s.HardStand.Print()
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

			var result map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
				fmt.Println("Error parsing JSON:", err)
				log.Fatalf("Error parsing JSON string: %v", err)
			}

            s.Playbook = result["playbook"].(string)
            s.Counts = parseIntSlice(result["counts"].([]interface{}))
            s.Bets = parseIntSlice(result["bets"].([]interface{}))
            s.Insurance = result["insurance"].(string)

            parseStringMap(result["soft-double"].(map[string]interface{}), s.SoftDouble)
            parseStringMap(result["hard-double"].(map[string]interface{}), s.HardDouble)
            parseStringMap(result["pair-split"].(map[string]interface{}), s.PairSplit)
            parseStringMap(result["soft-stand"].(map[string]interface{}), s.SoftStand)
            parseStringMap(result["hard-stand"].(map[string]interface{}), s.HardStand)

			return nil
		}
	}
	return fmt.Errorf("No matching strategy found")
}

func (s *Strategy) GetBet(seenCards *[13]int) int {
	return s.getTrueCount(seenCards, s.getRunningCount(seenCards)) * constants.TrueCountBet
}

func (s *Strategy) GetInsurance(seenCards *[13]int) bool {
    trueCount := s.getTrueCount(seenCards, s.getRunningCount(seenCards))
    return s.processValue(s.Insurance, trueCount, false)
}

func (s *Strategy) GetDouble(seenCards *[13]int, total int, soft bool, up *cards.Card) bool {
    trueCount := s.getTrueCount(seenCards, s.getRunningCount(seenCards))
    if (soft) {
        return s.processValue(s.SoftDouble.GetValueByTotal(total, up.Offset), trueCount, false)
    }
    return s.processValue(s.HardDouble.GetValueByTotal(total, up.Offset), trueCount, false)
}

func (s *Strategy) GetSplit(seenCards *[13]int, pair, up *cards.Card) bool {
    trueCount := s.getTrueCount(seenCards, s.getRunningCount(seenCards))
    return s.processValue(s.PairSplit.GetValue(pair.Key, up.Offset), trueCount, false)
}

func (s *Strategy) GetStand(seenCards *[13]int, total int, soft bool, up *cards.Card) bool {
    trueCount := s.getTrueCount(seenCards, s.getRunningCount(seenCards))
    if (soft) {
        return s.processValue(s.SoftStand.GetValueByTotal(total, up.Offset), trueCount, false)
    }
    return s.processValue(s.HardStand.GetValueByTotal(total, up.Offset), trueCount, false)
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
		return int(float64(runningCount) / (float64(unseen) / float64(constants.TrueCounTMultiplier)))
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

func parseIntSlice(data []interface{}) []int {
	result := make([]int, len(data))
	for i, v := range data {
		result[i] = int(v.(float64))
	}
	return result
}

func parseStringMap(data map[string]interface{}, chart *Chart) {
	for key, val := range data {
		parseStringSlice(val.([]interface{}), key, chart)
	}
}

func parseStringSlice(data []interface{}, key string, chart *Chart) {
	for i, v := range data {
		chart.Insert(key, i, v.(string))
	}
}

