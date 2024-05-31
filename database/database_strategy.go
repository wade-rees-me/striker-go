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

func StrategyUpdate(id, hand, payload string) (*dynamodb.ScanOutput, error) {
	svc := dynamodb.New(CreateSession())
    input := &dynamodb.UpdateItemInput{
		TableName: aws.String("StrikerStrategy"),
        Key: map[string]*dynamodb.AttributeValue{
            "StrategyDecksRules": {
                S: aws.String(id),
            },
            "Hand": {
                S: aws.String(hand),
            },
        },
		ExpressionAttributeNames: map[string]*string{
			"#P": aws.String("Map"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p": {
				S: aws.String(payload),
			}
		},
        UpdateExpression: aws.String("SET #P = :p"),
        ReturnValues:     aws.String("UPDATED_NEW"),
	}
    return svc.UpdateItem(input)
}

/*
func updateItem(sess *session.Session, tableName string, id string, name string, value string) error {
    svc := dynamodb.New(sess)

    // Create the item to update
    item := Item{
        ID:    id,
        Name:  name,
        Value: value,
    }

    // Convert the item to a map of AttributeValues
    av, err := dynamodbattribute.MarshalMap(item)
    if err != nil {
        return fmt.Errorf("failed to marshal item: %v", err)
    }

    // Create the input for the update item
    input := &dynamodb.UpdateItemInput{
        TableName: aws.String(tableName),
        Key: map[string]*dynamodb.AttributeValue{
            "ID": {
                S: aws.String(id),
            },
        },
        ExpressionAttributeNames: map[string]*string{
            "#N": aws.String("Name"),
            "#V": aws.String("Value"),
        },
        ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
            ":n": {
                S: aws.String(name),
            },
            ":v": {
                S: aws.String(value),
            },
        },
        UpdateExpression: aws.String("SET #N = :n, #V = :v"),
        ReturnValues:     aws.String("UPDATED_NEW"),
    }
    input := &dynamodb.UpdateItemInput{
		TableName: aws.String("StrikerStrategy"),
        Key: map[string]*dynamodb.AttributeValue{
            "StrategyDecksRules": {
                S: aws.String(id),
            },
        },
		ExpressionAttributeNames: map[string]*string{
			"#P": aws.String("Map"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p": {
				S: aws.String(payload),
			}
		},
        UpdateExpression: aws.String("SET #P = :p"),
        ReturnValues:     aws.String("UPDATED_NEW"),
	}

    // Update the item
    _, err = svc.UpdateItem(input)
    if err != nil {
        return fmt.Errorf("failed to update item: %v", err)
    }

    return nil
}
*/
