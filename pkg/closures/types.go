package closures

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type doc goquery.Document

// Closure represents a beach and/or road closure notice from the Cameron County SpaceX site.
type Closure struct {
	ClosureType string
	Date        string
	Time        string
	Start       time.Time
	End         time.Time
	Status      string
	Expires     int64
}

func (c Closure) String() string {
	return fmt.Sprintf("%s - %s", c.Date, c.Time)
}
