package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

var twitterParamsPath = "/boca-chica-bot/prod/"

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

	if os.Getenv("TWITTER_ENVIRONMENT") == "test" {
		twitterParamsPath = "/boca-chica-bot/test/"
	}
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, e events.DynamoDBEvent) error {
	log.Debugf("Event: %+v\n", e)
	for _, record := range e.Records {
		fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

		switch record.EventName {
		case string(events.DynamoDBOperationTypeInsert):
			fmt.Println("INSERT")
		case string(events.DynamoDBOperationTypeModify):
			fmt.Println("MODIFY")
		case string(events.DynamoDBOperationTypeRemove):

			for name, value := range record.Change.OldImage {
				if value.DataType() == events.DataTypeNumber {
					fmt.Printf("Attribute name: %s, value: %s\n", name, value.Number())
				} else {
					fmt.Printf("Attribute name: %s, value: %s\n", name, value.String())
				}
			}
		}
	}

	// TODO: create tweet string from event INSERT, MODIFY, or REMOVE

	/*
		if err != nil {
			switch err.(type) {
			case *db.ItemUnchangedError:
				log.Debugf("%s - Closure: %s", err.Error(), c)
			default:
				log.Errorf("%s - Closure: %s", err.Error(), c)
			}
		} else if existingClosure != nil {
			// if there was an existing closure in db and an attribute changed (e.g. status
			// changed from "Scheduled" to "Cancelled")
			if c.Status == closures.CancelledStatus {
				tweet = fmt.Sprintf("Beach closure for %s - has been cancelled.\n%s\n#spacex #starship", c, closures.SiteURL)
			} else {
				tweet = fmt.Sprintf("Beach closure status change:\n%s - %s\n%s\n#spacex #starship", c, c.Status, closures.SiteURL)
			}
			tweets = append(tweets, tweet)
		} else {
			// existingClosure is nil (meaning new addition to db)
			tweet = fmt.Sprintf("New beach closure scheduled:\n%s - %s\n%s\n#spacex #starship",
				c.ClosureType, c, closures.SiteURL)
			tweets = append(tweets, tweet)
		}
	*/

	// tweet := ""
	// return handleTweet(tweet)
	return nil
}

/*
func handleTweet(tweet string) error {
	if disable := os.Getenv("DISABLE_TWEETS"); disable != "" && disable != "false" {
		log.Debugf("DISABLE_TWEETS env var enabled, skipping publishing of tweet: %v", tweet)
		return nil
	}

	pClient := param.NewClient()
	params, err := pClient.GetParams(twitterParamsPath)
	if err != nil {
		return fmt.Errorf("error retrieving Twitter API creds from parameter store: %v", err)
	}

	c := &twitter.Credentials{
		ConsumerKey:    params["twitter_consumer_key"],
		ConsumerSecret: params["twitter_consumer_secret"],
		AccessToken:    params["twitter_access_token"],
		AccessSecret:   params["twitter_access_secret"],
	}

	client, err := twitter.GetClient(c)
	if err != nil {
		return fmt.Errorf("error getting twitter client: %v", err)
	}

	// log.Debug(client.Verify())

	log.Debugf("Tweet length: %d\n", len(tweet))
	log.Infof("Tweeting: %s\n", tweet)
	createdAt, err := client.Tweet(tweet)
	if err != nil {
		return err
	}
	log.Debugf("Tweet created at %s", createdAt)

	return nil
}
*/
