package usecase

import (
	"final/internal/domain"
	"final/internal/repository"
	"final/internal/usecase/implementation"
)

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
	CreatePVZ(*domain.PVZ) (*domain.PVZ, error)
	GetPVZInfo(string, string, string, string) ([]domain.PVZWithReceptions, error)
}

type ReceptionUsecase interface {
	CreateReception(*domain.Reception) (*domain.Reception, error)
	CloseReception(*string) (*domain.Reception, error)
}

type ProductUsecase interface {
	CreateProduct(*domain.InputProduct) (*domain.Product, error)
	DeleteLastProductForPVZ(*string) error
}

type UserUsecase interface {
	Register(*domain.InputUser) (*domain.User, error)
	Login(*domain.InputUser) (*domain.User, error)
}
