package main

import "time"

var TIME_FORMAT = "2006-01-02T15:04"

// return local time
func ParseTime(input string) time.Time {
	t, _ := time.ParseInLocation(TIME_FORMAT, input, time.Local)
	return t
}

func IsLate(duration time.Time) bool {
	r := time.Now().Compare(duration)

	switch r {
	case -1:
		return false
	default:
		return true
	}
}
