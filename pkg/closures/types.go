package closures

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type doc goquery.Document

// Closure represents a beach and/or road closure notice from the Cameron County SpaceX site.
type Closure struct {
	ClosureType  string
	Date         string
	RawTimeRange string
	Time         int64
	TimeType     string
	Status       string
}

func (c Closure) String() string {
	return fmt.Sprintf("%s - %s (%s)", c.Date, c.RawTimeRange, c.TimeType)
}
