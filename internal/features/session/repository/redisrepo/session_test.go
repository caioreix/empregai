package redisrepo_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-api/internal/core/session"
	"go-api/internal/features/session/repository/redisrepo"
	"go-api/pkg/config"
)

func TestSessionRepository_CreateSession(t *testing.T) {
	sessRepository := setupTest(t)

	t.Run("Success", func(t *testing.T) {
		userUUID := uuid.New()

		s, err := sessRepository.CreateSession(context.Background(), userUUID)
		assert.NoError(t, err)
		assert.NotEmpty(t, s)
	})
}

func TestSessionRepository_GetSessionByID(t *testing.T) {
	sessRepository := setupTest(t)

	t.Run("Success", func(t *testing.T) {
		userUUID := uuid.New()

		createdSess, err := sessRepository.CreateSession(context.Background(), userUUID)
		assert.NoError(t, err)
		assert.NotEqual(t, createdSess, "")

		s, err := sessRepository.GetSessionByID(context.Background(), createdSess)
		assert.NoError(t, err)
		assert.NotEqual(t, s, "")
	})
}

func TestSessionRepository_DeleteByID(t *testing.T) {
	sessRepository := setupTest(t)

	t.Run("Success", func(t *testing.T) {
		userUUID := uuid.New()

		err := sessRepository.DeleteByID(context.Background(), userUUID.String())
		assert.NoError(t, err)
	})
}

func setupTest(t *testing.T) session.Repository {
	t.Helper()

	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	cfg := &config.Config{Session: config.Session{
		BasePrefix: "api-session:",
		Duration:   10,
	}}

	sessRepository := redisrepo.NewSessionRepository(client, cfg)
	return sessRepository
}
