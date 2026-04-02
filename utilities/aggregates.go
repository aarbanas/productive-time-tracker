package utilities

import (
	"slices"

	"github.com/aarbanas/productive-time-tracker/api"
)

func AbsenceMinutes(c *api.Client) (int32, error) {
	after, before := PreviousMonthBounds()
	bookings, events, err := c.GetBookingsWithEvents(after, before)
	if err != nil {
		return 0, err
	}

	eventIDs := make([]string, 0, len(events))
	for _, e := range events {
		eventIDs = append(eventIDs, e.ID)
	}

	totalAbsenceMinutes := 0
	for _, booking := range bookings {
		if booking.Relationships.Event.Data == nil {
			continue
		}
		eventID := booking.Relationships.Event.Data.ID
		if slices.Contains(eventIDs, eventID) {
			totalAbsenceMinutes += booking.Attributes.TotalTime
		}
	}

	return int32(totalAbsenceMinutes), nil
}

func TimeEntriesMinutes(c *api.Client) (int32, error) {
	after, before := PreviousMonthBounds()
	entries, err := c.GetTimeEntries(after, before)
	if err != nil {
		return 0, err
	}

	total := 0
	for _, entry := range entries {
		total += entry.Attributes.Time
	}

	return int32(total), nil
}
