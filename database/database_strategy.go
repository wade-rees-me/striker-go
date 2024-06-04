package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	StrategyTableName = "StrikerStrategy"
)

type dbStrategyPayload struct {
}

type dbStrategy struct {
	Style     string
	Epoch     int64
	timestamp string
	Payload   string
}

func StrategyInsert(style, strategy string) error {
	return nil
}

func StrategyUpdate(style, strategy string) error {
	// Inset if not in table
	return nil
}

func StrategyGet(style, strategy string) error {
	return nil
}

func StrategyScan(style, strategy string) (*dynamodb.ScanOutput, error) {
	databaseLog.Debug(fmt.Sprintf("Scan items from Strategy table: %v", StrategyTableName))
	values := map[string]*dynamodb.AttributeValue{
		":t": {
			S: aws.String(style),
		},
	}
	names := map[string]*string{
		"#S": aws.String("Style"),
	}
	filter := "#S = :t"
	return ScanItemsFromTable(StrategyTableName, values, names, filter)
}

/*
func (s *DBSimulationTable) setTime() {
	s.Epoch = time.Now().Unix()
	timeFromEpoch := time.Unix(s.Epoch, 0)
	s.Timestamp = timeFromEpoch.Format("2006-01-02 15:04:05")
}
*/
