package utilities

import (
	"github.com/aarbanas/productive-time-tracker/api"
)

func ReportMinutes(c *api.ProductiveClient, currentMonth bool) (int32, error) {
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
		if attrs, ok := report.Attributes.Get(); ok {
			totalWorkedMinutes += attrs.WorkedTime.Or(0)
			totalScheduledMinutes += attrs.ScheduledTime.Or(0)
			totalAbsenceMinutes += attrs.EventTime.Or(0)
		}
	}

	remainingMinutes := totalScheduledMinutes - (totalWorkedMinutes + totalAbsenceMinutes)

	return int32(remainingMinutes), nil
}
