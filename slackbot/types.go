package slackbot

type UserConfig struct {
	Token    string `json:"token"`
	OrgID    string `json:"orgId"`
	MinHours int    `json:"minHours"`
}
