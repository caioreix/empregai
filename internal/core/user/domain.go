package user

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-api/pkg/utils"
)

type Handlers interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
	Logout() gin.HandlerFunc
	Update() gin.HandlerFunc
	Delete() gin.HandlerFunc
	GetUserByID() gin.HandlerFunc
	GetUsers() gin.HandlerFunc
	GetMe() gin.HandlerFunc
}

type Repository interface {
	Register(ctx context.Context, user *Model) (*Model, error)
	Update(ctx context.Context, user *Model) (*Model, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID uuid.UUID) (*Model, error)
	FindByEmail(ctx context.Context, email string) (*Model, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*List, error)
}

type UseCase interface {
	Register(ctx context.Context, user *Model) (*Token, error)
	Login(ctx context.Context, email, password string) (*Token, error)
	Update(ctx context.Context, user *Model) (*Model, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID uuid.UUID) (*Model, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*List, error)
}
