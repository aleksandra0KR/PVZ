package implementation

import (
	"final/internal/domain"
	"final/internal/repository"
)

type ProductUseCase struct {
	productRepository repository.ProductRepository
}

func NewProductUseCase(productRepository repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{productRepository: productRepository}
}

func (uc *ProductUseCase) CreateProduct(inputProduct *domain.InputProduct) (*domain.Product, error) {
	return uc.productRepository.CreateProduct(inputProduct)
}

func (uc *ProductUseCase) DeleteLastProductForPVZ(pvzId *string) error {
	return uc.productRepository.DeleteLastProductForPVZ(pvzId)
}
