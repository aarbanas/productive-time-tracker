package config

// UserConfig holds per-user Productive API credentials and preferences.
type UserConfig struct {
	Token                string `json:"token"`
	OrgID                string `json:"orgId"`
	MinHours             int    `json:"minHours"`
	ScheduleNotification bool   `json:"scheduleNotification,omitempty"`
}
