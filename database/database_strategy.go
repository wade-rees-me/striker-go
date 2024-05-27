package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	StrategyTableName = "StrikerStrategy"
)

func StrategyScan(strategy, rules, decks string) (*dynamodb.ScanOutput, error) {
	svc := dynamodb.New(CreateSession())
	input := &dynamodb.ScanInput{
		TableName: aws.String("StrikerStrategy"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				S: aws.String(strategy),
			},
			":r": {
				S: aws.String(rules),
			},
			":d": {
				S: aws.String(decks),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#R": aws.String("Rules"),
		},
		FilterExpression: aws.String("Strategy = :s and #R = :r and Decks = :d"),
	}
	return svc.Scan(input)
}
