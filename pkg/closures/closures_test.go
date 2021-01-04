package closures

import (
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func newTime(t string) *time.Time {
	o, _ := time.Parse(timeLayout, t)
	return &o
}

func TestParseTimeRange(t *testing.T) {
	tests := []struct {
		name          string
		timeRange     string
		expectedStart *time.Time
		expectedEnd   *time.Time
		errAssertion  assert.ErrorAssertionFunc
	}{
		{"valid range", "9:00 am to 9:00 pm", newTime("9:00 am"), newTime("9:00 pm"), assert.NoError},
		{"valid range 2", "8:00 am to 5:00 pm", newTime("8:00 am"), newTime("5:00 pm"), assert.NoError},
		{"invalid range", "8:00 am to to 5:00 pm", nil, nil, assert.Error},
		{"invalid range 2", "faketime to faketime", nil, nil, assert.Error},
		{"invalid range 3", "8:00 am to faketime", nil, nil, assert.Error},
		{"invalid range 3", "", nil, nil, assert.Error},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualStart, actualEnd, actualErr := parseTimeRange(tt.timeRange)

			tt.errAssertion(t, actualErr)
			assert.Equal(t, tt.expectedStart, actualStart, "start times should match")
			assert.Equal(t, tt.expectedEnd, actualEnd, "end times should match")
		})
	}
}

func TestScrape(t *testing.T) {
	mockSite, err := ioutil.ReadFile("./testdata/closures-table.html")
	if err != nil {
		t.Fatal(err)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name         string
		responder    httpmock.Responder
		assertion    assert.ValueAssertionFunc
		errAssertion assert.ErrorAssertionFunc
	}{
		{"200 response", httpmock.NewBytesResponder(200, mockSite), assert.NotNil, assert.NoError},
		{"GET request err", httpmock.NewErrorResponder(errors.New("mock http client error")), assert.Nil, assert.Error},
		{"404 response", httpmock.NewBytesResponder(404, nil), assert.Nil, assert.Error},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.RegisterResponder("GET", SiteURL, tt.responder)
			document, err := scrapeClosuresSite()
			tt.assertion(t, document)
			tt.errAssertion(t, err)
		})

		httpmock.Reset()
	}
}
