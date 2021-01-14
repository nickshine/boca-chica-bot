package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nickshine/boca-chica-bot/internal/param"
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
	var messages []string
	for _, record := range e.Records {
		log.Debugf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

		var image map[string]events.DynamoDBAttributeValue

		if len(record.Change.NewImage) > 0 {
			image = record.Change.NewImage
		} else if len(record.Change.OldImage) > 0 {
			image = record.Change.OldImage
		} else {
			return fmt.Errorf("Invalid DynamoDBEvent: %v", e)
		}

		closureType := image["ClosureType"].String()
		date := image["Date"].String()
		rawTimeRange := image["RawTimeRange"].String()
		timeRangeStatus := closures.TimeRangeStatus(image["TimeRangeStatus"].String())
		closureStatus := closures.ClosureStatus(image["ClosureStatus"].String())

		switch record.EventName {
		// An INSERT event means a new Closure has been added.
		case string(events.DynamoDBOperationTypeInsert):
			if closureStatus == closures.ClosureStatusCanceled {
				log.Debugf("Closure Status of '%s' on 'INSERT', skipping publish", closureStatus)
				return nil
			}

			messages = append(messages, fmt.Sprintf("New closure scheduled:\n%s - %s - %s\n%s",
				closureType, date, rawTimeRange, closures.SiteURL))
		// A MODIFY event means an existing Closure has been changed.
		case string(events.DynamoDBOperationTypeModify):

			oldRawTimeRange := record.Change.OldImage["RawTimeRange"].String()
			oldTimeRangeStatus := closures.TimeRangeStatus(record.Change.OldImage["TimeRangeStatus"].String())
			oldClosureStatus := closures.ClosureStatus(record.Change.OldImage["ClosureStatus"].String())

			if closureStatus == oldClosureStatus && closureStatus == closures.ClosureStatusCanceled {
				log.Debugf("Closure is cancelled, skipping publish")
				return nil
			} else if timeRangeStatus != oldTimeRangeStatus && closureStatus != closures.ClosureStatusCanceled {
				switch timeRangeStatus {
				case closures.TimeRangeStatusActive:
					messages = append(messages, fmt.Sprintf("Closure for %s - %s has started.\n%s",
						date, rawTimeRange, closures.SiteURL))
				case closures.TimeRangeStatusExpired:
					messages = append(messages, fmt.Sprintf("Closure for %s - %s has ended.\n%s",
						date, rawTimeRange, closures.SiteURL))
				}
			} else if rawTimeRange != oldRawTimeRange && closureStatus != closures.ClosureStatusCanceled {
				messages = append(messages, fmt.Sprintf("Time window for the %s - %s closure has changed to %s.\n%s",
					date, oldRawTimeRange, rawTimeRange, closures.SiteURL))
			} else if closureStatus != oldClosureStatus {
				messages = append(messages, fmt.Sprintf("Status change for the %s closure: %s.\n%s",
					date, closureStatus, closures.SiteURL))
			} else {
				messages = append(messages, fmt.Sprintf("Closure status change:\n%s - %s - %s\n%s",
					date, rawTimeRange, closureStatus, closures.SiteURL))
			}
		}
	}

	if len(messages) == 0 {
		return nil
	} else if disable := os.Getenv("DISABLE_PUBLISH"); disable != "" && disable != "false" {
		log.Debugf("DISABLE_PUBLISH env var enabled, skipping publishing: %v", messages)
		return nil
	}

	pClient := param.NewClient()
	params, err := pClient.GetParams(paramsPath)
	if err != nil {
		return fmt.Errorf("error retrieving Twitter/Discord API creds from parameter store: %v", err)
	}

	err = handleTweets(params, messages)
	if err != nil {
		log.Info(err)
	}
	err = handleDiscord(params, messages)
	if err != nil {
		log.Info(err)
	}
	return nil
}
