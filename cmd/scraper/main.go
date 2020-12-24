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

	if os.Getenv("TWITTER_ENVIRONMENT") == "test" {
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
		if time.Now().Unix() < c.Expires {
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
			}
		}
	}

	return nil
}
