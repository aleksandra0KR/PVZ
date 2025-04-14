package usecase

import (
	"context"
	"final/internal/domain"
	"final/internal/repository"
	"final/internal/usecase/implementation"
)

//go:generate mockgen -source=usecase.go -destination=mocks/mock.go
type Usecase struct {
	PVZUsecase       PVZUsecase
	ReceptionUsecase ReceptionUsecase
	ProductUsecase   ProductUsecase
	UserUsecase      UserUsecase
}

func NewUseCase(repository *repository.Repository) *Usecase {
	return &Usecase{
		PVZUsecase:       implementation.NewPVZUseCase(repository.PvzRepository),
		ReceptionUsecase: implementation.NewReceptionUseCase(repository.ReceptionRepository),
		ProductUsecase:   implementation.NewProductUseCase(repository.ProductRepository),
		UserUsecase:      implementation.NewUserUseCase(repository.UserRepository),
	}
}

type PVZUsecase interface {
	CreatePVZ(context.Context, *domain.PVZ) (*domain.PVZ, error)
	GetPVZInfo(context.Context, string, string, string, string) ([]domain.PVZWithReceptions, error)
}

type ReceptionUsecase interface {
	CreateReception(context.Context, *domain.Reception) (*domain.Reception, error)
	CloseReception(context.Context, *string) (*domain.Reception, error)
}

type ProductUsecase interface {
	CreateProduct(context.Context, *domain.InputProduct) (*domain.Product, error)
	DeleteLastProductForPVZ(context.Context, *string) error
}

type UserUsecase interface {
	Register(context.Context, *domain.InputUser) (*domain.User, error)
	Login(context.Context, *domain.InputUser) (*domain.User, error)
}
