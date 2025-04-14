package implementation

import (
	"context"
	"final/internal/domain"
	"final/internal/repository"
)

type ReceptionUseCase struct {
	receptionRepository repository.ReceptionRepository
}

func NewReceptionUseCase(receptionRepository repository.ReceptionRepository) *ReceptionUseCase {
	return &ReceptionUseCase{receptionRepository: receptionRepository}
}

func (uc *ReceptionUseCase) CreateReception(ctx context.Context, reception *domain.Reception) (*domain.Reception, error) {
	if reception.PvzId == nil {
		return nil, domain.ErrInvalidInputData
	}
	return uc.receptionRepository.CreateReception(ctx, reception)
}

func (uc *ReceptionUseCase) CloseReception(ctx context.Context, pvzID *string) (*domain.Reception, error) {
	if pvzID == nil || *pvzID == "" {
		return nil, domain.ErrInvalidInputData
	}
	return uc.receptionRepository.CloseReception(ctx, pvzID)
}
