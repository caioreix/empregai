package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"go-api/internal/core/session"
	"go-api/pkg/config"
)

type sessionRepository struct {
	conn *redis.Client
	cfg  *config.Config
}

// Session redis repository constructor
func NewSessionRepository(c *redis.Client, cfg *config.Config) session.Repository {
	return &sessionRepository{
		conn: c,
		cfg:  cfg,
	}
}

// CreateSession in redis
func (r *sessionRepository) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	s := &session.Session{
		UserID: userID,
	}

	s.SessionID = uuid.New().String()
	sessionKey := r.newKey(s.SessionID)

	sessBytes, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	err = r.conn.Set(ctx, sessionKey, sessBytes, r.cfg.Session.Duration).Err()
	if err != nil {
		return "", err
	}
	return sessionKey, nil
}

// GetSessionByID in redis
func (r *sessionRepository) GetSessionByID(ctx context.Context, sessionID string) (*session.Session, error) {
	sessBytes, err := r.conn.Get(ctx, sessionID).Bytes()
	if err != nil {
		return nil, err
	}

	sess := &session.Session{}
	if err = json.Unmarshal(sessBytes, &sess); err != nil {
		return nil, err
	}
	return sess, nil
}

// DeleteByID in redis
func (r *sessionRepository) DeleteByID(ctx context.Context, sessionID string) error {
	return r.conn.Del(ctx, sessionID).Err()
}

func (r *sessionRepository) newKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", r.cfg.Session.BasePrefix, sessionID)
}
