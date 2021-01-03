package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nickshine/boca-chica-bot/pkg/closures"
)

func buildPutInput(tablename string, c *closures.Closure) *dynamodb.PutItemInput {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"ClosureType": {
				S: aws.String(c.ClosureType),
			},
			"ClosureTypeSort": {
				S: aws.String(c.ClosureType + "#" + c.TimeType), // Primary Date#start, Primary Date#end
			},
			"Date": {
				S: aws.String(c.Date),
			},
			"RawTimeRange": {
				S: aws.String(c.RawTimeRange),
			},
			"Time": {
				N: aws.String(fmt.Sprint(c.Time)),
			},
			"TimeType": {
				S: aws.String(c.TimeType),
			},
			"Status": {
				S: aws.String(c.Status),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(tablename),
		ReturnValues:           aws.String("ALL_OLD"),
		// Only overwrite existing if something has changed
		ConditionExpression: aws.String("#Status <> :status OR #Time <> :time OR ClosureType <> :type OR RawTimeRange <> :rawTimeRange"), // nolint
		ExpressionAttributeNames: map[string]*string{
			"#Status": aws.String("Status"),
			"#Time":   aws.String("Time"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {
				S: aws.String(c.Status),
			},
			":time": {
				N: aws.String(fmt.Sprint(c.Time)),
			},
			":type": {
				S: aws.String(c.ClosureType),
			},
			":rawTimeRange": {
				S: aws.String(c.RawTimeRange),
			},
		},
	}

	return input
}

func buildDeleteInput(tablename string, c *closures.Closure) *dynamodb.DeleteItemInput {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"Date": {
				S: aws.String(c.Date),
			},
			"ClosureTypeSort": {
				S: aws.String(c.ClosureType + "#" + c.TimeType), // Primary Date#start, Primary Date#end
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
	}

	return input
}

func buildTimeQueryInput(tablename string, t time.Time) *dynamodb.QueryInput {
	input := &dynamodb.QueryInput{
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(tablename),
		KeyConditionExpression: aws.String("#Date = :date"),
		FilterExpression:       aws.String("#Time <= :time"),
		ExpressionAttributeNames: map[string]*string{
			"#Date": aws.String("Date"),
			"#Time": aws.String("Time"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":date": {
				S: aws.String(t.Format(closures.DateLayout)),
			},
			":time": {
				N: aws.String(fmt.Sprint(t.Unix())),
			},
		},
	}

	return input
}

func buildClosure(attributes map[string]*dynamodb.AttributeValue) (*closures.Closure, error) {
	if attributes == nil {
		return nil, nil
	}

	ct := aws.StringValue(attributes["ClosureType"].S)
	date := aws.StringValue(attributes["Date"].S)
	rawTimeRange := aws.StringValue(attributes["RawTimeRange"].S)
	timeString := aws.StringValue(attributes["Time"].N)
	timestamp, err := strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("problem parsing 'Time' attribute: %v", err)
	}
	timeType := aws.StringValue(attributes["TimeType"].S)
	status := aws.StringValue(attributes["Status"].S)

	c := &closures.Closure{
		ClosureType:  ct,
		Date:         date,
		RawTimeRange: rawTimeRange,
		Time:         timestamp,
		TimeType:     timeType,
		Status:       status,
	}

	return c, nil
}
