package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"go-api/internal/core/user"
	"go-api/pkg/utils"
)

type UserRepository struct {
	conn *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) user.Repository {
	return &UserRepository{
		conn: db,
	}
}

func (r *UserRepository) Register(ctx context.Context, usr *user.Model) (*user.Model, error) {
	u := &user.Model{}
	err := r.conn.QueryRowxContext(
		ctx,
		createUserQuery,
		usr.Email,
		usr.Password,
		usr.Role,
	).StructScan(u)

	return u, errors.Wrap(err, "UserRepository.Register.StructScan")
}

func (r *UserRepository) Update(ctx context.Context, usr *user.Model) (*user.Model, error) {
	u := &user.Model{}
	err := r.conn.GetContext(
		ctx,
		u,
		updateUserQuery,
		usr.ID,
		usr.Email,
		usr.Password,
		usr.Role,
	)

	return u, errors.Wrap(err, "UserRepository.Update.GetContext")
}

func (r *UserRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	res, err := r.conn.ExecContext(ctx, deleteUserQuery, userID)
	if err != nil {
		return errors.Wrap(err, "UserRepository.Delete.ExecContext")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "UserRepository.Delete.RowsAffected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "UserRepository.Delete.NoRows")
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID uuid.UUID) (*user.Model, error) {
	u := &user.Model{}
	err := r.conn.QueryRowxContext(
		ctx,
		getUserByIDQuery,
		userID,
	).StructScan(u)

	return u, errors.Wrap(err, "UserRepository.GetByID.StructScan")
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.Model, error) {
	u := &user.Model{}
	err := r.conn.QueryRowxContext(
		ctx,
		getUserByEmailQuery,
		email,
	).StructScan(u)

	return u, errors.Wrap(err, "UserRepository.FindByEmail.StructScan")
}

func (r *UserRepository) GetUsers(ctx context.Context, pagination *utils.PaginationQuery) (*user.List, error) {
	var totalCount int
	err := r.conn.GetContext(ctx, &totalCount, getUsersCountQuery)
	if err != nil {
		return nil, errors.Wrap(err, "UserRepository.GetUsers.GetContext")
	}

	users := make([]*user.Model, 0, pagination.GetSize())
	usersList := &user.List{
		TotalCount: totalCount,
		TotalPages: pagination.GetTotalPages(totalCount),
		Page:       pagination.GetPage(),
		Size:       pagination.GetSize(),
		HasMore:    pagination.GetHasMore(totalCount),
		Users:      &users,
	}

	if totalCount == 0 {
		return usersList, nil
	}

	err = r.conn.SelectContext(
		ctx,
		&users,
		getAllUsersQuery,
		pagination.GetOrderBy(),
		pagination.GetOffset(),
		pagination.GetLimit(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "UserRepository.GetUsers.SelectContext")
	}

	usersList.Users = &users
	return usersList, nil
}
