package temporal

import (
	"time"

	"google.golang.org/genproto/googleapis/type/datetime"
	"google.golang.org/protobuf/types/known/durationpb"
)

func timeToDateTime(t time.Time) *datetime.DateTime {
	dt := &datetime.DateTime{
		Year:    int32(t.Year()),
		Month:   int32(t.Month()),
		Day:     int32(t.Day()),
		Hours:   int32(t.Hour()),
		Minutes: int32(t.Minute()),
		Seconds: int32(t.Second()),
		Nanos:   int32(t.Second()),
	}
	zone, offset := t.Zone()
	if zone == "UTC" && offset == 0 {
		return dt
	}
	if zone != "" {
		dt.TimeOffset = &datetime.DateTime_TimeZone{TimeZone: &datetime.TimeZone{Id: zone}}
	} else {
		dt.TimeOffset = &datetime.DateTime_UtcOffset{UtcOffset: &durationpb.Duration{Seconds: int64(offset)}}
	}
	return dt
}

func dateTimeToTime(dt *datetime.DateTime) time.Time {
	zone := dt.GetTimeZone()
	offset := dt.GetUtcOffset()

	if zone != nil {
		return timeWithLocation(dt, time.FixedZone(zone.Id, offsetFromZone(zone.Id)))
	}
	if offset != nil {
		return timeWithLocation(dt, time.FixedZone("", int(offset.Seconds)))
	}
	return timeWithLocation(dt, time.UTC)
}

func timeWithLocation(dt *datetime.DateTime, loc *time.Location) time.Time {
	return time.Date(
		int(dt.Year),
		time.Month(dt.Month),
		int(dt.Day),
		int(dt.Hours),
		int(dt.Minutes),
		int(dt.Seconds),
		int(dt.Nanos),
		loc,
	)
}

func offsetFromZone(zone string) int {
	t := time.Now()
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return 0
	}
	t = t.In(loc)
	_, offset := t.Zone()
	return offset
}
