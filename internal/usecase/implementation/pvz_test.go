package implementation

import (
	"context"
	"final/internal/domain"
	"final/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreatePVZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockPVZRepository(ctrl)
	useCase := NewPVZUseCase(mockRepo)
	allowedCity := "Москва"
	wrongCity := "Париж"
	id := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	layout := "2006-01-02T15:04:05.000Z"
	str := "2025-04-12T20:10:25.102Z"
	dateTime, err := time.Parse(layout, str)
	if err != nil {
		log.Error(err)
	}

	ctx := context.Background()
	t.Run("Successful_Create_Pvz", func(t *testing.T) {
		pvz := &domain.PVZ{
			ID:               &id,
			RegistrationDate: &dateTime,
			City:             &allowedCity,
		}
		outputPVZ := &domain.PVZ{City: &allowedCity, ID: &id, RegistrationDate: &dateTime}

		mockRepo.EXPECT().CreatePVZ(ctx, pvz).Return(pvz, nil)
		result, err := useCase.CreatePVZ(ctx, pvz)
		assert.NoError(t, err)
		assert.Equal(t, outputPVZ, result)
	})

	t.Run("Failure_Invalid_City", func(t *testing.T) {
		pvz := &domain.PVZ{
			ID:               &id,
			RegistrationDate: &dateTime,
			City:             &wrongCity,
		}

		result, err := useCase.CreatePVZ(ctx, pvz)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidCity, err)
		assert.Nil(t, result)
	})

	t.Run("Failure_Nil_City", func(t *testing.T) {
		pvz := &domain.PVZ{
			ID:               &id,
			RegistrationDate: &dateTime,
		}

		result, err := useCase.CreatePVZ(ctx, pvz)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidInputData, err)
		assert.Nil(t, result)
	})
}

func TestGetPVZInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockPVZRepository(ctrl)
	useCase := NewPVZUseCase(mockRepo)

	allowedCity := "Москва"
	id := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	layout := "2006-01-02T15:04:05.000Z"
	str := "2025-04-12T20:10:25.102Z"
	dateTime, err := time.Parse(layout, str)
	if err != nil {
		log.Error(err)
	}

	ctx := context.Background()
	t.Run("Successful_Get_Pvz_Info", func(t *testing.T) {
		startDateStr := "2025-01-01T00:00:00Z"
		endDateStr := "2025-12-31T23:59:59Z"
		pageStr := "1"
		limitStr := "10"

		startDate, _ := time.Parse(time.RFC3339, startDateStr)
		endDate, _ := time.Parse(time.RFC3339, endDateStr)

		pvzInfo := []domain.PVZWithReceptions{
			{PVZ: domain.PVZ{
				ID:               &id,
				RegistrationDate: &dateTime,
				City:             &allowedCity,
			}},
			{PVZ: domain.PVZ{
				ID:               &id,
				RegistrationDate: &dateTime,
				City:             &allowedCity,
			}},
		}

		mockRepo.EXPECT().GetPVZInfo(ctx, &startDate, &endDate, 0, 10).Return(pvzInfo, nil)

		result, err := useCase.GetPVZInfo(ctx, startDateStr, endDateStr, pageStr, limitStr)
		assert.NoError(t, err)
		assert.Equal(t, pvzInfo, result)
	})

	t.Run("Failure_Invalid_Start_Date", func(t *testing.T) {
		startDateStr := "invalid-date"
		endDateStr := "2023-12-31T23:59:59Z"
		pageStr := "1"
		limitStr := "10"

		result, err := useCase.GetPVZInfo(ctx, startDateStr, endDateStr, pageStr, limitStr)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrGetPVZInfo, err)
		assert.Nil(t, result)
	})

	t.Run("Failure_Invalid_Invalid_End_Date", func(t *testing.T) {
		startDateStr := "2023-01-01T00:00:00Z"
		endDateStr := "invalid-date"
		pageStr := "1"
		limitStr := "10"

		result, err := useCase.GetPVZInfo(ctx, startDateStr, endDateStr, pageStr, limitStr)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrGetPVZInfo, err)
		assert.Nil(t, result)
	})

	t.Run("Failure_Invalid_Invalid_Page", func(t *testing.T) {
		startDateStr := "2023-01-01T00:00:00Z"
		endDateStr := "2023-12-31T23:59:59Z"
		pageStr := "invalid-page"
		limitStr := "10"

		result, err := useCase.GetPVZInfo(ctx, startDateStr, endDateStr, pageStr, limitStr)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrGetPVZInfo, err)
		assert.Nil(t, result)
	})

	t.Run("Failure_Invalid_Invalid_Limit", func(t *testing.T) {
		startDateStr := "2023-01-01T00:00:00Z"
		endDateStr := "2023-12-31T23:59:59Z"
		pageStr := "1"
		limitStr := "invalid-limit"

		result, err := useCase.GetPVZInfo(ctx, startDateStr, endDateStr, pageStr, limitStr)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrGetPVZInfo, err)
		assert.Nil(t, result)
	})
}
