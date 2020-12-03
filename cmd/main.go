package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nickshine/boca-chica-bot/internal/db"
	"github.com/nickshine/boca-chica-bot/internal/param"
	"github.com/nickshine/boca-chica-bot/internal/twitter"
	"github.com/nickshine/boca-chica-bot/pkg/closures"

	"go.uber.org/zap"
)

var log *zap.SugaredLogger

var (
	twitterParamsPath = "/boca-chica-bot/prod/"
	tablename         = "BocaChicaBot-closures"
)

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
		tablename = "BocaChicaBot-closures-test"
	}
}

func main() {
	lambda.Start(handler)
}

func handler() {
	cls, err := closures.Get()
	if err != nil {
		log.Fatal(err.Error())
	}

	var tweets []string

	for _, c := range cls {
		// don't bother putting expired closures in db as the TTL will remove them anyways
		if time.Now().Unix() < c.Expires {
			var tweet string
			existingClosure, err := db.Put(tablename, c)
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
		}
	}

	handleTweets(tweets)
}

func handleTweets(tweets []string) {
	if len(tweets) == 0 {
		return
	}
	if disable := os.Getenv("DISABLE_TWEETS"); disable != "" && disable != "false" {
		log.Debugf("DISABLE_TWEETS env var enabled, skipping publishing of tweets: %v", tweets)
		return
	}

	pClient := param.GetClient()
	params, err := pClient.GetParams(twitterParamsPath)
	if err != nil {
		log.Fatalf("error retrieving Twitter API creds from parameter store: %v", err)
	}

	c := &twitter.Credentials{
		ConsumerKey:    params["twitter_consumer_key"],
		ConsumerSecret: params["twitter_consumer_secret"],
		AccessToken:    params["twitter_access_token"],
		AccessSecret:   params["twitter_access_secret"],
	}

	// c := &twitter.Credentials{
	// 	ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
	// 	ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
	// 	AccessToken:    os.Getenv("TWITTER_ACCESS_TOKEN"),
	// 	AccessSecret:   os.Getenv("TWITTER_ACCESS_SECRET"),
	// }

	client, err := twitter.GetClient(c)
	if err != nil {
		log.Errorf("error getting twitter client: %v", err)
	}

	// log.Debug(client.Verify())

	for _, t := range tweets {
		log.Debugf("Tweet length: %d\n", len(t))
		log.Infof("Tweeting: %s\n", t)
		createdAt, err := client.Tweet(t)
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("Tweet created at %s", createdAt)
	}
}
