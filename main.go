package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nickshine/boca-chica-bot/closure"
	"github.com/nickshine/boca-chica-bot/internal/db"

	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func init() {
	var logger *zap.Logger
	if _, ok := os.LookupEnv("DEBUG"); ok {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()
	log = logger.Sugar()
}

func main() {

	/*
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
	*/

	cls, err := closure.Get()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, c := range cls {
		// don't bother putting expired closures in db as the TTL will remove them anyways
		if time.Now().Unix() < c.Expires {
			err := db.Put(c)
			if err != nil {
				switch err.(type) {
				case *db.ErrItemUnchanged:
					fmt.Println(err.Error())
				default:
					fmt.Println("unknown error")
				}
			}
		}
	}

	// db.Info()

}
