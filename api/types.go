package api

// TimeReport is a single time report row from the API.
type TimeReport struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Attributes TimeReportAttributes `json:"attributes"`
}

// TimeReportAttributes are the fields used for tracked-hours calculations.
type TimeReportAttributes struct {
	WorkedTime    float64 `json:"worked_time"`
	EventTime     float64 `json:"event_time"`
	ScheduledTime float64 `json:"scheduled_time"`
}

// TimeReportsResponse is the JSON:API payload for time reports.
type TimeReportsResponse struct {
	Data []TimeReport `json:"data"`
}
