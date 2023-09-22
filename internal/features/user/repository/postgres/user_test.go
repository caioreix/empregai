package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"go-api/internal/core/user"
	"go-api/internal/features/user/repository/postgres"
	"go-api/pkg/utils"
)

func TestUserRepository_Register(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, repo, mock, want, rows := setupTest(t)
		defer db.Close()

		mock.ExpectQuery(createUserQuery).
			WithArgs(want.Name, want.Email, want.Password).
			WillReturnRows(rows)

		got, err := repo.Register(context.TODO(), want)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestUserRepository_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, repo, mock, want, rows := setupTest(t)
		defer db.Close()

		mock.ExpectQuery(updateUserQuery).
			WithArgs(want.ID, want.Name, want.Email, want.Password).
			WillReturnRows(rows)

		got, err := repo.Update(context.TODO(), want)
		assert.NoError(t, err)
		assert.Equal(t, got, want)
	})
}

func TestUserRepository_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, repo, mock, want, _ := setupTest(t)
		defer db.Close()

		mock.ExpectExec(deleteUserQuery).
			WithArgs(want.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(context.TODO(), want.ID)
		assert.NoError(t, err)
	})

	t.Run("Fail with no rows", func(t *testing.T) {
		db, repo, mock, want, _ := setupTest(t)
		defer db.Close()

		mock.ExpectExec(deleteUserQuery).
			WithArgs(want.ID).
			WillReturnResult(sqlmock.NewResult(1, 0))

		err := repo.Delete(context.TODO(), want.ID)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, repo, mock, want, rows := setupTest(t)
		defer db.Close()

		mock.ExpectQuery(getUserByIDQuery).
			WithArgs(want.ID).
			WillReturnRows(rows)

		got, err := repo.GetByID(context.TODO(), want.ID)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, repo, mock, want, rows := setupTest(t)
		defer db.Close()

		mock.ExpectQuery(getUserByEmailQuery).
			WithArgs(want.Email).
			WillReturnRows(rows)

		got, err := repo.FindByEmail(context.TODO(), want.Email)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestUserRepository_GetUsers(t *testing.T) {
	pag := &utils.PaginationQuery{
		OrderBy: "",
		Page:    0,
		Size:    10,
	}

	t.Run("Success with users", func(t *testing.T) {
		db, repo, mock, _, rows := setupTest(t)
		defer db.Close()

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(getUsersCountQuery).WillReturnRows(totalRows)
		mock.ExpectQuery(getAllUsersQuery).
			WithArgs(pag.OrderBy, pag.Page, pag.Size).
			WillReturnRows(rows)

		got, err := repo.GetUsers(context.TODO(), pag)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(*got.Users))
	})

	t.Run("Success with no users", func(t *testing.T) {
		db, repo, mock, _, rows := setupTest(t)
		defer db.Close()

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(getUsersCountQuery).WillReturnRows(totalRows)
		mock.ExpectQuery(getAllUsersQuery).
			WithArgs(pag.OrderBy, pag.Page, pag.Size).
			WillReturnRows(rows)

		got, err := repo.GetUsers(context.TODO(), pag)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(*got.Users))
	})
}

func setupTest(t *testing.T) (*sql.DB, user.Repository, sqlmock.Sqlmock, *user.Raw, *sqlmock.Rows) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	dbx := sqlx.NewDb(db, "sqlmock")
	repo := postgres.NewUserRepository(dbx)

	want := &user.Raw{
		ID:        uuid.New(),
		Name:      "Fake Name",
		Email:     "fake@mail.com",
		Password:  "fake_password",
		LastLogin: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"email",
		"password",
		"created_at",
		"updated_at",
		"last_login",
	}).AddRow(
		want.ID,
		want.Name,
		want.Email,
		want.Password,
		want.CreatedAt,
		want.UpdatedAt,
		want.LastLogin,
	)

	return db, repo, mock, want, rows
}
