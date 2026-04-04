package slackbot

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aarbanas/productive-time-tracker/slackbot/config"
	"github.com/aarbanas/productive-time-tracker/slackbot/handlers"
)

// Env vars (all required): S3_BUCKET, S3_KEY, SLACK_BOT_TOKEN.
func requiredEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalf("missing required environment variable %s", name)
	}
	return v
}

// ServeSlackbot starts the HTTP server for Slack slash commands and interactions.
func ServeSlackbot() {
	s3Bucket := requiredEnv("S3_BUCKET")
	s3Key := requiredEnv("S3_KEY")
	botToken := requiredEnv("SLACK_BOT_TOKEN")

	ctx := context.Background()
	store, err := config.NewStore(ctx, s3Bucket, s3Key)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}
	store.Load(ctx)

	h := &handlers.Mux{Store: store, BotToken: botToken}
	http.HandleFunc("/command", h.SlashCommand)
	http.HandleFunc("/interactive", h.Interactive)

	log.Println("Server running on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
