package utilities

import (
	"time"
)

func PreviousMonthBounds() (firstDate, lastDate string) {
	now := time.Now()

	firstDayThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDayPrevMonth := firstDayThisMonth.AddDate(0, 0, -1)
	firstDayPrevMonth := time.Date(lastDayPrevMonth.Year(), lastDayPrevMonth.Month(), 1, 0, 0, 0, 0, time.UTC)

	return firstDayPrevMonth.Format("2006-01-02"), lastDayPrevMonth.Format("2006-01-02")
}

func CurrentMonthBounds() (firstDate, lastDate string) {
	now := time.Now()
	firstDayThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	firstOfNextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	lastOfCurrentMonth := firstOfNextMonth.AddDate(0, 0, -1)

	return firstDayThisMonth.Format("2006-01-02"), lastOfCurrentMonth.Format("2006-01-02")
}
