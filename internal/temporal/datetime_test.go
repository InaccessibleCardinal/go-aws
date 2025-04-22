package temporal

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/type/datetime"
	"google.golang.org/protobuf/types/known/durationpb"
)

func Test_timeToDateTime(t *testing.T) {
	type testCase struct {
		name   string
		timeIn time.Time
		want   *datetime.DateTime
	}

	for _, tc := range []testCase{
		{
			name:   "time_with_valid_time_zone",
			timeIn: time.Date(2020, 1, 21, 12, 1, 1, 1, time.FixedZone("America/Chicago", -18000)),
			want: &datetime.DateTime{
				Year: 2020, Month: 1, Day: 21, Hours: 12, Minutes: 1, Seconds: 1, Nanos: 1,
				TimeOffset: &datetime.DateTime_TimeZone{TimeZone: &datetime.TimeZone{Id: "America/Chicago"}},
			},
		},
		{
			name:   "time_with_offset_west",
			timeIn: time.Date(2020, 1, 21, 12, 1, 1, 1, time.FixedZone("", -14400)),
			want: &datetime.DateTime{
				Year: 2020, Month: 1, Day: 21, Hours: 12, Minutes: 1, Seconds: 1, Nanos: 1,
				TimeOffset: &datetime.DateTime_UtcOffset{UtcOffset: &durationpb.Duration{Seconds: -14400}},
			},
		},
		{
			name:   "time_with_offset_east",
			timeIn: time.Date(2020, 1, 21, 12, 1, 1, 1, time.FixedZone("", 14400)),
			want: &datetime.DateTime{
				Year: 2020, Month: 1, Day: 21, Hours: 12, Minutes: 1, Seconds: 1, Nanos: 1,
				TimeOffset: &datetime.DateTime_UtcOffset{UtcOffset: &durationpb.Duration{Seconds: 14400}},
			},
		},
		{
			name:   "time_in_utc",
			timeIn: time.Date(2020, 1, 21, 12, 1, 1, 1, time.UTC),
			want:   &datetime.DateTime{Year: 2020, Month: 1, Day: 21, Hours: 12, Minutes: 1, Seconds: 1, Nanos: 1},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := timeToDateTime(tc.timeIn)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_dateTimeToTime(t *testing.T) {
	type testCase struct {
		name       string
		datetimeIn *datetime.DateTime
		want       time.Time
	}

	for _, tc := range []testCase{
		{
			name: "datetime_with_time_zone",
			datetimeIn: &datetime.DateTime{
				Year: 2020, Month: 1, Day: 1, Hours: 12, Minutes: 1, Seconds: 23, Nanos: 1,
				TimeOffset: &datetime.DateTime_TimeZone{TimeZone: &datetime.TimeZone{Id: "America/New_York"}},
			},
			want: time.Date(2020, 1, 1, 12, 1, 23, 1, time.FixedZone("", -14400)),
		},
		{
			name: "datetime_with_time_zone_east",
			datetimeIn: &datetime.DateTime{
				Year: 2020, Month: 1, Day: 1, Hours: 12, Minutes: 1, Seconds: 23, Nanos: 1,
				TimeOffset: &datetime.DateTime_UtcOffset{UtcOffset: &durationpb.Duration{Seconds: 7200}},
			},
			want: time.Date(2020, 1, 1, 12, 1, 23, 1, time.FixedZone("", 7200)),
		},
		{
			name: "datetime_with_time_zone_west",
			datetimeIn: &datetime.DateTime{
				Year: 2020, Month: 1, Day: 1, Hours: 12, Minutes: 1, Seconds: 23, Nanos: 1,
				TimeOffset: &datetime.DateTime_UtcOffset{UtcOffset: &durationpb.Duration{Seconds: -18000}},
			},
			want: time.Date(2020, 1, 1, 12, 1, 23, 1, time.FixedZone("", -18000)),
		},
		{
			name: "datetime_in_utc",
			datetimeIn: &datetime.DateTime{
				Year: 2020, Month: 1, Day: 1, Hours: 12, Minutes: 1, Seconds: 23, Nanos: 1,
			},
			want: time.Date(2020, 1, 1, 12, 1, 23, 1, time.UTC),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := dateTimeToTime(tc.datetimeIn)
			assertSameDate(t, tc.want, got)
		})
	}
}

func Test_offsetFromTimeZone(t *testing.T) {
	type testCase struct {
		name     string
		zoneName string
		want     int
	}

	for _, tc := range []testCase{
		{
			name:     "time_zone_known",
			zoneName: "America/New_York",
			want:     -14400,
		},
		{
			name:     "time_zone_unknown",
			zoneName: "Fake/Place",
			want:     0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := offsetFromZone(tc.zoneName)
			assert.Equal(t, tc.want, got)
		})
	}
	// tz := "America/New_York"
	// want := -14400
	// got := offsetFromZone(tz)
	// assert.Equal(t, want, got)
}

func assertSameDate(t *testing.T, wantTime, gotTime time.Time) {
	assert.Equal(t, wantTime.Year(), gotTime.Year())
	assert.Equal(t, wantTime.Month(), gotTime.Month())
	assert.Equal(t, wantTime.Day(), gotTime.Day())
	assert.Equal(t, wantTime.Hour(), gotTime.Hour())
	assert.Equal(t, wantTime.Minute(), gotTime.Minute())
	assert.Equal(t, wantTime.Second(), gotTime.Second())

	_, wantOffset := wantTime.Zone()
	_, gotOffset := gotTime.Zone()
	assert.Equal(t, wantOffset, gotOffset)
}

func Test_LocalBug(t *testing.T) { //21600
	tm := time.Date(2020, 1, 1, 1, 1, 1, 1, time.FixedZone("", -21600))

	s := tm.Format(time.RFC3339Nano)

	fmt.Println("s:", s)

	parsed, _ := time.Parse(time.RFC3339Nano, s)

	fmt.Println("parsed:", parsed)
}
