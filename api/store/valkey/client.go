package valkey

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/valkey-io/valkey-go"
)

type Client struct {
	client valkey.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	addr := getEnvOrDefault("VALKEY_ADDR", "localhost:6379")
	password := os.Getenv("VALKEY_PASSWORD")

	opts := valkey.ClientOption{
		InitAddress: []string{addr},
		Password:    password,
	}

	client, err := valkey.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to valkey: %w", err)
	}

	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		return nil, fmt.Errorf("failed to ping valkey: %w", err)
	}

	return &Client{client: client}, nil
}

func (c *Client) Close() {
	c.client.Close()
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type SessionStore struct {
	client *Client
	prefix string
}

func NewSessionStore(client *Client) *SessionStore {
	return &SessionStore{
		client: client,
		prefix: "session:",
	}
}

func (s *SessionStore) key(sessionID string) string {
	return s.prefix + sessionID
}

func (s *SessionStore) Set(ctx context.Context, sessionID string, data []byte, ttl time.Duration) error {
	return s.client.client.Do(ctx,
		s.client.client.B().Set().Key(s.key(sessionID)).Value(string(data)).Ex(ttl).Build(),
	).Error()
}

func (s *SessionStore) Exists(ctx context.Context, sessionID string) (bool, error) {
	result, err := s.client.client.Do(ctx,
		s.client.client.B().Exists().Key(s.key(sessionID)).Build(),
	).AsInt64()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (s *SessionStore) Get(ctx context.Context, sessionID string) ([]byte, error) {
	result, err := s.client.client.Do(ctx,
		s.client.client.B().Get().Key(s.key(sessionID)).Build(),
	).AsBytes()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (s *SessionStore) Delete(ctx context.Context, sessionID string) error {
	return s.client.client.Do(ctx,
		s.client.client.B().Del().Key(s.key(sessionID)).Build(),
	).Error()
}
