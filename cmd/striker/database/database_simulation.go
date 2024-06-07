package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	SimulationTableName = "StrikerSimulation"
)

type dbSimulationPayload struct {
}

type DBSimulationTable struct {
	Style     string
	Guid      string
	Target    string
	Hostname  string
	Status    string
	Epoch     int64
	Timestamp string
	Payload   string
}

func (s *DBSimulationTable) SimulationUpdate() error {
	databaseLog.Debug(fmt.Sprintf("Update item in Simulation table: %v", SimulationTableName))
	s.setTime()

	key := map[string]*dynamodb.AttributeValue{
		"Style": {
			S: aws.String(s.Style),
		},
		"Guid": {
			S: aws.String(s.Guid),
		},
	}
	updateExpression := "SET #S = :s, #P = :p"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":s": {
			S: aws.String(s.Status),
		},
		":p": {
			S: aws.String(s.Payload),
		},
	}
	expressionAttributeNames := map[string]*string{
		"#S": aws.String("Status"),
		"#P": aws.String("Payload"),
	}

	return UpdateItemInTable(SimulationTableName, key, updateExpression, expressionAttributeValues, expressionAttributeNames)
}

func (s *DBSimulationTable) SimulationInsert() error {
	databaseLog.Debug(fmt.Sprintf("Insert from Simulation table: %v", SimulationTableName))
	s.setTime()
	av, err := dynamodbattribute.MarshalMap(s)
	if err != nil {
		return err
	}
	return InsertItemIntoTable(SimulationTableName, av)
}

func (s *DBSimulationTable) SimulationDelete() error {
	databaseLog.Debug(fmt.Sprintf("Delete from Simulation table: %v", SimulationTableName))
	key := map[string]*dynamodb.AttributeValue{
		"Style": {
			S: aws.String(strings.ToLower(s.Style)),
		},
		"Guid": {
			S: aws.String(s.Guid),
		},
	}
	return DeleteItemFromTable(SimulationTableName, key)
}

func (s *DBSimulationTable) setTime() {
	s.Epoch = time.Now().Unix()
	timeFromEpoch := time.Unix(s.Epoch, 0)
	s.Timestamp = timeFromEpoch.Format("2006-01-02 15:04:05")
}
