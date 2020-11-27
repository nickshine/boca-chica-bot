package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const tablename = "BocaChicaBot-closures"

var svc *dynamodb.DynamoDB

func init() {
	sess := session.Must(session.NewSession())
	svc = dynamodb.New(sess)
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
