package implementation

import (
	"final/internal/domain"
	"final/internal/repository"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type PVZUseCase struct {
	pvzRepository repository.PVZRepository
}

func NewPVZUseCase(pvzRepository repository.PVZRepository) *PVZUseCase {
	return &PVZUseCase{pvzRepository: pvzRepository}
}

func (uc *PVZUseCase) CreatePVZ(pvz *domain.PVZ) (*domain.PVZ, error) {
	err := isValidCity(pvz.City)
	if err != nil {
		return nil, err
	}
	return uc.pvzRepository.CreatePVZ(pvz)
}

func (uc *PVZUseCase) GetPVZInfo(startDateStr, endDateStr, pageStr, limitStr string) ([]domain.PVZWithReceptions, error) {
	var startDate, endDate *time.Time
	if startDateStr != "" {
		startDateTime, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			log.Error(err)
			return nil, domain.ErrGetPVZInfo
		}
		startDate = &startDateTime
	}

	if endDateStr != "" {
		endDateTime, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			log.Error(err)
			return nil, domain.ErrGetPVZInfo
		}
		endDate = &endDateTime
	}

	page := 1
	if pageStr != "" {
		val, err := strconv.Atoi(pageStr)
		if err != nil {
			log.Error(err)
			return nil, domain.ErrGetPVZInfo
		}
		page = val
	}

	limit := 10
	if limitStr != "" {
		val, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Error(err)
			return nil, domain.ErrGetPVZInfo
		}
		limit = val
	}

	offset := (page - 1) * limit

	pvzInfo, err := uc.pvzRepository.GetPVZInfo(startDate, endDate, offset, limit)
	if err != nil {
		return nil, domain.ErrGetPVZInfo
	}
	return pvzInfo, nil
}

func isValidCity(city *string) error {
	if city == nil {
		return domain.ErrInvalidInputData
	}
	var allowedCities = map[string]struct{}{
		"Москва":          {},
		"Санкт-Петербург": {},
		"Казань":          {},
	}
	_, ok := allowedCities[*city]
	if !ok {
		return domain.ErrInvalidCity
	}
	return nil
}
