package timeof

import (
	"strconv"
	"strings"
	"time"
)

var zoneStr = time.Now().Format("-0700")

// FIXME: 2033-05-18T03:33:19Z  Unix(1999999999) This function will faild to parse time.Time.
// TimeOf convert input string into time.Time struct.
// supported formats:
// (date) 20060102, 2006-01-02
//
//	(date+time) 20060102150405,2006-01-02-15-04-05,2006-01-02T15:04:05
//
// (timestamp in seconds) 1588776655
// (timestamp in miliseconds) 1707122646123
func TimeOf(input string) (time.Time, bool) {
	t := time.Time{}
	if input == "" {
		return t, false
	}
	var err error
	tryParse := func(layout string) bool {

		str := input + zoneStr
		t, err = time.Parse(layout+"-0700", str)
		return err == nil
	}
	// firstly try to parse input into date or datetime format.
	// should be fixed before 2100...
	if strings.HasPrefix(input, "20") {
		if tryParse("20060102") || tryParse("2006-01-02") ||
			tryParse("20060102150405") || tryParse("2006-01-02-15-04-05") || tryParse("2006-01-02T15:04:05") {
			return t, true
		}
		t, err = time.Parse(input, time.RFC3339)
		if err == nil {
			return t, true
		}
	}
	// then, try to parse the time as timestamp in seconds.
	// this have bugs for parsing time before 2001-09-09...
	if len(input) == 10 {
		ts, err := strconv.ParseInt(input, 10, 64)
		if err == nil {
			return time.Unix(ts, 0), true
		}
	} else if len(input) == 13 {
		ts, err := strconv.ParseInt(input, 10, 64)
		if err == nil {
			return time.UnixMilli(ts), true
		}
	}

	// invalid input.
	return t, false
}
