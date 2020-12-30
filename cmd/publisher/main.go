package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nickshine/boca-chica-bot/internal/discord"
	"github.com/nickshine/boca-chica-bot/internal/param"
	"github.com/nickshine/boca-chica-bot/internal/twitter"
	"github.com/nickshine/boca-chica-bot/pkg/closures"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

var paramsPath = "/boca-chica-bot/prod/"

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

	if os.Getenv("AWS_ENVIRONMENT") == "test" {
		paramsPath = "/boca-chica-bot/test/"
	}
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, e events.DynamoDBEvent) error {
	log.Debugf("Event: %+v\n", e)
	var tweets []string
	for _, record := range e.Records {
		fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

		var image map[string]events.DynamoDBAttributeValue

		if record.Change.NewImage != nil {
			image = record.Change.NewImage
		} else if record.Change.OldImage != nil {
			image = record.Change.OldImage
		} else {
			return fmt.Errorf("Invalid DynamoDBEvent: %v", e)
		}

		closureType := image["ClosureType"].String()
		date := image["Date"].String()
		rawTimeRange := image["RawTimeRange"].String()
		timeType := image["TimeType"].String()
		status := image["Status"].String()

		switch record.EventName {
		// An INSERT event means a new Closure has been added
		case string(events.DynamoDBOperationTypeInsert):
			// Each closure has two entries, a 'start' type and 'end' type.  Only tweet a new closure once (on 'start' type).
			if timeType == closures.TimeTypeEnd {
				log.Debugf("Closure TimeType of '%s' on 'INSERT', skipping publish", timeType)
				return nil
			} else if status == closures.CancelledStatus {
				log.Debugf("Closure Status of '%s' on 'INSERT', skipping publish", status)
				return nil
			}

			tweets = append(tweets, fmt.Sprintf("New closure scheduled:\n%s - %s - %s\n%s\n#spacex #starship",
				closureType, date, rawTimeRange, closures.SiteURL))
		case string(events.DynamoDBOperationTypeModify):
			// Each closure has two entries, a 'start' type and 'end' type.  Only tweet a closure update once (on 'start' type).
			if timeType == closures.TimeTypeEnd {
				log.Debugf("Closure TimeType of '%s' on 'MODIFY', skipping publish", timeType)
				return nil
			}
			if status == closures.CancelledStatus {
				tweets = append(tweets, fmt.Sprintf("Closure for %s - %s has been cancelled.\n%s\n#spacex #starship",
					date, rawTimeRange, closures.SiteURL))
			} else {
				tweets = append(tweets, fmt.Sprintf("Closure status change:\n%s - %s - %s\n%s\n#spacex #starship",
					date, rawTimeRange, status, closures.SiteURL))
			}
		// A REMOVE event means the closure has expired (time range started or ended)
		case string(events.DynamoDBOperationTypeRemove):
			if status != closures.ScheduledStatus {
				log.Debugf("Closure Status of '%s' on 'REMOVE', skipping publish", status)
				return nil
			}
			if timeType == closures.TimeTypeStart {
				tweets = append(tweets, fmt.Sprintf("Closure for %s - %s has started.\n%s\n#spacex #starship",
					date, rawTimeRange, closures.SiteURL))
			} else if timeType == closures.TimeTypeEnd {
				tweets = append(tweets, fmt.Sprintf("Closure for %s - %s has ended.\n%s\n#spacex #starship",
					date, rawTimeRange, closures.SiteURL))
			}
		}
	}

	if len(tweets) != 0 {
		return nil
	} else if disable := os.Getenv("DISABLE_TWEETS"); disable != "" && disable != "false" {
		log.Debugf("DISABLE_TWEETS env var enabled, skipping publishing of tweets: %v", tweets)
		return nil
	}

	pClient := param.NewClient()
	params, err := pClient.GetParams(paramsPath)
	if err != nil {
		return fmt.Errorf("error retrieving Twitter/Discord API creds from parameter store: %v", err)
	}

	err = handleTweets(params, tweets)
	if err != nil {
		log.Error(err)
	}
	err = handleDiscord(params, tweets)
	return err
}

func handleDiscord(params map[string]string, messages []string) error {
	c := &discord.Credentials{
		Token: params["discord_bot_token"],
	}

	discord, err := discord.GetSession(c)
	if err != nil {
		return fmt.Errorf("error getting twitter client: %v", err)
	}

	errors := discord.Send(messages)
	for _, e := range errors {
		log.Error(e)
	}

	return nil
}

func handleTweets(params map[string]string, tweets []string) error {
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

	for _, t := range tweets {
		log.Debugf("Tweet length: %d\n", len(t))
		log.Infof("Tweeting: %s\n", t)
		createdAt, err := client.Tweet(t)
		if err != nil {
			return err
		}
		log.Debugf("Tweet created at %s", createdAt)

	}
	return nil
}
