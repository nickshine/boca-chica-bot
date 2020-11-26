package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"go.uber.org/zap"
)

const (
	closureSiteURL    = "https://www.cameroncounty.us/spacex/"
	closureDateLayout = "Jan 2, 2006"
	closureLocation   = "America/Chicago"
	closureTimeLayout = "3:04 pm"
)

var log *zap.SugaredLogger

type doc goquery.Document

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

	document, err := scrape(closureSiteURL)
	if err != nil {
		log.Fatalf("Failed to scrape Cameron County SpaceX page", err)
	}

	closures, err := document.getClosures()
	if err != nil {
		log.Fatalf("Failed to get closures", err)
	}

	for _, c := range closures {
		fmt.Printf("%+v\n", c)

	}
}

type closure struct {
	closureType string
	dateString  string
	timeString  string
	start       time.Time
	end         time.Time
	status      string
}

func (c closure) String() string {
	return fmt.Sprintf("%s: %v\tEnd: %v\tStatus: %s\n", c.closureType, c.start, c.end, c.status)
}

func (d *doc) getClosures() ([]*closure, error) {
	var closures []*closure

	location, err := time.LoadLocation(closureLocation)
	if err != nil {
		return nil, err
	}

	d.Find("table tbody > tr").Each(func(i int, s *goquery.Selection) {

		// html, _ := s.Html()
		// fmt.Printf("row:\n%s\n", html)

		cells := s.Find("td")

		if cells.Length() != 4 {
			log.Fatalf("table format changed: row does not have 4 cells: cell count: %d", cells.Length())
		}

		dateString := cells.Get(1).FirstChild.Data
		date, err := time.Parse(closureDateLayout, dateString)
		if err != nil {
			log.Fatalf("date format changed from 'Jan 2, 2006' to %s", cells.Get(1).FirstChild.Data)
		}

		timeString := cells.Get(2).FirstChild.Data
		startTime, endTime, err := parseTimeRange(timeString)
		if err != nil {
			log.Fatalf("failed to parse time range: %s", err)
		}

		startDate := time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, location)
		endDate := time.Date(date.Year(), date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, location)

		// fmt.Printf("start date: %s\n", startDate)
		// fmt.Printf("end date: %s\n", endDate)

		closures = append(closures,
			&closure{
				closureType: cells.Get(0).FirstChild.Data,
				dateString:  dateString,
				timeString:  timeString,
				start:       startDate,
				end:         endDate,
				status:      cells.Get(3).FirstChild.Data,
			},
		)
	})

	return closures, nil
}

// timeRange format: "9:00 am to 9:00 pm"
func parseTimeRange(timeRange string) (*time.Time, *time.Time, error) {

	times := strings.Split(timeRange, "to")

	if len(times) != 2 {
		return nil, nil, fmt.Errorf("Date range format has changed from '9:00am to 9:00pm' to %s", timeRange)
	}

	start, err := time.Parse(closureTimeLayout, strings.TrimSpace(times[0]))
	if err != nil {
		return nil, nil, err
	}
	end, err := time.Parse(closureTimeLayout, strings.TrimSpace(times[1]))
	if err != nil {
		return nil, nil, err
	}

	return &start, &end, nil

}

// Scrape Cameron County SpaceX page for road closures
func scrape(url string) (*doc, error) {

	client := &http.Client{Timeout: time.Second * 10}

	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return (*doc)(document), nil
}
