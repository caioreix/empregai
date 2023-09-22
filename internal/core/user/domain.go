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
	// FindByName() gin.HandlerFunc
	GetUsers() gin.HandlerFunc
	GetMe() gin.HandlerFunc
	// GetCSRFToken() gin.HandlerFunc
}

type Repository interface {
	Register(ctx context.Context, user *Raw) (*Raw, error)
	Update(ctx context.Context, user *Raw) (*Raw, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID uuid.UUID) (*Raw, error)
	// FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*List, error)
	FindByEmail(ctx context.Context, email string) (*Raw, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*List, error)
}

type UseCase interface {
	Register(ctx context.Context, user *Raw) (*Token, error)
	Login(ctx context.Context, email, password string) (*Token, error)
	Update(ctx context.Context, user *Raw) (*Raw, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID uuid.UUID) (*Raw, error)
	// FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*List, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*List, error)
}
