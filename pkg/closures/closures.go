// Package closures retrieves beach and road closures related to SpaceX Starship testing in Boca
// Chica, TX.
package closures

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
	"golang.org/x/net/html"
)

const (
	// SiteURL is the website publishing the Boca Chica Beach and Road closures.
	SiteURL = "https://www.cameroncountytx.gov/spacex/"

	// ClosureTypePrimary represents a primary closure.
	ClosureTypePrimary ClosureType = "Primary Date"

	// ClosureTypeSecondary represents a secondary (backup) closure.
	ClosureTypeSecondary ClosureType = "Secondary Date"

	// ClosureStatusCanceled represents the canceled Boca Chica Beach status display text.
	ClosureStatusCanceled ClosureStatus = "Closure Canceled"

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
	timeLayoutAlt   = "Jan 2 - 3:04 pm"
)

// Get pulls the current beach/road closures from https://www.cameroncountytx.gov/spacex/.
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
			continue
		}

		cellData := cells.Map(func(i int, c *goquery.Selection) string {
			return c.Text()
		})

		rawClosureType := cellData[0]
		rawDateString := cellData[1]
		rawTimeRangeString := cellData[2]
		rawClosureStatus := cellData[3]

		if rawClosureType == "" {
			// row is malformed, skip
			continue
		}
		closureType := ClosureType(strings.TrimSpace(rawClosureType))
		dateString := strings.TrimSpace(rawDateString)
		var date time.Time
		date, err = dateparse.ParseAny(dateString)
		if err != nil {
			// try handling misspelled weekday
			parts := strings.Split(dateString, ",")
			if len(parts) == 3 {
				dateString = strings.TrimSpace(strings.Join(parts[1:], ","))
				date, err = dateparse.ParseAny(dateString)
				if err != nil {
					// skip malformed closure
					continue
				}
			} else {
				// skip malformed closure
				continue
			}
		}

		// reset dateString to formated 'Monday, Jan 2, 2006' for primary key consistency
		dateString = date.Format(DateLayout)
		rawTimeRange := strings.TrimSpace(rawTimeRangeString)
		// try to sanitize/consolidate time range variations a bit
		rawTimeRange = strings.ReplaceAll(rawTimeRange, ".", "")
		rawTimeRange = strings.ToLower(rawTimeRange)

		startTime, endTime, err := parseTimeRange(rawTimeRange)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time range: %s", err)
		}

		startDate := time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, location)
		endDate := time.Date(date.Year(), date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, location)

		closureStatus := ClosureStatus(strings.TrimSpace(rawClosureStatus))

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

	start, err := time.Parse(timeLayout, strings.Trim(strings.TrimSpace(times[0]), "."))
	if err != nil {
		return nil, nil, err
	}

	end, err := time.Parse(timeLayout, strings.Trim(strings.TrimSpace(times[1]), "."))
	if err != nil {
		// try alternate timeLayout
		end, err = time.Parse(timeLayoutAlt, strings.TrimSpace(times[1]))
		// fallback to midnight
		if err != nil {
			end, err = time.Parse(timeLayout, "11:59 pm")
			if err != nil {
				return &start, nil, err
			}
		}
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

// IsCanceled checks if a closure status is 'Status Canceled' or 'Testing concluded' or 'Testing
// ended', forgiving for alternate spellings and capitilization differences.
func IsCanceled(status ClosureStatus) bool {
	s := strings.ToLower(string(status))
	return strings.Contains(s, "cancel") || strings.Contains(s, "conclude") || strings.Contains(s, "ended")
}
