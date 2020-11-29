package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nickshine/boca-chica-bot/pkg/closures"
)

func buildPutInput(c *closures.Closure) *dynamodb.PutItemInput {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"ClosureType": {
				S: aws.String(c.ClosureType),
			},
			"Date": {
				S: aws.String(c.Date),
			},
			"Time": {
				S: aws.String(c.Time),
			},
			"Start": {
				S: aws.String(c.Start.Format(time.RFC3339)),
			},
			"End": {
				S: aws.String(c.End.Format(time.RFC3339)),
			},
			"Status": {
				S: aws.String(c.Status),
			},
			"Expires": {
				N: aws.String(fmt.Sprint(c.Expires)),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(tablename),
		ReturnValues:           aws.String("ALL_OLD"),
		// Only overwrite existing if something has changed
		ConditionExpression: aws.String("ClosureType <> :type OR #Status <> :status OR #Start <> :start OR #End <> :end"),
		ExpressionAttributeNames: map[string]*string{
			"#Status": aws.String("Status"),
			"#Start":  aws.String("Start"),
			"#End":    aws.String("End"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":end": {
				S: aws.String(c.End.Format(time.RFC3339)),
			},
			":start": {
				S: aws.String(c.Start.Format(time.RFC3339)),
			},
			":status": {
				S: aws.String(c.Status),
			},
			":type": {
				S: aws.String(c.ClosureType),
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
	timeRange := aws.StringValue(attributes["Time"].S)
	startString := aws.StringValue(attributes["Start"].S)
	start, err := time.Parse(time.RFC3339, startString)
	if err != nil {
		return nil, fmt.Errorf("problem parsing 'Start' attribute: %v", err)
	}
	endString := aws.StringValue(attributes["End"].S)
	end, err := time.Parse(time.RFC3339, endString)
	if err != nil {
		return nil, fmt.Errorf("problem parsing 'End' attribute: %v", err)
	}
	status := aws.StringValue(attributes["Status"].S)
	expiresString := aws.StringValue(attributes["Expires"].N)
	expires, err := strconv.ParseInt(expiresString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("problem parsing 'Expires' attribute: %v", err)
	}

	c := &closures.Closure{
		ClosureType: ct,
		Date:        date,
		Time:        timeRange,
		Start:       start,
		End:         end,
		Status:      status,
		Expires:     expires,
	}

	return c, nil
}
