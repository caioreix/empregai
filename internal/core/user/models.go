package user

import (
	"context"
	"time"

	"github.com/google/uuid"

	"go-api/pkg/apierrors"
	"go-api/pkg/utils"
)

// Model model store user data
type Model struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	LastLogin time.Time `json:"last_login" db:"last_login"`
}

// List model store user pages
type List struct {
	TotalCount int       `json:"total_count"`
	TotalPages int       `json:"total_pages"`
	Page       int       `json:"page"`
	Size       int       `json:"size"`
	HasMore    bool      `json:"has_more"`
	Users      *[]*Model `json:"users"`
}

// Token model store user token
type Token struct {
	User  *Model `json:"user"`
	Token string `json:"token"`
}

// CtxKey is a key used for the User object in the context
type CtxKey struct{}

func (u *Model) HashPassword() error {
	hash, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = hash

	return nil
}

func (u *Model) ComparePassword(password string) bool {
	return utils.ComparePassword(u.Password, password)
}

func (u *Model) Sanitize() {
	u.Password = ""
}

// Get user from context
func GetUserFromCtx(ctx context.Context) (*Model, error) {
	user, ok := ctx.Value(CtxKey{}).(*Model)
	if !ok {
		return nil, apierrors.Unauthorized()
	}

	return user, nil
}
