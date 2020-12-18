package db

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Client is a dynamodb client with some convenience receiver functions added.
type Client dynamodb.DynamoDB

// ItemUnchangedError is an implementation of the error interface used to signify items in the db
// exist and have no changes.
type ItemUnchangedError struct {
	message string
}

// NewItemUnchangedError returns a new instance of ItemUnchangedError.
func NewItemUnchangedError() *ItemUnchangedError {
	return &ItemUnchangedError{
		message: "Item exists and is unchanged",
	}
}

func (e *ItemUnchangedError) Error() string {
	return e.message
}
