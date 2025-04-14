package postgres

import (
	"context"
	"database/sql"
	"final/internal/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegister(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Error(err)
		}
	}(db)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	email := "test"
	password := "password"
	role := "employee"
	userId := "1"

	ctx := context.Background()
	t.Run("Successful_Registration", func(t *testing.T) {
		inputUser := &domain.User{
			Email:    &email,
			Password: &password,
			Role:     &role,
		}

		mock.ExpectQuery("SELECT id, email, password, role FROM users").
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "role"}))

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(email, password, role).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userId))

		user, err := repo.Register(ctx, inputUser)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userId, *user.ID)
		assert.Equal(t, email, *user.Email)
		assert.Equal(t, password, *user.Password)
		assert.Equal(t, role, *user.Role)
	})

	t.Run("Failure_User_Already_Exists", func(t *testing.T) {
		inputUser := &domain.User{
			Email:    &email,
			Password: &password,
			Role:     &role,
		}

		mock.ExpectQuery("SELECT id, email, password, role FROM users").
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "role"}).AddRow(userId, email, password, role))

		user, err := repo.Register(ctx, inputUser)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrUserWithSuchEmailAlreadyExists, err)
		assert.Nil(t, user)
	})

	t.Run("Failure_DB_Error", func(t *testing.T) {
		inputUser := &domain.User{
			Email:    &email,
			Password: &password,
			Role:     &role,
		}

		mock.ExpectQuery("SELECT id, email, password, role FROM users").
			WithArgs(email).
			WillReturnError(sql.ErrConnDone)

		user, err := repo.Register(ctx, inputUser)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrCreateUser, err)
		assert.Nil(t, user)
	})
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Error(err)
		}
	}(db)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	email := "test"
	password := "password"
	role := "user"
	userId := "1"

	ctx := context.Background()
	t.Run("Successful_Get_User_By_Email", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, password, role FROM users").
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "role"}).AddRow(userId, email, password, role))

		user, err := repo.GetUserByEmail(ctx, email)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userId, *user.ID)
		assert.Equal(t, email, *user.Email)
		assert.Equal(t, password, *user.Password)
		assert.Equal(t, role, *user.Role)
	})

	t.Run("Failure_User_Not_Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, password, role FROM users").
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByEmail(ctx, email)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrFindUser, err)
		assert.Nil(t, user)
	})

	t.Run("Failure_DB_Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, password, role FROM users").
			WithArgs(email).
			WillReturnError(sql.ErrConnDone)

		user, err := repo.GetUserByEmail(ctx, email)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrFindUser, err)
		assert.Nil(t, user)
	})
}
