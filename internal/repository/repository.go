package repository

import (
	"final/internal/domain"
	"final/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
	"time"
)

type Repository struct {
	PvzRepository       PVZRepository
	ReceptionRepository ReceptionRepository
	ProductRepository   ProductRepository
	UserRepository      UserRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	receptionRepository := postgres.NewReceptionRepository(db)
	return &Repository{
		PvzRepository:       postgres.NewPvzRepository(db),
		ReceptionRepository: receptionRepository,
		ProductRepository:   postgres.NewProductRepository(db, receptionRepository),
		UserRepository:      postgres.NewUserRepository(db),
	}
}

type PVZRepository interface {
	CreatePVZ(*domain.PVZ) (*domain.PVZ, error)
	GetPVZInfo(*time.Time, *time.Time, int, int) ([]domain.PVZWithReceptions, error)
}

type ReceptionRepository interface {
	CreateReception(*domain.Reception) (*domain.Reception, error)
	CloseReception(*string) (*domain.Reception, error)
}

type ProductRepository interface {
	CreateProduct(*domain.InputProduct) (*domain.Product, error)
	DeleteLastProductForPVZ(*string) error
}

type UserRepository interface {
	Register(*domain.User) (*domain.User, error)
	GetUserByEmail(string) (*domain.User, error)
}
