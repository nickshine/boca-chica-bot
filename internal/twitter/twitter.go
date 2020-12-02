package twitter

import (
	"errors"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Client is an alias for twitter.Client with some additional receiver functions for convenience.
type Client twitter.Client

// GetClient returns a twitter client given proper credentials.
func GetClient(c *Credentials) (*Client, error) {
	if c.ConsumerKey == "" || c.ConsumerSecret == "" || c.AccessToken == "" || c.AccessSecret == "" {
		return nil, errors.New("consumer key/secret and access token/secret required")
	}

	config := oauth1.NewConfig(c.ConsumerKey, c.ConsumerSecret)
	token := oauth1.NewToken(c.AccessToken, c.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	return (*Client)(client), nil
}

// Tweet is a convenience function wrapping the twitter client call to send a tweet.
//
// A successful tweet will return the CreatedAt timestamp.
func (c *Client) Tweet(s string) (string, error) {
	tweet, _, err := c.Statuses.Update(s, nil)
	return tweet.CreatedAt, err
}

// Verify returns user info using the provided credentials to validate they work.
func (c *Client) Verify() string {
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	user, _, _ := c.Accounts.VerifyCredentials(verifyParams)
	return fmt.Sprintf("User's Account:\n%+v\n", user)
}
