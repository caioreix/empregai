package user

import (
	"context"
	"time"

	"github.com/google/uuid"

	"go-api/pkg/apierrors"
	"go-api/pkg/utils"
)

// Raw model store user data
type Raw struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	LastLogin time.Time `json:"last_login" db:"last_login"`
}

// List model store user pages
type List struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Users      *[]*Raw `json:"users"`
}

// Token model store user token
type Token struct {
	User  *Raw   `json:"user"`
	Token string `json:"token"`
}

// CtxKey is a key used for the User object in the context
type CtxKey struct{}

func (u *Raw) HashPassword() error {
	hash, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = hash

	return nil
}

func (u *Raw) ComparePassword(password string) bool {
	return utils.ComparePassword(u.Password, password)
}

func (u *Raw) Sanitize() {
	u.Password = ""
}

// Get user from context
func GetUserFromCtx(ctx context.Context) (*Raw, error) {
	user, ok := ctx.Value(CtxKey{}).(*Raw)
	if !ok {
		return nil, apierrors.Unauthorized()
	}

	return user, nil
}
