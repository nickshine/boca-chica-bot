package db

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/nickshine/boca-chica-bot/closure"
)

const tablename = "BocaChicaBot-closures"

var svc *dynamodb.DynamoDB

func init() {
	sess := session.Must(session.NewSession())
	svc = dynamodb.New(sess)
}

func buildPutInput(c *closure.Closure) *dynamodb.PutItemInput {
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
		// ConditionExpression:    aws.String("attribute_not_exists(#D) AND attribute_not_exists(#T) OR #S <> :status OR ClosureType <> :type"),
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

func Put(c *closure.Closure) {

	if time.Now().Unix() > c.Expires {
		// don't bother adding expired closures to table since the TTL will remove them anyways
		return
	}

	input := buildPutInput(c)

	// check if closure already exists. If so, figure out if anything has changed

	res, err := svc.PutItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeTransactionConflictException:
				fmt.Println(dynamodb.ErrCodeTransactionConflictException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())

			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	// changed := res.Attributes
	// fmt.Printf("\nchanged: %+v\n", changed)

	fmt.Println(res)

}

func Info() {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tablename),
	}

	res, err := svc.DescribeTable(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(res)

}
