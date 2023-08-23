package session

import (
	"context"

	"github.com/google/uuid"
)

// Repository session interface
type Repository interface {
	CreateSession(ctx context.Context, userID uuid.UUID) (string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*Session, error)
	DeleteByID(ctx context.Context, sessionID string) error
}

// UseCase session interface
type UseCase interface {
	CreateSession(ctx context.Context, userID uuid.UUID) (string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*Session, error)
	DeleteByID(ctx context.Context, sessionID string) error
}
