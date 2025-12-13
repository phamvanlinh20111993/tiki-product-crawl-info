package util

import "time"

func timeToString(time time.Time, pattern string) string {
	if pattern == "" {
		pattern = "2006-01-02 15:04:05"
	}
	return time.Format(pattern)
}

var TimeToString = timeToString
