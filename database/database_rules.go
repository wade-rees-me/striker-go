package database

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	RulesTableName = "StrikerRules"
)

type TableRules struct {
	Rules             string
	Decks             string
	HitSoft17         bool
	Surrender         bool
	DoubleAnyTwoCards bool
	DoubleAfterSplit  bool
	ResplitAces       bool
	HitSplitAces      bool
}

func RulesSelect(rules, decks string) (*TableRules, error) {
	item := TableRules{}
	svc := dynamodb.New(CreateSession())
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(RulesTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Rules": {
				S: aws.String(strings.ToLower(rules)),
			},
			"Decks": {
				S: aws.String(strings.ToLower(decks)),
			},
		},
	})
	if err != nil {
		return &item, err
	}

	if result.Item == nil {
		return &item, errors.New("Could not find table rules")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	return &item, err
}
