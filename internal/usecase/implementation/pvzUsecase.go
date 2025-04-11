package implementation

import (
	"errors"
	"final/internal/domain"
	"final/internal/repository"
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
	if !isValidCity(*pvz.City) {
		return nil, errors.New("invalid city")
	}
	return uc.pvzRepository.CreatePVZ(pvz)
}

func isValidCity(city string) bool {
	var allowedCities = map[string]struct{}{
		"Москва":          {},
		"Санкт-Петербург": {},
		"Казань":          {},
	}
	_, ok := allowedCities[city]
	return ok
}

func (uc *PVZUseCase) GetPVZInfo(startDateStr, endDateStr, pageStr, limitStr string) ([]domain.PVZWithReceptions, error) {
	var startDate, endDate *time.Time
	if startDateStr != "" {
		startDateTime, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return nil, err
		}
		startDate = &startDateTime
	}

	if endDateStr != "" {
		endDateTime, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return nil, err
		}
		endDate = &endDateTime
	}

	page := 1
	if pageStr != "" {
		val, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, err
		}
		page = val
	}

	limit := 10
	if limitStr != "" {
		val, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, err
		}
		limit = val
	}

	offset := (page - 1) * limit

	pvzInfo, err := uc.pvzRepository.GetPVZInfo(startDate, endDate, offset, limit)
	if err != nil {
		return nil, err
	}
	return pvzInfo, nil
}
