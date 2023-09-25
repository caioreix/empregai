package middleware

import (
	"go-api/internal/core/session"
	"go-api/internal/core/user"
	"go-api/pkg/config"
	"go-api/pkg/logger"
)

type Manager struct {
	cfg       *config.Config
	log       logger.Logger
	userUC    user.UseCase
	sessionUC session.UseCase
}

func New(
	cfg *config.Config,
	logger logger.Logger,
	userUC user.UseCase,
	sessionUC session.UseCase,
) *Manager {
	return &Manager{
		cfg:       cfg,
		log:       logger,
		userUC:    userUC,
		sessionUC: sessionUC,
	}
}
