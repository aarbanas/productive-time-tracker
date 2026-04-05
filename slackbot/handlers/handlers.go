package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/aarbanas/productive-time-tracker/slackbot/config"
	"github.com/aarbanas/productive-time-tracker/slackbot/service"
	"github.com/aarbanas/productive-time-tracker/slackbot/views"
	"github.com/slack-go/slack"
)

// Mux holds dependencies for Slack HTTP endpoints.
type Mux struct {
	Store    *config.Store
	BotToken string
}

func (m *Mux) SlashCommand(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userID := r.FormValue("user_id")
	text := r.FormValue("text")
	args := strings.Fields(text)

	if len(args) == 0 {
		views.WriteFirstTimeHelp(w)
		return
	}

	switch args[0] {
	case "configure":
		m.handleConfigure(w, r)
	case "check":
		m.handleCheck(w, userID)
	case "help":
		views.WriteFirstTimeHelp(w)
	default:
		views.WritePlainText(w, "Unknown command. Use 'help' for more information.'")
	}
}

func (m *Mux) handleConfigure(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	triggerID := r.FormValue("trigger_id")
	userID := r.FormValue("user_id")

	var initialScheduleNotification bool
	if c, ok := m.Store.Get(userID); ok {
		initialScheduleNotification = c.ScheduleNotification
	}

	modal := views.ConfigurationModal(initialScheduleNotification)

	slackClient := slack.New(m.BotToken)
	_, err := slackClient.OpenView(triggerID, modal)
	if err != nil {
		log.Println("Error opening modal:", err)
	}

	w.WriteHeader(http.StatusOK)
}

func (m *Mux) handleCheck(w http.ResponseWriter, userID string) {
	cfg, ok := m.Store.Get(userID)
	if !ok {
		views.WriteFirstTimeHelp(w)
		return
	}

	response, err := service.MonthlySummaryForUser(cfg)
	if err != nil {
		views.WritePlainText(w, fmt.Sprintf("Error: %v", err))
		return
	}

	views.WriteJSON(w, response)
}

func (m *Mux) Interactive(w http.ResponseWriter, r *http.Request) {
	payload := r.FormValue("payload")

	var callback slack.InteractionCallback
	if err := json.Unmarshal([]byte(payload), &callback); err != nil {
		http.Error(w, "bad payload", http.StatusBadRequest)
		return
	}

	if callback.Type == slack.InteractionTypeViewSubmission {
		m.handleModalSubmit(callback)
	}

	w.WriteHeader(http.StatusOK)
}

func (m *Mux) handleModalSubmit(callback slack.InteractionCallback) {
	userID := callback.User.ID
	values := callback.View.State.Values

	token := values["token_block"]["token_input"].Value
	orgID := values["org_block"]["org_input"].Value
	hoursStr := values["hours_block"]["hours_input"].Value
	scheduleNotification := false
	if v, ok := values["schedule_block"]["schedule_checkbox"]; ok && len(v.SelectedOptions) > 0 {
		scheduleNotification = true
	}

	minHours := 8
	if h, err := strconv.Atoi(hoursStr); err == nil {
		minHours = h
	}

	m.Store.Set(userID, config.UserConfig{
		Token:                token,
		OrgID:                orgID,
		MinHours:             minHours,
		ScheduleNotification: scheduleNotification,
	})
	m.Store.SaveAsync()

	log.Println("Saved config for user:", userID)
}
