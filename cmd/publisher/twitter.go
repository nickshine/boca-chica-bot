package main

import (
	"fmt"

	"github.com/nickshine/boca-chica-bot/internal/twitter"
)

func handleTweets(params map[string]string, messages []string) error {
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

	for _, m := range messages {
		log.Debugf("Tweet length: %d\n", len(m))
		log.Infof("Tweeting: %s\n", m)
		createdAt, err := client.Tweet(m + "\n#spacex #starship")
		if err != nil {
			log.Info(err)
		}
		log.Debugf("Tweet created at %q", createdAt)

	}
	return nil
}
