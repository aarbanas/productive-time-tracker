package utilities

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Holiday struct {
	StartDate string `json:"startDate"`
}

var (
	holidayMu    sync.Mutex
	holidayCache map[int]map[string]bool // year → YYYY-MM-DD set
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

func RequiredWorkingMinutesPreviousMonth(minHours int) (int32, *time.Time, *time.Time) {
	return requiredWorkingMinutesPreviousMonthAt(time.Now(), minHours)
}

func requiredWorkingMinutesPreviousMonthAt(now time.Time, minHours int) (int32, *time.Time, *time.Time) {
	firstDayThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDayPrevMonth := firstDayThisMonth.AddDate(0, 0, -1)
	firstDayPrevMonth := time.Date(lastDayPrevMonth.Year(), lastDayPrevMonth.Month(), 1, 0, 0, 0, 0, time.UTC)

	workingDays := 0
	for d := firstDayPrevMonth; !d.After(lastDayPrevMonth); d = d.AddDate(0, 0, 1) {
		if w := d.Weekday(); w != time.Saturday && w != time.Sunday {
			workingDays++
		}
	}

	return int32(workingDays * minHours * 60), &firstDayPrevMonth, &lastDayPrevMonth
}

func fetchHolidaysForYear(year int) (map[string]bool, error) {
	url := fmt.Sprintf(
		"https://openholidaysapi.org/PublicHolidays?countryIsoCode=HR&validFrom=%d-01-01&validTo=%d-12-31",
		year, year,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close holidays response body:", err)
			return
		}
	}(resp.Body)

	var holidays []Holiday
	if err := json.NewDecoder(resp.Body).Decode(&holidays); err != nil {
		return nil, err
	}

	holidayMap := make(map[string]bool)
	for _, h := range holidays {
		holidayMap[h.StartDate] = true
	}

	return holidayMap, nil
}

func holidaysForYear(year int) (map[string]bool, error) {
	holidayMu.Lock()
	defer holidayMu.Unlock()
	if holidayCache != nil {
		if m, ok := holidayCache[year]; ok {
			return m, nil
		}
	} else {
		holidayCache = make(map[int]map[string]bool)
	}
	m, err := fetchHolidaysForYear(year)
	if err != nil {
		return nil, err
	}
	holidayCache[year] = m
	return m, nil
}

func lastWorkingDay(year int, month time.Month) (*time.Time, error) {
	return lastWorkingDayInLoc(year, month, time.Local)
}

func lastWorkingDayInLoc(year int, month time.Month, loc *time.Location) (*time.Time, error) {
	firstNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, loc)
	lastDay := firstNextMonth.AddDate(0, 0, -1)

	var holidays map[string]bool
	var loadedYear int

	for {
		if lastDay.Year() != loadedYear {
			var err error
			holidays, err = holidaysForYear(lastDay.Year())
			if err != nil {
				fmt.Println("Failed to fetch holidays:", err)
				return nil, err
			}
			loadedYear = lastDay.Year()
		}

		dateStr := lastDay.Format("2006-01-02")

		isWeekend := lastDay.Weekday() == time.Saturday || lastDay.Weekday() == time.Sunday
		isHoliday := holidays[dateStr]

		if !isWeekend && !isHoliday {
			chosen := lastDay
			return &chosen, nil
		}

		lastDay = lastDay.AddDate(0, 0, -1)
	}
}

func parseExecutionTime(exec string) (int, int, error) {
	t, err := time.Parse("15:04", exec)
	if err != nil {
		return 0, 0, err
	}
	return t.Hour(), t.Minute(), nil
}

func NextRunTime(now time.Time, executionTime string) (time.Time, error) {
	year, month, _ := now.Date()
	hour, minute, err := parseExecutionTime(executionTime)
	if err != nil {
		return time.Time{}, err
	}

	runDay, err := lastWorkingDay(year, month)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(runDay.Year(), runDay.Month(), runDay.Day(), hour, minute, 0, 0, time.Local), nil
}

// NextRunLastWorkingDayUTC returns the instant on the last working day of the UTC calendar month
// that contains now, at hour:minute in UTC (Croatian public holidays applied to those dates).
func NextRunLastWorkingDayUTC(now time.Time, hour, minute int) (time.Time, error) {
	y, m, _ := now.UTC().Date()
	runDay, err := lastWorkingDayInLoc(y, m, time.UTC)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(runDay.Year(), runDay.Month(), runDay.Day(), hour, minute, 0, 0, time.UTC), nil
}

func SameScheduleMinute(a, b time.Time) bool {
	a = a.In(time.Local).Truncate(time.Minute)
	b = b.In(time.Local).Truncate(time.Minute)
	return a.Equal(b)
}
