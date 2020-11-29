package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/nickshine/boca-chica-bot/pkg/closures"
)

const tablename = "BocaChicaBot-closures"

var svc *dynamodb.DynamoDB

func init() {
	sess := session.Must(session.NewSession())
	svc = dynamodb.New(sess)
}

// Put will insert a road and/or beach closure notice in to the db.
//
// If the closure already exists in the db but something has changed (e.g. 'Status' has changed from
// 'Scheduled' to 'Cancelled'), then the new closure overwrites existing. The existing closure is
// returned in this case, otherwise nil.
//
// Closures are automatically deleted from the db table when their 'Expires' attribute becomes older
// than the current time (See DynamoDB Managed TTL).
func Put(c *closures.Closure) (*closures.Closure, error) {
	input := buildPutInput(c)

	res, err := svc.PutItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return nil, NewItemUnchangedError()
			default:
				return nil, fmt.Errorf("put item failure: %v", aerr.Error())
			}
		}

		return nil, fmt.Errorf("unknown put item failure: %v", err.Error())
	}

	// res.Attributes will contain the existing closure if an attribute has changed.
	// See the ConditionalExpression on PutItem input.
	// New closures added to db will return nil.
	closure, err := buildClosure(res.Attributes)
	if err != nil {
		return nil, err
	}

	return closure, nil
}

/*
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
*/
