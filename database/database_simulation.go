// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	SimulationTableName = "StrikerSimulation"
)

type dbSimulation struct {
	Target    string
	Guid      string
	Hostname  string
	Strategy  string
	Rules     string
	Decks     string
	Epoch     int64
	Timestamp string
	Duration  string
	Payload   string
}

func SimulationInsert(target, guid, host, strategy, rules, decks string, epoch int64, duration, payload string) error {
	item := dbSimulation{
		Target:    target,
		Guid:      guid,
		Hostname:  host,
		Strategy:  strategy,
		Rules:     rules,
		Decks:     decks,
		Epoch:     epoch,
		Timestamp: time.Unix(epoch, 0).Format("2006-01-02 15:04:05.000"),
		Duration:  duration,
		Payload:   payload,
	}
fmt.Println(item)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	svc := dynamodb.New(CreateSession())
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(SimulationTableName),
	}

	_, err = svc.PutItem(input)
	return err
}

func SimulationDelete(guid string) error {
	svc := dynamodb.New(CreateSession())
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Guid": {
				S: aws.String(guid),
			},
		},
		TableName: aws.String(SimulationTableName),
	}

	_, err := svc.DeleteItem(input)
	return err
}
