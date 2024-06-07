package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	RulesTableName = "StrikerRules"
)

type DBRulesPayload struct {
	HitSoft17         bool
	Surrender         bool
	DoubleAnyTwoCards bool
	DoubleAfterSplit  bool
	ResplitAces       bool
	HitSplitAces      bool
}

type DBRulesTable struct {
	Style     string
	Epoch     int64
	Timestamp string
	Payload   string
}

func (r *DBRulesTable) Insert() error {
	databaseLog.Debug(fmt.Sprintf("Insert from Rules table: %v", RulesTableName))
	r.setTime()
	key := map[string]*dynamodb.AttributeValue{
		"Style": {
			S: aws.String(strings.ToLower(r.Style)),
		},
	}
	return InsertItemIntoTable(RulesTableName, key)
}

func (r *DBRulesTable) Update() error {
	databaseLog.Debug(fmt.Sprintf("Update from Rules table: %v", RulesTableName))
	r.setTime()
	//key := map[string]*dynamodb.AttributeValue{
	//"Style": {
	//S: aws.String(strings.ToLower(r.Style)),
	//},
	//}
	//return UpdateItemInTable(RulesTableName, key)
	return nil
}

func (r *DBRulesTable) Delete() error {
	databaseLog.Debug(fmt.Sprintf("Delete from Rules table: %v", RulesTableName))
	r.setTime()
	key := map[string]*dynamodb.AttributeValue{
		"Style": {
			S: aws.String(strings.ToLower(r.Style)),
		},
	}
	return DeleteItemFromTable(RulesTableName, key)
}

func (r *DBRulesTable) Select() (*DBRulesPayload, error) {
	databaseLog.Debug(fmt.Sprintf("Select from Rules table: %v", RulesTableName))
	payload := new(DBRulesPayload)
	key := map[string]*dynamodb.AttributeValue{
		"Style": {
			S: aws.String(strings.ToLower(r.Style)),
		},
	}
	databaseLog.Debug(fmt.Sprintf("  Using key: %s", key))
	result, err := SelectItemFromTable(RulesTableName, key)
	if err != nil || result.Item == nil {
		return payload, errors.New("could not find table rules")
	}
	databaseLog.Debug(fmt.Sprintf("  Results: %v", result.Item))
	if err = dynamodbattribute.UnmarshalMap(result.Item, r); err != nil {
		return payload, errors.New("could not parse table rules")
	}
	err = json.Unmarshal([]byte(r.Payload), &payload)
	databaseLog.Debug(fmt.Sprintf("  Results: %v", *payload))
	return payload, err
}

func (r *DBRulesTable) setTime() {
	r.Epoch = time.Now().Unix()
	timeFromEpoch := time.Unix(r.Epoch, 0)
	r.Timestamp = timeFromEpoch.Format("2006-01-02 15:04:05")
}
