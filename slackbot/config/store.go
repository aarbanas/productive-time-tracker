package config

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"maps"
	"sync"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Store struct {
	mu     sync.RWMutex
	users  map[string]UserConfig
	client *s3.Client
	bucket string
	key    string
}

func NewStore(ctx context.Context, bucket, key string) (*Store, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &Store{
		users:  make(map[string]UserConfig),
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
		key:    key,
	}, nil
}

func (s *Store) Load(ctx context.Context) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &s.key,
	})
	if err != nil {
		log.Println("No existing config found in S3 (starting fresh)")
		return
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Println("Failed to close S3 config body:", err)
		}
	}(out.Body)

	var data map[string]UserConfig
	if err := json.NewDecoder(out.Body).Decode(&data); err != nil {
		log.Println("Failed to decode S3 config:", err)
		return
	}

	s.mu.Lock()
	s.users = data
	s.mu.Unlock()

	log.Println("Configs loaded from S3")
}

func (s *Store) Save(ctx context.Context) {
	s.mu.RLock()
	data, err := json.MarshalIndent(s.users, "", "  ")
	s.mu.RUnlock()

	if err != nil {
		log.Println("Failed to marshal configs:", err)
		return
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &s.key,
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		log.Println("Failed to upload to S3:", err)
		return
	}

	log.Println("Configs saved to S3")
}

func (s *Store) SaveAsync() {
	go func() {
		s.Save(context.Background())
	}()
}

func (s *Store) Get(userID string) (UserConfig, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.users[userID]
	return c, ok
}

func (s *Store) Set(userID string, cfg UserConfig) {
	s.mu.Lock()
	s.users[userID] = cfg
	s.mu.Unlock()
}

func (s *Store) Snapshot() map[string]UserConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]UserConfig, len(s.users))
	maps.Copy(out, s.users)
	return out
}
