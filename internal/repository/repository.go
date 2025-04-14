package repository

import (
	"context"
	"final/internal/domain"
	"final/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
	"time"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

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
	CreatePVZ(context.Context, *domain.PVZ) (*domain.PVZ, error)
	GetPVZInfo(context.Context, *time.Time, *time.Time, int, int) ([]domain.PVZWithReceptions, error)
	GetAllPVZ(context.Context) ([]*domain.PVZ, error)
}

type ReceptionRepository interface {
	CreateReception(context.Context, *domain.Reception) (*domain.Reception, error)
	CloseReception(context.Context, *string) (*domain.Reception, error)
}

type ProductRepository interface {
	CreateProduct(context.Context, *domain.InputProduct) (*domain.Product, error)
	DeleteLastProductForPVZ(context.Context, *string) error
}

type UserRepository interface {
	Register(context.Context, *domain.User) (*domain.User, error)
	GetUserByEmail(context.Context, string) (*domain.User, error)
}
