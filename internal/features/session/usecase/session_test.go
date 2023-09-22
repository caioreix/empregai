package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go-api/internal/core/session"
	sessionmock "go-api/internal/core/session/mocks"
	"go-api/internal/features/session/usecase"
	"go-api/pkg/config"
)

func TestSessionUC_CreateSession(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repoMock, sessionUC := setupTest(t)
		ctx := context.TODO()
		sess := &session.Session{
			SessionID: uuid.NewString(),
			UserID:    uuid.New(),
		}

		repoMock.On("CreateSession", ctx, sess.UserID).
			Return(sess.SessionID, nil).
			Once()

		got, err := sessionUC.CreateSession(ctx, sess.UserID)
		assert.NoError(t, err)
		assert.Equal(t, sess.SessionID, got)
	})

	t.Run("Fail", func(t *testing.T) {
		repoMock, sessionUC := setupTest(t)
		ctx := context.TODO()
		sess := &session.Session{
			SessionID: uuid.NewString(),
			UserID:    uuid.New(),
		}

		repoMock.On("CreateSession", ctx, sess.UserID).
			Return("", errors.New("fake_error")).
			Once()

		got, err := sessionUC.CreateSession(ctx, sess.UserID)
		assert.Error(t, err)
		assert.Empty(t, got)
	})
}

func TestSessionUC_GetSessionByID(t *testing.T) {
	repoMock, sessionUC := setupTest(t)

	t.Run("Success", func(t *testing.T) {
		ctx := context.TODO()
		sess := &session.Session{
			SessionID: uuid.NewString(),
			UserID:    uuid.New(),
		}

		repoMock.On("GetSessionByID", ctx, sess.SessionID).
			Return(sess, nil).
			Once()

		got, err := sessionUC.GetSessionByID(ctx, sess.SessionID)
		assert.NoError(t, err)
		assert.Equal(t, sess, got)
	})

	t.Run("Fail", func(t *testing.T) {
		ctx := context.TODO()
		sess := &session.Session{
			SessionID: uuid.NewString(),
			UserID:    uuid.New(),
		}

		repoMock.On("GetSessionByID", ctx, sess.SessionID).
			Return(nil, errors.New("fake_error")).
			Once()

		got, err := sessionUC.GetSessionByID(ctx, sess.SessionID)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestSessionUC_DeleteByID(t *testing.T) {
	repoMock, sessionUC := setupTest(t)

	t.Run("Success", func(t *testing.T) {
		ctx := context.TODO()
		sess := &session.Session{
			SessionID: uuid.NewString(),
			UserID:    uuid.New(),
		}

		repoMock.On("DeleteByID", ctx, sess.SessionID).
			Return(nil).
			Once()

		err := sessionUC.DeleteByID(ctx, sess.SessionID)
		assert.NoError(t, err)
	})

	t.Run("Fail", func(t *testing.T) {
		ctx := context.TODO()
		sess := &session.Session{
			SessionID: uuid.NewString(),
			UserID:    uuid.New(),
		}

		repoMock.On("DeleteByID", ctx, sess.SessionID).
			Return(errors.New("fake_error")).
			Once()

		err := sessionUC.DeleteByID(ctx, sess.SessionID)
		assert.Error(t, err)
	})
}

func setupTest(t *testing.T) (*sessionmock.Repository, session.UseCase) {
	t.Helper()

	mockRepo := sessionmock.NewRepository(t)
	cfg := &config.Config{}
	sessUC := usecase.NewSessionUseCase(mockRepo, cfg)

	return mockRepo, sessUC
}
