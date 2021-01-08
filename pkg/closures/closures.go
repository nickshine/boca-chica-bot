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
	// SiteURL is the website publishing the Boca Chica Beach and Road closures.
	SiteURL = "https://www.cameroncounty.us/spacex/"

	// ClosureTypePrimary represents a primary closure.
	ClosureTypePrimary ClosureType = "Primary Date"

	// ClosureTypeSecondary represents a secondary (backup) closure.
	ClosureTypeSecondary ClosureType = "Secondary Date"

	// ClosureStatusCancelled represents the cancelled Boca Chica Beach status display text.
	ClosureStatusCancelled ClosureStatus = "Closure Cancelled"

	// ClosureStatusScheduled represents the scheduled Boca Chica Beach status display text.
	ClosureStatusScheduled ClosureStatus = "Closure Scheduled"

	// TimeRangeStatusActive notates when the current time is within a closure time range.
	TimeRangeStatusActive TimeRangeStatus = "active"

	// TimeRangeStatusPending notates when the current time is before a closure time range.
	TimeRangeStatusPending TimeRangeStatus = "pending"

	// TimeRangeStatusExpired notates when the current time is after a closure time range.
	TimeRangeStatusExpired TimeRangeStatus = "expired"

	// DateLayout represents the current date layout posted for each Closure.
	DateLayout = "Monday, Jan 2, 2006"
)

const (
	closureLocation = "America/Chicago"
	timeLayout      = "3:04 pm"
)

// Get pulls the current beach/road closures from https://www.cameroncounty.us/spacex/.
func Get() ([]*Closure, error) {
	document, err := scrapeClosuresSite()
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

		closureType := ClosureType(strings.TrimSpace(cells.Get(0).FirstChild.Data))
		dateString := strings.TrimSpace(cells.Get(1).FirstChild.Data)
		date, err := time.Parse(DateLayout, dateString)
		if err != nil {
			return nil, fmt.Errorf("date format changed from 'Monday, Jan 2, 2006' to '%s'", cells.Get(1).FirstChild.Data)
		}

		// reset dateString to formated 'Monday, Jan 2, 2006' for primary key consistency
		dateString = date.Format(DateLayout)
		rawTimeRange := strings.TrimSpace(cells.Get(2).FirstChild.Data)
		startTime, endTime, err := parseTimeRange(rawTimeRange)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time range: %s", err)
		}

		startDate := time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, location)
		endDate := time.Date(date.Year(), date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, location)

		closureStatus := ClosureStatus(strings.TrimSpace(cells.Get(3).FirstChild.Data))

		zone, _ := startDate.Zone()
		rawTimeRangeWithZone := fmt.Sprintf("%s (%s)", rawTimeRange, zone)

		var timeRangeStatus TimeRangeStatus
		// give a little leeway on 'active' status since lambda schedule starts at 8am CST, which is
		// one of the usual closure start times.
		if time.Now().Add(30 * time.Second).Before(startDate) {
			timeRangeStatus = TimeRangeStatusPending
		} else if time.Now().Before(endDate) {
			timeRangeStatus = TimeRangeStatusActive
		} else {
			timeRangeStatus = TimeRangeStatusExpired
		}

		// create a Closure for the start and end of the time range
		closures = append(closures,
			&Closure{
				ClosureType:     closureType,
				Date:            dateString,
				RawTimeRange:    rawTimeRangeWithZone,
				TimeRangeStatus: timeRangeStatus,
				ClosureStatus:   closureStatus,
				Expires:         endDate.Add(2 * time.Hour).Unix(),
			},
		)
	}

	return closures, nil
}

// parseTimeRange parses a string with two times in the format "9:00 am to 9:00 pm".
//
// A start and end Time struct is returned, or an error.
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
func scrapeClosuresSite() (*doc, error) {
	client := &http.Client{Timeout: time.Second * 10}

	res, err := client.Get(SiteURL)
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
