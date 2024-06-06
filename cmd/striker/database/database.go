// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	//"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/wade-rees-me/striker-go/logger"
)

const (
	debugFileName = "_database.log"
)

var databaseLog = logger.NewLogger(os.Stdout, os.Stdout, os.Stderr, ioutil.Discard)

func init() {
	databaseLog.OpenDebugFile(debugFileName)
	databaseLog.Debug("Starting Striker-Database")
}

func Finish() {
	databaseLog.CloseDebugFile()
}

func InsertItemIntoTable(tableName string, av map[string]*dynamodb.AttributeValue) error {
	databaseLog.Debug(fmt.Sprintf("Insert into table: %s", tableName))
	databaseLog.Debug(fmt.Sprintf("Item: %v", av))
	svc := dynamodb.New(CreateSession())
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}
	_, err := svc.PutItem(input)
	return err
}

func UpdateItemInTable(tableName string, key map[string]*dynamodb.AttributeValue, updateExpression string, expressionAttributeValues map[string]*dynamodb.AttributeValue, expressionAttributeNames map[string]*string) error {
	databaseLog.Debug(fmt.Sprintf("Update item in table: %s", tableName))
	databaseLog.Debug(fmt.Sprintf("Using key: %v", key))
	svc := dynamodb.New(CreateSession())
	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ExpressionAttributeNames:  expressionAttributeNames,
	}
	_, err := svc.UpdateItem(input)
	return err
}

func DeleteItemFromTable(tableName string, key map[string]*dynamodb.AttributeValue) error {
	databaseLog.Debug(fmt.Sprintf("Delete item from table: %s", tableName))
	databaseLog.Debug(fmt.Sprintf("Using key: %v", key))
	svc := dynamodb.New(CreateSession())
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}
	_, err := svc.DeleteItem(input)
	return err
}

func SelectItemFromTable(tableName string, key map[string]*dynamodb.AttributeValue) (*dynamodb.GetItemOutput, error) {
	databaseLog.Debug(fmt.Sprintf("Select from table: %s", tableName))
	databaseLog.Debug(fmt.Sprintf("Using key: %s", key))
	svc := dynamodb.New(CreateSession())
	return svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})
}

func ScanItemsFromTable(tableName string, values map[string]*dynamodb.AttributeValue, names map[string]*string, filter string) (*dynamodb.ScanOutput, error) {
	databaseLog.Debug(fmt.Sprintf("Scan items from table: %s", tableName))
	databaseLog.Debug(fmt.Sprintf("Using values: %s", values))
	databaseLog.Debug(fmt.Sprintf("Using names: %v", names))
	databaseLog.Debug(fmt.Sprintf("Using filter: %s", filter))
	svc := dynamodb.New(CreateSession())
	input := &dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeValues: values,
		ExpressionAttributeNames:  names,
		FilterExpression:          aws.String(filter),
	}
	return svc.Scan(input)
}

// Initialize a session that the SDK will use to load
// credentials from the shared credentials file ~/.aws/credentials
// and region from the shared configuration file ~/.aws/config.
func CreateSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}
