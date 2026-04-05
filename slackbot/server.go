package slackbot

import (
	"log"
	"net/http"

	"github.com/aarbanas/productive-time-tracker/slackbot/config"
	"github.com/aarbanas/productive-time-tracker/slackbot/handlers"
	"github.com/aarbanas/productive-time-tracker/slackbot/scheduler"
)

// ServeSlackbot starts the HTTP server for Slack slash commands and interactions.
func ServeSlackbot(store *config.Store, botToken string) {

	h := &handlers.Mux{Store: store, BotToken: botToken}
	http.HandleFunc("/command", h.SlashCommand)
	http.HandleFunc("/interactive", h.Interactive)

	scheduler.Start(store, botToken)

	log.Println("Server running on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
