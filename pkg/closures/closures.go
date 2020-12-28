// Package closures retrieves beach and road closures related to SpaceX Starship testing in Boca
// Chica, TX.
package closures

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const (
	// CancelledStatus represents the cancelled Boca Chica Beach status display text.
	CancelledStatus = "Closure Cancelled"

	// ScheduledStatus represents the scheduled Boca Chica Beach status display text.
	ScheduledStatus = "Closure Scheduled"

	// TimeTypeStart represents a Closure's beginning time in a posted time range.
	TimeTypeStart = "start"

	// TimeTypeEnd represents a Closure's ending time in a posted time range.
	TimeTypeEnd = "end"

	// SiteURL is the website publishing the Boca Chica Beach and Road closures.
	SiteURL = "https://www.cameroncounty.us/spacex/"
)

const (
	dateLayout      = "Monday, Jan 2, 2006"
	closureLocation = "America/Chicago"
	timeLayout      = "3:04 pm"
)

// Get pulls the current beach/road closures from https://www.cameroncounty.us/spacex/.
func Get() ([]*Closure, error) {
	document, err := scrape(SiteURL)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape Cameron County SpaceX page: %w", err)
	}

	closures, err := document.getClosures()
	if err != nil {
		return nil, fmt.Errorf("failed to get closures: %w", err)
	}

	return closures, nil
}

func (d *doc) getClosures() ([]*Closure, error) {
	var closures []*Closure

	location, err := time.LoadLocation(closureLocation)
	if err != nil {
		return nil, err
	}

	rows := d.Find("table tbody > tr") // no .Each function callback in order to return errors

	for _, row := range rows.Nodes {

		sel := &goquery.Selection{
			Nodes: []*html.Node{row},
		}

		cells := sel.Find("td")
		if cells.Length() != 4 {
			return nil, fmt.Errorf("table format changed: row does not have 4 cells: cell count: %d", cells.Length())
		}

		closureType := strings.TrimSpace(cells.Get(0).FirstChild.Data)
		dateString := strings.TrimSpace(cells.Get(1).FirstChild.Data)
		date, err := time.Parse(dateLayout, dateString)
		if err != nil {
			return nil, fmt.Errorf("date format changed from 'Monday, Jan 2, 2006' to '%s'", cells.Get(1).FirstChild.Data)
		}

		// reset dateString to formated 'Jan 2, 2006' for primary key consistency
		dateString = date.Format(dateLayout)
		rawTimeRange := strings.TrimSpace(cells.Get(2).FirstChild.Data)
		startTime, endTime, err := parseTimeRange(rawTimeRange)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time range: %s", err)
		}

		startDate := time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, location)
		endDate := time.Date(date.Year(), date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, location)

		status := strings.TrimSpace(cells.Get(3).FirstChild.Data)

		zone, _ := startDate.Zone()
		rawTimeRangeWithZone := fmt.Sprintf("%s (%s)", rawTimeRange, zone)

		// create a Closure for the start and end of the time range
		closures = append(closures,
			&Closure{
				ClosureType:  closureType,
				Date:         dateString,
				RawTimeRange: rawTimeRangeWithZone,
				Time:         startDate.Unix(),
				TimeType:     TimeTypeStart,
				Status:       status,
			},
			&Closure{
				ClosureType:  closureType,
				Date:         dateString,
				RawTimeRange: rawTimeRangeWithZone,
				Time:         endDate.Unix(),
				TimeType:     TimeTypeEnd,
				Status:       status,
			},
		)
	}

	return closures, nil
}

// timeRange format: "9:00 am to 9:00 pm"
func parseTimeRange(timeRange string) (*time.Time, *time.Time, error) {
	times := strings.Split(timeRange, "to")

	if len(times) != 2 {
		return nil, nil, fmt.Errorf("date range format has changed from '9:00am to 9:00pm' to %s", timeRange)
	}

	start, err := time.Parse(timeLayout, strings.TrimSpace(times[0]))
	if err != nil {
		return nil, nil, err
	}
	end, err := time.Parse(timeLayout, strings.TrimSpace(times[1]))
	if err != nil {
		return nil, nil, err
	}

	return &start, &end, nil
}

// Scrape Cameron County SpaceX page for road Closures
func scrape(url string) (*doc, error) {
	client := &http.Client{Timeout: time.Second * 10}

	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close() //nolint:errcheck
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
