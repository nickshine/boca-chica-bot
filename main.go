package main

import (
	"fmt"
	"os"

	"github.com/nickshine/boca-chica-bot/closure"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	log = logger.Sugar()

	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		log.Fatal("Consumer Key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	user, _, _ := client.Accounts.VerifyCredentials(verifyParams)
	fmt.Printf("User's Account:\n%+v\n", user)

	cls, err := closure.Get()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, c := range cls {
		fmt.Print(c)
	}
}
