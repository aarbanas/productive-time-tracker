package slackbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/aarbanas/productive-time-tracker/api"
	"github.com/aarbanas/productive-time-tracker/utilities"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	userConfigs = map[string]UserConfig{}
	mu          sync.RWMutex

	s3Client *s3.Client
	bucket   = "productive-time-tracker"
	key      = "configs.json"
)

func initS3() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
}

func loadConfigsFromS3() {
	out, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		log.Println("No existing config found in S3 (starting fresh)")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Failed to close S3 config body:", err)
		}
	}(out.Body)

	var data map[string]UserConfig
	if err := json.NewDecoder(out.Body).Decode(&data); err != nil {
		log.Println("Failed to decode S3 config:", err)
		return
	}

	mu.Lock()
	userConfigs = data
	mu.Unlock()

	log.Println("Configs loaded from S3")
}

func saveConfigsToS3() {
	mu.RLock()
	data, err := json.MarshalIndent(userConfigs, "", "  ")
	mu.RUnlock()

	if err != nil {
		log.Println("Failed to marshal configs:", err)
		return
	}

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		log.Println("Failed to upload to S3:", err)
		return
	}

	log.Println("Configs saved to S3")
}

func respond(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := fmt.Fprint(w, msg)
	if err != nil {
		return
	}
}

func runYourExistingLogic(cfg UserConfig) (string, error) {
	client := api.NewClient(cfg.Token, cfg.OrgID)

	totalAbsenceMinutes, err := utilities.AbsenceMinutes(client)
	if err != nil {
		return "", fmt.Errorf("failed to calculate absence minutes: %v", err)
	}

	totalTimeEntriesMinutes, err := utilities.TimeEntriesMinutes(client)
	if err != nil {
		return "", fmt.Errorf("failed to calculate time entries minutes: %v\n", err)
	}

	totalMinutes := totalAbsenceMinutes + totalTimeEntriesMinutes
	requiredMinutes := utilities.RequiredWorkingMinutesPreviousMonth()

	var b strings.Builder
	_, err = fmt.Fprintf(&b, "Total absence hours: %d\n", totalAbsenceMinutes/60)
	if err != nil {
		return "", err
	}
	_, err = fmt.Fprintf(&b, "Total time entries hours: %d\n", totalTimeEntriesMinutes/60)
	if err != nil {
		return "", err
	}
	_, err = fmt.Fprintf(&b, "Required hours to track for previous month: %d\n", requiredMinutes/60)
	if err != nil {
		return "", err
	}
	_, err = fmt.Fprintf(&b, "Total hours tracked: %d\n", totalMinutes/60)
	if err != nil {
		return "", err
	}

	if totalMinutes < requiredMinutes {
		_, err = fmt.Fprintf(&b, "You are %d minutes behind schedule.\n", requiredMinutes-totalMinutes)
		if err != nil {
			return "", err
		}
	} else {
		b.WriteString("Great job! You are on track!\n")
	}

	return b.String(), nil
}

func slashHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userID := r.FormValue("user_id")
	text := r.FormValue("text")
	args := strings.Fields(text)

	if len(args) == 0 {
		respond(w, "Usage: configure <token> <org_id> [hours] | check")
		return
	}

	switch args[0] {
	case "configure":
		handleConfigure(w, userID, args[1:])
	case "check":
		handleCheck(w, userID)
	default:
		respond(w, "Unknown command. Use 'configure' or 'check'")
	}
}

func handleConfigure(w http.ResponseWriter, userID string, args []string) {
	if len(args) < 2 {
		respond(w, "Usage: configure <token> <org_id> [hours]")
		return
	}

	token := args[0]
	orgID := args[1]
	minHours := 8

	if len(args) >= 3 {
		if h, err := strconv.Atoi(args[2]); err == nil {
			minHours = h
		}
	}

	mu.Lock()
	userConfigs[userID] = UserConfig{
		Token:    token,
		OrgID:    orgID,
		MinHours: minHours,
	}
	mu.Unlock()

	// Persist to S3
	go saveConfigsToS3()

	respond(w, "Configuration saved ✅")
}

func handleCheck(w http.ResponseWriter, userID string) {
	mu.RLock()
	cfg, ok := userConfigs[userID]
	mu.RUnlock()

	if !ok {
		respond(w, "You need to configure first")
		return
	}

	result, err := runYourExistingLogic(cfg)
	if err != nil {
		respond(w, fmt.Sprintf("Error: %v", err))
		return
	}
	respond(w, result)
}

func ServeSlackbot() {
	initS3()
	loadConfigsFromS3()

	http.HandleFunc("/command", slashHandler)

	log.Println("Server running on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
