package service

import (
	"fmt"

	"github.com/aarbanas/productive-time-tracker/api"
	"github.com/aarbanas/productive-time-tracker/slackbot/config"
	"github.com/aarbanas/productive-time-tracker/slackbot/views"
	"github.com/aarbanas/productive-time-tracker/utilities"
)

// MonthlySummaryForUser loads Productive data and returns Slack block-kit JSON for the summary.
func MonthlySummaryForUser(cfg config.UserConfig) (map[string]interface{}, error) {
	client := api.NewClient(cfg.Token, cfg.OrgID)

	totalAbsenceMinutes, err := utilities.AbsenceMinutes(client)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate absence minutes: %w", err)
	}

	totalTimeEntriesMinutes, err := utilities.TimeEntriesMinutes(client)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate time entries minutes: %w", err)
	}

	totalMinutes := totalAbsenceMinutes + totalTimeEntriesMinutes
	requiredMinutes, firstDayPrevMonth, lastDayPrevMonth := utilities.RequiredWorkingMinutesPreviousMonth(cfg.MinHours)

	return views.MonthlySummary(
		totalAbsenceMinutes,
		totalTimeEntriesMinutes,
		requiredMinutes,
		totalMinutes,
		firstDayPrevMonth,
		lastDayPrevMonth,
	), nil
}
