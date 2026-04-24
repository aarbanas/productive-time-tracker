package utilities

import (
	"github.com/aarbanas/productive-time-tracker/api"
)

func ReportMinutes(c *api.Client, currentMonth bool) (int32, error) {
	var after, before string
	if currentMonth {
		after, before = CurrentMonthBounds()
	} else {
		after, before = PreviousMonthBounds()
	}

	reports, err := c.GetTimeReport(after, before)
	if err != nil {
		return 0, err
	}

	var totalWorkedMinutes float64
	var totalScheduledMinutes float64
	var totalAbsenceMinutes float64

	for _, report := range reports {
		totalWorkedMinutes += report.Attributes.WorkedTime
		totalScheduledMinutes += report.Attributes.ScheduledTime
		totalAbsenceMinutes += report.Attributes.EventTime
	}

	remainingMinutes := totalScheduledMinutes - (totalWorkedMinutes + totalAbsenceMinutes)

	return int32(remainingMinutes), nil
}
