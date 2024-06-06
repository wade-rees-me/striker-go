package simulation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wade-rees-me/striker-go/arguments"
	"github.com/wade-rees-me/striker-go/constants"
	"github.com/wade-rees-me/striker-go/database"
	"github.com/wade-rees-me/striker-go/logger"
	"github.com/wade-rees-me/striker-go/tables"
	"github.com/wade-rees-me/striker-go/utilities"
)

type SimulationReports struct {
	Reports []string
}

type SimulationParameters struct {
	Guid          string
	Timestamp     int64
	Target        string
	Strategy      string
	Rules         string
	Decks         string
	Tables        int
	Rounds        int
	BlackjackPays string
	Penetration   float64
}

type Simulation struct {
	Name       string
	Guid       string
	Hostname   string
	Duration   time.Duration
	ElapsedTime string
	Year       int
	Month      int
	Day        int
	Parameters SimulationParameters
	TableRules *database.DBRulesPayload
	tableList  []tables.Table
}

func NewSimulation() *Simulation {
	t := time.Now()
	s := new(Simulation)
	s.Year = t.Year()
	s.Month = int(t.Month())
	s.Day = t.Day()
	s.Name = fmt.Sprintf("go-striker-%4d_%02d_%02d_%012d", s.Year, s.Month, s.Day, t.Unix())
	s.Guid = uuid.New().String()
	s.Hostname = getHostname()

	s.Parameters.Guid = uuid.New().String()
	s.Parameters.Timestamp = time.Now().UTC().Unix()
	s.Parameters.Target = constants.StrikerWhoAmI

	s.Parameters.Strategy = arguments.CLStrategy.Get()
	s.Parameters.Rules = arguments.CLRules.Get()
	s.Parameters.Decks = arguments.CLTable.Get()

	s.Parameters.Tables = int(math.Max(math.Min(float64(arguments.CLSimulation.Tables), float64(constants.MaxNumberOfTables)), float64(constants.MinNumberOfTables)))
	s.Parameters.Rounds = int(math.Max(math.Min(float64(arguments.CLSimulation.Rounds), float64(constants.MaxNumberOfRounds)), float64(constants.MinNumberOfRounds)))
	s.Parameters.BlackjackPays = arguments.CLSimulation.BlackjackPays
	s.Parameters.Penetration = arguments.CLSimulation.Penetration

	s.TableRules = getTableRules(s.Parameters.Rules, s.Parameters.Decks)

	strategy := s.newPlayStrategy()
	for i := 1; i <= s.Parameters.Tables; i++ {
		s.tableList = append(s.tableList, *s.getTable(i, tables.NewPlayer(s.TableRules, strategy, s.Parameters.BlackjackPays)))
	}

	return s
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func getTableRules(rules, decks string) *database.DBRulesPayload {
	rulesTable := new(database.DBRulesTable)
	rulesTable.Style = fmt.Sprintf("%s-%s", rules, decks)
	tr, err := rulesTable.Select()
	if err != nil {
		panic(fmt.Sprintf("Cannot get table rules: %v", err))
	}
	return tr
}

func (s *Simulation) getTable(tableNumber int, player *tables.Player) *tables.Table {
	decks := 1

	if strings.ToLower(s.Parameters.Decks) == "double" {
		decks = 2
	}
	if strings.ToLower(s.Parameters.Decks) == "multi" {
		decks = 6
	}

	return tables.NewTable(s.TableRules, player, tableNumber, decks, s.Parameters.Penetration, s.Name, s.Guid)
}

func (s *Simulation) RunSimulation() {
	var wg sync.WaitGroup
	var start = time.Now()

	logger.Log.Info(fmt.Sprintf("Simulation %v, started at %v", s.Name, start))
	wg.Add(len(s.tableList))
	for i := range s.tableList {
		t := &s.tableList[i]
		go t.Session(&wg, s.Parameters.Rounds)
	}
	wg.Wait()

	end := time.Now()
	s.Duration = time.Since(start).Round(time.Second)
	s.ElapsedTime = s.Duration.String()
	logger.Log.Info(fmt.Sprintf("Simulation %v, ended at %v, total elapsed time: %s", s.Name, end, s.Duration))
}

func (s *Simulation) newPlayStrategy() *tables.PlayStrategy {
	ps := new(tables.PlayStrategy)

	results, err := database.StrategyScan(fmt.Sprintf("%s-%s-%s", strings.ToLower(s.Parameters.Rules), strings.ToLower(s.Parameters.Decks), strings.ToLower(s.Parameters.Strategy)), "HardDouble")
	if err != nil {
		panic("Failed to scan table:")
	}

	for _, item := range results.Items {
		if handAttr, ok := item["Strategy"]; ok && handAttr.S != nil {
			hand := *handAttr.S
			switch hand {
			case "HardDouble":
				ps.HardDouble = unmarshalStrategyMap(item["Payload"].S)
			case "SoftDouble":
				ps.SoftDouble = unmarshalStrategyMap(item["Payload"].S)
			case "PairSplit":
				ps.PairSplit = unmarshalStrategyMap(item["Payload"].S)
			case "HardStand":
				ps.HardStand = unmarshalStrategyMap(item["Payload"].S)
			case "SoftStand":
				ps.SoftStand = unmarshalStrategyMap(item["Payload"].S)
			default:
				panic("Unknown hand type:")
			}
		} else {
			panic(fmt.Sprintf("Missing hand type: %v", item))
		}
	}

	if ps.HardDouble == nil || ps.SoftDouble == nil || ps.PairSplit == nil || ps.HardStand == nil || ps.SoftStand == nil {
		panic("Missing hand type:")
	}

	//fmt.Println(s.HardDouble)
	//fmt.Println(s.SoftDouble)
	//fmt.Println(s.PairSplit)
	//fmt.Println(s.HardStand)
	//fmt.Println(s.SoftStand)

	return ps
}

func unmarshalStrategyMap(data *string) map[int][]string {
	c := map[int][]string{}

	err := json.Unmarshal([]byte(*data), &c)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Unmarshal: %v", err))
	}
	return c
}

func (s *Simulation) getReport() string {
	report := new(SimulationReports)
	for i := range s.tableList {
		t := &s.tableList[i]
		r := *t.GetReport()
		b, err := json.Marshal(r)
		if err != nil {
			panic(err)
		}
		report.Reports = append(report.Reports, string(b))
		fmt.Printf("\n%v\n", utilities.JsonPrettyPrint(string(b)))
	}
	b, err := json.Marshal(report)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (s *Simulation) PrintSimulation() {
	var out bytes.Buffer
	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Indent(&out, []byte(string(b)), "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Simulation:\n%v\n", out.String())
}
