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
	err := isValidType(inputProduct.Type)
	if err != nil {
		return nil, err
	}
	if inputProduct.PvzId == nil {
		return nil, domain.ErrEmptyPvzID
	}
	return uc.productRepository.CreateProduct(inputProduct)
}

func (uc *ProductUseCase) DeleteLastProductForPVZ(pvzId *string) error {
	if pvzId == nil || *pvzId == "" {
		return domain.ErrEmptyPvzID
	}
	return uc.productRepository.DeleteLastProductForPVZ(pvzId)
}

func isValidType(productType *string) error {
	if productType == nil {
		return domain.ErrInvalidInputData
	}
	var allowedProductTypes = map[string]struct{}{
		"электроника": {},
		"одежда":      {},
		"обувь":       {},
	}
	_, ok := allowedProductTypes[*productType]
	if !ok {
		return domain.ErrInvalidProductType
	}
	return nil
}
