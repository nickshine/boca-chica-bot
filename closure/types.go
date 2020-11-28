package closure

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type doc goquery.Document

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
	return fmt.Sprintf("%s - %v to %v - %s\n", c.ClosureType, c.Start, c.End, c.Status)
}
