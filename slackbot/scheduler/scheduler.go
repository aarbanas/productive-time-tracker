package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aarbanas/productive-time-tracker/slackbot/config"
	"github.com/aarbanas/productive-time-tracker/slackbot/service"
	"github.com/aarbanas/productive-time-tracker/utilities"
	"github.com/slack-go/slack"
)

// Start launches a background ticker (every minute) that sends scheduled monthly summaries by DM.
func Start(store *config.Store, botToken string) {
	api := slack.New(botToken)
	s := &scheduler{store: store, api: api, sent: make(map[string]struct{})}
	go s.loop()
}

type scheduler struct {
	store *config.Store
	api   *slack.Client
	mu    sync.Mutex
	sent  map[string]struct{}
}

func (s *scheduler) loop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.runOnce(time.Now())
	}
}

func (s *scheduler) runOnce(now time.Time) {
	for userID, cfg := range s.store.Snapshot() {
		if !cfg.ScheduleNotification {
			continue
		}

		nextRun, err := utilities.NextRunLastWorkingDayUTC(now, 12, 0)
		if err != nil {
			log.Printf("scheduler: user %s next run: %v", userID, err)
			continue
		}

		if !utilities.SameScheduleMinute(now, nextRun) {
			continue
		}

		key := userID + "|" + nextRun.Format("2006-01-02T15:04")
		s.mu.Lock()
		if _, ok := s.sent[key]; ok {
			s.mu.Unlock()
			continue
		}
		s.sent[key] = struct{}{}
		s.mu.Unlock()

		if err := s.sendDM(userID, cfg); err != nil {
			log.Printf("scheduler: DM user %s: %v", userID, err)
			s.mu.Lock()
			delete(s.sent, key)
			s.mu.Unlock()
		}
	}
}

func (s *scheduler) sendDM(userID string, cfg config.UserConfig) error {
	ch, _, _, err := s.api.OpenConversation(&slack.OpenConversationParameters{
		Users:    []string{userID},
		ReturnIM: true,
	})
	if err != nil {
		return err
	}
	if ch == nil {
		return errors.New("slack: conversations.open returned no channel")
	}

	payload, err := service.MonthlySummaryForUser(cfg)
	if err != nil {
		return err
	}
	raw, ok := payload["blocks"]
	if !ok {
		return errors.New("monthly summary: missing blocks")
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("marshal blocks: %w", err)
	}
	var blocks slack.Blocks
	if err := blocks.UnmarshalJSON(b); err != nil {
		return fmt.Errorf("parse slack blocks: %w", err)
	}
	_, _, err = s.api.PostMessage(ch.ID, slack.MsgOptionBlocks(blocks.BlockSet...))
	return err
}
