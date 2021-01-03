package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nickshine/boca-chica-bot/internal/db"
	"github.com/nickshine/boca-chica-bot/pkg/closures"

	"go.uber.org/zap"
)

var (
	log      *zap.SugaredLogger
	dbClient *db.Client
)

var tablename = "BocaChicaBot-closures"

func init() {
	var logger *zap.Logger
	debug := os.Getenv("DEBUG")
	if debug != "false" && debug != "" {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync() // nolint:errcheck
	log = logger.Sugar()

	dbClient = db.NewClient()

	if os.Getenv("AWS_ENVIRONMENT") == "test" {
		tablename = "BocaChicaBot-closures-test"
	}
}

func main() {
	lambda.Start(handler)
}

func handler() error {
	cls, err := closures.Get()
	if err != nil {
		return err
	}

	for _, c := range cls {
		// don't bother putting expired closures in db as the TTL will remove them anyways
		if time.Now().Unix() < c.Time {
			existingClosure, err := dbClient.Put(tablename, c)
			if err != nil {
				switch err.(type) {
				case *db.ItemUnchangedError:
					log.Debugf("%s - Closure: %s", err.Error(), c)
				default:
					return fmt.Errorf("%s - Closure: %s", err.Error(), c)
				}
			} else if existingClosure != nil {
				// if there was an existing closure in db that was overwritten
				log.Debugf("Closure changed: %s", c)
			} else {
				log.Debugf("New closure added: %s", c)
			}
		} else {
			log.Debugf("Closure expired, skipping database call: %s", c)
		}
	}

	return handleRemovals()
}

// handleRemovals explicitly removes closures with Timestamps close to current time.
//
// The Time attribute on each item in the database has a TTL for expiration, but it is only
// guaranteed to be removed within 48 hours, and the 'REMOVE' event is intended to happen at or
// close to the timestamp in order to send the 'start' and 'end' closure notifications at the right
// time. This explicit removal is a workaround around the eventual removal limitation of DynamoDB
// TTLs.
func handleRemovals() error {
	// This will only return results when the function is ran the day of a Closure (partition key)
	// close to the Timestamp of the Closure item (filter expression). Most of the time this will
	// return nil.
	cls, err := dbClient.QueryByTime(tablename, time.Now().Add(2*time.Minute))
	if err != nil {
		return err
	}

	for _, c := range cls {
		log.Debugf("Closure to be removed: %s", c)
		err := dbClient.RemoveClosure(tablename, c)
		if err != nil {
			log.Error("problem removing closure: %v", err)
		}
	}

	return nil
}
