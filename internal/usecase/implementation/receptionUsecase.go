package implementation

import (
	"final/internal/domain"
	"final/internal/repository"
)

type ReceptionUseCase struct {
	receptionRepository repository.ReceptionRepository
}

func NewReceptionUseCase(receptionRepository repository.ReceptionRepository) *ReceptionUseCase {
	return &ReceptionUseCase{receptionRepository: receptionRepository}
}

func (uc *ReceptionUseCase) CreateReception(reception *domain.Reception) (*domain.Reception, error) {
	if reception.PvzId == nil {
		return nil, domain.ErrInvalidInputData
	}
	return uc.receptionRepository.CreateReception(reception)
}

func (uc *ReceptionUseCase) CloseReception(pvzID *string) (*domain.Reception, error) {
	if pvzID == nil || *pvzID == "" {
		return nil, domain.ErrInvalidInputData
	}
	return uc.receptionRepository.CloseReception(pvzID)
}
