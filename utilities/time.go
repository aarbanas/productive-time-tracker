package utilities

import (
	"time"
)

func PreviousMonthBounds() (firstDate, lastDate string) {
	return previousMonthBoundsAt(time.Now())
}

func previousMonthBoundsAt(now time.Time) (firstDate, lastDate string) {
	firstDayThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDayPrevMonth := firstDayThisMonth.AddDate(0, 0, -1)
	firstDayPrevMonth := time.Date(lastDayPrevMonth.Year(), lastDayPrevMonth.Month(), 1, 0, 0, 0, 0, time.UTC)

	return firstDayPrevMonth.Format("2006-01-02"), lastDayPrevMonth.Format("2006-01-02")
}

func RequiredWorkingMinutesPreviousMonth(minHours int) (int32, time.Time, time.Time) {
	return requiredWorkingMinutesPreviousMonthAt(time.Now(), minHours)
}

func requiredWorkingMinutesPreviousMonthAt(now time.Time, minHours int) (int32, time.Time, time.Time) {
	firstDayThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDayPrevMonth := firstDayThisMonth.AddDate(0, 0, -1)
	firstDayPrevMonth := time.Date(lastDayPrevMonth.Year(), lastDayPrevMonth.Month(), 1, 0, 0, 0, 0, time.UTC)

	workingDays := 0
	for d := firstDayPrevMonth; !d.After(lastDayPrevMonth); d = d.AddDate(0, 0, 1) {
		if w := d.Weekday(); w != time.Saturday && w != time.Sunday {
			workingDays++
		}
	}

	return int32(workingDays * minHours * 60), firstDayPrevMonth, lastDayPrevMonth
}
