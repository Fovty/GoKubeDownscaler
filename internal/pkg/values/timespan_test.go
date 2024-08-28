package values

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var zeroTime = time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)

func TestParseRelativeTimeSpan(t *testing.T) {
	tests := []struct {
		name           string
		timespanString string
		wantResult     *relativeTimeSpan
		wantErr        bool
	}{
		{
			name:           "valid",
			timespanString: "Mon-Fri 07:00-16:00 UTC",
			wantResult: &relativeTimeSpan{
				timezone:    time.UTC,
				weekdayFrom: time.Monday,
				weekdayTo:   time.Friday,
				timeFrom:    zeroTime.Add(7 * time.Hour),
				timeTo:      zeroTime.Add(16 * time.Hour),
			},
			wantErr: false,
		},
		{
			name:           "reverse",
			timespanString: "Sat-Sun 20:00-06:00 UTC",
			wantResult: &relativeTimeSpan{
				timezone:    time.UTC,
				weekdayFrom: time.Saturday,
				weekdayTo:   time.Sunday,
				timeFrom:    zeroTime.Add(20 * time.Hour),
				timeTo:      zeroTime.Add(6 * time.Hour),
			},
			wantErr: false,
		},
		{
			name:           "invalid TZ",
			timespanString: "Mon-Fri 07:00-16:00 Invalid",
			wantResult:     nil,
			wantErr:        true,
		},
		{
			name:           "invalid Time",
			timespanString: "Mon-Fri 03:00-04:0 UTC",
			wantResult:     nil,
			wantErr:        true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotResult, gotErr := parseRelativeTimeSpan(test.timespanString)
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
			assert.Equal(t, test.wantResult, gotResult)
		})
	}
}

func TestRelativeTimeSpan_isWeekdayInRange(t *testing.T) {
	tests := []struct {
		name       string
		timespan   relativeTimeSpan
		weekday    time.Weekday
		wantResult bool
	}{
		{
			name:       "in range",
			timespan:   relativeTimeSpan{weekdayFrom: time.Monday, weekdayTo: time.Friday},
			weekday:    time.Wednesday,
			wantResult: true,
		},
		{
			name:       "from in range",
			timespan:   relativeTimeSpan{weekdayFrom: time.Monday, weekdayTo: time.Friday},
			weekday:    time.Monday,
			wantResult: true,
		},
		{
			name:       "to in range",
			timespan:   relativeTimeSpan{weekdayFrom: time.Monday, weekdayTo: time.Friday},
			weekday:    time.Friday,
			wantResult: true,
		},
		{
			name:       "reverse in range",
			timespan:   relativeTimeSpan{weekdayFrom: time.Saturday, weekdayTo: time.Sunday},
			weekday:    time.Saturday,
			wantResult: true,
		},
		{
			name:       "reverse out of range",
			timespan:   relativeTimeSpan{weekdayFrom: time.Saturday, weekdayTo: time.Sunday},
			weekday:    time.Monday,
			wantResult: false,
		},
		{
			name:       "out of range",
			timespan:   relativeTimeSpan{weekdayFrom: time.Monday, weekdayTo: time.Friday},
			weekday:    time.Saturday,
			wantResult: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotResult := test.timespan.isWeekdayInRange(test.weekday)
			assert.Equal(t, test.wantResult, gotResult)
		})
	}
}

func TestRelativeTimeSpan_isTimeOfDayInRange(t *testing.T) {
	tests := []struct {
		name       string
		timespan   relativeTimeSpan
		timeOfDay  time.Time
		wantResult bool
	}{
		{
			name:       "in range",
			timespan:   relativeTimeSpan{timeFrom: zeroTime.Add(6 * time.Hour), timeTo: zeroTime.Add(20 * time.Hour)},
			timeOfDay:  zeroTime.Add(16 * time.Hour),
			wantResult: true,
		},
		{
			name:       "to out of range",
			timespan:   relativeTimeSpan{timeFrom: zeroTime.Add(6 * time.Hour), timeTo: zeroTime.Add(20 * time.Hour)},
			timeOfDay:  zeroTime.Add(20 * time.Hour),
			wantResult: false,
		},
		{
			name:       "reverse in range",
			timespan:   relativeTimeSpan{timeFrom: zeroTime.Add(18 * time.Hour), timeTo: zeroTime.Add(4 * time.Hour)},
			timeOfDay:  zeroTime.Add(3 * time.Hour),
			wantResult: true,
		},
		{
			name:       "reverse to out of range",
			timespan:   relativeTimeSpan{timeFrom: zeroTime.Add(18 * time.Hour), timeTo: zeroTime.Add(4 * time.Hour)},
			timeOfDay:  zeroTime.Add(4 * time.Hour),
			wantResult: false,
		},
		{
			name:       "from in range",
			timespan:   relativeTimeSpan{timeFrom: zeroTime.Add(6 * time.Hour), timeTo: zeroTime.Add(20 * time.Hour)},
			timeOfDay:  zeroTime.Add(6 * time.Hour),
			wantResult: true,
		},
		{
			name:       "reverse from in range",
			timespan:   relativeTimeSpan{timeFrom: zeroTime.Add(18 * time.Hour), timeTo: zeroTime.Add(4 * time.Hour)},
			timeOfDay:  zeroTime.Add(18 * time.Hour),
			wantResult: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotResult := test.timespan.isTimeOfDayInRange(test.timeOfDay)
			assert.Equal(t, test.wantResult, gotResult)
		})
	}
}

func TestGetTimeOfDay(t *testing.T) {
	tests := []struct {
		name       string
		time       time.Time
		wantResult time.Time
	}{
		{
			name:       "utc",
			time:       time.Date(2024, time.April, 12, 10, 20, 0, 0, time.UTC),
			wantResult: time.Date(0, time.January, 1, 10, 20, 0, 0, time.UTC),
		},
		{
			name:       "not utc",
			time:       time.Date(2024, time.April, 12, 10, 20, 0, 0, time.FixedZone("UTC+2", 2*int(time.Hour/time.Second))),
			wantResult: time.Date(0, time.January, 1, 10, 20, 0, 0, time.FixedZone("UTC+2", 2*int(time.Hour/time.Second))),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotResult := getTimeOfDay(test.time)
			assert.Equal(t, test.wantResult, gotResult)
		})
	}
}

func TestAbsoluteTimeSpan_isTimeInSpan(t *testing.T) {
	tests := []struct {
		name       string
		timespan   absoluteTimeSpan
		time       time.Time
		wantResult bool
	}{
		{
			name: "in range",
			timespan: absoluteTimeSpan{
				from: time.Date(2024, time.June, 3, 20, 0, 0, 0, time.UTC),
				to:   time.Date(2024, time.June, 10, 6, 0, 0, 0, time.UTC),
			},
			time:       time.Date(2024, time.June, 9, 12, 34, 2, 152, time.UTC),
			wantResult: true,
		},
		{
			name: "out of range",
			timespan: absoluteTimeSpan{
				from: time.Date(2024, time.November, 1, 22, 0, 0, 0, time.UTC),
				to:   time.Date(2024, time.November, 22, 5, 0, 0, 0, time.UTC),
			},
			time:       time.Date(2024, time.December, 5, 2, 30, 0, 0, time.UTC),
			wantResult: false,
		},
		{
			name: "from in range",
			timespan: absoluteTimeSpan{
				from: time.Date(2024, time.November, 1, 22, 0, 0, 0, time.UTC),
				to:   time.Date(2024, time.November, 22, 5, 0, 0, 0, time.UTC),
			},
			time:       time.Date(2024, time.November, 1, 22, 0, 0, 0, time.UTC),
			wantResult: true,
		},
		{
			name: "to out of range",
			timespan: absoluteTimeSpan{
				from: time.Date(2024, time.November, 1, 22, 0, 0, 0, time.UTC),
				to:   time.Date(2024, time.November, 22, 5, 0, 0, 0, time.UTC),
			},
			time:       time.Date(2024, time.November, 22, 5, 0, 0, 0, time.UTC),
			wantResult: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotResult := test.timespan.isTimeInSpan(test.time)
			assert.Equal(t, test.wantResult, gotResult)
		})
	}
}

func TestOverlappingTimespans(t *testing.T) {
	tests := []struct {
		name         string
		span1        TimeSpan
		span2        TimeSpan
		wantedResult bool
	}{
		{
			name: "rel rel overlap",
			span1: relativeTimeSpan{
				timezone:    time.UTC,
				weekdayFrom: time.Monday,
				weekdayTo:   time.Friday,
				timeFrom:    zeroTime.Add(8 * time.Hour),
				timeTo:      zeroTime.Add(20 * time.Hour),
			},
			span2: relativeTimeSpan{
				timezone:    time.UTC,
				weekdayFrom: time.Monday,
				weekdayTo:   time.Monday,
				timeFrom:    zeroTime.Add(12 * time.Hour),
				timeTo:      zeroTime.Add(18 * time.Hour),
			},
			wantedResult: true,
		},
		{
			name: "rel rel dont overlap",
			span1: relativeTimeSpan{
				timezone:    time.UTC,
				weekdayFrom: time.Monday,
				weekdayTo:   time.Friday,
				timeFrom:    zeroTime.Add(8 * time.Hour),
				timeTo:      zeroTime.Add(20 * time.Hour),
			},
			span2: relativeTimeSpan{
				timezone:    time.UTC,
				weekdayFrom: time.Tuesday,
				weekdayTo:   time.Tuesday,
				timeFrom:    zeroTime.Add(22 * time.Hour),
				timeTo:      zeroTime.Add(24 * time.Hour),
			},
			wantedResult: false,
		},
		{
			name: "abs abs overlap",
			span1: absoluteTimeSpan{ // all of January 2024
				from: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),  // from 1st of January
				to:   time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC), // to 1st of February
			},
			span2: absoluteTimeSpan{ // from the 10th of Janurary 2024 until the end of the 19th of Janurary 2024
				from: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC), // from 10th of January
				to:   time.Date(2024, time.January, 20, 0, 0, 0, 0, time.UTC), // to 20th of January
			},
			wantedResult: true,
		},
		{
			name: "abs abs dont overlap",
			span1: absoluteTimeSpan{ // all of January 2024
				from: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),  // from 1st of January
				to:   time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC), // to 1st of February
			},
			span2: absoluteTimeSpan{ // all of February 2024
				from: time.Date(2024, time.February, 10, 0, 0, 0, 0, time.UTC), // from 1st of February
				to:   time.Date(2024, time.March, 20, 0, 0, 0, 0, time.UTC),    // to 1st of March
			},
			wantedResult: false,
		},
		{
			name: "rel abs overlap",
			span1: relativeTimeSpan{
				timezone:    time.UTC,
				weekdayFrom: time.Monday,
				weekdayTo:   time.Friday,
				timeFrom:    zeroTime.Add(8 * time.Hour),
				timeTo:      zeroTime.Add(20 * time.Hour),
			},
			span2: absoluteTimeSpan{ // all of January 2024
				from: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),  // from Monday 1st of January
				to:   time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC), // to Thursday 1st of Feburary
			},
			wantedResult: true,
		},
		{
			name: "rel abs dont overlap",
			span1: relativeTimeSpan{
				timezone:    time.UTC,
				weekdayFrom: time.Monday,
				weekdayTo:   time.Friday,
				timeFrom:    zeroTime.Add(8 * time.Hour),
				timeTo:      zeroTime.Add(20 * time.Hour),
			},
			span2: absoluteTimeSpan{ // the entire day on 6th of January 2024 (Saturday)
				from: time.Date(2024, time.January, 6, 0, 0, 0, 0, time.UTC), // from Saturday 6st of January
				to:   time.Date(2024, time.January, 7, 0, 0, 0, 0, time.UTC), // to Sunday 7st of January
			},
			wantedResult: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotResult := areTimespanOverlapped(test.span1, test.span2)
			assert.Equal(t, test.wantedResult, gotResult)
		})
	}
}
