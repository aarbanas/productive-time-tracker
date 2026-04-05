package views

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/slack-go/slack"
)

func ConfigurationModal(initialScheduleNotification bool) slack.ModalViewRequest {
	scheduleOpt := slack.NewOptionBlockObject(
		"schedule_on",
		slack.NewTextBlockObject(slack.PlainTextType, "Last working day of each month at 12:00 UTC", false, false),
		nil,
	)
	scheduleCheckbox := slack.NewCheckboxGroupsBlockElement("schedule_checkbox", scheduleOpt)
	if initialScheduleNotification {
		scheduleCheckbox.InitialOptions = []*slack.OptionBlockObject{scheduleOpt}
	}

	return slack.ModalViewRequest{
		Type:       slack.VTModal,
		Title:      slack.NewTextBlockObject(slack.PlainTextType, "Configuration", false, false),
		Submit:     slack.NewTextBlockObject(slack.PlainTextType, "Save", false, false),
		Close:      slack.NewTextBlockObject(slack.PlainTextType, "Cancel", false, false),
		CallbackID: "config_modal",
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewInputBlock(
					"token_block",
					slack.NewTextBlockObject(slack.PlainTextType, "Token", false, false),
					nil,
					slack.NewPlainTextInputBlockElement(
						slack.NewTextBlockObject(slack.PlainTextType, "Enter token", false, false),
						"token_input",
					),
				),
				slack.NewInputBlock(
					"org_block",
					slack.NewTextBlockObject(slack.PlainTextType, "Org ID", false, false),
					nil,
					slack.NewPlainTextInputBlockElement(
						slack.NewTextBlockObject(slack.PlainTextType, "Enter org ID", false, false),
						"org_input",
					),
				),
				slack.NewInputBlock(
					"hours_block",
					slack.NewTextBlockObject(slack.PlainTextType, "Min hours/day", false, false),
					nil,
					slack.NewPlainTextInputBlockElement(
						slack.NewTextBlockObject(slack.PlainTextType, "e.g. 8", false, false),
						"hours_input",
					),
				),
				slack.NewInputBlock(
					"schedule_block",
					slack.NewTextBlockObject(slack.PlainTextType, "Monthly DM reminder", false, false),
					nil,
					scheduleCheckbox,
				),
			},
		},
	}
}

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
