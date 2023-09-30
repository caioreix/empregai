package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"go-api/internal/core/user"
	"go-api/pkg/apierrors"
	"go-api/pkg/config"
	"go-api/pkg/token"
	"go-api/pkg/utils"
)

type userUseCase struct {
	cfg  *config.Config
	repo user.Repository
}

func NewUserUseCase(cfg *config.Config, repo user.Repository) user.UseCase {
	return &userUseCase{cfg: cfg, repo: repo}
}

func (uc *userUseCase) Register(ctx context.Context, usr *user.Model) (*user.Token, error) {
	err := usr.HashPassword()
	if err != nil {
		return nil, err
	}

	createdUser, err := uc.repo.Register(ctx, usr)
	if err != nil {
		return nil, err
	}
	createdUser.Sanitize()

	tokenDuration := 60 * time.Minute
	jwt, err := token.GenerateJWT(createdUser.Email, createdUser.ID.String(), tokenDuration, uc.cfg)
	if err != nil {
		return nil, err
	}

	return &user.Token{
		User:  createdUser,
		Token: jwt,
	}, nil
}

func (uc *userUseCase) Login(ctx context.Context, email, password string) (*user.Token, error) {
	foundUser, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !foundUser.ComparePassword(password) {
		return nil, apierrors.Unauthorized()
	}

	foundUser.Sanitize()

	tokenDuration := 60 * time.Minute
	jwt, err := token.GenerateJWT(foundUser.Email, foundUser.ID.String(), tokenDuration, uc.cfg)
	if err != nil {
		return nil, err
	}

	return &user.Token{
		User:  foundUser,
		Token: jwt,
	}, nil
}

func (uc *userUseCase) Update(ctx context.Context, usr *user.Model) (*user.Model, error) {
	if usr.Password != "" {
		err := usr.HashPassword()
		if err != nil {
			return nil, err
		}
	}

	updatedUser, err := uc.repo.Update(ctx, usr)
	if err != nil {
		return nil, err
	}

	updatedUser.Sanitize()

	return updatedUser, nil
}

func (uc *userUseCase) Delete(ctx context.Context, userID uuid.UUID) error {
	err := uc.repo.Delete(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *userUseCase) GetByID(ctx context.Context, userID uuid.UUID) (*user.Model, error) {
	usr, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	usr.Sanitize()

	return usr, nil
}

func (uc *userUseCase) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*user.List, error) {
	list, err := uc.repo.GetUsers(ctx, pq)
	if err != nil {
		return nil, err
	}

	users := *list.Users

	for _, usr := range users {
		usr.Sanitize()
	}

	return list, nil
}
