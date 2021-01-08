package db

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nickshine/boca-chica-bot/pkg/closures"
)

func buildPutInput(tablename string, c *closures.Closure) *dynamodb.PutItemInput {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"ClosureType": {
				S: aws.String(string(c.ClosureType)),
			},
			"Date": {
				S: aws.String(c.Date),
			},
			"RawTimeRange": {
				S: aws.String(c.RawTimeRange),
			},
			"TimeRangeStatus": {
				S: aws.String(string(c.TimeRangeStatus)),
			},
			"ClosureStatus": {
				S: aws.String(string(c.ClosureStatus)),
			},
			"Expires": {
				N: aws.String(fmt.Sprint(c.Expires)),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(tablename),
		ReturnValues:           aws.String("ALL_OLD"),
		// Only overwrite existing if something has changed
		ConditionExpression: aws.String("ClosureStatus <> :closureStatus OR TimeRangeStatus <> :timeRangeStatus OR ClosureType <> :type OR RawTimeRange <> :rawTimeRange"), // nolint
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":closureStatus": {
				S: aws.String(string(c.ClosureStatus)),
			},
			":timeRangeStatus": {
				S: aws.String(string(c.TimeRangeStatus)),
			},
			":type": {
				S: aws.String(string(c.ClosureType)),
			},
			":rawTimeRange": {
				S: aws.String(c.RawTimeRange),
			},
		},
	}

	return input
}

func buildClosure(attributes map[string]*dynamodb.AttributeValue) (*closures.Closure, error) {
	if attributes == nil {
		return nil, nil
	}

	closureType := closures.ClosureType(aws.StringValue(attributes["ClosureType"].S))
	date := aws.StringValue(attributes["Date"].S)
	rawTimeRange := aws.StringValue(attributes["RawTimeRange"].S)
	timeRangeStatus := closures.TimeRangeStatus(aws.StringValue(attributes["TimeRangeStatus"].S))
	closureStatus := closures.ClosureStatus(aws.StringValue(attributes["ClosureStatus"].S))
	expiresString := aws.StringValue(attributes["Expires"].N)
	expires, err := strconv.ParseInt(expiresString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("problem parsing 'Expires' attribute: %v", err)
	}

	c := &closures.Closure{
		ClosureType:     closureType,
		Date:            date,
		RawTimeRange:    rawTimeRange,
		TimeRangeStatus: timeRangeStatus,
		ClosureStatus:   closureStatus,
		Expires:         expires,
	}

	return c, nil
}
