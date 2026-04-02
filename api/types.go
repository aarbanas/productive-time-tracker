package api

// TimeEntry is a single time entry from the Productive API.
type TimeEntry struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Date           string `json:"date"`
		Time           int    `json:"time"`
		BillableTime   int    `json:"billable_time"`
		RecognizedTime int    `json:"recognized_time"`
		Note           string `json:"note"`
		Approved       bool   `json:"approved"`
		Submitted      bool   `json:"submitted"`
	} `json:"attributes"`
}

// TimeEntriesResponse is the JSON:API payload for time entries.
type TimeEntriesResponse struct {
	Data []TimeEntry `json:"data"`
}

// Event is an included event resource from the bookings API.
type Event struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name                     string      `json:"name"`
		EventTypeID              int         `json:"event_type_id"`
		IconID                   string      `json:"icon_id"`
		ColorID                  string      `json:"color_id"`
		ArchivedAt               interface{} `json:"archived_at"`
		LimitationTypeID         int         `json:"limitation_type_id"`
		SyncPersonalIntegrations bool        `json:"sync_personal_integrations"`
		HalfDayBookings          bool        `json:"half_day_bookings"`
		Description              interface{} `json:"description"`
		AbsenceType              string      `json:"absence_type"`
	} `json:"attributes"`
}

// Booking is a single booking from the API.
type Booking struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		StartedOn string `json:"started_on"`
		EndedOn   string `json:"ended_on"`
		TotalTime int    `json:"total_time"`
		Approved  bool   `json:"approved"`
		Note      string `json:"note"`
	} `json:"attributes"`
	Relationships struct {
		Event struct {
			Data *struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"data"`
		} `json:"event"`
	} `json:"relationships"`
}

// BookingsResponse is the JSON:API payload for bookings with included events.
type BookingsResponse struct {
	Data   []Booking `json:"data"`
	Events []Event   `json:"included"`
}
