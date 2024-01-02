package internal

import "time"

func getYearBounds(year int) (time.Time, time.Time) {
	loc, _ := time.LoadLocation("Local")
	firstDay := time.Date(year, time.January, 1, 0, 0, 0, 0, loc)
	lastDay := time.Date(year, time.December, 31, 23, 59, 59, 999999999, loc)
	return firstDay, lastDay
}
