package http_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sessionmock "go-api/internal/core/session/mocks"
	"go-api/internal/core/user"
	usermock "go-api/internal/core/user/mocks"
	userhttp "go-api/internal/features/user/delivery/http"
	"go-api/pkg/config"
	"go-api/pkg/token"
)

func TestUserHandler_Register(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			cfg, ctx, rw, userUC, sessionUC, h = setupTest(t)
			usr                                = &user.Raw{
				Name:     "Fake Name",
				Email:    "fake@mail.com",
				Password: "fake_password",
			}
			userID   = uuid.New()
			usrToken = &user.Token{
				User: &user.Raw{
					ID:        userID,
					Name:      usr.Name,
					Email:     usr.Email,
					Password:  "",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					LastLogin: time.Now(),
				},
			}
			sessionID = "fake_session_id"
		)

		jwt, err := token.GenerateJWT(usr.Email, userID.String(), 60*time.Minute, cfg)
		assert.NoError(t, err)
		usrToken.Token = jwt

		setupRequest(t, ctx, http.MethodPost, usr)

		userUC.On("Register", ctx, usr).
			Return(usrToken, nil).
			Once()

		sessionUC.On("CreateSession", ctx, userID).
			Return(sessionID, nil).
			Once()

		handlerFunc := h.Register()
		handlerFunc(ctx)

		assert.Equal(t, http.StatusCreated, rw.Result().StatusCode)
		plainToken, err := json.Marshal(usrToken)
		assert.NoError(t, err)
		assert.Equal(t, string(plainToken), rw.Body.String())
	})
}

func setupTest(t *testing.T) (*config.Config, *gin.Context, *httptest.ResponseRecorder, *usermock.UseCase, *sessionmock.UseCase, user.Handlers) {
	t.Helper()

	userUC := usermock.NewUseCase(t)
	sessionUC := sessionmock.NewUseCase(t)

	cfg := &config.Config{
		Server: config.Server{
			JWTSecret: "fake_secret",
		},
		Session: config.Session{
			Duration: 10 * time.Second,
		},
	}

	userHandlers := userhttp.NewUserHandler(cfg, userUC, sessionUC)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	return cfg, ctx, w, userUC, sessionUC, userHandlers
}

func setupRequest(t *testing.T, ctx *gin.Context, method string, v any) {
	t.Helper()

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(v)
	require.NoError(t, err)
	require.NotNil(t, buf)

	ctx.Request = &http.Request{
		Method: method,
		Header: make(http.Header),
	}

	switch method {
	case http.MethodPost:
		ctx.Request.Header.Set("Content-Type", "application/json")
	}

	jsonValue, err := json.Marshal(v)
	require.NoError(t, err)

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonValue))
}
