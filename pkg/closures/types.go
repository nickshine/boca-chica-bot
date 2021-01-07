package closures

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type doc goquery.Document

// TimeRangeStatus notates the status of a closure time range relative to current time: pending,
// active, or expired.
type TimeRangeStatus string

// ClosureStatus represents the official status of the closure as reported on https://www.cameroncounty.us/spacex/.
type ClosureStatus string

// ClosureType represents the type of Closure reported.
type ClosureType string

// Closure represents a beach and/or road closure notice from the Cameron County SpaceX site.
type Closure struct {
	ClosureType     ClosureType
	Date            string
	RawTimeRange    string
	TimeRangeStatus TimeRangeStatus
	ClosureStatus   ClosureStatus
	Expires         int64
}

func (c Closure) String() string {
	return fmt.Sprintf("%s - %s", c.Date, c.RawTimeRange)
}
