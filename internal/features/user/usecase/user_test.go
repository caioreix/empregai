package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go-api/internal/core/user"
	usermock "go-api/internal/core/user/mocks"
	"go-api/internal/features/user/usecase"
	"go-api/pkg/config"
	"go-api/pkg/utils"
)

func TestUserUseCase_Register(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, mock, uc := setupTest(t)

		usr := &user.Raw{
			Name:     "Fake Name",
			Password: "fake_password",
			Email:    "fake@mail.com",
		}

		mock.On("Register", ctx, usr).
			Return(usr, nil).
			Once()

		got, err := uc.Register(ctx, usr)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

func TestUserUseCase_Login(t *testing.T) {
	t.Run("Success with sanitize data", func(t *testing.T) {
		ctx, mock, uc := setupTest(t)

		password := "fake_password"

		usr := &user.Raw{
			Password: password,
			Email:    "fake@mail.com",
		}

		err := usr.HashPassword()
		assert.NoError(t, err)

		mock.On("FindByEmail", ctx, usr.Email).
			Return(usr, nil).
			Once()

		got, err := uc.Login(ctx, usr.Email, password)
		assert.NoError(t, err, usr)
		assert.NotNil(t, got)
		assert.Empty(t, got.User.Password)
	})
}

func TestUserUseCase_Update(t *testing.T) {
	t.Run("Success with password", func(t *testing.T) {
		ctx, mock, uc := setupTest(t)

		usr := &user.Raw{
			Name:     "Fake Name",
			Password: "fake_password",
			Email:    "fake@mail.com",
		}

		mock.On("Update", ctx, usr).
			Return(usr, nil).
			Once()

		got, err := uc.Update(ctx, usr)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Empty(t, got.Password)
	})

	t.Run("Success with no password", func(t *testing.T) {
		ctx, mock, uc := setupTest(t)

		usr := &user.Raw{
			Name:  "Fake Name",
			Email: "fake@mail.com",
		}

		mock.On("Update", ctx, usr).
			Return(usr, nil).
			Once()

		got, err := uc.Update(ctx, usr)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Empty(t, got.Password)
	})
}

func TestUserUseCase_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, mock, uc := setupTest(t)

		id := uuid.New()

		mock.On("Delete", ctx, id).
			Return(nil).
			Once()

		err := uc.Delete(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("Fail", func(t *testing.T) {
		ctx, mock, uc := setupTest(t)

		id := uuid.New()

		mock.On("Delete", ctx, id).
			Return(errors.New("fake_err")).
			Once()

		err := uc.Delete(ctx, id)
		assert.Error(t, err)
	})
}

func TestUserUseCase_GetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, mock, uc := setupTest(t)

		usr := &user.Raw{
			ID:       uuid.New(),
			Name:     "Fake Name",
			Password: "fake_password",
			Email:    "fake@mail.com",
		}

		mock.On("GetByID", ctx, usr.ID).
			Return(usr, nil).
			Once()

		got, err := uc.GetByID(ctx, usr.ID)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Empty(t, got.Password)
	})
}

func TestUserUseCase_GetUsers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, mock, uc := setupTest(t)

		pq := &utils.PaginationQuery{
			OrderBy: "",
			Page:    0,
			Size:    10,
		}

		usr := &user.Raw{
			ID:       uuid.New(),
			Name:     "Fake Name",
			Password: "fake_password",
			Email:    "fake@mail.com",
		}

		list := &user.List{
			Users: &[]*user.Raw{
				{
					ID:       usr.ID,
					Name:     usr.Name,
					Password: "password",
					Email:    usr.Email,
				},
				{
					ID:       usr.ID,
					Name:     usr.Name,
					Password: "password",
					Email:    usr.Email,
				},
			},
		}

		mock.On("GetUsers", ctx, pq).
			Return(list, nil).
			Once()

		got, err := uc.GetUsers(ctx, pq)
		assert.NoError(t, err)
		users := *got.Users
		for _, u := range users {
			assert.Empty(t, u.Password)
		}
	})
}

func setupTest(t *testing.T) (context.Context, *usermock.Repository, user.UseCase) {
	t.Helper()

	cfg := &config.Config{
		Server: config.Server{
			JWTSecret: "fake_secret",
		},
	}

	repo := usermock.NewRepository(t)
	uc := usecase.NewUserUseCase(cfg, repo)

	return context.TODO(), repo, uc
}
