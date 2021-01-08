package main

import (
	"errors"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func handleTweets(params map[string]string, messages []string) error {
	c := &TwitterCredentials{
		ConsumerKey:    params["twitter_consumer_key"],
		ConsumerSecret: params["twitter_consumer_secret"],
		AccessToken:    params["twitter_access_token"],
		AccessSecret:   params["twitter_access_secret"],
	}

	client, err := getTwitterClient(c)
	if err != nil {
		return fmt.Errorf("error getting twitter client: %v", err)
	}

	// log.Debug(verify(client))

	for _, m := range messages {
		log.Debugf("Tweet length: %d\n", len(m))
		log.Infof("Tweeting: %s\n", m)
		createdAt, err := tweet(client, m+"\n#spacex #starship")
		if err != nil {
			log.Info(err)
		}
		log.Debugf("Tweet created at %q", createdAt)

	}
	return nil
}

// TwitterCredentials stores the secrets/tokens/keys necessary for authentication with the Twitter API.
type TwitterCredentials struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

// getTwitterClient returns a twitter client given proper credentials.
func getTwitterClient(c *TwitterCredentials) (*twitter.Client, error) {
	if c.ConsumerKey == "" || c.ConsumerSecret == "" || c.AccessToken == "" || c.AccessSecret == "" {
		return nil, errors.New("consumer key/secret and access token/secret required")
	}

	config := oauth1.NewConfig(c.ConsumerKey, c.ConsumerSecret)
	token := oauth1.NewToken(c.AccessToken, c.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	return client, nil
}

// Tweet is a convenience function wrapping the twitter client call to send a tweet.
//
// A successful tweet will return the CreatedAt timestamp.
func tweet(client *twitter.Client, s string) (string, error) {
	tweet, _, err := client.Statuses.Update(s, nil)
	return tweet.CreatedAt, err
}

/*
// verify returns user info using the provided credentials to validate they work.
func verify(client *twitter.Client) string {
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	user, _, _ := client.Accounts.VerifyCredentials(verifyParams)
	return fmt.Sprintf("User's Account:\n%+v\n", user)
}
*/
