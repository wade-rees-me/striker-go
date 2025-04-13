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
	Id         string `json:"_id"`
	Playbook   string `json:"playbook"`
	Counts     []int  `json:"counts"`
	Insurance  string `json:"insurance"`
	SoftDouble *Chart
	HardDouble *Chart
	PairSplit  *Chart
	SoftStand  *Chart
	HardStand  *Chart

	NumberOfCards int
	JsonResponse  map[string]interface{}
}

func NewStrategy(decks, strategy string, numberOfDecks int) *Strategy {
	s := &Strategy{NumberOfCards: numberOfDecks * constants.NumberOfCardsInDeck}

	s.SoftDouble = NewChart("Soft Double")
	s.HardDouble = NewChart("Hard Double")
	s.PairSplit = NewChart("Pair Split")
	s.SoftStand = NewChart("Soft Stand")
	s.HardStand = NewChart("Hard Stand")

	if strategy != "mimic" {
		err := s.fetchJson("http://" + constants.ChartsUrl + "/" + decks + "/" + strategy)
		if err != nil {
			log.Fatalf("Error fetching JSON: %v", err)
		}
		s.fetchTable(decks, strategy)

		if false {
			s.SoftDouble.Print()
			s.HardDouble.Print()
			s.PairSplit.Print()
			s.SoftStand.Print()
			s.HardStand.Print()
			s.PrintCounts()
		}
	}
	return s
}

func (s *Strategy) fetchJson(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	bodyString = constants.UnescapeJSON(bodyString)
	bodyString = constants.StripQuotes(bodyString)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	err = json.Unmarshal([]byte(bodyString), &s.JsonResponse)
	return err
}

func (s *Strategy) fetchTable(decks, strategy string) {
	if playbook, ok := s.JsonResponse["playbook"].(string); ok {
		s.Playbook = playbook
	}

	s.Insurance = s.JsonResponse["insurance"].(string)
	s.Counts = parseIntSlice(s.JsonResponse["counts"].([]interface{}))
	s.Counts = append([]int{0, 0}, s.Counts...)

	parseStringMap(s.JsonResponse["soft-double"].(map[string]interface{}), s.SoftDouble)
	parseStringMap(s.JsonResponse["hard-double"].(map[string]interface{}), s.HardDouble)
	parseStringMap(s.JsonResponse["pair-split"].(map[string]interface{}), s.PairSplit)
	parseStringMap(s.JsonResponse["soft-stand"].(map[string]interface{}), s.SoftStand)
	parseStringMap(s.JsonResponse["hard-stand"].(map[string]interface{}), s.HardStand)
}

func (s *Strategy) GetBet(seenCards *[cards.MAXIMUM_CARD_VALUE + 1]int) int {
	return s.getTrueCount(seenCards, s.getRunningCount(seenCards)) * constants.TrueCountBet
}

func (s *Strategy) GetInsurance(seenCards *[cards.MAXIMUM_CARD_VALUE + 1]int) bool {
	trueCount := s.getTrueCount(seenCards, s.getRunningCount(seenCards))
	return s.processValue(s.Insurance, trueCount, false)
}

func (s *Strategy) GetDouble(seenCards *[cards.MAXIMUM_CARD_VALUE + 1]int, total int, soft bool, up *cards.Card) bool {
	trueCount := s.getTrueCount(seenCards, s.getRunningCount(seenCards))
	if soft {
		return s.processValue(s.SoftDouble.GetValueByTotal(total, up.Value), trueCount, false)
	}
	return s.processValue(s.HardDouble.GetValueByTotal(total, up.Value), trueCount, false)
}

func (s *Strategy) GetSplit(seenCards *[cards.MAXIMUM_CARD_VALUE + 1]int, pair, up *cards.Card) bool {
	trueCount := s.getTrueCount(seenCards, s.getRunningCount(seenCards))
	return s.processValue(s.PairSplit.GetValue(pair.Key, up.Value), trueCount, false)
}

func (s *Strategy) GetStand(seenCards *[cards.MAXIMUM_CARD_VALUE + 1]int, total int, soft bool, up *cards.Card) bool {
	trueCount := s.getTrueCount(seenCards, s.getRunningCount(seenCards))
	if soft {
		return s.processValue(s.SoftStand.GetValueByTotal(total, up.Value), trueCount, false)
	}
	return s.processValue(s.HardStand.GetValueByTotal(total, up.Value), trueCount, false)
}

func (s *Strategy) getRunningCount(seenCards *[cards.MAXIMUM_CARD_VALUE + 1]int) int {
	running := 0
	for i, count := range s.Counts {
		running += count * seenCards[i]
	}
	return running
}

func (s *Strategy) getTrueCount(seenCards *[cards.MAXIMUM_CARD_VALUE + 1]int, runningCount int) int {
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

// Print prints the entire chart to the console
func (s *Strategy) PrintCounts() {
	fmt.Println("Counts")
	fmt.Println("--------------------2-----3-----4-----5-----6-----7-----8-----9-----X-----A---")
	fmt.Printf("     ")
	for _, value := range s.Counts {
		fmt.Printf("%4d, ", value)
	}
	fmt.Println()
	fmt.Println("------------------------------------------------------------------------------")
}
