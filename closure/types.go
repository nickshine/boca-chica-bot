package closure

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type doc goquery.Document

type Closure struct {
	ClosureType string
	DateString  string
	TimeString  string
	Start       time.Time
	End         time.Time
	Status      string
}

func (c Closure) String() string {
	return fmt.Sprintf("%s: %v\tEnd: %v\tStatus: %s\n", c.ClosureType, c.Start, c.End, c.Status)
}
