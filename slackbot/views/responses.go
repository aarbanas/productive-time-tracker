package views

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// FirstTimeHelp returns the Slack slash-command JSON for onboarding.
func FirstTimeHelp() map[string]interface{} {
	return map[string]interface{}{
		"response_type": "ephemeral",
		"blocks": []map[string]interface{}{
			{
				"type": "section",
				"text": map[string]string{
					"type": "mrkdwn",
					"text": "👋 *Welcome to SlackBot!*",
				},
			},
			{
				"type": "section",
				"text": map[string]string{
					"type": "mrkdwn",
					"text": "*Getting started:*\n1. Run `/tracker configure`\n2. Enter your token & org ID\n3. Run `/tracker check`",
				},
			},
			{
				"type": "context",
				"elements": []map[string]string{
					{
						"type": "mrkdwn",
						"text": "Need help? Just run the command again anytime.",
					},
				},
			},
		},
	}
}

// MonthlySummary builds the Slack JSON for the monthly check result.
func MonthlySummary(
	totalAbsenceMinutes int32,
	totalTimeEntriesMinutes int32,
	requiredMinutes int32,
	totalMinutes int32,
	firstDayPrevMonth time.Time,
	lastDayPrevMonth time.Time,
) map[string]interface{} {
	absence := totalAbsenceMinutes / 60
	entries := totalTimeEntriesMinutes / 60
	required := requiredMinutes / 60
	total := totalMinutes / 60

	var statusText string
	if totalMinutes < requiredMinutes {
		diff := (requiredMinutes - totalMinutes) / 60
		statusText = fmt.Sprintf("⚠️ *You are %d hours behind schedule.*", diff)
	} else {
		extra := (totalMinutes - requiredMinutes) / 60
		statusText = fmt.Sprintf("✅ *Great job!* You are %d hours ahead 🚀", extra)
	}

	return map[string]interface{}{
		"response_type": "ephemeral",
		"blocks": []map[string]interface{}{
			{
				"type": "header",
				"text": map[string]string{
					"type": "plain_text",
					"text": fmt.Sprintf("📊 Monthly Summary (%s - %s)", firstDayPrevMonth.Format("02.01.2006"), lastDayPrevMonth.Format("02.01.2006")),
				},
			},
			{"type": "divider"},
			{
				"type": "section",
				"text": map[string]string{
					"type": "mrkdwn",
					"text": fmt.Sprintf(
						"*📝 Time entries:* %dh\n\n*⏸ Absence:* %dh\n\n*🎯 Required:* %dh\n\n*⏱ Total tracked:* %dh",
						entries,
						absence,
						required,
						total,
					),
				},
			},
			{
				"type": "section",
				"text": map[string]string{
					"type": "mrkdwn",
					"text": statusText,
				},
			},
		},
	}
}

// WriteJSON sets Content-Type and encodes v as JSON.
func WriteJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// WritePlainText responds with text/plain.
func WritePlainText(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = fmt.Fprint(w, msg)
}

// WriteFirstTimeHelp writes the onboarding payload as JSON.
func WriteFirstTimeHelp(w http.ResponseWriter) {
	WriteJSON(w, FirstTimeHelp())
}
