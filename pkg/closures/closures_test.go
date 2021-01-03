package closures

import (
	"testing"
	"time"

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
		{"invalid range 2", "", nil, nil, assert.Error},
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
