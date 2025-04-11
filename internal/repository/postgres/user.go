package postgres

import (
	"database/sql"
	"errors"
	"final/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Register(user *domain.User) (*domain.User, error) {
	userGetByEmail, err := r.GetUserByEmail(*user.Email)
	if err != nil {
		return nil, err
	} else if userGetByEmail != nil {
		return nil, errors.New("user with such email already exists")
	}

	query := `INSERT INTO users (email, password, role)
              VALUES ($1, $2, $3)
              RETURNING id`
	row := r.db.QueryRow(query, user.Email, user.Password, user.Role)
	err = row.Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	query := `SELECT id, email, password, role FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)
	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}
