package db

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/nickshine/boca-chica-bot/pkg/closures"
)

// NewClient returns a new Client
func NewClient() *Client {
	sess := session.Must(session.NewSession())
	ddb := dynamodb.New(sess)

	return (*Client)(ddb)
}

// Put will insert a road and/or beach closure notice in to the db.
//
// If the closure already exists in the db but something has changed (e.g. 'Status' has changed from
// 'Scheduled' to 'Cancelled'), then the new closure overwrites existing. The existing closure is
// returned in this case, otherwise nil.
//
// Closures are automatically deleted from the db table when their 'Expires' attribute becomes older
// than the current time (See DynamoDB Managed TTL).
func (client *Client) Put(tablename string, c *closures.Closure) (*closures.Closure, error) {
	input := buildPutInput(tablename, c)

	res, err := (*dynamodb.DynamoDB)(client).PutItem(input)
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

// RemoveClosure removes the given Closure from the db table.
func (client *Client) RemoveClosure(tablename string, c *closures.Closure) error {
	input := buildDeleteInput(tablename, c)
	_, err := (*dynamodb.DynamoDB)(client).DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
}

// QueryByTime returns a slice of Closures with timestamps close to the passed in time.
//
// This will filter items with TTL'd Timestamps <= the given time.
//
func (client *Client) QueryByTime(tablename string, t time.Time) ([]*closures.Closure, error) {
	input := buildTimeQueryInput(tablename, t)

	res, err := (*dynamodb.DynamoDB)(client).Query(input)
	if err != nil {
		return nil, fmt.Errorf("db query error: %v", err)
	}

	// res.Items contains the items in the table matching the query.
	var closures []*closures.Closure
	for _, item := range res.Items {
		closure, err := buildClosure(item)
		if err != nil {
			return closures, err
		}
		closures = append(closures, closure)
	}

	return closures, nil
}
