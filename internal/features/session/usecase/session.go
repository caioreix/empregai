package usecase

import (
	"context"

	"github.com/google/uuid"

	"go-api/internal/core/session"
	"go-api/pkg/config"
)

// sessionUC struct
type sessionUC struct {
	repo session.Repository
	cfg  *config.Config
}

// NewSessionUseCase constructor
func NewSessionUseCase(sessionRepo session.Repository, cfg *config.Config) session.UseCase {
	return &sessionUC{repo: sessionRepo, cfg: cfg}
}

// CreateSession usecase
func (u *sessionUC) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	return u.repo.CreateSession(ctx, userID)
}

// DeleteByID usecase
func (u *sessionUC) DeleteByID(ctx context.Context, sessionID string) error {
	return u.repo.DeleteByID(ctx, sessionID)
}

// GetSessionByID usecase
func (u *sessionUC) GetSessionByID(ctx context.Context, sessionID string) (*session.Session, error) {
	return u.repo.GetSessionByID(ctx, sessionID)
}
