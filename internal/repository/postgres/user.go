package postgres

import (
	"final/internal/domain"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Register(user *domain.User) (*domain.User, error) {
	userGetByEmail, err := r.GetUserByEmail(*user.Email)
	if userGetByEmail != nil && err != nil {
		return nil, domain.ErrCreateUser
	} else if userGetByEmail != nil {
		return nil, domain.ErrUserWithSuchEmailAlreadyExists
	}

	query := `INSERT INTO users (email, password, role)
              VALUES ($1, $2, $3)
              RETURNING id`
	row := r.db.QueryRow(query, user.Email, user.Password, user.Role)
	err = row.Scan(&user.ID)
	if err != nil {
		log.Error(err)
		return nil, domain.ErrCreateUser
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	query := `SELECT id, email, password, role FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		log.Error(err)
		return nil, domain.ErrFindUser
	}
	return &user, nil
}
