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
	ProcessingTableName = "StrikerProcessing"
)

type dbProcessing struct {
	Target    string
	Guid      string
	Hostname  string
	Epoch     int64
	Timestamp string
	Payload   string
}

func ProcessingInsert(target, guid, host string, epoch int64, payload string) error {
	item := dbProcessing{
		Target:    target,
		Guid:      guid,
		Hostname:  host,
		Epoch:     epoch,
		Timestamp: time.Unix(epoch, 0).Format("2006-01-02 15:04:05.000"),
		Payload:   payload,
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	svc := dynamodb.New(CreateSession())
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(ProcessingTableName),
	}

	_, err = svc.PutItem(input)
	return err
}

func ProcessingDelete(guid string) error {
	svc := dynamodb.New(CreateSession())
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Guid": {
				S: aws.String(guid),
			},
		},
		TableName: aws.String(ProcessingTableName),
	}

	_, err := svc.DeleteItem(input)
	return err
}
