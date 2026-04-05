package main

import (
	"context"
	"log"
	"os"

	"github.com/aarbanas/productive-time-tracker/slackbot"
	"github.com/aarbanas/productive-time-tracker/slackbot/config"
)

func requiredEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalf("missing required environment variable %s", name)
	}
	return v
}

func main() {
	s3Bucket := requiredEnv("S3_BUCKET")
	s3Key := requiredEnv("S3_KEY")
	botToken := requiredEnv("SLACK_BOT_TOKEN")

	ctx := context.Background()
	store, err := config.NewStore(ctx, s3Bucket, s3Key)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}
	store.Load(ctx)

	slackbot.ServeSlackbot(store, botToken)
}
